import copy
import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch
from synapse_axon import Axon

ROOT = Path(__file__).resolve().parents[3]
DEMO = ROOT / "examples/sdk_memory_axon.toml"

class ConfigurationTests(unittest.TestCase):
    def create(self, content):
        temp = tempfile.TemporaryDirectory()
        self.addCleanup(temp.cleanup)
        path = Path(temp.name) / "axon.toml"
        path.write_text(content)
        return path

    def test_invalid_configuration_fails_before_client_creation(self):
        source = DEMO.read_text()
        invalid = [
            source.replace('axon.card.v1', 'axon.card.v2'),
            source.replace('["memory", "logs", "controls"]', '["ghost"]'),
            source.replace('ttl = 6', 'ttl = 0'),
            source.replace('max = 100', 'max = "bad"'),
            source.replace('max_items = 5', 'max_items = -1'),
            source.replace('default = 45', 'default = "bad"'),
            source.replace('id = "recover"', 'id = "bad/id"'),
            source + '\nunknown = "field"\n',
            source + '\n[[components.controls.items]]\nid = "recover"\nlabel = "Duplicate"\n',
        ]
        for text in invalid:
            with self.subTest(text=text[-80:]), patch('synapse_axon.client.mqtt.Client', side_effect=AssertionError('network client created')):
                with self.assertRaises(ValueError):
                    Axon(self.create(text), token="private")

    def test_mapping_and_typed_runtime_updates(self):
        axon = Axon(DEMO, token="private")
        payload = axon.snapshot()
        self.assertEqual(payload['api_version'], 'v1')
        self.assertEqual(payload['components']['controls']['items'][0]['action_id'], 'recover')
        self.assertNotIn('meta', payload)
        for component in payload['components'].values():
            self.assertNotIn('default', component)
        before = copy.deepcopy(payload)
        for value in [True, 'bad', float('nan'), -1, 101]:
            with self.assertRaises(ValueError):
                axon.components['memory'].update(value)
        self.assertEqual(axon.snapshot(), before)
        axon.components['memory'].update(75)
        self.assertEqual(axon.snapshot()['components']['memory']['value'], 75)
        for i in range(7):
            axon.components['logs'].update(str(i))
        self.assertEqual(axon.snapshot()['components']['logs']['value'], ['2','3','4','5','6'])
        axon.components['logs'].update(['snapshot'])
        self.assertEqual(axon.snapshot()['components']['logs']['value'], ['snapshot'])
        with self.assertRaises(ValueError):
            axon.components['logs'].update(['good', 1])
        with self.assertRaises(ValueError):
            axon.on_action('undeclared')
        with patch('synapse_axon.client.mqtt.Client', side_effect=AssertionError('network client created')):
            with self.assertRaises(ValueError):
                axon.start()
        before = axon.snapshot()
        with self.assertRaises(ValueError):
            axon.components['logs'].update('x' * (1024*1024))
        self.assertEqual(axon.snapshot(), before)
        # Heartbeat payloads are full, stable array snapshots.
        self.assertEqual(json.dumps(axon.snapshot()), json.dumps(axon.snapshot()))

if __name__ == '__main__':
    unittest.main()
