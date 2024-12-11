# Kismet Zero

As of right now, this is a work-in-progress. The goal here is to run
local LLM inference using Go without the FFI overhead by leveraging
llama.cpp project as a static library.

## Installation

The installation is rough because this is a very early prototype that I've been
doing locally. Make sure to clone this repository with all submodules included
which should pull the llama.cpp project for you. The `build-llama-go.sh` shell
will run the CMake projects and copy the resulting `.dylib` library for you. This
is because this was built for Mac/Darwin, but you can adjust it for your own
platform.

The model file is loaded via an environment variable "KISMET_MODEL" which should
point to the GFML model location. Using an application like Ollama can make this
easier to download, otherwise use the instructions listed at llama.cpp in order
to find and convert models.

## Attributions

The [llama-go.cpp](llama-go.cpp/) code was copied from [github.com/kelindar/search](https://github.com/kelindar/search/blob/main),
which is MIT licensed.
