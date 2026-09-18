#!/usr/bin/env python3

# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

import argparse
import re
from pathlib import Path


# Repository order used for generated navigation weights.
DEFAULT_REPO_ORDER = [
    "dse.clib",
    "dse.fmi",
    "dse.modelc",
    "dse.network",
    "dse.schemas",
    "dse.sdp",
    "dse.standards",
]


REPO_WEIGHT_STEP = 1000
FILE_WEIGHT_STEP = 10


def parse_args():
    parser = argparse.ArgumentParser(
        description="Generate missing Markdown navigation metadata."
    )

    parser.add_argument(
        "--repos-dir",
        "-r",
        default="build/repos",
        help="Directory containing copied repositories",
    )

    parser.add_argument(
        "--repo-order",
        default=",".join(DEFAULT_REPO_ORDER),
        help="Comma-separated repository order",
    )

    return parser.parse_args()


def read_text(path: Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        return path.read_text(
            encoding="utf-8",
            errors="ignore",
        )


def get_front_matter_end(lines):
    """
    Return the line index of the closing front-matter delimiter.

    Example:

        ---
        title: Example
        ---
        Content

    Returns the index of the second '---'.
    """
    if not lines or lines[0].strip() != "---":
        return None

    for index in range(1, len(lines)):
        if lines[index].strip() == "---":
            return index

    return None


def has_weight(lines, end_index):
    if end_index is None:
        return False

    for line in lines[1:end_index]:
        if re.match(
            r"^\s*weight\s*:",
            line,
            re.IGNORECASE,
        ):
            return True

    return False


def add_weight(text: str, weight: int) -> str:
    """
    Add a generated weight.

    If front matter does not exist, create it.

    If front matter exists, insert the weight before
    the closing delimiter.
    """
    lines = text.splitlines(keepends=True)

    end_index = get_front_matter_end(lines)

    # No front matter: create it.
    if end_index is None:
        return (
            "---\n"
            f"weight: {weight}\n"
            "---\n\n"
            + text
        )

    # Preserve the existing line-ending style.
    newline = "\n"

    if lines[end_index].endswith("\r\n"):
        newline = "\r\n"

    lines.insert(
        end_index,
        f"weight: {weight}{newline}",
    )

    return "".join(lines)


def process_repo(repo_dir: Path, repo_index: int) -> int:
    """
    Add missing weights to Markdown files in one repository.

    Explicit weights are preserved.
    README.md is ignored.
    """
    doc_dir = repo_dir / "doc"

    if not doc_dir.is_dir():
        print(
            f"[SKIP] {repo_dir.name}: "
            "no doc/ directory"
        )
        return 0
    repo_base = (repo_index + 1) * REPO_WEIGHT_STEP

    # Use a deterministic order inside each repository.
    markdown_files = sorted(
        doc_dir.rglob("*.md"),
        key=lambda path: (
            path.relative_to(doc_dir)
            .as_posix()
            .lower()
        ),
    )

    updated_count = 0

    print(
        f"[INFO] Processing {repo_dir.name} "
        f"(base weight: {repo_base})"
    )

    position = 1

    for path in markdown_files:

        if path.name.lower() == "readme.md":
            continue

        text = read_text(path)
        lines = text.splitlines(keepends=True)

        front_matter_end = get_front_matter_end(lines)

        if has_weight(lines, front_matter_end):
            position += 1
            continue

        weight = (
            repo_base
            + (position * FILE_WEIGHT_STEP)
        )

        updated = add_weight(
            text,
            weight,
        )

        if updated != text:
            path.write_text(
                updated,
                encoding="utf-8",
            )

            updated_count += 1

            relative = path.relative_to(repo_dir)

            print(
                f"[UPDATE] {repo_dir.name}/{relative} "
                f"-> weight: {weight}"
            )

        position += 1

    return updated_count


def main():
    args = parse_args()

    repos_dir = Path(args.repos_dir)

    repo_order = [
        item.strip()
        for item in args.repo_order.split(",")
        if item.strip()
    ]

    if not repos_dir.exists():
        print(
            f"[WARNING] Repository directory does not exist: "
            f"{repos_dir}"
        )
        return 0

    print(
        f"[INFO] Updating metadata in: {repos_dir}"
    )

    total_updated = 0

    for repo_index, repo_name in enumerate(repo_order):

        repo_dir = repos_dir / repo_name

        if not repo_dir.is_dir():
            print(
                f"[SKIP] {repo_name}: "
                "repository not found"
            )
            continue

        total_updated += process_repo(
            repo_dir,
            repo_index,
        )

    print()

    print(
        f"[SUCCESS] Metadata update complete: "
        f"{total_updated} files updated."
    )


if __name__ == "__main__":
    main()