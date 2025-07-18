package main

import (
	"context"
	"fmt"

	"github.com/chris-pikul/kismet-zero/exfer"
)

func main() {
	// Test the generation method
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fragments := make(chan exfer.GenerateFragment)

	err := exfer.OllamaGenerateAsync(ctx, &exfer.GenerateOptions{
		Model:  "llama3.2",
		Prompt: "Why is the sky blue?",
	}, fragments)
	if err != nil {
		fmt.Print(err)
		return
	}

	for frag := range fragments {
		if frag.Done {
			break
		}
		fmt.Print(frag.Content)
	}
	fmt.Print("\nDONE\n")
}
