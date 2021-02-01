import sys
from dataclasses import dataclass, field
from typing import Optional, Dict

from .daemon import Daemon


def _get_imported_packages() -> Dict[str, str]:
    """
    Returns a list of packages that have been imported, as a {name: version} dict.
    """
    try:
        # Should we vendor pkg_resources? See https://github.com/replicate/replicate/issues/350
        import pkg_resources
    except ImportError:
        print("Could not import setuptools/pkg_resources, not tracking package versions")
        # console.warn(
        # "Could not import setuptools/pkg_resources, not tracking package versions"
        # )
        return {}
    result = {}
    for d in pkg_resources.working_set:
        if _is_imported(d.key):
            result[d.key] = d.version
    return result


def _is_imported(module_name):
    return module_name in sys.modules


def _get_python_version():
    return ".".join([str(x) for x in sys.version_info[:3]])


@dataclass
class Provider:
    """
    Cloud compute provider.

    Instantiated from .latchrc.
    """


class Job:
    """
    A unit of computation (ie. script) operating on an arbitrary cluster.
    """

    def __init__(self):
        self.python_version: float = _get_python_version()
        self.python_packages: Optional[Dict[str,
                                            str]] = _get_imported_packages()
        self.provider: Provider = ""
        self.script: str = "/test/"
        self._daemon_instance: Daemon = Daemon(self)

    def _daemon(self) -> Daemon:
        if self._daemon_instance is None:
            self._daemon_instance = Daemon(self)
        return self._daemon_instance

    def ping(self):
        return self._daemon().ping()

    def launch(self):
        return self._daemon().launch_job(self)


def init(
    data: Optional[str] = None,
    out: Optional[str] = None
):
    """
    Create and launch a new job.
    """
    # print(Job().ping())
    print(Job().launch())
