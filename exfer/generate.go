package exfer

import "time"

// GenerateOptions are the options provided for a completion request.
type GenerateOptions struct {
	Model             string  `json:"model"`
	SystemPrompt      string  `json:"systemPrompt"`
	Prompt            string  `json:"prompt"`
	Seed              *int    `json:"seed"`
	ContextSize       uint    `json:"contextSize"`
	PredictSize       uint    `json:"predictSize"`
	RepetitionPenalty float32 `json:"repetitionPenalty"`
	Temperature       float32 `json:"temperature"`
	TopK              *int    `json:"topK"`
	TopP              float32 `json:"topP"`
	MinP              float32 `json:"minP"`
}

// GenerateResult is the combined results from a synchronous completion request.
type GenerateResult struct {
	Model             string        `json:"model"`
	Content           string        `json:"content"`
	TokenCount        int           `json:"tokens"`
	TotalDuration     time.Duration `json:"totalDuration"`
	InferenceDuration time.Duration `json:"inferenceDuration"`
}

// GenerateFragment is a partial token fragment result from an asynchronous streaming
// completion request.
type GenerateFragment struct {
	Done    bool   `json:"done"`
	Content string `json:"content"`
}
