#!/usr/bin/env bash
set -euo pipefail

COV_DIR="coverage"
TMP_DIR="${COV_DIR}/tmp"

mkdir -p "${COV_DIR}"
rm -rf "${TMP_DIR}"
mkdir -p "${TMP_DIR}"

echo "Collecting packages with test files..."
PKGS=$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

FIRST_MODE=""

for PKG in $PKGS; do
    case "$PKG" in
        */mocks)
            echo "Skipping generated/mock package ${PKG}"
            continue
            ;;
    esac
    SAFE=$(echo "$PKG" | tr '/:' '__')
    OUT="${TMP_DIR}/${SAFE}.out"
    echo "Running tests for ${PKG}"
    go test -coverprofile="${OUT}" "${PKG}"
    if [[ -z "${FIRST_MODE}" && -f "${OUT}" ]]; then
        FIRST_MODE=$(head -1 "${OUT}")
    fi
done

MERGED="${COV_DIR}/coverage.out"

if [[ -z "${FIRST_MODE}" ]]; then
    echo "No eligible packages with tests were found"
    exit 1
fi

echo "${FIRST_MODE}" > "${MERGED}"

for f in "${TMP_DIR}"/*.out; do
    [[ -f "$f" ]] || continue
    tail -n +2 "$f" >> "${MERGED}"
done

rm -rf "${TMP_DIR}"
echo "Merged coverage written to ${MERGED}"
