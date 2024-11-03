// Copied from github.com/kelindar/search, original copyright is as follows:
//
// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
package llama

import "fmt"

// Vectorizer represents a loaded LLM/Embedding model.
type Vectorizer struct {
	handle uintptr
	n_embd int32
	pool   *pool[*Context]
}

// NewVectorizer creates a new vectorizer model from the given model file.
func NewVectorizer(modelPath string, gpuLayers int) (*Vectorizer, error) {
	handle := load_model(modelPath, uint32(gpuLayers))
	if handle == 0 {
		return nil, fmt.Errorf("failed to load model (%s)", modelPath)
	}

	model := &Vectorizer{
		handle: handle,
		n_embd: embed_size(handle),
	}

	// Initialize the context pool to reduce allocations
	model.pool = newPool(16, func() *Context {
		return model.Context(0)
	})
	return model, nil
}

// Close closes the model and releases any resources associated with it.
func (m *Vectorizer) Close() error {
	free_model(m.handle)
	m.handle = 0
	m.pool.Close()
	return nil
}

// Context creates a new context of the given size.
func (m *Vectorizer) Context(size int) *Context {
	return &Context{
		parent: m,
		handle: load_context(m.handle, uint32(size), true),
	}
}

// EmbedText embeds the given text using the model.
func (m *Vectorizer) EmbedText(text string) ([]float32, error) {
	ctx := m.pool.Get()
	defer m.pool.Put(ctx)
	return ctx.EmbedText(text)
}
