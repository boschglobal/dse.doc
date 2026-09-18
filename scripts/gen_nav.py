#!/usr/bin/env python3

# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

import argparse
import os
import re
from pathlib import Path

TOP_NAV = [
    ("Home", "index.md"),
    ("About", "about"),
    ("Guide", "user"),
    ("Doc", "doc"),
    ("Api", "apis"),
]


def parse_args():
    parser = argparse.ArgumentParser(
        description="Generate Zensical navigation from staged site docs."
    )
    parser.add_argument(
        "--input",
        "-i",
        default="zensical.toml",
        help="Path to base TOML config",
    )
    parser.add_argument(
        "--output",
        "-o",
        default="build/zensical.generated.toml",
        help="Path to output TOML config",
    )
    parser.add_argument(
        "--site-dir",
        "-s",
        default="build/staged",
        help="Staged docs directory",
    )
    parser.add_argument(
        "--site-output-dir",
        default="build/site",
        help="Built site output directory",
    )
    return parser.parse_args()


def strip_quotes(value: str) -> str:
    value = value.strip()
    if len(value) >= 2 and (
        (value[0] == '"' and value[-1] == '"')
        or (value[0] == "'" and value[-1] == "'")
    ):
        return value[1:-1]
    return value


def read_front_matter(md_path: Path) -> dict:
    meta = {}

    try:
        text = md_path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        text = md_path.read_text(encoding="utf-8", errors="ignore")

    lines = text.splitlines()

    if not lines or lines[0].strip() != "---":
        return meta

    end_idx = None

    for idx in range(1, len(lines)):
        if lines[idx].strip() == "---":
            end_idx = idx
            break

    if end_idx is None:
        return meta

    for line in lines[1:end_idx]:
        match = re.match(
            r"^([A-Za-z_][A-Za-z0-9_-]*)\s*:\s*(.*)$",
            line,
        )
        if not match:
            continue

        key = match.group(1)
        raw_value = match.group(2).strip()
        value = strip_quotes(raw_value)
        meta[key] = value

    return meta


def parse_weight(value: str) -> int:
    if value is None:
        return 10_000

    match = re.match(r"^\s*(-?\d+)", str(value))

    if match:
        return int(match.group(1))

    return 10_000


# Known proper casing and acronym overrides for directory labels.
# Filesystem names remain unchanged; this only controls how they are displayed
# in the generated navigation.
_CUSTOM_LABELS = {
    "sdp": "SDP",
    "csv": "CSV",
    "lua": "Lua",
    "fmi": "FMI",
    "fmu": "FMU",
    "modelc": "ModelC",
    "modelc-fmu": "ModelC FMU",
    "modelc_fmu": "ModelC FMU",
    "ncodec": "NCodec",
    "mcl": "MCL",
    "ast": "AST",
    "dsl": "DSL",
    "lsp": "LSP",
    "cmocka": "CMocka",
    "simmock": "SimMock",
    "testscript": "Testscript",
    "c4": "C4 Diagrams",
    "c4-diagrams": "C4 Diagrams",
    "c4_diagrams": "C4 Diagrams",
    "network-bus": "Network Bus",
    "network_bus": "Network Bus",
}


def directory_label(name: str) -> str:
    key = name.strip().lower()

    if key in _CUSTOM_LABELS:
        return _CUSTOM_LABELS[key]

    return re.sub(r"[._\-]+", " ", name).title()


def fallback_title(md_path: Path) -> str:
    return md_path.stem


def nav_title(md_path: Path, meta: dict) -> str:
    for key in ("linkTitle", "title"):
        if key in meta and str(meta[key]).strip():
            return str(meta[key]).strip()

    return fallback_title(md_path)


def collect_markdown(root: Path):
    if not root.exists():
        return []

    files = []

    for p in sorted(root.rglob("*.md"), key=lambda x: str(x).lower()):
        if any(part.startswith(".") for part in p.parts):
            continue

        files.append(p)

    return files

def get_effective_weight(directory: Path) -> int:
    """
    Return the effective navigation weight for a directory.

    If the directory has index.md, its weight is used.

    If the directory has no index.md, the lowest weight
    found anywhere below that directory is used.

    This allows category-only directories such as:
        user/guides/
        user/envars/
        user/models/

    to be positioned according to their contents.
    """

    index_path = directory / "index.md"

    # Directory has its own landing page.
    if index_path.exists():
        meta = read_front_matter(index_path)
        return parse_weight(meta.get("weight"))

    weights = []

    for path in directory.iterdir():
        if path.name.startswith("."):
            continue

        if path.is_dir():
            child_weight = get_effective_weight(path)

            if child_weight < 10_000:
                weights.append(child_weight)

        elif (
            path.is_file()
            and path.suffix.lower() == ".md"
            and path.name.lower() not in ("index.md", "readme.md")
        ):
            meta = read_front_matter(path)

            if "weight" in meta:
                weight = parse_weight(meta.get("weight"))

                if weight < 10_000:
                    weights.append(weight)

    if weights:
        return min(weights)

    return 10_000


