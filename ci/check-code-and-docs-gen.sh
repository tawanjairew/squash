#!/bin/bash

set -ex

protoc --version

if [ ! -f .gitignore ]; then
  echo "_output" > .gitignore
fi

git init
git add .
git commit -m "set up dummy repo for diffing" -q

git clone https://github.com/solo-io/solo-kit /workspace/gopath/src/github.com/solo-io/solo-kit

make update-deps
make pin-repos

PATH=/workspace/gopath/bin:$PATH

set +e

make generatecode -B > /dev/null
if [[ $? -ne 0 ]]; then
  echo "Code generation failed"
  exit 1;
fi
if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Generating code produced a non-empty diff"
  echo "Try running 'make update-deps generatecode -B' then re-pushing."
  git status --porcelain
  git diff | cat
  exit 1;
fi
