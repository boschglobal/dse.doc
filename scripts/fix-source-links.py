#!/usr/bin/env python3

# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

import argparse
import html
import re
from pathlib import Path


GITHUB_REPO_BASE = "https://github.com/boschglobal"


def parse_args():
    parser = argparse.ArgumentParser(
        description="Fix Zensical Edit/View links using the generated source map."
    )

    parser.add_argument(
        "site_dir",
        type=Path,
        help="Generated Zensical site directory.",
    )

    parser.add_argument(
        "source_map",
        type=Path,
        help="Source map generated while staging documentation.",
    )

    return parser.parse_args()


def load_source_map(source_map: Path):
    """
    Load mappings of:

        staged_path<TAB>repo<TAB>source_path

    Example:

        user/models/csv/index.md    dse.modelc    doc/content/docs/user/models/csv/index.md
    """

    mappings = {}

    if not source_map.is_file():
        raise FileNotFoundError(
            f"Source map not found: {source_map}"
        )

    for line_number, line in enumerate(
        source_map.read_text(encoding="utf-8").splitlines(),
        start=1,
    ):
        line = line.strip()

        if not line or line.startswith("#"):
            continue

        parts = line.split("\t")

        if len(parts) != 3:
            raise ValueError(
                f"Invalid source-map entry at line {line_number}: {line}"
            )

        staged_path, repo, source_path = parts

        mappings[staged_path.replace("\\", "/")] = (
            repo,
            source_path.replace("\\", "/"),
        )

    return mappings


def github_url(repo: str, source_path: str, action: str) -> str:
    if action == "raw":
        action = "blob"

    return (
        f"https://github.com/"
        f"boschglobal/{repo}/{action}/main/{source_path}"
    )

def normalise_source_path(path: str) -> str:
    """
    Convert a Zensical source URL path into the staged markdown path
    used by source-map.tsv.

    Example:

        user/models/csv/index.md
        user/models/csv/

    both resolve to:

        user/models/csv/index.md
    """

    path = path.lstrip("/")

    if path.endswith("/"):
        path = f"{path}index.md"

    return path


def patch_html_file(html_file: Path, mappings) -> int:
    content = html_file.read_text(encoding="utf-8")
    replacements = 0

    pattern = re.compile(
        r"https://github\.com/boschglobal/dse\.doc/"
        r"(edit|raw)/[^/]+/staged/([^\"'?#\s]+)"
    )

    def replace(match):
        nonlocal replacements

        action = match.group(1)
        staged_path = normalise_source_path(match.group(2))

        mapping = mappings.get(staged_path)

        if mapping is None:
            return match.group(0)

        repo, source_path = mapping

        replacements += 1

        return github_url(
            repo,
            source_path,
            action,
        )

    updated = pattern.sub(replace, content)

    if updated != content:
        html_file.write_text(
            updated,
            encoding="utf-8",
        )

    return replacements


def main():
    args = parse_args()

    mappings = load_source_map(args.source_map)

    html_files = list(
        args.site_dir.rglob("*.html")
    )

    total_replacements = 0

    for html_file in html_files:
        replacements = patch_html_file(
            html_file,
            mappings,
        )

        if replacements:
            print(
                f"Updated {replacements} link(s): "
                f"{html_file}"
            )

        total_replacements += replacements

    print(
        f"Fixed {total_replacements} GitHub Edit/View link(s) "
        f"across {len(html_files)} HTML file(s)."
    )


if __name__ == "__main__":
    main()