def build_tree(section_dir: Path, site_root: Path):

    entries = []

    for path in section_dir.iterdir():
        if path.name.startswith("."):
            continue

        if path.is_dir():
            index_path = path / "index.md"

            if index_path.exists():
                meta = read_front_matter(index_path)

                title = nav_title(
                    index_path,
                    meta,
                )

                weight = parse_weight(
                    meta.get("weight")
                )

            else:
                title = directory_label(path.name)

                weight = get_effective_weight(
                    path
                )

            entries.append({
                "kind": "directory",
                "path": path,
                "title": title,
                "weight": weight,
            })

        elif (
            path.is_file()
            and path.suffix.lower() == ".md"
            and path.name.lower() not in (
                "index.md",
                "readme.md",
            )
        ):
            meta = read_front_matter(path)

            title = nav_title(
                path,
                meta,
            )

            weight = parse_weight(
                meta.get("weight")
            )

            entries.append({
                "kind": "file",
                "path": path,
                "title": title,
                "weight": weight,
            })


    entries.sort(
        key=lambda entry: (
            entry["weight"],
            entry["title"].casefold(),
        )
    )

    children = []

    for entry in entries:
        path = entry["path"]
        title = entry["title"]

        if entry["kind"] == "file":
            rel = path.relative_to(
                site_root
            ).as_posix()

            children.append({
                title: rel
            })

            continue

        subtree = build_tree(
            path,
            site_root,
        )

        index_path = path / "index.md"


        if index_path.exists():
            landing_rel = index_path.relative_to(
                site_root
            ).as_posix()

            directory_children = [
                {
                    nav_title(
                        index_path,
                        read_front_matter(index_path),
                    ): landing_rel
                }
            ]

            directory_children.extend(subtree)

            children.append({
                title: directory_children
            })


        elif subtree:
            children.append({
                title: subtree
            })

    return children


def dump_toml_inline(obj):
    if isinstance(obj, list):
        return "[" + ", ".join(
            dump_toml_inline(i) for i in obj
        ) + "]"

    if isinstance(obj, dict):
        parts = []

        for k, v in obj.items():
            parts.append(
                f'"{k}" = {dump_toml_inline(v)}'
            )

        return "{" + ", ".join(parts) + "}"

    if isinstance(obj, str):
        escaped = (
            obj
            .replace("\\", "\\\\")
            .replace('"', '\\"')
        )
        return f'"{escaped}"'

    if isinstance(obj, int):
        return str(obj)

    return '""'


def remove_existing_nav(toml_text: str) -> str:
    match = re.search(
        r"(?m)^\s*nav\s*=\s*\[",
        toml_text,
    )

    if not match:
        return toml_text

    start = match.start()
    bracket_start = toml_text.find("[", match.start())

    if bracket_start == -1:
        return toml_text

    depth = 0
    i = bracket_start
    in_string = False
    escaped = False

    while i < len(toml_text):
        ch = toml_text[i]

        if in_string:
            if escaped:
                escaped = False
            elif ch == "\\":
                escaped = True
            elif ch == '"':
                in_string = False

        else:
            if ch == '"':
                in_string = True
            elif ch == "[":
                depth += 1
            elif ch == "]":
                depth -= 1

                if depth == 0:
                    end = i + 1

                    while (
                        end < len(toml_text)
                        and toml_text[end] in " \t"
                    ):
                        end += 1

                    if (
                        end < len(toml_text)
                        and toml_text[end] == "\n"
                    ):
                        end += 1

                    return (
                        toml_text[:start]
                        + toml_text[end:]
                    )

        i += 1

    return toml_text


def ensure_docs_dir(toml_text: str, docs_dir: str) -> str:
    docs_dir = docs_dir.replace("\\", "/")

    if re.search(
        r'(?m)^\s*docs_dir\s*=\s*"',
        toml_text,
    ):
        return re.sub(
            r'(?m)^(\s*docs_dir\s*=\s*")[^"]*("\s*)$',
            rf'\1{docs_dir}\2',
            toml_text,
            count=1,
        )

    return toml_text + f'\ndocs_dir = "{docs_dir}"\n'


def ensure_site_dir(toml_text: str, site_dir: str) -> str:
    site_dir = site_dir.replace("\\", "/")

    if re.search(
        r'(?m)^\s*site_dir\s*=\s*"',
        toml_text,
    ):
        return re.sub(
            r'(?m)^(\s*site_dir\s*=\s*")[^"]*("\s*)$',
            rf'\1{site_dir}\2',
            toml_text,
            count=1,
        )

    return toml_text + f'\nsite_dir = "{site_dir}"\n'


def relative_path(from_path: Path, to_path: Path) -> str:
    try:
        return os.path.relpath(
            to_path,
            start=from_path.parent,
        ).replace("\\", "/")
    except ValueError:
        return to_path.as_posix()


