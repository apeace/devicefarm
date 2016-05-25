#!/bin/bash
#
# release.sh
#
# Creates a Github draft release for the current tag.
# This should be run via ./release.sh from the project root.
#
# Requirements:
#  - go get github.com/mitchellh/gox
#  - go get github.com/tcnksm/ghr
#  - set GITHUB_TOKEN (see https://help.github.com/articles/creating-an-access-token-for-command-line-use/)
#
# Optional:
#  - set CIRCLE_TAG to specify the tag for the release. if unset, it will default to the most recent tag.
#
# Note: We should NOT use OS X builds that are cross-compiled from linux.
# See docs/development.md and this issue for more info: https://github.com/golang/go/issues/6376

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# get either CIRCLE_TAG, or the most recent tag
TAG=${CIRCLE_TAG:-$(git log --simplify-by-decoration --oneline --pretty=format:"%d" --all | grep -Eo "tag: [^),\s]+" | sed "s/tag: //" | head -n 1)}

echo ">> Creating release $TAG"

./dist.sh $TAG

ghr -t $GITHUB_TOKEN -u ride -r devicefarm --draft $TAG dist/

echo ">> Done"
