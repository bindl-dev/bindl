#!/usr/bin/env bash

set -e
set -o pipefail

dir=$(dirname ${0})

pushd $dir &>/dev/null
  tar -czf myprogram.tar.gz myprogram
popd &>/dev/null
