#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

install_hook() {
  local name="$1"
  local src="$REPO_ROOT/scripts/${name}.sh"
  local dst="$HOOKS_DIR/$name"

  if [ ! -f "$src" ]; then
    echo "WARNING: $src not found, skipping"
    return
  fi

  cp "$src" "$dst"
  chmod +x "$dst"
  echo "Installed: $dst"
}

install_hook pre-commit

echo "Done. Git hooks installed."
