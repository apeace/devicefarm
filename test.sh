#!/bin/bash
#
# test.sh

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

go vet

echo "mode: atomic" > coverage.out

PACKAGES=$(go list ./... | sed -E -e "s/^github.com\/ride\/devicefarm\/([^\/]+)$/\1/" -e 'tx' -e 'd' -e ':x')
for package in $PACKAGES
do
  echo ">> package $package"
  go test -race -v -cover -coverprofile="$package.out" github.com/ride/devicefarm/$package
  cat "$package.out" | grep -v "mode:" >> coverage.out
  rm "$package.out"
done

ARTIFACTS=${CIRCLE_ARTIFACTS:-""}
go tool cover -html="coverage.out" -o="coverage.html"
rm coverage.out
if [ ! -z "$ARTIFACTS" ]; then
  mv "coverage.html" "$ARTIFACTS/"
fi
