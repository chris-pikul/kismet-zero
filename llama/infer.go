package llama

import "errors"

// Response returns the response from the LLM inference last ran.
func Response() string {
	return get_last_output()
}

// InferSync runs the inference for the given prompt using the LLM model that has
// been initialized, if one hasn't than an error is returned. Additionally, any
// errors that happen during inference will be returned.
//
// Any parameters needed by the model, context, or sampling must be set prior
// to calling this.
//
// This is the synchronous version, meaning that it blocks until the entire
// output is generated.
func InferSync(prompt string) (string, error) {
	if !inited {
		return "", errors.New("library not initialized")
	}

	if !modelLoaded {
		return "", errors.New("no model loaded")
	}

	if !infer_sync(prompt) {
		return "", LastError()
	}

	return Response(), nil
}
