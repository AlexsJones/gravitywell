#!/usr/bin/env bash
set -x

git add VERSION
git commit -m "bumping version to $(cat version)"
git push origin master
git tag -a $(cat VERSION) -m "Automated release version: $(cat VERSION)"
git push origin $(cat VERSION)
export VERSION=$(cat version)
goreleaser --rm-dist