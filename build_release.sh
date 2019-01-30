#!/usr/bin/env bash
set -x

git tag -a $(cat VERSION) -m "Automated release version: $(cat VERSION)"
git push origin $(cat VERSION)
goreleaser