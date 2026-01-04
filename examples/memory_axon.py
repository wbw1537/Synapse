import time
import json
import psutil
import paho.mqtt.client as mqtt

# Configuration
BROKER = "localhost"
PORT = 1883
SERVICE_ID = "memory-sidecar-python"
TOKEN = "synapse-secret"
TTL = 10

def get_payload(simulated_percent=None):
    # Get memory stats
    mem = psutil.virtual_memory()
    
    # Use simulated value if provided, otherwise use real value
    if simulated_percent is not None:
        val = simulated_percent
    else:
        val = mem.percent
    
    # Derive status for status_indicator
    status_key = "normal"
    if val > 90:
        status_key = "critical"
    elif val > 70:
        status_key = "high"

    return {
        "api_version": "v1",
        "auth_token": TOKEN,
        "id": SERVICE_ID,
        "name": "Memory Monitor",
        "group": "System",
        "status": "online",
        "ttl": TTL,
        "markdown_docs": """# Memory Monitor Runbook

This Axon monitors the system memory usage and reports it to Synapse.

## Troubleshooting
If you receive a **Critical** alert (>90% usage):
1. Check for memory-hungry processes: `top -o %MEM`
2. Identify if it's a leak or expected load.
3. Use the **Drop Cache** action below if the system is sluggish.

## Data Sources
- **psutil**: Used for cross-platform memory statistics.
- **Log Stream**: Captures every sample event for audit.
""",
        "widgets": [
            # 1. Stat (Existing logic preserved)
            {
                "id": "mem_stat",
                "type": "stat",
                "label": "Memory Usage",
                "value": val,
                "unit": "%",
                "copyable": True,
                "monitors": [
                    {
                        "condition": "value > 90",
                        "severity": "critical",
                        "message": "Memory Critical: >90% usage detected!"
                    }
                ]
            },
            # 2. Gauge
            {
                "id": "mem_gauge",
                "type": "gauge",
                "label": "Load Visual",
                "value": val,
                "min": 0,
                "max": 100,
                "unit": "%",
                "thresholds": {
                    "70": "warning",
                    "90": "danger"
                }
            },
            # 3. Status Indicator
            {
                "id": "health_status",
                "type": "status_indicator",
                "label": "Health State",
                "value": status_key,
                "mapping": {
                    "normal": {"text": "Healthy", "color": "success", "icon": "check"},
                    "high": {"text": "High Load", "color": "warning", "icon": "activity"},
                    "critical": {"text": "Critical", "color": "danger", "icon": "alert", "animate": True}
                }
            },
            # 4. Log Stream
            {
                "id": "mem_logs",
                "type": "log_stream",
                "label": "Activity Log",
                "value": f"[{time.strftime('%H:%M:%S')}] Sampled memory usage at {val}%",
                "max_items": 5
            },
            # 5. Action Group
            {
                "id": "ops_actions",
                "type": "action_group",
                "label": "Operations",
                "items": [
                    {
                        "action_id": "drop_cache",
                        "label": "Drop Cache",
                        "style": "primary",
                        "confirm": False
                    },
                    {
                        "action_id": "restart_process",
                        "label": "Restart",
                        "style": "danger",
                        "confirm": True
                    }
                ]
            },
            # 6. Link
            {
                "id": "docs_link",
                "type": "link",
                "label": "Documentation",
                "text": "Wiki",
                "uri": "https://en.wikipedia.org/wiki/Random-access_memory"
            }
        ]
    }

def on_connect(client, userdata, flags, rc, properties=None):
    if rc == 0:
        print("Connected to MQTT Broker!")
        # Subscribe to command topic for this service
        command_topic = f"synapse/v1/command/{SERVICE_ID}"
        client.subscribe(command_topic)
        print(f"Subscribed to {command_topic}")
    else:
        print(f"Failed to connect, return code {rc}")

def on_message(client, userdata, msg):
    try:
        print(f"Received command on {msg.topic}")
        payload = json.loads(msg.payload.decode())
        action_id = payload.get("action_id")
        
        if action_id == "drop_cache":
            print(">>> ACTION: Dropping Cache... (Simulated)")
            # Simulated effect: reduce memory usage
            # In a real scenario, this would run `sync; echo 3 > /proc/sys/vm/drop_caches`
        elif action_id == "restart_process":
            print(">>> ACTION: Restarting Process... (Simulated)")
        else:
            print(f">>> WARNING: Unknown action '{action_id}' requested.")
            
    except Exception as e:
        print(f"Error handling message: {e}")

# Initialize MQTT Client (using CallbackAPIVersion.VERSION2 for newer paho-mqtt)
try:
    client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
except AttributeError:
    # Fallback for older paho-mqtt versions
    client = mqtt.Client()

client.on_connect = on_connect
client.on_message = on_message

print(f"Connecting to {BROKER}:{PORT}...")
try:
    client.connect(BROKER, PORT, 60)
except Exception as e:
    print(f"Connection error: {e}")
    exit(1)

# Start loop
client.loop_start()

start_time = time.time()

try:
    print(f"Starting {SERVICE_ID}...")
    print("Scenario: 10s Normal -> 5s Critical -> Normal")
    
    while True:
        elapsed = time.time() - start_time
        
        # Determine value based on elapsed time
        if elapsed < 10:
            current_val = 45 # Normal
            phase = "Normal (Initial)"
        elif 10 <= elapsed < 15:
            current_val = 95 # Critical
            phase = "CRITICAL TEST"
        else:
            current_val = 50 # Normal
            phase = "Normal (Resumed)"
            
        payload = get_payload(simulated_percent=current_val)
        topic = f"synapse/v1/discovery/{SERVICE_ID}"
        
        print(f"[{time.strftime('%H:%M:%S')}] Phase: {phase} | Publishing Value: {payload['widgets'][0]['value']}")
        client.publish(topic, json.dumps(payload))
        
        # Sleep for half of TTL
        time.sleep(TTL / 2)
        
except KeyboardInterrupt:
    print("\nShutting down...")
    # Send offline status
    payload = get_payload()
    payload["status"] = "offline"
    client.publish(f"synapse/v1/discovery/{SERVICE_ID}", json.dumps(payload))
    client.loop_stop()
    client.disconnect()