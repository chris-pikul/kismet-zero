// Copyright 2024 Chris Pikul. All rights reserved.
//
// Based on the kelindar/search project, but rewritten for inference instead of
// embeddings.
#include "llama-go.h"

#include <vector>
#include "common.h"

extern "C"
{
    LLAMA_API void init_library(uint8_t numa)
    {
        common_params params;
        params.numa = (ggml_numa_strategy)numa;
        common_init();

        llama_backend_init();
        llama_numa_init(params.numa);
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

    LLAMA_API model_ptr load_model(const char *path, const uint32_t n_gpu_layers)
    {
        llama_model_params mparams = llama_model_default_params();
        mparams.n_gpu_layers = n_gpu_layers;

        return llama_load_model_from_file(path, mparams);
    }

    LLAMA_API void free_model(model_ptr model)
    {
        llama_free_model(model);
    }

    LLAMA_API bool complete(model_ptr model, const char *prompt, const char *outGenerated)
    {
        int n_predict = 32;

        // Tokenize the prompt
        const int n_prompt = -llama_tokenize(model, prompt, strlen(prompt), NULL, 0, true, true);

        std::vector<llama_token> prompt_tokens(n_prompt);
        if (llama_tokenize(model, prompt, strlen(prompt), prompt_tokens.data(), prompt_tokens.size(), true, true) < 0)
            return false;

        // Generate a context for it
        llama_context_params cparams = llama_context_default_params();
        cparams.n_ctx = n_prompt + n_predict - 1;
        cparams.n_batch = n_prompt;
        cparams.no_perf = true;

        context_ptr ctx = llama_new_context_with_model(model, cparams);
        if (ctx == NULL)
            return false;

        // Initialize the sampler
        llama_sampler_chain_params sparams = llama_sampler_chain_default_params();
        sparams.no_perf = true;

        llama_sampler *smpl = llama_sampler_chain_init(sparams);

        // Prepare batch for the prompt
        llama_batch batch = llama_batch_get_one(prompt_tokens.data(), prompt_tokens.size());

        // Start inference loop
        std::string output;
        llama_token ntoken;
        char buf[128];
        std::string sbuf;
        int len;
        for (int n_pos = 0; n_pos + batch.n_tokens < n_prompt + n_predict;)
        {
            if (llama_decode(ctx, batch))
                return false;

            n_pos += batch.n_tokens;

            // Sample next token
            {
                ntoken = llama_sampler_sample(smpl, ctx, -1);
                if (llama_token_is_eog(model, ntoken))
                    break;

                // Save the current token piece
                len = llama_token_to_piece(model, ntoken, buf, sizeof(buf), 0, true);
                if (len < 0)
                    return false;
                output.append(buf, len);

                // Prepare next batch
                batch = llama_batch_get_one(&ntoken, 1);
            }
        }

        // Free resources
        llama_sampler_free(smpl);
        llama_free(ctx);

        // Fill output
        outGenerated = output.c_str();

        return true;
    }
}