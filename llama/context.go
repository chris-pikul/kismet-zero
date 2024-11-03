// Copied from github.com/kelindar/search, original copyright is as follows:
//
// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
package llama

import (
	"fmt"
	"sync/atomic"
)

// Context represents a context for embedding text using the model.
type Context struct {
	parent *Vectorizer
	handle uintptr
	tokens atomic.Uint64
}

// Close closes the context and releases any resources associated with it.
func (ctx *Context) Close() error {
	free_context(ctx.handle)
	ctx.handle = 0
	return nil
}

// Tokens returns the number of tokens processed by the context.
func (ctx *Context) Tokens() uint {
	return uint(ctx.tokens.Load())
}

// EmbedText embeds the given text using the model.
func (ctx *Context) EmbedText(text string) ([]float32, error) {
	switch {
	case ctx.handle == 0 || ctx.parent.handle == 0:
		return nil, fmt.Errorf("context is not initialized")
	case ctx.parent.n_embd <= 0:
		return nil, fmt.Errorf("model does not support embedding")
	}

	out := make([]float32, ctx.parent.n_embd)
	tok := uint32(0)
	ret := embed_text(ctx.handle, text, out, &tok)
	ctx.tokens.Add(uint64(tok))
	switch ret {
	case 0:
		return out, nil
	case 1:
		return nil, fmt.Errorf("number of tokens (%d) exceeds batch size", tok)
	case 2:
		return nil, fmt.Errorf("last token in the prompt is not SEP")
	case 3:
		return nil, fmt.Errorf("failed to decode/encode text")
	default:
		return nil, fmt.Errorf("failed to embed text (code=%d)", ret)
	}
}
