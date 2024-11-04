package main

import (
	"fmt"
	"kismet/llama"
	"os"
	"time"
)

const modelPath = "/Users/chris/.ollama/models/blobs/sha256-00e1317cbf74d901080d7100f57580ba8dd8de57203072dc6f668324ba545f29"

func main() {
	start := time.Now()

	var model llama.Model
	model.Params = llama.NewParams()
	model.GPULayers = 99

	if err := model.Load(modelPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output, err := model.InferSync("Write a short story about llamas.")
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
