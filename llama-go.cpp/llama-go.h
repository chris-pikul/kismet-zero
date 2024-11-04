// Copyright 2024 Chris Pikul. All rights reserved.
//
// Based on the kelindar/search project, but rewritten for inference instead of
// embeddings.
#include "llama.h"

typedef struct llama_model *model_ptr;
typedef struct llama_context *context_ptr;

extern "C"
{
    LLAMA_API bool did_error();
    LLAMA_API const char *get_last_error();
    LLAMA_API void clear_error();
    LLAMA_API const char *get_last_output();
    LLAMA_API void clear_output();
    LLAMA_API void set_log_level(uint8_t n);

    LLAMA_API void set_predict(int32_t n);
    LLAMA_API void set_context_size(int32_t n);
    LLAMA_API void set_batch(int32_t n);
    LLAMA_API void set_ubatch(int32_t n);
    LLAMA_API void set_draft_size(int32_t n);
    LLAMA_API void set_parallel(int32_t n);
    LLAMA_API void set_gpu_layers(int32_t n);

    LLAMA_API void set_seed(int32_t n);
    LLAMA_API void set_min_keep(int32_t n);
    LLAMA_API void set_top_k(int32_t n);
    LLAMA_API void set_top_p(float n);
    LLAMA_API void set_typical_p(float n);
    LLAMA_API void set_min_p(float n);
    LLAMA_API void set_temp(float n);

    LLAMA_API void init_library(uint8_t numa);
    LLAMA_API void free_library();

    LLAMA_API bool load_model(const char *path);
    LLAMA_API void free_model();

    LLAMA_API bool infer_sync(const char *prompt);
}
