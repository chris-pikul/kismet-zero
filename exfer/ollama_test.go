package exfer

import (
	"context"
	"testing"
)

func TestOllamaGenerate(t *testing.T) {
	resp, err := OllamaGenerate(&GenerateOptions{
		Model:       "llama3.2",
		Prompt:      "Why is the sky blue?",
		PredictSize: 10,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Response: %s", resp.Content)
}

func TestOllamaGenerateAsync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fragments := make(chan GenerateFragment)

	err := OllamaGenerateAsync(ctx, &GenerateOptions{
		Model:       "llama3.2",
		Prompt:      "Why is the sky blue?",
		PredictSize: 10,
	}, fragments)
	if err != nil {
		t.Error(err)
		return
	}

	for frag := range fragments {
		if frag.Done {
			break
		}
		t.Logf("Frag: %s", frag.Content)
	}
}
