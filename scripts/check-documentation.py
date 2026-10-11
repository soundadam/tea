#!/usr/bin/env python3
"""Validate local Markdown links."""

from __future__ import annotations

import re
import sys
from pathlib import Path
from urllib.parse import unquote

ROOT = Path(__file__).resolve().parent.parent
MARKDOWN_FILES = sorted(
    path
    for pattern in ("*.md", "*.mdx")
    for path in ROOT.rglob(pattern)
    if ".build" not in path.parts and ".git" not in path.parts
)
LINK_PATTERN = re.compile(r"(?<!!)\[[^\]]+\]\(([^)]+)\)")


def check_local_links() -> list[str]:
    errors: list[str] = []
    for document in MARKDOWN_FILES:
        text = document.read_text(encoding="utf-8")
        for raw_target in LINK_PATTERN.findall(text):
            target = raw_target.strip().split()[0].strip("<>")
            if target.startswith(("http://", "https://", "mailto:", "#", "/")):
                continue
            path_text = unquote(target.split("#", 1)[0].split("?", 1)[0])
            if not path_text:
                continue
            resolved = (document.parent / path_text).resolve()
            if not resolved.exists():
                errors.append(
                    f"{document.relative_to(ROOT)}: missing local link target {target!r}"
                )
    return errors


def main() -> int:
    errors = check_local_links()
    if errors:
        for error in errors:
            print(f"documentation check: {error}", file=sys.stderr)
        return 1
    print(
        f"documentation check passed: {len(MARKDOWN_FILES)} Markdown/MDX files; "
        "local links resolve"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
