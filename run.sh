#!/usr/bin/env bash

if [ "$(uname)" == "Darwin" ]; then
  export DYLD_LIBRARY_PATH="./build/lib"
fi

./build.sh && build/bin/start-server