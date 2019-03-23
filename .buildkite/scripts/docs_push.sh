#!/bin/bash

if [[ "$BUILDKITE_BRANCH" != "master" ]]; then
  echo "nothing to do"
  exit 0
fi

git checkout docs
