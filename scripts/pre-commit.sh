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
total=$(go tool cover -func=coverage/coverage.out | awk '/^total:/ {gsub("%", "", $3); print $3}')
min=95
ok=$(awk -v t="$total" -v m="$min" 'BEGIN { print (t+0 >= m+0) ? "1" : "0" }')
if [ "$ok" != "1" ]; then
  echo "ERROR: coverage ${total}% is below minimum ${min}%"
  exit 1
fi
echo "OK: coverage ${total}%"
