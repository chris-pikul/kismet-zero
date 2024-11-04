// Copyright 2024 Chris Pikul. All rights reserved.
//
// Based on the kelindar/search project, but rewritten for inference instead of
// embeddings.
#include "llama-go.h"

#include <vector>
#include <algorithm>

static bool g_errored = false;
static std::string g_error;
static std::string g_output;

static int32_t g_predict = -1;
static int32_t g_ctx = 4096;
static int32_t g_batch = 2048;
static int32_t g_ubatch = 512;
static int32_t g_draft = 5;
static int32_t g_parallel = 1;
static int32_t g_gpu_layers = -1;

static int32_t g_seed = 0xFFFFFFFF;
static int32_t g_min_keep = 0;
static int32_t g_top_k = 40;
static float g_top_p = 0.95f;
static float g_typ_p = 1.0f;
static float g_min_p = 0.05f;
static float g_temp = 0.80f;

bool error(const std::string &msg)
{
    g_error.assign(msg);
    g_errored = true;
    return false;
}

extern "C"
{
    LLAMA_API bool did_error()
    {
        return g_errored;
    }

    LLAMA_API const char *get_last_error()
    {
        return g_error.c_str();
    }

    LLAMA_API void clear_error()
    {
        g_errored = false;
        g_error.clear();
    }

    LLAMA_API const char *get_last_output()
    {
        return g_output.c_str();
    }

    LLAMA_API void clear_output()
    {
        g_output.clear();
    }

    LLAMA_API void set_predict(int32_t n)
    {
        g_predict = n;
    }

    LLAMA_API void set_context_size(int32_t n)
    {
        g_ctx = n;
    }

    LLAMA_API void set_batch(int32_t n)
    {
        g_batch = n;
    }

    LLAMA_API void set_ubatch(int32_t n)
    {
        g_ubatch = n;
    }

    LLAMA_API void set_draft_size(int32_t n)
    {
        g_draft = n;
    }

    LLAMA_API void set_parallel(int32_t n)
    {
        g_parallel = n;
    }

    LLAMA_API void set_gpu_layers(int32_t n)
    {
        g_gpu_layers = n;
    }

    LLAMA_API void set_seed(int32_t n)
    {
        g_seed = n;
    }

    LLAMA_API void set_min_keep(int32_t n)
    {
        g_min_keep = n;
    }

    LLAMA_API void set_top_k(int32_t n)
    {
        g_top_k = n;
    }

    LLAMA_API void set_top_p(float n)
    {
        g_top_p = n;
    }

    LLAMA_API void set_typical_p(float n)
    {
        g_typ_p = n;
    }

    LLAMA_API void set_min_p(float n)
    {
        g_min_p = n;
    }

    LLAMA_API void set_temp(float n)
    {
        g_temp = n;
    }

    LLAMA_API void init_library(uint8_t numa)
    {
        llama_backend_init();
        llama_numa_init((ggml_numa_strategy)numa);
    }

    LLAMA_API void init_logging(ggml_log_level level)
    {
        auto level_p = new ggml_log_level;
        *level_p = level;

        llama_log_set([](ggml_log_level lvl, const char *text, void *user_data)
                      {
            ggml_log_level inLevel = *(ggml_log_level*)user_data;
            if (lvl < inLevel) return;

            fputs(text, stderr);
            fflush(stderr); }, level_p);
    }

    LLAMA_API void free_library()
    {
        llama_backend_free();
    }

    LLAMA_API model_ptr load_model(const char *path)
    {
        llama_model_params mparams = llama_model_default_params();
        mparams.n_gpu_layers = g_gpu_layers;

        return llama_load_model_from_file(path, mparams);
    }

    LLAMA_API void free_model(model_ptr model)
    {
        llama_free_model(model);
    }

    LLAMA_API bool infer_sync(model_ptr model, const char *prompt)
    {
        clear_error();

        // Tokenize the prompt
        const int n_prompt = -llama_tokenize(model, prompt, strlen(prompt), NULL, 0, true, true);
        std::vector<llama_token> prompt_tokens(n_prompt);
        if (llama_tokenize(model, prompt, strlen(prompt), prompt_tokens.data(), prompt_tokens.size(), true, true) < 0)
            return error("failed to tokenize prompt");

        const int32_t n_predict = g_predict > 0 ? std::min(n_prompt + g_predict - 1, g_ctx) : g_ctx;

        // Create a new context
        llama_context_params cparams = llama_context_default_params();
        cparams.n_ctx = n_predict;
        cparams.n_batch = g_batch;
        cparams.n_ubatch = g_ubatch;
        cparams.no_perf = true;

        context_ptr ctx = llama_new_context_with_model(model, cparams);
        if (ctx == NULL)
            return error("failed to create new context");
        const int n_ctx = llama_n_ctx(ctx);

        // Initialize the sampler
        llama_sampler_chain_params sparams = llama_sampler_chain_default_params();
        sparams.no_perf = true;
        llama_sampler *smpl = llama_sampler_chain_init(sparams);
        llama_sampler_chain_add(smpl, llama_sampler_init_top_k(g_top_k));
        llama_sampler_chain_add(smpl, llama_sampler_init_top_p(g_top_p, g_min_keep));
        llama_sampler_chain_add(smpl, llama_sampler_init_typical(g_typ_p, g_min_keep));
        llama_sampler_chain_add(smpl, llama_sampler_init_min_p(g_min_p, 1));
        llama_sampler_chain_add(smpl, llama_sampler_init_temp(g_temp));
        llama_sampler_chain_add(smpl, llama_sampler_init_dist(g_seed));

        // Prepare batch for the prompt
        llama_batch batch = llama_batch_get_one(prompt_tokens.data(), prompt_tokens.size());

        // Start inference loop
        clear_output();

        std::string output;
        llama_token next;
        char buf[128];
        int len, n_ctx_used;
        int32_t decode_res;
        while (true)
        {
            n_ctx_used = llama_get_kv_cache_used_cells(ctx);
            if (n_ctx_used + batch.n_tokens > n_predict)
                break;

            decode_res = llama_decode(ctx, batch);
            if (decode_res < 0)
                return error("error decoding");
            else if (decode_res > 0)
                return error("failed to find KV slot for batch");

            // Sample next token
            next = llama_sampler_sample(smpl, ctx, -1);
            if (llama_token_is_eog(model, next))
                break;

            // Save the current token piece
            len = llama_token_to_piece(model, next, buf, sizeof(buf), 0, true);
            if (len < 0)
                return error("failed to convert token to piece");
            output.append(buf, len);

            // Prepare next batch
            batch = llama_batch_get_one(&next, 1);
        }

        // Copy output into saved buffer
        g_output.assign(output);

        // Free resources
        llama_sampler_free(smpl);
        llama_free(ctx);

        return true;
    }
}