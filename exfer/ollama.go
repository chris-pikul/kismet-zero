package exfer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chris-pikul/kismet-zero/utils"
)

const DefaultOllamaHost = "localhost"
const DefaultOllamaPort uint16 = 11434

// Ollama implements the [Provider] interface to allow for local Ollama inference.
type Ollama struct {
	Host   string `json:"host"`
	Port   uint16 `json:"port"`
	UseSSL bool   `json:"useSSL"`
}

// Address returns the connection URL for this provider.
func (p Ollama) Address() string {
	return fmt.Sprintf(`http%s://%s:%d`, utils.If(p.UseSSL, "s", ""), p.Host, p.Port)
}

// URL returns an HTTP request URL combining the connection settings and the given
// path segments which will be combined using slashes.
func (p Ollama) URL(path string, query map[string]any) string {
	if len(path) == 0 {
		path = "/"
	} else if path[0] != '/' {
		path = "/" + path
	}
	if query != nil {
		return fmt.Sprintf("%s%s?%s", p.Address(), path, utils.QueryString(query))
	}
	return fmt.Sprintf("%s%s", p.Address(), path)
}

// constructs a JSON marshalled request body for Ollama using the options
func (p Ollama) requestBody(opts *GenerateOptions, stream bool) ([]byte, error) {
	req := map[string]any{
		"model":  opts.Model,
		"prompt": opts.Prompt,
		"stream": stream,
	}
	if opts.SystemPrompt != "" {
		req["system"] = opts.SystemPrompt
	}

	options := map[string]any{}
	if opts.Seed != nil {
		options["seed"] = opts.Seed
	}
	if opts.ContextSize > 0 {
		options["num_ctx"] = opts.ContextSize
	}
	if opts.PredictSize > 0 {
		options["num_predict"] = opts.PredictSize
	}
	if opts.RepetitionPenalty > 0 {
		options["repeat_penalty"] = opts.RepetitionPenalty
	}
	if opts.Temperature > 0 {
		options["temperature"] = opts.Temperature
	}
	if opts.TopK != nil {
		options["top_k"] = opts.TopK
	}
	if opts.TopP > 0 {
		options["top_p"] = opts.TopP
	}
	if opts.MinP > 0 {
		options["min_p"] = opts.MinP
	}

	if len(options) > 0 {
		req["options"] = options
	}

	// Marshal as JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return body, nil
}

func (p Ollama) Generate(opts *GenerateOptions) (GenerateResult, error) {
	var result GenerateResult

	// Request payload for Ollama
	body, err := p.requestBody(opts, false)
	if err != nil {
		return result, err
	}

	url := p.URL("/api/generate", nil)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return result, fmt.Errorf("http POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, respData)
	}

	var apiResp struct {
		Model           string `json:"model"`
		Response        string `json:"response"`
		TotalDurationNS uint64 `json:"total_duration"`
		EvalDurationNS  uint64 `json:"eval_duration"`
		EvalTokens      int    `json:"eval_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return result, fmt.Errorf("decode error: %w", err)
	}

	result.Model = apiResp.Model
	result.Content = apiResp.Response
	result.TotalDuration = utils.NanoSecondDuration(apiResp.TotalDurationNS)
	result.InferenceDuration = utils.NanoSecondDuration(apiResp.EvalDurationNS)
	result.TokenCount = apiResp.EvalTokens

	return result, nil
}

func (p Ollama) GenerateAsync(ctx context.Context, opts *GenerateOptions, fragmentChan chan<- GenerateFragment) error {
	body, err := p.requestBody(opts, true)
	if err != nil {
		return fmt.Errorf("request body error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL("/api/generate", nil), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST error: %w", err)
	}

	go func() {
		defer resp.Body.Close()
		defer close(fragmentChan)

		if resp.StatusCode != http.StatusOK {
			msg, _ := io.ReadAll(resp.Body)
			fragmentChan <- GenerateFragment{
				Content: fmt.Sprintf("[error] status %d: %s", resp.StatusCode, msg),
				Done:    true,
			}
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var part struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}

			if err := decoder.Decode(&part); err != nil {
				fragmentChan <- GenerateFragment{
					Content: fmt.Sprintf("[decode error] %v", err),
					Done:    true,
				}
				return
			}

			fragmentChan <- GenerateFragment{
				Content: part.Response,
				Done:    part.Done,
			}

			if part.Done {
				return
			}
		}
	}()

	return nil
}

// DefaultOllama is an instance of [Ollama] setup with the default host and port
// configuration for local inference.
var DefaultOllama = Ollama{
	Host: DefaultOllamaHost,
	Port: DefaultOllamaPort,
}

// OllamaGenerate uses [DefaultOllama] to perform a synchronous generate request.
func OllamaGenerate(opts *GenerateOptions) (GenerateResult, error) {
	return DefaultOllama.Generate(opts)
}

// OllamaGenerateAsync uses [DefaultOllama] to perform an asynchronous streaming generation
// request.
func OllamaGenerateAsync(ctx context.Context, opts *GenerateOptions, fragmentChan chan<- GenerateFragment) error {
	return DefaultOllama.GenerateAsync(ctx, opts, fragmentChan)
}
