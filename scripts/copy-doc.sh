#!/usr/bin/env bash
set -euo pipefail

# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

BUILD_DIR="${1:-build/repos}"
SITE_DIR="${2:-build/staged}"
SOURCE_MAP="${3:-build/source-map.tsv}"

LOCAL_CONTENT_DIR="content"

copy_tree_contents() {
    local src="$1"
    local dst="$2"
    local repo="$3"
    local repo_root="$4"

    [[ -d "$src" ]] || return 0

    mkdir -p "$dst"

    while IFS= read -r -d '' file; do
        # If a directory contains an explicit index.md, do not stage
        # README.md from the same directory. This prevents Zensical
        # from treating README.md as an alternative section index.
        if [[ "$(basename "$file")" == "README.md" ]] &&
           [[ -f "$(dirname "$file")/index.md" ]]; then
            continue
        fi

        local relative_path
        local staged_path
        local source_path

        relative_path="${file#"$src"/}"
        staged_path="$dst/$relative_path"

        mkdir -p "$(dirname "$staged_path")"
        cp -a "$file" "$staged_path"

        staged_path="${staged_path#"$SITE_DIR"/}"

        source_path="${file#"$repo_root"/}"

        printf '%s\t%s\t%s\n' \
            "$staged_path" \
            "$repo" \
            "$source_path" >> "$SOURCE_MAP"

    done < <(find "$src" -type f -print0)
}

stage_from_content_root() {
    local content_root="$1"
    local repo="$2"
    local repo_root="$3"
    local source_path

    [[ -d "$content_root" ]] || return 0

    if [[ -f "$content_root/index.md" ]]; then
        cp -a \
            "$content_root/index.md" \
            "$SITE_DIR/index.md"

        source_path="${content_root#"$repo_root"/}/index.md"

        printf '%s\t%s\t%s\n' \
            "index.md" \
            "$repo" \
            "$source_path" >> "$SOURCE_MAP"
    fi

    copy_tree_contents \
        "$content_root/docs/user" \
        "$SITE_DIR/user" \
        "$repo" \
        "$repo_root"

    copy_tree_contents \
        "$content_root/apis" \
        "$SITE_DIR/apis" \
        "$repo" \
        "$repo_root"

    copy_tree_contents \
        "$content_root/about" \
        "$SITE_DIR/about" \
        "$repo" \
        "$repo_root"

    if [[ -d "$content_root/docs" ]]; then
        mkdir -p "$SITE_DIR/doc"

        if [[ -f "$content_root/docs/index.md" ]]; then
            cp -a \
                "$content_root/docs/index.md" \
                "$SITE_DIR/doc/index.md"

            source_path="${content_root#"$repo_root"/}/docs/index.md"

            printf '%s\t%s\t%s\n' \
                "doc/index.md" \
                "$repo" \
                "$source_path" >> "$SOURCE_MAP"
        fi

        for d in "$content_root/docs"/*; do
            [[ -d "$d" ]] || continue
            [[ "$(basename "$d")" == "user" ]] && continue

            local dirname
            dirname="$(basename "$d")"

            copy_tree_contents \
                "$d" \
                "$SITE_DIR/doc/$dirname" \
                "$repo" \
                "$repo_root"
        done
    fi
}


rm -f "$SITE_DIR/index.md"

rm -rf \
    "$SITE_DIR/user" \
    "$SITE_DIR/apis" \
    "$SITE_DIR/doc" \
    "$SITE_DIR/about"

mkdir -p "$SITE_DIR"

rm -f "$SOURCE_MAP"
mkdir -p "$(dirname "$SOURCE_MAP")"

printf '# staged_path\trepository\tsource_path\n' > "$SOURCE_MAP"


if [[ -d "$LOCAL_CONTENT_DIR" ]]; then
    echo "Staging local documentation"

    stage_from_content_root \
        "$LOCAL_CONTENT_DIR" \
        "dse.doc" \
        "."
fi



for repo_content in "$BUILD_DIR"/*/doc/content; do
    [[ -d "$repo_content" ]] || continue

    repo_dir="$(dirname "$(dirname "$repo_content")")"
    repo="$(basename "$repo_dir")"

    echo "Staging from $repo_content"

    stage_from_content_root \
        "$repo_content" \
        "$repo" \
        "$repo_dir"
done


if [[ -f "$LOCAL_CONTENT_DIR/stylesheets/extra.css" ]]; then
    mkdir -p "$SITE_DIR/stylesheets"

    cp -a \
        "$LOCAL_CONTENT_DIR/stylesheets/extra.css" \
        "$SITE_DIR/stylesheets/extra.css"
fi

echo "Staged docs into $SITE_DIR"
echo "Generated source map: $SOURCE_MAP"
