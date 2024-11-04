package llama

import "math"

var (
	predictSize int32   = -1
	contextSize int32   = 4096
	batch       int32   = 2048
	uBatch      int32   = 512
	draft       int32   = 5
	parallel    int32   = 1
	gpuLayers   int32   = -1
	seed        int32   = math.MaxInt32
	minKeep     int32   = 0
	topK        int32   = 40
	topP        float32 = 0.95
	typicalP    float32 = 1
	minP        float32 = 0.05
	temperature float32 = 0.80
)

func PredictSize() int32 {
	return predictSize
}
func SetPredictSize(n int32) {
	predictSize = n
	set_predict(n)
}

func ContextSize() int32 {
	return contextSize
}
func SetContextSize(n int32) {
	contextSize = n
	set_context_size(n)
}

func Batch() int32 {
	return batch
}
func SetBatch(n int32) {
	batch = n
	set_batch(n)
}

func UBatch() int32 {
	return uBatch
}
func SetUBatch(n int32) {
	uBatch = n
	set_ubatch(n)
}

func Draft() int32 {
	return draft
}
func SetDraft(n int32) {
	draft = n
	set_draft_size(n)
}

func Parallel() int32 {
	return parallel
}
func SetParallel(n int32) {
	parallel = n
	set_parallel(n)
}

func GPULayers() int32 {
	return gpuLayers
}
func SetGPULayers(n int32) {
	gpuLayers = n
	set_gpu_layers(n)
}

func Seed() int32 {
	return seed
}
func SetSeed(n int32) {
	seed = n
	set_seed(n)
}

func MinKeep() int32 {
	return minKeep
}
func SetMinKeep(n int32) {
	minKeep = n
	set_min_keep(n)
}

func TopK() int32 {
	return topK
}
func SetTopK(n int32) {
	topK = n
	set_top_k(n)
}

func TopP() float32 {
	return topP
}
func SetTopP(n float32) {
	topP = n
	set_top_p(n)
}

func TypicalP() float32 {
	return typicalP
}
func SetTypicalP(n float32) {
	typicalP = n
	set_typical_p(n)
}

func MinP() float32 {
	return minP
}
func SetMinP(n float32) {
	minP = n
	set_min_p(n)
}

func Temperature() float32 {
	return temperature
}
func SetTemperature(n float32) {
	temperature = n
	set_temp(n)
}
