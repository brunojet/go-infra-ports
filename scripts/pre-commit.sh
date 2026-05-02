#!/usr/bin/env bash
set -euo pipefail

echo "==> [pre-commit] go fmt"
go fmt ./...

echo "==> [pre-commit] golangci-lint"
golangci-lint run ./...

echo "==> [pre-commit] go test"
go test ./...

echo "==> [pre-commit] coverage gate (min 95%)"
bash ./scripts/coverage.sh
COV_FILE="coverage/coverage.out"
if [[ ! -f "${COV_FILE}" ]]; then
  echo "WARN: coverage file not found; skipping coverage gate"
  exit 0
fi
total=$(go tool cover -func="${COV_FILE}" | awk '/^total:/ {gsub("%", "", $3); print $3}')
min=95
ok=$(awk -v t="$total" -v m="$min" 'BEGIN { print (t+0 >= m+0) ? "1" : "0" }')
if [ "$ok" != "1" ]; then
  echo "ERROR: coverage ${total}% is below minimum ${min}%"
  exit 1
fi
echo "OK: coverage ${total}%"
