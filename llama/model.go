package llama

import "errors"

var modelLoaded = false

// IsModelLoaded returns true if the program believes the model is loaded. Because
// the actual model memory happens over FFI anything can happen.
func IsModelLoaded() bool {
	return modelLoaded
}

// LoadModel initializes a model at the given path and loads it into memory. This
// returns an error if the model fails to load. The [GPULayers] parameter must
// be set prior to calling this. The model file must adhere to llama.cpp standards,
// use GGLF/GGUF format.
func LoadModel(path string) error {
	if !inited {
		return errors.New("library not initialized")
	}

	if !load_model(path) {
		return LastError()
	}
	modelLoaded = true

	return nil
}

// FreeModel releases the model memory. There isn't really a reason to call this
// if the model will exist for the life of the program.
func FreeModel() {
	free_model()
	modelLoaded = false
}
