"""
@generated by mypy-protobuf.  Do not edit manually!
isort:skip_file
"""
import builtins
import google.protobuf.descriptor
import google.protobuf.internal.containers
import google.protobuf.message
import typing
import typing_extensions

DESCRIPTOR: google.protobuf.descriptor.FileDescriptor = ...

class PingRequest(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor = ...

    def __init__(self,
        ) -> None: ...
global___PingRequest = PingRequest

class PingReply(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor = ...
    success: builtins.bool = ...

    def __init__(self,
        *,
        success : builtins.bool = ...,
        ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal[u"success",b"success"]) -> None: ...
global___PingReply = PingReply

class LaunchJobRequest(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor = ...

    @property
    def job(self) -> global___Job: ...

    def __init__(self,
        *,
        job : typing.Optional[global___Job] = ...,
        ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal[u"job",b"job"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal[u"job",b"job"]) -> None: ...
global___LaunchJobRequest = LaunchJobRequest

class LaunchJobReply(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor = ...
    success: builtins.bool = ...

    def __init__(self,
        *,
        success : builtins.bool = ...,
        ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal[u"success",b"success"]) -> None: ...
global___LaunchJobReply = LaunchJobReply

class Job(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor = ...
    class PythonPackagesEntry(google.protobuf.message.Message):
        DESCRIPTOR: google.protobuf.descriptor.Descriptor = ...
        key: typing.Text = ...
        value: typing.Text = ...

        def __init__(self,
            *,
            key : typing.Text = ...,
            value : typing.Text = ...,
            ) -> None: ...
        def ClearField(self, field_name: typing_extensions.Literal[u"key",b"key",u"value",b"value"]) -> None: ...

    script: typing.Text = ...
    pythonVersion: typing.Text = ...

    @property
    def pythonPackages(self) -> google.protobuf.internal.containers.ScalarMap[typing.Text, typing.Text]: ...

    def __init__(self,
        *,
        script : typing.Text = ...,
        pythonPackages : typing.Optional[typing.Mapping[typing.Text, typing.Text]] = ...,
        pythonVersion : typing.Text = ...,
        ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal[u"pythonPackages",b"pythonPackages",u"pythonVersion",b"pythonVersion",u"script",b"script"]) -> None: ...
global___Job = Job
