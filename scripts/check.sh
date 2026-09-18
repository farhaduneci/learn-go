#!/usr/bin/env bash
# Verifies the repo: formatting, vet, and every reference solution under -race.
# Stubs are expected to fail, so they are only checked for formatting.
set -euo pipefail
cd "$(dirname "$0")/.."

unformatted=$(gofmt -l .)
if [[ -n "$unformatted" ]]; then
  echo "gofmt needed:"; echo "$unformatted"; exit 1
fi

go vet -tags solution ./...
go test -race -tags solution -timeout 120s ./...
echo "all checks passed"
