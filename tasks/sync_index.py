#!/usr/bin/env python3
"""Validate stable task contracts and derive their read-only JSON index."""
import argparse
from datetime import date
import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parent
REPO = ROOT.parent
TIERS = {
    "T0-bug-fix": "Bug fixes",
    "T1-feature": "Features",
    "T2-quality": "Quality",
    "T3-infrastructure": "Infrastructure",
}
STATES = ("backlog", "ready", "in-progress", "blocked", "done", "obsoleted")
PRIORITIES = ("high", "medium", "low")


def section(text, name):
    match = re.search(r"^## " + re.escape(name) + r"\s*\n(.*?)(?=^## |\Z)", text, re.M | re.S)
    return match.group(1).strip() if match else ""


def scan(root=ROOT):
    records, errors = [], []
    repo = root.parent
    for path in sorted(root.rglob("*.md")):
        rel = path.relative_to(root)
        if rel.parts[0] not in TIERS:
            if len(rel.parts) > 1 or path.name not in ("README.md", "template.md", "migration.md"):
                errors.append(f"{rel}: unknown task category or unexpected root task")
            continue
        label = str(path.relative_to(repo))
        if len(rel.parts) != 2:
            errors.append(f"{label}: task paths must be stable category/slug.md")
        if not re.fullmatch(r"[a-z0-9]+(?:-[a-z0-9]+)*", path.stem):
            errors.append(f"{label}: invalid filename slug")
        text = path.read_text(encoding="utf-8")
        pairs = re.findall(r"^\*\*([\w-]+):\*\*[ \t]*(.*)$", text, re.M)
        fields = {}
        for key, value in pairs:
            key = key.lower()
            if key in fields:
                errors.append(f"{label}: duplicate field {key}")
            fields[key] = value.strip()
        for key in ("id", "status", "created", "priority", "depends-on"):
            if not fields.get(key):
                errors.append(f"{label}: missing {key}")
        tid = fields.get("id", "")
        if not re.fullmatch(r"SYN-\d{3,}", tid):
            errors.append(f"{label}: invalid ID {tid!r}")
        status = fields.get("status", "")
        if status not in STATES:
            errors.append(f"{label}: unknown status {status!r}")
        if fields.get("priority") not in PRIORITIES:
            errors.append(f"{label}: unknown priority")
        try:
            date.fromisoformat(fields.get("created", ""))
        except ValueError:
            errors.append(f"{label}: invalid created date")
        title = re.search(r"^# (.+)$", text, re.M)
        if not title:
            errors.append(f"{label}: missing task title")
        for name in ("Goal", "Acceptance", "Verification"):
            if not section(text, name):
                errors.append(f"{label}: missing {name} content")
        criteria = re.findall(r"^- \[([ xX])\] (.+)$", section(text, "Acceptance"), re.M)
        if not criteria:
            errors.append(f"{label}: acceptance needs observable checklist items")
        if status in ("ready", "in-progress", "blocked", "done"):
            if re.search(r"\b(TODO|TBD|placeholder)\b", section(text, "Goal") + section(text, "Acceptance"), re.I):
                errors.append(f"{label}: unresolved contract placeholders")
        execution = section(text, "Execution")
        for key in (("next",) if status == "in-progress" else ("next", "blocker", "resume-when") if status == "blocked" else ()):
            if not fields.get(key) or not execution:
                errors.append(f"{label}: {status} requires Execution and {key}")
        if status == "done":
            if any(mark == " " for mark, _ in criteria):
                errors.append(f"{label}: done has unchecked acceptance")
            if not fields.get("evidence") or not section(text, "Verification"):
                errors.append(f"{label}: done requires concrete Evidence in Verification")
            elif "**Evidence:**" not in section(text, "Verification"):
                errors.append(f"{label}: Evidence belongs in Verification")
        if status == "obsoleted" and not fields.get("reason"):
            errors.append(f"{label}: obsoleted requires Reason")
        deps = [v.strip() for v in fields.get("depends-on", "-").split(",") if v.strip() and v.strip() != "-"]
        if len(deps) != len(set(deps)):
            errors.append(f"{label}: duplicate dependency")
        records.append({"id": tid, "name": path.stem, "title": title.group(1) if title else path.stem,
                        "file": label, "category": rel.parts[0], "status": status,
                        "priority": fields.get("priority", ""), "created": fields.get("created", ""),
                        "depends_on": deps, "description": " ".join(section(text, "Goal").split())})
    by_id = {}
    for task in records:
        if task["id"] in by_id:
            errors.append(f"duplicate ID {task['id']}: {task['file']}")
        by_id[task["id"]] = task
    for task in records:
        for dep in task["depends_on"]:
            if dep not in by_id:
                errors.append(f"{task['id']}: unknown dependency {dep}")
            elif task["status"] in ("in-progress", "done") and by_id[dep]["status"] != "done":
                errors.append(f"{task['id']}: active/completed dependency {dep} is not done")
    visiting, visited = set(), set()
    def visit(tid):
        if tid in visiting:
            errors.append(f"dependency cycle through {tid}")
            return
        if tid in visited or tid not in by_id:
            return
        visiting.add(tid)
        for dep in by_id[tid]["depends_on"]:
            visit(dep)
        visiting.remove(tid)
        visited.add(tid)
    for tid in by_id:
        visit(tid)
    return records, errors


def build_index(records):
    return {"schema_version": 1, "tiers": {
        tier: {"label": label, "tasks": {status: sorted(
            [t for t in records if t["category"] == tier and t["status"] == status],
            key=lambda t: (PRIORITIES.index(t["priority"]), t["id"])) for status in STATES}}
        for tier, label in TIERS.items()}}


def run(command, root=ROOT):
    records, errors = scan(root)
    if errors:
        for error in errors:
            print(f"ERROR: {error}", file=sys.stderr)
        return 1
    expected = build_index(records)
    index_path = root / "index.json"
    if command == "update":
        serialized = json.dumps(expected, ensure_ascii=False, indent=2) + "\n"
        temporary = index_path.with_suffix(".json.tmp")
        temporary.write_text(serialized, encoding="utf-8")
        temporary.replace(index_path)
        print(f"Updated tasks/index.json: {len(records)} tasks")
    elif command == "check":
        try:
            actual = json.loads(index_path.read_text(encoding="utf-8"))
        except (OSError, ValueError):
            print("ERROR: missing/invalid task index; run update", file=sys.stderr)
            return 1
        if actual != expected:
            print("ERROR: task index drift; run update", file=sys.stderr)
            return 1
        print(f"Tasks OK: {len(records)} contracts, metadata and dependencies")
    elif command == "next":
        by_id = {t["id"]: t for t in records}
        ready = sorted([t for t in records if t["status"] == "ready" and all(
            by_id[d]["status"] == "done" for d in t["depends_on"])],
            key=lambda t: (PRIORITIES.index(t["priority"]), t["id"]))
        for task in ready:
            print(f"{task['id']} [{task['priority']}] {task['title']} — {task['file']}")
        if not ready:
            print("No dependency-satisfied ready tasks")
    return 0


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", nargs="?", default="check", choices=("check", "update", "next"))
    return run(parser.parse_args().command)


if __name__ == "__main__":
    sys.exit(main())
