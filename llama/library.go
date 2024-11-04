package llama

var inited = false

// IsInitialized returns true if the underlying llama.cpp library is initialized.
func IsInitialized() bool {
	return inited
}

// Initialize starts up the underlying llama.cpp library over the FFI bridge.
//
// Internally this package uses [purego] to do this over dynamic library symbols
// instead of build time CGO.
func Initialize(numa NUMASetting) {
	init_library(numa)
	inited = true
}

// FreeLibrary frees the memory of the underlying llama.cpp library over the
// FFI bridge. There isn't much reason to call this really.
func FreeLibrary() {
	free_library()
	inited = false
}
