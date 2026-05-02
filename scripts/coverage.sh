#!/usr/bin/env bash
set -euo pipefail

COV_DIR="coverage"
MERGED="${COV_DIR}/coverage.out"

mkdir -p "${COV_DIR}"

PKGS=$(go list ./... 2>/dev/null | grep -v '/mocks' || true)

if [[ -z "${PKGS}" ]]; then
    echo "No packages found; writing empty coverage file"
    echo "mode: set" > "${MERGED}"
    echo "Coverage written to ${MERGED} (no packages)"
    exit 0
fi

COVPKGS=$(echo "${PKGS}" | tr '\n' ',' | sed 's/,$//')

echo "Running tests with coverage across all packages..."
# shellcheck disable=SC2086
go test -coverprofile="${MERGED}" -coverpkg="${COVPKGS}" ${PKGS}

echo "Coverage written to ${MERGED}"
