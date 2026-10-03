#!/usr/bin/env python3
"""Stage maintained Markdown for MkDocs without copying the repository."""
import argparse
import json
import os
import re
import shutil
from pathlib import Path
from urllib.parse import quote, unquote, urlsplit

ROOT = Path(__file__).resolve().parents[1]
REPOSITORY = "https://github.com/wbw1537/Synapse"
LINK = re.compile(r'(!?\[[^\]\n]*\]\()(<[^>\n]+>|[^\s)]+)(\s+(?:"[^"\n]*"|\x27[^\x27\n]*\x27))?(\))')


def selected_sources(root):
    catalog = json.loads((root / "docs/index.json").read_text())
    sources = {"README.md": "index.md"}
    for entry in catalog["documents"]:
        if entry["status"] == "current":
            relative = entry["file"]
            sources.setdefault(relative, relative)
    for relative in sources:
        relative_path = Path(relative)
        path = root / relative_path
        if (relative_path.is_absolute() or ".." in relative_path.parts or
                path.suffix != ".md" or path.resolve() != path or not path.is_file()):
            raise ValueError(f"invalid published source: {relative}")
    return sources


def rewrite_links(text, source, sources, root):
    def rewrite(match):
        target = match[2].strip("<>")
        url = urlsplit(target)
        if url.scheme or url.netloc or not url.path:
            return match[0]
        path = (root / source).parent / unquote(url.path)
        path = path.resolve()
        if not path.is_relative_to(root) or not path.exists():
            raise ValueError(f"missing/outside link from {source}: {target}")
        relative = path.relative_to(root).as_posix()
        if relative in sources:
            mapped = Path(sources[relative])
            destination = os.path.relpath(mapped, Path(sources[source]).parent).replace(os.sep, "/")
        else:
            kind = "tree" if path.is_dir() else "blob"
            destination = f"{REPOSITORY}/{kind}/main/{quote(relative)}"
        if url.query:
            destination += "?" + url.query
        if url.fragment:
            destination += "#" + url.fragment
        return match[1] + destination + (match[3] or "") + match[4]

    # Example Markdown in fenced blocks is content, not a navigation link.
    output, fence = [], None
    for line in text.splitlines(keepends=True):
        marker = re.match(r"^\s{0,3}(`{3,}|~{3,})", line)
        if marker:
            token = marker[1]
            if fence is None:
                fence = token
            elif token[0] == fence[0] and len(token) >= len(fence):
                fence = None
            output.append(line)
        else:
            output.append(line if fence else LINK.sub(rewrite, line))
    return "".join(output)


def stage(root=ROOT):
    root = root.resolve()
    destination = root / ".build/docs"
    if destination.resolve() != destination:
        raise ValueError("documentation staging path must not contain symlinks")
    sources = selected_sources(root)
    # Validate/render everything before replacing a previous generated tree.
    rendered = {}
    for source, output in sources.items():
        text = rewrite_links((root / source).read_text(), source, sources, root)
        text += f"\n\n---\n[View Markdown source]({REPOSITORY}/blob/main/{quote(source)})\n"
        rendered[output] = text
    if destination.exists():
        shutil.rmtree(destination)
    for output, text in rendered.items():
        path = destination / output
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(text)
    return destination, len(rendered)


def main():
    argparse.ArgumentParser(description=__doc__).parse_args()
    destination, count = stage()
    print(f"Staged {count} maintained Markdown pages in {destination}")


if __name__ == "__main__":
    main()
