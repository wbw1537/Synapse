"""MQTT lifecycle and typed state for a reference Axon."""
import copy
import json
import logging
import os
import queue
import threading

import paho.mqtt.client as mqtt
from .config import load_config, require, validate_value

log = logging.getLogger(__name__)

class Component:
    def __init__(self, axon, cid):
        self._axon, self.id = axon, cid

    def update(self, value):
        self._axon.update(self.id, value)

class Axon:
    def __init__(self, config_path, *, token=None, host=None, port=None):
        self._state, self._actions = load_config(config_path)
        self._token = token or os.environ.get("SYNAPSE_AUTH_TOKEN")
        require(isinstance(self._token, str) and bool(self._token), "supply token or SYNAPSE_AUTH_TOKEN")
        self.host = host or os.environ.get("SYNAPSE_MQTT_HOST", "localhost")
        require(isinstance(self.host,str) and bool(self.host), "invalid MQTT host")
        self.port = int(port or os.environ.get("SYNAPSE_MQTT_PORT", "1883").lstrip(":"))
        require(1 <= self.port <= 65535, "invalid MQTT port")
        self.components = {cid: Component(self, cid) for cid in self._state["components"]}
        self._lock = threading.RLock()
        self._connected = threading.Event()
        self._stop = threading.Event()
        self._handlers = {}
        self._commands = queue.Queue(maxsize=32)
        self._client = None
        self._threads = []
        self._subscription = None

    @property
    def connected(self):
        return self._connected.is_set()

    def snapshot(self):
        with self._lock:
            return dict(copy.deepcopy(self._state), auth_token=self._token)

    def on_action(self, action_id):
        require(action_id in self._actions, "action must be declared in TOML")
        def bind(handler):
            require(callable(handler), "action handler must be callable")
            self._handlers[action_id] = handler
            return handler
        return bind

    def update(self, component_id, value):
        with self._lock:
            component = self._state["components"][component_id]
            validate_value(component, value)
            if component["type"] == "log_stream":
                if isinstance(value, str):
                    value = component["value"] + ([value] if value else [])
                value = value[-component["max_items"]:]
            previous = component["value"]
            component["value"] = copy.deepcopy(value)
            try:
                self.publish()
            except ValueError:
                component["value"] = previous
                raise

    def set_status(self, status, message=""):
        require(status in ("online", "warning", "error", "offline"), "invalid status")
        require(isinstance(message, str), "message must be string")
        with self._lock:
            self._state["status"], self._state["message"] = status, message
            self.publish()

    def publish(self):
        with self._lock:
            payload = json.dumps(self.snapshot(), allow_nan=False)
            require(len(payload.encode()) <= 1024*1024, "snapshot exceeds 1 MiB")
            if self._client is None or not self.connected:
                return None
            return self._client.publish(f"synapse/v1/discovery/{self._state['id']}", payload, qos=1)

    def _on_connect(self, client, userdata, flags, reason_code, properties):
        self._connected.clear()
        if reason_code.is_failure:
            log.warning("Axon connection refused")
            return
        rc, self._subscription = client.subscribe(f"synapse/v1/command/{self._state['id']}", qos=1)
        if rc != mqtt.MQTT_ERR_SUCCESS:
            log.warning("Axon command subscription failed")

    def _on_subscribe(self, client, userdata, mid, reasons, properties):
        if mid != self._subscription or self._stop.is_set():
            return
        if not reasons or any(reason.is_failure for reason in reasons):
            log.warning("Axon command subscription refused")
            return
        self._connected.set()
        self.publish()

    def _on_disconnect(self, client, userdata, flags, reason_code, properties):
        self._connected.clear()

    def _on_message(self, client, userdata, message):
        if message.topic != f"synapse/v1/command/{self._state['id']}" or len(message.payload) > 4096:
            return
        try:
            command = json.loads(message.payload)
            if not isinstance(command, dict):
                return
            action_id = command.get("action_id")
            if not isinstance(action_id, str) or action_id not in self._handlers:
                return
            self._commands.put_nowait(action_id)
        except (ValueError, queue.Full):
            log.warning("Invalid command or full action queue")

    def _heartbeat(self):
        while not self._stop.wait(self._state["ttl"] / 2):
            try:
                self.publish()
            except (ValueError, RuntimeError):
                log.warning("Heartbeat publication failed")

    def _dispatch(self):
        while not self._stop.is_set():
            try:
                action_id = self._commands.get(timeout=0.2)
            except queue.Empty:
                continue
            try:
                self._handlers[action_id]()
            except Exception:
                # Handler exceptions can contain credentials; do not log their text.
                log.warning("Action handler failed")
            finally:
                self._commands.task_done()

    def start(self, timeout=10):
        require(self._client is None, "Axon already started")
        require(not any(thread.is_alive() for thread in self._threads), "previous action handler is still running")
        self.publish()
        require(self._actions <= self._handlers.keys(), "bind all declared action handlers before start")
        with self._lock:
            if self._state["status"] == "offline":
                self._state["status"] = "online"
        self._stop.clear()
        while not self._commands.empty():
            self._commands.get_nowait()
            self._commands.task_done()
        self._client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2, client_id=self._state["id"])
        self._client.username_pw_set("axon", self._token)
        self._client.on_connect = self._on_connect
        self._client.on_subscribe = self._on_subscribe
        self._client.on_disconnect = self._on_disconnect
        self._client.on_message = self._on_message
        self._client.reconnect_delay_set(min_delay=1, max_delay=10)
        self._client.max_queued_messages_set(32)
        try:
            self._client.connect(self.host, self.port, keepalive=max(5, self._state["ttl"]))
            self._client.loop_start()
            if not self._connected.wait(timeout):
                raise TimeoutError("Axon registration subscription timed out")
            self._threads = [threading.Thread(target=target, daemon=True) for target in (self._heartbeat, self._dispatch)]
            for thread in self._threads:
                thread.start()
        except Exception:
            self.stop(offline=False)
            raise
        return self

    def stop(self, *, offline=True):
        self._stop.set()
        if self._client is None:
            return
        if offline and self.connected:
            with self._lock:
                self._state["status"] = "offline"
            try:
                publication = self.publish()
                if publication is not None:
                    publication.wait_for_publish(timeout=3)
            except (ValueError, RuntimeError):
                log.warning("Offline publication unavailable; core TTL will expire")
        self._client.disconnect()
        self._client.loop_stop()
        self._connected.clear()
        for thread in self._threads:
            if thread is not threading.current_thread():
                thread.join(timeout=3)
        self._threads = [thread for thread in self._threads if thread.is_alive()]
        self._client = None

    def __enter__(self):
        return self.start()

    def __exit__(self, *args):
        self.stop()
