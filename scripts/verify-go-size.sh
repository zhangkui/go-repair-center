#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

mapfile -t files < <(find "$ROOT_DIR" -type f -name '*.go' \
  ! -name '*_test.go' \
  ! -path "$ROOT_DIR/tests/*" \
  ! -path "$ROOT_DIR/vendor/*" \
  ! -path "$ROOT_DIR/migrations/*" | sort)

line_count=0
for file in "${files[@]}"; do
  count="$(wc -l < "$file")"
  line_count=$((line_count + count))
  printf '%6s  %s\n' "$count" "${file#$ROOT_DIR/}"
done

printf '\nFILES=%s\nLINES=%s\n' "${#files[@]}" "$line_count"