def is_absolute_or_url(path_value: str) -> bool:
    return bool(
        re.match(
            r"^[A-Za-z][A-Za-z0-9+.-]*://",
            path_value,
        )
    ) or path_value.startswith("/")


def rewrite_config_path_line(
    toml_text: str,
    key: str,
) -> str:


    pattern = rf'(?m)^(\s*{re.escape(key)}\s*=\s*)("[^"]*")\s*$'

    def replace(match):
        original_value = match.group(2).strip('"')

        if original_value.startswith(("http://", "https://", "/")):
            rewritten_value = original_value
        else:
            rewritten_value = f"../{original_value}"

        return f'{match.group(1)}"{rewritten_value}"'

    return re.sub(pattern, replace, toml_text)

def rewrite_relative_project_paths(
    toml_text: str,
    input_config_path: Path,
    output_config_path: Path,
) -> str:

    toml_text = rewrite_config_path_line(
        toml_text,
        "custom_dir",
    )

    return toml_text


def build_nav(site_root: Path):
    nav = []

    home_path = site_root / "index.md"

    if home_path.exists():
        nav.append({
            "Home": "index.md"
        })

    for section_title, section_dir_name in TOP_NAV[1:]:
        section_dir = site_root / section_dir_name

        if not section_dir.exists():
            continue

        index_path = section_dir / "index.md"

        if index_path.exists():
            landing_rel = index_path.relative_to(
                site_root
            ).as_posix()

            subtree = build_tree(
                section_dir,
                site_root
            )

            directory_children = [
                {
                    nav_title(
                        index_path,
                        read_front_matter(index_path),
                    ): landing_rel
                }
            ]

            directory_children.extend(subtree)

            nav.append({
                section_title: directory_children
            })
            continue

        # Sections without an index.md use the generated tree.
        section_tree = build_tree(
            section_dir,
            site_root
        )

        if section_tree:
            nav.append({
                section_title: section_tree
            })

    return nav

def insert_nav_into_project_section(toml_text: str, nav_tree) -> str:
    nav_line = (
        "nav = "
        + dump_toml_inline(nav_tree)
        + "\n"
    )

    project_match = re.search(
        r"(?m)^\[project\]\s*$",
        toml_text,
    )

    if not project_match:
        return "[project]\n" + nav_line + "\n" + toml_text

    project_start = project_match.end()

    next_section = re.search(
        r"(?m)^\[[^\[]",
        toml_text[project_start:],
    )

    if next_section:
        insert_pos = project_start + next_section.start()
    else:
        insert_pos = len(toml_text)

    # Remove any existing nav inside the project section.
    project_content = toml_text[project_start:insert_pos]

    project_content = re.sub(
        r"(?ms)^\s*nav\s*=\s*.*?(?=^\s*\[[^\[]|\Z)",
        "\n",
        project_content,
        count=1,
    )

    return (
        toml_text[:project_start]
        + project_content.rstrip()
        + "\n"
        + nav_line
        + toml_text[insert_pos:]
    )

def main():
    args = parse_args()

    input_config_path = Path(args.input)
    output_config_path = Path(args.output)
    site_root = Path(args.site_dir)
    site_output_root = Path(args.site_output_dir)

    output_config_path.parent.mkdir(
        parents=True,
        exist_ok=True,
    )

    site_root.mkdir(
        parents=True,
        exist_ok=True,
    )

    print(
        f"[INFO] Scanning staged docs: {site_root}"
    )

    nav_tree = build_nav(site_root)

    base_content = ""

    if input_config_path.exists():
        print(
            "[INFO] Reading base configuration from: "
            f"{input_config_path}"
        )

        base_content = input_config_path.read_text(
            encoding="utf-8"
        )

    else:
        print(
            f"[WARNING] Base config "
            f"'{input_config_path}' not found. "
            "Creating generated config from nav only."
        )

    cleaned = rewrite_relative_project_paths(
        base_content,
        input_config_path,
        output_config_path,
    )

    cleaned = remove_existing_nav(cleaned)

    docs_dir_for_config = relative_path(
        output_config_path,
        site_root.resolve(),
    )

    site_dir_for_config = relative_path(
        output_config_path,
        site_output_root.resolve(),
    )

    cleaned = ensure_docs_dir(
        cleaned,
        docs_dir_for_config,
    )

    cleaned = ensure_site_dir(
        cleaned,
        site_dir_for_config,
    )

    cleaned = insert_nav_into_project_section(
        cleaned,
        nav_tree,
    )

    print(
        f"[INFO] Writing generated configuration to: "
        f"{output_config_path}"
    )

    output_config_path.write_text(
        cleaned,
        encoding="utf-8",
    )

    print(
        "[SUCCESS] Navigation configuration "
        "generated successfully."
    )


if __name__ == "__main__":
    main()