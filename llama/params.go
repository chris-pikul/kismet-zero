package llama

import "math"

type Params struct {
	PredictSize int32
	ContextSize int32
	BatchSize   int32
	UBatchSize  int32
	DraftSize   int32
	NumParallel int32
	GPULayers   int32
	Seed        int32
	MinKeep     int32
	TopK        int32
	TopP        float32
	TypicalP    float32
	MinP        float32
	Temperature float32
}

func (p Params) Flush() {
	set_predict(p.PredictSize)
	set_context_size(p.ContextSize)
	set_batch(p.BatchSize)
	set_ubatch(p.UBatchSize)
	set_draft_size(p.DraftSize)
	set_parallel(p.NumParallel)
	set_gpu_layers(p.GPULayers)
	set_seed(p.Seed)
	set_min_keep(p.MinKeep)
	set_top_k(p.TopK)
	set_top_p(p.TopP)
	set_typical_p(p.TypicalP)
	set_min_p(p.MinP)
	set_temp(p.Temperature)
}

func NewParams() *Params {
	return &Params{
		PredictSize: -1,
		ContextSize: 4096,
		BatchSize:   2048,
		UBatchSize:  512,
		DraftSize:   5,
		NumParallel: 1,
		GPULayers:   -1,
		Seed:        math.MaxInt32,
		MinKeep:     0,
		TopK:        40,
		TopP:        0.95,
		TypicalP:    1.0,
		MinP:        0.05,
		Temperature: 0.80,
	}
}
