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
    
    return {
        "api_version": "v1",
        "auth_token": TOKEN,
        "id": SERVICE_ID,
        "name": "Memory Monitor",
        "group": "System",
        "status": "online",
        "ttl": TTL,
        "widgets": [
            {
                "type": "stat",
                "label": "Memory Usage (%)",
                "value": val,
                "monitors": [
                    {
                        "condition": "value > 90",
                        "severity": "critical",
                        "message": "Memory Critical: >90% usage detected!"
                    }
                ]
            }
        ]
    }

def on_connect(client, userdata, flags, rc, properties=None):
    if rc == 0:
        print("Connected to MQTT Broker!")
    else:
        print(f"Failed to connect, return code {rc}")

# Initialize MQTT Client (using CallbackAPIVersion.VERSION2 for newer paho-mqtt)
try:
    client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
except AttributeError:
    # Fallback for older paho-mqtt versions
    client = mqtt.Client()

client.on_connect = on_connect

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