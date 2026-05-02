#!/usr/bin/env bash
set -euo pipefail

COV_DIR="coverage"
MERGED="${COV_DIR}/coverage.out"

mkdir -p "${COV_DIR}"

ALL_PKGS=$(go list ./... 2>/dev/null | grep -Ev '/mocks($|/)' || true)
TEST_PKGS=$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... 2>/dev/null | grep -Ev '/mocks($|/)' || true)

if [[ -z "${ALL_PKGS}" ]]; then
    echo "No packages found; writing empty coverage file"
    echo "mode: set" > "${MERGED}"
    echo "Coverage written to ${MERGED} (no packages)"
    exit 0
fi

if [[ -z "${TEST_PKGS}" ]]; then
    echo "No packages with tests found; writing empty coverage file"
    echo "mode: set" > "${MERGED}"
    echo "Coverage written to ${MERGED} (no test packages)"
    exit 0
fi

COVPKGS=$(echo "${ALL_PKGS}" | tr '\n' ',' | sed 's/,$//')

echo "Running tests for packages with tests and coverage across all non-mock packages..."
# shellcheck disable=SC2086
go test -coverprofile="${MERGED}" -coverpkg="${COVPKGS}" ${TEST_PKGS}

echo "Coverage written to ${MERGED}"
