#!/bin/bash
#
# dist.sh VERSION
#
# Builds a distribution with the given version number.
# This should be run via ./dist.sh from the project root.
#
# Requirements:
#  - pass VERSION argument to the script
#  - go get github.com/mitchellh/gox
#
# Note: We should NOT use OS X builds that are cross-compiled from linux.
# See docs/development.md and this issue for more info: https://github.com/golang/go/issues/6376

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

TAG=${1:-""}

if [[ -z $TAG ]]; then
    echo "usage: ./dist.sh VERSION";
    exit 1;
fi

# build for osx and linux
mkdir -p dist
gox \
    -ldflags "-X main.Version $TAG" \
    -output "dist/devicefarm_{{.OS}}_{{.Arch}}" \
    -osarch "darwin/386" \
    -osarch "darwin/amd64" \
    -osarch "linux/386" \
    -osarch "linux/amd64"
