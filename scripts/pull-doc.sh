#!/usr/bin/env bash
set -euo pipefail

# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

# Usage:
#   ./scripts/pull-doc.sh repos.txt
#   ./scripts/pull-doc.sh owner/repo another-owner/another-repo
#   ./scripts/pull-doc.sh https://github.com/owner/repo.git

OUT_DIR="${REPOS_DIR:-build/repos}"

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <repos.txt | repo1 repo2 ...>"
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

read_repos() {
  if [[ $# -eq 1 && -f "$1" ]]; then
    grep -vE '^\s*(#|$)' "$1"
  else
    printf '%s\n' "$@"
  fi
}

repo_name_from_ref() {
  local ref="$1"

  ref="${ref%/}"
  ref="${ref%.git}"
  basename "$ref"
}

repo_url_from_ref() {
  local ref="$1"

  if [[ "$ref" =~ ^https?:// || "$ref" =~ ^git@ ]]; then
    echo "$ref"
  else
    echo "https://github.com/${ref}.git"
  fi
}

mkdir -p "$OUT_DIR"

while IFS= read -r repo_ref; do
  repo_name="$(repo_name_from_ref "$repo_ref")"
  repo_url="$(repo_url_from_ref "$repo_ref")"

  clone_dir="$tmp_dir/$repo_name"
  target_dir="$OUT_DIR/$repo_name/doc"

  echo "Fetching doc from $repo_ref..."

  git clone \
    --depth 1 \
    --single-branch \
    --filter=blob:none \
    --sparse \
    "$repo_url" \
    "$clone_dir"

  git -C "$clone_dir" sparse-checkout set --cone doc

  if [[ ! -d "$clone_dir/doc" ]]; then
    echo "Warning: $repo_ref has no doc/ directory, skipping."
    continue
  fi

  rm -rf "$target_dir"
  mkdir -p "$(dirname "$target_dir")"

  cp -a "$clone_dir/doc" "$target_dir"

  echo "Downloaded to $target_dir"
done < <(read_repos "$@")

echo "Done."
