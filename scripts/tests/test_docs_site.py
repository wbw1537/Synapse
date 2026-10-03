import json
from pathlib import Path
import sys
import tempfile
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from docs_site import stage

class DocumentationSiteTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        (self.root / "docs").mkdir()
        (self.root / "sdk/python").mkdir(parents=True)
        (self.root / "tasks").mkdir()
        (self.root / "README.md").write_text('[Protocol](docs/protocol.md)')
        (self.root / "docs/protocol.md").write_text('[SDK](../sdk/python/README.md)\n[Home](../README.md)\n[Old](old.md)\n[Task](../tasks/work.md)\n```md\n[Example](does-not-exist.md)\n```\n')
        (self.root / "sdk/python/README.md").write_text('[Protocol](../../docs/protocol.md)')
        (self.root / "docs/old.md").write_text('Historical content')
        (self.root / "tasks/work.md").write_text('Task content')
        (self.root / ".env").write_text('SECRET=must-not-publish')
        self.catalog = {"documents":[
            {"file":"docs/protocol.md","status":"current"},
            {"file":"sdk/python/README.md","status":"current"},
            {"file":"docs/old.md","status":"historical"}]}
        self.save_catalog()

    def save_catalog(self):
        (self.root / "docs/index.json").write_text(json.dumps(self.catalog))

    def test_selected_pages_and_source_links(self):
        output, count = stage(self.root)
        self.assertEqual(count, 3)
        self.assertEqual({str(p.relative_to(output)) for p in output.rglob('*') if p.is_file()},
                         {'index.md','docs/protocol.md','sdk/python/README.md'})
        protocol = (output / 'docs/protocol.md').read_text()
        self.assertIn('[SDK](../sdk/python/README.md)', protocol)
        self.assertIn('[Home](../index.md)', protocol)
        self.assertIn('https://github.com/wbw1537/Synapse/blob/main/docs/old.md', protocol)
        self.assertIn('https://github.com/wbw1537/Synapse/blob/main/tasks/work.md', protocol)
        self.assertIn('[Example](does-not-exist.md)', protocol)
        self.assertNotIn('must-not-publish', ''.join(p.read_text() for p in output.rglob('*.md')))
        (output / 'stale.md').write_text('removed document')
        stage(self.root)
        self.assertFalse((output / 'stale.md').exists())

    def test_invalid_links_fail_without_destroying_previous_build(self):
        output, _ = stage(self.root)
        before = (output / 'docs/protocol.md').read_text()
        (self.root / 'docs/protocol.md').write_text('[Broken](missing.md)')
        with self.assertRaises(ValueError):
            stage(self.root)
        self.assertEqual((output / 'docs/protocol.md').read_text(), before)

    def test_invalid_catalog_paths_cannot_publish_or_overwrite_sources(self):
        for path in ['.env', str(self.root / 'docs/protocol.md'), '../outside.md']:
            with self.subTest(path=path):
                self.catalog['documents'].append({'file':path,'status':'current'})
                self.save_catalog()
                with self.assertRaises(ValueError):
                    stage(self.root)
                self.catalog['documents'].pop()
        (self.root / 'docs/link.md').symlink_to(self.root / 'docs/protocol.md')
        self.catalog['documents'].append({'file':'docs/link.md','status':'current'})
        self.save_catalog()
        with self.assertRaises(ValueError):
            stage(self.root)

if __name__ == '__main__':
    unittest.main()
