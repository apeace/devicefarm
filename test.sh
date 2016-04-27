#!/bin/bash
#
# test.sh
#
# Runs unit tests for all sub-packages, but not the main package.
# This should be run via ./test.sh from the project root.

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# get project name from Circle env var if it exists, otherwise
# take the name of the current directory
PROJECT_NAME=${CIRCLE_PROJECT_REPONAME:-${PWD##*/}}

# start a coverage file. the "mode: atomic" is output by each
# run of `go test -cover`, but we only want it to appear once
echo "mode: atomic" > coverage.out

# get every sub-package, excluding the vendor directory,
# and excluding the main (top-level) package.
PACKAGES=$(go list ./... | grep -v "vendor" | sed -E -e "s/.*$PROJECT_NAME\/([^\/]+)$/\1/" -e 'tx' -e 'd' -e ':x')

# if user ran as `./test.sh package`, only test the given package
if [ ! -z $1 ]; then
  PACKAGES=$1
fi

for package in $PACKAGES
do
  echo ">> package $package"
  # cannot use vet because of this "unfortunate" issue https://github.com/golang/go/issues/9171
  #go vet ./$package
  go test -race -v -cover -coverprofile="$package.out" ./$package
  cat "$package.out" | grep -v "mode:" >> coverage.out
  rm "$package.out"
done

# put a coverage.html in this directory, or copy it to Circle
# artifacts if we are on Circle
ARTIFACTS=${CIRCLE_ARTIFACTS:-""}
go tool cover -html="coverage.out" -o="coverage.html"
rm coverage.out
if [ ! -z "$ARTIFACTS" ]; then
  mv "coverage.html" "$ARTIFACTS/"
fi
