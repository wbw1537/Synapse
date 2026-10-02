"""Behavioral checks for task integrity, installation and viewer exposure."""
from contextlib import redirect_stdout, redirect_stderr
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest

REPO = Path(__file__).resolve().parents[2]


def module(name, path):
    spec = importlib.util.spec_from_file_location(name, path)
    result = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(result)
    return result


tasks = module("workspace_tasks", REPO / "tasks/sync_index.py")
installer = module("workspace_installer", REPO / "scripts/install_skills.py")
viewer = module("workspace_viewer", REPO / "scripts/tasks_server.py")


class TaskTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name) / "tasks"
        (self.root / "T0-bug-fix").mkdir(parents=True)

    def write(self, name="one", tid="SYN-101", state="ready", deps="-", acceptance="[ ]", extra=""):
        path = self.root / "T0-bug-fix" / f"{name}.md"
        path.write_text(f"""# Test outcome

**ID:** {tid}
**Status:** {state}
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** {deps}

## Goal

Keep observable service state correct.

## Acceptance

- {acceptance} An open page receives the persisted state.

## Verification

Compare an HTTP snapshot and open page after stopping heartbeats.
{extra}
""")
        return path

    def run_quiet(self, command):
        with redirect_stdout(io.StringIO()), redirect_stderr(io.StringIO()):
            return tasks.run(command, self.root)

    def test_index_detects_metadata_drift(self):
        path = self.write()
        self.assertEqual(self.run_quiet("update"), 0)
        self.assertEqual(self.run_quiet("check"), 0)
        path.write_text(path.read_text().replace("**Priority:** high", "**Priority:** low"))
        self.assertEqual(self.run_quiet("check"), 1)

    def test_unknown_status_does_not_destroy_valid_index(self):
        path = self.write()
        self.run_quiet("update")
        index = (self.root / "index.json").read_bytes()
        path.write_text(path.read_text().replace("**Status:** ready", "**Status:** implemented"))
        self.assertEqual(self.run_quiet("update"), 1)
        self.assertEqual((self.root / "index.json").read_bytes(), index)

    def test_cycles_unknown_dependencies_and_duplicate_ids_fail(self):
        self.write(deps="SYN-102")
        self.assertTrue(any("unknown dependency" in e for e in tasks.scan(self.root)[1]))
        self.write("two", "SYN-102", deps="SYN-101")
        self.assertTrue(any("cycle" in e for e in tasks.scan(self.root)[1]))
        self.write("two", "SYN-101")
        self.assertTrue(any("duplicate ID" in e for e in tasks.scan(self.root)[1]))

    def test_done_requires_checked_acceptance_and_evidence(self):
        self.write(state="done")
        self.assertTrue(any("unchecked acceptance" in e for e in tasks.scan(self.root)[1]))
        self.write(state="done", acceptance="[x]")
        self.assertTrue(any("concrete Evidence" in e for e in tasks.scan(self.root)[1]))
        self.write(state="done", acceptance="[x]", extra="**Evidence:** Runtime comparison passed on a recorded revision.")
        self.assertEqual(tasks.scan(self.root)[1], [])

    def test_blocked_requires_recovery_condition(self):
        self.write(state="blocked")
        self.assertTrue(any("resume-when" in e for e in tasks.scan(self.root)[1]))
        self.write(state="blocked", extra="\n## Execution\n\n**Next:** Run broker smoke.\n**Blocker:** Required environment unavailable.\n**Resume-when:** Disposable broker is reachable.")
        self.assertEqual(tasks.scan(self.root)[1], [])

    def test_misplaced_task_is_not_silently_omitted(self):
        path = self.write()
        path.rename(self.root / "one.md")
        self.assertTrue(any("unexpected root task" in e for e in tasks.scan(self.root)[1]))

    def test_next_excludes_unsatisfied_dependencies(self):
        self.write(state="backlog")
        self.write("two", "SYN-102", deps="SYN-101")
        output = io.StringIO()
        with redirect_stdout(output):
            self.assertEqual(tasks.run("next", self.root), 0)
        self.assertIn("No dependency-satisfied", output.getvalue())
        self.write(state="done", acceptance="[x]", extra="**Evidence:** Recorded smoke passed.")
        output = io.StringIO()
        with redirect_stdout(output):
            tasks.run("next", self.root)
        self.assertIn("SYN-102", output.getvalue())
        self.assertNotIn("SYN-101 [", output.getvalue())


class InstallerTests(unittest.TestCase):
    def test_install_is_repeatable_and_refuses_conflicts(self):
        with tempfile.TemporaryDirectory() as tmp, redirect_stdout(io.StringIO()):
            base = Path(tmp)
            global_root, links = base / "user-agents", base / "local-skills"
            installer.install(global_root, links)
            installer.install(global_root, links)
            installer.install(global_root, links, checking=True)
            self.assertTrue((links / "planning/SKILL.md").is_file())
            (links / "review").unlink()
            (links / "review").mkdir()
            marker = links / "review/user-file"
            marker.write_text("preserve")
            with self.assertRaises(ValueError):
                installer.install(global_root, links)
            self.assertEqual(marker.read_text(), "preserve")

    def test_global_copy_is_distinct_and_tampering_detected(self):
        with tempfile.TemporaryDirectory() as tmp, redirect_stdout(io.StringIO()):
            base = Path(tmp)
            installer.install(base / "agents", base / "links")
            copied = (base / "links/planning/SKILL.md").resolve()
            self.assertFalse(copied.is_relative_to(REPO))
            copied.write_text("changed")
            with self.assertRaises(ValueError):
                installer.install(base / "agents", base / "links", checking=True)
            with self.assertRaises(ValueError):
                installer.install(base / "agents", base / "links")


class ViewerTests(unittest.TestCase):
    def request(self, path, body=True):
        # Exercise the actual route/file policy without requiring a listening socket.
        handler = object.__new__(viewer.Handler)
        handler.path = path
        handler.wfile = io.BytesIO()
        result = {"status": None, "headers": {}}
        handler.send_error = lambda code: result.update(status=code)
        handler.send_response = lambda code: result.update(status=code)
        handler.send_header = lambda name, value: result["headers"].update({name: value})
        handler.end_headers = lambda: None
        handler.serve(body)
        result["body"] = handler.wfile.getvalue()
        return result

    def test_allowed_task_data_and_head(self):
        response = self.request("/tasks/index.json?ts=123")
        self.assertEqual(response["status"], 200)
        self.assertEqual(json.loads(response["body"])["schema_version"], 1)
        self.assertEqual(self.request("/")["status"], 200)
        self.assertEqual(self.request("/docs/protocol.md")["status"], 200)
        head = self.request("/tasks/index.json", False)
        self.assertEqual(head["status"], 200)
        self.assertEqual(head["body"], b"")

    def test_unrelated_data_and_traversal_denied(self):
        for path in ("/.env", "/synapse.db", "/internal/models/service.go", "/.git/config",
                     "/tasks/../.env", "/tasks/%2e%2e/.env", "/tasks/sync_index.py"):
            with self.subTest(path=path):
                self.assertEqual(self.request(path)["status"], 404)


if __name__ == "__main__":
    unittest.main()
