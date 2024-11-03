#!/usr/bin/env bash

# Copied in as per their instructions from:
# https://github.com/kelindar/search/blob/main/README.md#compile-on-linux

mkdir -p build && cd build
cmake -DBUILD_SHARED_LIBS=ON -DLLAMA_BUILD_COMMON=ON -DCMAKE_BUILD_TYPE=Release -DCMAKE_CXX_COMPILER=g++ -DCMAKE_C_COMPILER=gcc ..
cmake --build . --config Release