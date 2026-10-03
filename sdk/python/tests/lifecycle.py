"""Observable SDK lifecycle against an isolated, already-running core."""
import json
import os
import threading
import time
import urllib.request
from pathlib import Path
import paho.mqtt.client as mqtt
from synapse_axon import Axon

ROOT = Path(__file__).resolve().parents[3]
BASE = os.environ.get("SYNAPSE_HTTP_URL", "http://127.0.0.1:18080")
ADMIN = os.environ["SYNAPSE_ADMIN_TOKEN"]

def request(path, method="GET"):
    req = urllib.request.Request(BASE + path, method=method, headers={"Authorization":"Bearer " + ADMIN})
    with urllib.request.urlopen(req, timeout=3) as response:
        body = response.read()
        return json.loads(body) if body and method == "GET" else None

def wait_for(predicate, timeout=12):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        try:
            result = predicate()
            if result:
                return result
        except (urllib.error.HTTPError, urllib.error.URLError):
            pass
        time.sleep(.1)
    raise AssertionError("Lifecycle condition timed out")

def stored():
    return request('/api/v1/services/sdk-memory')

axon = Axon(ROOT / 'examples/sdk_memory_axon.toml')
handled = threading.Event()
started = threading.Event()
release = threading.Event()
@axon.on_action('recover')
def recover():
    started.set()
    if not release.wait(8):
        return
    axon.components['memory'].update(45)
    axon.components['logs'].update('Recovery callback completed')
    handled.set()

try:
    axon.start()
    wait_for(lambda: stored()['status'] == 'online')
    axon.components['logs'].update('One event')
    axon.components['memory'].update(95)
    state = wait_for(lambda: (s := stored())['components']['memory']['value'] == 95 and s)
    last_seen = state['last_seen']
    wait_for(lambda: stored()['last_seen'] != last_seen)
    assert stored()['components']['logs']['value'] == ['One event'], 'heartbeat duplicated logs'
    # Force an actual broker disconnect using a disposable duplicate identity.
    duplicate = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2, client_id='sdk-memory')
    duplicate.username_pw_set('axon', os.environ['SYNAPSE_AUTH_TOKEN'])
    duplicate.connect(axon.host, axon.port)
    duplicate.loop_start()
    wait_for(lambda: not axon.connected)
    duplicate.disconnect()
    duplicate.loop_stop()
    wait_for(lambda: axon.connected)
    request('/api/v1/services/sdk-memory/actions/recover', 'POST')
    assert started.wait(5), 'callback missing after reconnect/resubscribe'
    heartbeat_before = stored()['last_seen']
    wait_for(lambda: stored()['last_seen'] != heartbeat_before, timeout=5)
    assert not handled.is_set(), 'callback did not wait for release'
    release.set()
    assert handled.wait(5), 'callback failed to complete'
    wait_for(lambda: stored()['components']['memory']['value'] == 45)
    assert stored()['components']['logs']['value'] == ['One event','Recovery callback completed']
    axon.stop()
    wait_for(lambda: stored()['status'] == 'offline')
    assert not axon.connected and not any(t.is_alive() for t in axon._threads), 'shutdown retained worker'
    axon.start()
    wait_for(lambda: stored()['status'] == 'online')
    axon.stop()
    wait_for(lambda: stored()['status'] == 'offline')
    print('SDK lifecycle passed: registration, typed updates, heartbeat, idempotent logs, forced reconnect/resubscribe, actual callback and offline shutdown')
finally:
    release.set()
    axon.stop()
