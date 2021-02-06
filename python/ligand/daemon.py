import subprocess
import os
import atexit
import tempfile
import functools
import threading
import sys

import grpc
from google.rpc import status_pb2, error_details_pb2

from .servicepb.latch_pb2_grpc import DaemonStub
from .servicepb import latch_pb2 as pb
from . import exceptions

DAEMON_BINARY = os.path.join(os.path.dirname(__file__), "bin/ligand-daemon")


def _is_status_detail(x):
    return hasattr(x, "key") and x.key == "grpc-status-details-bin"


def _start_wrapped_pipe(pipe, writer):
    def wrap_pipe(pipe, writer):
        with pipe:
            for line in iter(pipe.readline, b""):
                writer.write(line)
                writer.flush()

    # if writer is normal sys.std{out,err}, it can't
    # write bytes directly.
    # see https://stackoverflow.com/a/908440/135797
    if hasattr(writer, "buffer"):
        writer = writer.buffer

    thread = threading.Thread(target=wrap_pipe, args=[pipe, writer], daemon=True)
    thread.start()
    return thread


def _handle_exception(code, details):
    if code == "DOES_NOT_EXIST":
        return exceptions.DoesNotExist(details)
    if code == "READ_ERROR":
        return exceptions.ReadError(details)
    if code == "WRITE_ERROR":
        return exceptions.WriteError(details)
    if code == "REPOSITORY_CONFIGURATION_ERROR":
        return exceptions.RepositoryConfigurationError(details)
    if code == "INCOMPATIBLE_REPOSITORY_VERSION":
        return exceptions.IncompatibleRepositoryVersion(details)
    if code == "CORRUPTED_REPOSITORY_SPEC":
        return exceptions.CorruptedRepositorySpec(details)
    if code == "CONFIG_NOT_FOUND":
        return exceptions.ConfigNotFound(details)


def _get_status_code(e, details):
    metadata = e.trailing_metadata()
    status_md = [x for x in metadata if _is_status_detail(x)]
    if status_md:
        for md in status_md:
            st = status_pb2.Status()
            st.MergeFromString(md.value)
        if st.details:
            val = error_details_pb2.ErrorInfo()
            st.details[0].Unpack(val)
            return val.reason
    return None


def handle_error(f):
    @functools.wraps(f)
    def wrapped(*args, **kwargs):
        try:
            return f(*args, **kwargs)
        except grpc.RpcError as e:
            code, name = e.code().value
            details = e.details()
            if name == "internal":
                status_code = _get_status_code(e, details)
                if status_code:
                    raise handle_exception(status_code, details)
            raise Exception(details)

    return wrapped


class Daemon:
    """
    todo... grpc communication process
    """

    def __init__(self, job, socket_path=None):
        self.job = job

        if socket_path is None:
            # create a new temporary file just to get a free name.
            # the Go GRPC server will create the file.
            f = tempfile.NamedTemporaryFile(
                prefix="ligand-daemon-", suffix=".sock", delete=False
            )
            self.socket_path = f.name
            f.close()
        else:
            self.socket_path = socket_path

        # the Go GRPC server will fail to start if the socket file
        # already exists.
        os.unlink(self.socket_path)

        cmd = [DAEMON_BINARY]

        # Init daemon, communicating with flags passed to binary
        # if self.project.repository:
        #     cmd += ["-R", self.project.repository]
        # if self.project.directory:
        #     cmd += ["-D", self.project.directory]
        # if debug:
        #     cmd += ["-v"]
        cmd.append(self.socket_path)
        self.process = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE
        )

        # need to wrap stdout and stderr for this to work in jupyter
        # notebooks. jupyter redefines sys.std{out,err} as custom
        # writers that eventually write the output to the notebook.
        self.stdout_thread = _start_wrapped_pipe(self.process.stdout, sys.stdout)
        self.stderr_thread = _start_wrapped_pipe(self.process.stderr, sys.stderr)

        atexit.register(self.cleanup)
        self.channel = grpc.insecure_channel("unix://" + self.socket_path)
        self.stub = DaemonStub(self.channel)

        TIMEOUT_SEC = 15
        grpc.channel_ready_future(self.channel).result(timeout=TIMEOUT_SEC)

    def cleanup(self):
        if self.process.poll() is None:  # check if process is still running:
            # the sigterm handler in the daemon process waits for any in-progress uploads etc. to finish.
            # the sigterm handler also deletes the socket file
            self.process.terminate()
            self.process.wait()

            # need to join these threads to avoid "could not acquire lock" error
            self.stdout_thread.join()
            self.stderr_thread.join()
        self.channel.close()

    @handle_error
    def ping(self):
        return self.stub.Ping(pb.PingRequest())

    @handle_error
    def launch_job(self, job):
        print(job.script, job.python_packages, job.python_version)
        pb_job = pb.Job(
            script=job.script,
            pythonPackages=job.python_packages,
            pythonVersion=job.python_version,
        )

        return self.stub.LaunchJob(
            pb.LaunchJobRequest(
                job=pb_job,
            ),
        )

    @handle_error
    def experiment_is_running(self, experiment_id: str) -> str:
        ret = self.stub.GetExperimentStatus(
            pb.GetExperimentStatusRequest(experimentID=experiment_id)
        )
        return ret.status == pb.GetExperimentStatusReply.Status.RUNNING
