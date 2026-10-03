"""Run after installing sdk/python; all measurements and actions are simulated."""
import time
from pathlib import Path
from synapse_axon import Axon

axon = Axon(Path(__file__).with_suffix(".toml"))

@axon.on_action("recover")
def recover():
    axon.components["memory"].update(45)
    axon.components["logs"].update("Simulated recovery callback completed")
    print("Simulated recovery callback completed", flush=True)

try:
    with axon:
        axon.components["logs"].update("Demo started")
        time.sleep(5)
        axon.components["memory"].update(95)
        axon.components["logs"].update("Simulated high memory; click recovery")
        while True:
            time.sleep(1)
except KeyboardInterrupt:
    pass
