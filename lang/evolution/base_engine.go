package evolution

import (
	"math/rand/v2"
)

// BaseEngine provides common functionality for all evolution engines.
// This eliminates the duplication of RNG initialization and configuration
// that was repeated across 8 different engine types.
type BaseEngine struct {
	rng    *rand.Rand
	config EvolutionConfig
}

// NewBaseEngine creates a new base engine with the given configuration.
// This consolidates the duplicate constructor pattern found across all engines.
func NewBaseEngine(config EvolutionConfig) BaseEngine {
	rng := rand.New(rand.NewPCG(uint64(config.Seed), 0))
	return BaseEngine{rng: rng, config: config}
}

// GetRNG returns the random number generator for the engine.
func (be *BaseEngine) GetRNG() *rand.Rand {
	return be.rng
}

// GetConfig returns the evolution configuration for the engine.
func (be *BaseEngine) GetConfig() EvolutionConfig {
	return be.config
}

// GetSeed returns the seed value from the configuration.
func (be *BaseEngine) GetSeed() int64 {
	return be.config.Seed
}

// GetRandomFloat32 returns a random float32 value between 0.0 and 1.0.
func (be *BaseEngine) GetRandomFloat32() float32 {
	return be.rng.Float32()
}

// GetRandomFloat64 returns a random float64 value between 0.0 and 1.0.
func (be *BaseEngine) GetRandomFloat64() float64 {
	return be.rng.Float64()
}

// GetRandomInt returns a random int value.
func (be *BaseEngine) GetRandomInt() int {
	return be.rng.Int()
}

// GetRandomInt64 returns a random int64 value.
func (be *BaseEngine) GetRandomInt64() int64 {
	return be.rng.Int64()
}
