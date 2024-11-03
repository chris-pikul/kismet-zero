// Copyright 2024 Chris Pikul. All rights reserved.
//
// Based on the kelindar/search project, but rewritten for inference instead of
// embeddings.
#include "llama.h"

typedef struct llama_model *model_ptr;
typedef struct llama_context *context_ptr;

extern "C"
{
    LLAMA_API void init_library(uint8_t numa);
    LLAMA_API void init_logging(ggml_log_level level);
    LLAMA_API void free_library();

    LLAMA_API model_ptr load_model(const char *path, const uint32_t n_gpu_layers);
    LLAMA_API void free_model(model_ptr model);

    LLAMA_API bool complete(model_ptr model, const char *prompt, const char *outGenerated);
}
