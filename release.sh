#!/bin/bash
#
# release.sh
#
# Creates a Github release for the current tag.
# This should be run via ./release.sh from the project root.
#
# Requirements:
#  - go get github.com/mitchellh/gox
#  - go get github.com/tcnksm/ghr
#  - set GITHUB_TOKEN (see https://help.github.com/articles/creating-an-access-token-for-command-line-use/)
#
# Optional:
#  - set CIRCLE_TAG to specify the tag for the release. if unset, it will default to the most recent tag.

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# get either CIRCLE_TAG, or the most recent tag
TAG=${CIRCLE_TAG:-$(git log --simplify-by-decoration --oneline --pretty=format:"%d" --all | grep -Eo "tag: [^),\s]+" | sed "s/tag: //" | head -n 1)}

echo ">> Creating release $TAG"

# build for osx and linux
mkdir -p dist
gox \
    -ldflags "-X main.Version $TAG" \
    -output "dist/devicefarm_{{.OS}}_{{.Arch}}" \
    -osarch "darwin/386" \
    -osarch "darwin/amd64" \
    -osarch "linux/386" \
    -osarch "linux/amd64"

ghr -t $GITHUB_TOKEN -u ride -r devicefarm --draft $TAG dist/

echo ">> Done"
