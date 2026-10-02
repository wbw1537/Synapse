#!/usr/bin/env python3
"""Install versioned portable skills and link local project skills without clobbering files."""
import argparse
import hashlib
from pathlib import Path
import shutil
import sys

ROOT = Path(__file__).resolve().parents[1]
SHARED = {"planning": "project-planner", "development": "agent-developer", "review": "review"}
LOCAL = ("dev-testing", "operation", "implement-axon")


def digest():
    sha = hashlib.sha256()
    for name in SHARED:
        for file in sorted((ROOT / "skill" / name).rglob("*")):
            if file.is_file():
                sha.update(str(file.relative_to(ROOT / "skill")).encode())
                sha.update(b"\0" + file.read_bytes() + b"\0")
    return sha.hexdigest()[:16]


def safe_link(path, target, permitted):
    if path.is_symlink():
        existing = path.resolve()
        if existing != target.resolve() and not any(existing.is_relative_to(base.resolve()) for base in permitted):
            raise ValueError(f"refusing unrelated link: {path}")
    elif path.exists():
        raise ValueError(f"refusing existing non-symlink: {path}")


def install(global_root, links_root, checking=False):
    bundle_root = global_root / "bundles/synapse-workflow"
    bundle = bundle_root / digest()
    pairs = [(links_root / name, bundle / name) for name in SHARED]
    pairs += [(global_root / "skills" / skill_name, bundle / folder) for folder, skill_name in SHARED.items()]
    pairs += [(links_root / name, ROOT / "skill" / name) for name in LOCAL]
    for link, target in pairs:
        safe_link(link, target, [bundle_root, ROOT / "skill"])
    if checking:
        for name in SHARED:
            source = ROOT / "skill" / name
            dest = bundle / name
            actual = {str(p.relative_to(dest)): p.read_bytes() for p in dest.rglob("*") if p.is_file()}
            expected = {str(p.relative_to(source)): p.read_bytes() for p in source.rglob("*") if p.is_file()}
            if expected != actual:
                raise ValueError(f"bundle content differs: {name}; reinstall")
        for link, target in pairs:
            if not link.is_symlink() or link.resolve() != target.resolve() or not target.is_dir():
                raise ValueError(f"missing/stale link: {link}; reinstall")
        print(f"Skills OK: pinned bundle {bundle.name}, {len(pairs)} links")
        return
    # All destination conflicts are checked before any links are changed.
    for name in SHARED:
        source, dest = ROOT / "skill" / name, bundle / name
        if dest.exists():
            actual = {str(p.relative_to(dest)): p.read_bytes() for p in dest.rglob("*") if p.is_file()}
            expected = {str(p.relative_to(source)): p.read_bytes() for p in source.rglob("*") if p.is_file()}
            if actual != expected:
                raise ValueError(f"immutable bundle modified: {dest}")
        else:
            dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.copytree(source, dest)
    for link, target in pairs:
        link.parent.mkdir(parents=True, exist_ok=True)
        if link.is_symlink() and link.resolve() == target.resolve():
            continue
        if link.is_symlink():
            link.unlink()
        link.symlink_to(target)
    print(f"Installed shared bundle: {bundle}")
    print(f"Linked {len(SHARED) + len(LOCAL)} skills in {links_root}")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--global-root", type=Path, default=Path.home() / ".agents")
    parser.add_argument("--links-root", type=Path, default=ROOT / ".agents/skills")
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    try:
        install(args.global_root.expanduser().absolute(), args.links_root.absolute(), args.check)
    except (OSError, ValueError) as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
