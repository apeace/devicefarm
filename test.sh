#!/bin/bash
#
# test.sh

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

echo "mode: atomic" > devicefarm.out

PACKAGES=$(go list ./... | sed -E -e "s/^github.com\/ride\/devicefarm\/([^\/]+)$/\1/" -e 'tx' -e 'd' -e ':x')
for package in $PACKAGES
do
  echo ">> package $package"
  go vet github.com/ride/devicefarm/$package
  godep go test -race -v -cover -coverprofile="$package.out" github.com/ride/devicefarm/$package
  cat "$package.out" | grep -v "mode:" >> devicefarm.out
  rm "$package.out"
done

ARTIFACTS=${CIRCLE_ARTIFACTS:-""}
go tool cover -html="devicefarm.out" -o="devicefarm.html"
rm devicefarm.out
if [ ! -z "$ARTIFACTS" ]; then
  mv "devicefarm.html" "$ARTIFACTS/"
fi
