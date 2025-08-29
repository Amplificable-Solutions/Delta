import os
import sys
from pathlib import Path

import pytest

from .network import setup_delta, setup_custom_delta, setup_geth

dir = os.path.dirname(os.path.realpath(__file__))
sys.path.append(dir + "/protobuf")


def pytest_configure(config):
    config.addinivalue_line("markers", "unmarked: fallback mark for unmarked tests")
    config.addinivalue_line("markers", "slow: marks tests as slow")
    config.addinivalue_line("markers", "gravity: gravity bridge test cases")
    config.addinivalue_line("markers", "ica: marks ica tests")
    config.addinivalue_line("markers", "upgrade: marks upgrade tests")
    config.addinivalue_line("markers", "ibc: marks default ibc tests")
    config.addinivalue_line("markers", "ibc_rly_evm: marks ibc_rly_evm tests")
    config.addinivalue_line("markers", "ibc_rly_gas: marks ibc relayer gas tests")
    config.addinivalue_line("markers", "ibc_timeout: marks ibc timeout tests")
    config.addinivalue_line("markers", "ibc_update_client: marks ibc updateclient test")
    config.addinivalue_line("markers", "gov: marks gov related tests")
    config.addinivalue_line("markers", "gas: marks gas related tests")


def pytest_collection_modifyitems(items, config):
    for item in items:
        if not any(item.iter_markers()):
            item.add_marker("unmarked")


@pytest.fixture(scope="session")
def suspend_capture(pytestconfig):
    """
    used to pause in testing

    Example:
    ```
    def test_simple(suspend_capture):
        with suspend_capture:
            # read user input
            print(input())
    ```
    """

    class SuspendGuard:
        def __init__(self):
            self.capmanager = pytestconfig.pluginmanager.getplugin("capturemanager")

        def __enter__(self):
            self.capmanager.suspend_global_capture(in_=True)

        def __exit__(self, _1, _2, _3):
            self.capmanager.resume_global_capture()

    yield SuspendGuard()


@pytest.fixture(scope="session", params=[True])
def delta(request, tmp_path_factory):
    enable_indexer = request.param
    if enable_indexer:
        path = tmp_path_factory.mktemp("indexer")
        yield from setup_custom_delta(
            path, 27000, Path(__file__).parent / "configs/enable-indexer.jsonnet"
        )
    else:
        path = tmp_path_factory.mktemp("delta")
        yield from setup_delta(path, 26650)


@pytest.fixture(scope="session")
def geth(tmp_path_factory):
    path = tmp_path_factory.mktemp("geth")
    yield from setup_geth(path, 8545)


@pytest.fixture(scope="session", params=["delta", "geth", "delta-ws"])
def cluster(request, delta, geth):
    """
    run on both delta and geth
    """
    provider = request.param
    if provider == "delta":
        yield delta
    elif provider == "geth":
        yield geth
    elif provider == "delta-ws":
        delta_ws = delta.copy()
        delta_ws.use_websocket()
        yield delta_ws
    else:
        raise NotImplementedError
