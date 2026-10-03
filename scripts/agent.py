#!/usr/bin/env python3
"""Synapse workspace checks and proportional build verification."""
import argparse
import json
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "tasks"))
import sync_index


def links(path, root=ROOT):
    root = root.resolve()
    public = root / "website/docs"
    text = path.read_text(encoding="utf-8")
    text = re.sub(r"^```.*?^```[^\n]*", "", text, flags=re.M | re.S)
    errors = []
    for target in re.findall(r"!?\[[^\]]*\]\(([^)]+)\)", text):
        target = target.strip().split(" ", 1)[0].strip("<>")
        if not target or target.startswith("#") or re.match(r"[a-zA-Z][\w+.-]*:", target):
            continue
        file = target.split("#", 1)[0]
        resolved = (root / file.lstrip("/") if file.startswith("/") else path.parent / file).resolve()
        if not resolved.is_relative_to(root) or not resolved.exists():
            errors.append(f"{path.relative_to(root)}: broken/outside link {target}")
        elif path.is_relative_to(public) and not resolved.is_relative_to(public):
            errors.append(f"{path.relative_to(root)}: public link leaves website/docs: {target}")
    return errors


def check():
    result = sync_index.run("check")
    errors, maintained = [], [ROOT / "AGENTS.md", ROOT / "AGENT.md", ROOT / "README.md"]
    try:
        catalog = json.loads((ROOT / "docs/index.json").read_text())
        seen = set()
        for entry in catalog["documents"]:
            path = (ROOT / entry["file"]).resolve()
            if entry["file"] in seen:
                errors.append(f"duplicate document: {entry['file']}")
            seen.add(entry["file"])
            if entry["status"] not in ("current", "proposed", "historical"):
                errors.append(f"unknown document status: {entry['file']}")
            if not path.is_relative_to(ROOT) or not path.is_file():
                errors.append(f"missing/outside document: {entry['file']}")
            elif entry["status"] == "current":
                maintained.append(path)
            if not entry.get("summary"):
                errors.append(f"missing document summary: {entry['file']}")
        for directory in (ROOT / "docs", ROOT / "website/docs"):
            for path in directory.rglob("*.md"):
                if str(path.relative_to(ROOT)) not in seen:
                    errors.append(f"uncatalogued document: {path.relative_to(ROOT)}")
                if path.is_relative_to(ROOT / "website/docs"):
                    maintained.append(path)
    except (OSError, ValueError, KeyError, TypeError) as exc:
        errors.append(f"invalid docs/index.json: {exc}")
    maintained += list((ROOT / "tasks").rglob("*.md"))
    for directory in sorted((ROOT / "skill").iterdir()):
        if not directory.is_dir():
            continue
        skill = directory / "SKILL.md"
        if not skill.is_file():
            errors.append(f"missing skill entrypoint: {directory.name}")
            continue
        text = skill.read_text()
        if not re.match(r"\A---\nname: .+\ndescription: .+\n---\n", text):
            errors.append(f"invalid skill frontmatter: {skill.relative_to(ROOT)}")
        maintained += list(directory.rglob("*.md"))
    for path in set(maintained):
        if path.is_file():
            errors += links(path)
        else:
            errors.append(f"missing maintained file: {path.relative_to(ROOT)}")
    for error in errors:
        print(f"ERROR: {error}", file=sys.stderr)
    if not errors:
        print("Docs/skills OK: catalog, entrypoints and maintained file links")
    return 1 if errors or result else 0


def doctor():
    checks = {"Python 3.10+": sys.version_info >= (3, 10), "Go": bool(shutil.which("go")),
              "Node": bool(shutil.which("node")), "npm": bool(shutil.which("npm")),
              "frontend dependencies": (ROOT / "web/node_modules").is_dir(),
              "frontend entry asset": (ROOT / "web/dist/index.html").is_file()}
    for label, ok in checks.items():
        print(f"{'OK' if ok else 'MISSING'} {label}")
    for name in ("planning", "development", "review", "dev-testing", "operation", "implement-axon"):
        path = ROOT / ".agents/skills" / name / "SKILL.md"
        checks[f"installed skill {name}"] = path.is_file()
        print(f"{'OK' if path.is_file() else 'MISSING'} installed skill {name}")
    print("For dependencies: npm ci --prefix web; for skills: python3 scripts/install_skills.py")
    return 0 if all(checks.values()) else 1


def command(args, cwd=ROOT):
    print("Running:", " ".join(args), flush=True)
    subprocess.run(args, cwd=cwd, check=True)


def verify(scope):
    if scope in ("workspace", "all"):
        if check():
            return 1
        command([sys.executable, "-m", "unittest", "discover", "-s", "scripts/tests", "-v"])
        if scope == "workspace":
            return 0
    if scope in ("web", "all") or not (ROOT / "web/dist/index.html").is_file():
        command(["npm", "run", "build"], ROOT / "web")
    if scope in ("core", "all"):
        command(["go", "vet", "./..."])
        command(["go", "test", "./..."])
    if scope == "all":
        with tempfile.TemporaryDirectory(prefix="synapse-build-") as tmp:
            command(["go", "build", "-o", str(Path(tmp) / "synapse"), "./cmd/synapse"])
    print("Verification passed; build checks do not prove runtime acceptance")
    return 0


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="action", required=True)
    sub.add_parser("check")
    sub.add_parser("doctor")
    v = sub.add_parser("verify")
    v.add_argument("--scope", choices=("workspace", "core", "web", "all"), default="all")
    args = parser.parse_args()
    try:
        return check() if args.action == "check" else doctor() if args.action == "doctor" else verify(args.scope)
    except (OSError, subprocess.CalledProcessError) as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
