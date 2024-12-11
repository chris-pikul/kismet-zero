package main

import (
	"fmt"
	"kismet/llama"
	"os"
	"time"
)

var modelPath = os.Getenv("KISMET_MODEL")

func main() {
	// Initialize the library
	llama.SetLogLevel(llama.GGML_LOG_LEVEL_WARN)
	llama.Initialize(llama.GGML_NUMA_STRATEGY_DISTRIBUTE)
	llama.SetGPULayers(99)

	// Load the model
	if err := llama.LoadModel(modelPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run the inference and return the output
	start := time.Now()
	output, err := llama.InferSync("Write a short story about llamas.")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("Generated output length %d\n", len(output))
	fmt.Printf("Time to completion: %s\n", time.Since(start).String())
	fmt.Println("Output:")
	fmt.Print(output)
	fmt.Print("\n\n")
}
