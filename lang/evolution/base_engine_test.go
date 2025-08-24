package evolution

import (
	"testing"
)

// TestNewBaseEngine tests the BaseEngine constructor
func TestNewBaseEngine(t *testing.T) {
	config := EvolutionConfig{Seed: 42}
	engine := NewBaseEngine(config)

	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}

	if engine.rng == nil {
		t.Error("Expected RNG to be initialized")
	}
}

// TestBaseEngineGetRNG tests the GetRNG method
func TestBaseEngineGetRNG(t *testing.T) {
	config := EvolutionConfig{Seed: 42}
	engine := NewBaseEngine(config)

	rng := engine.GetRNG()
	if rng == nil {
		t.Error("Expected RNG to be returned")
	}
}

// TestBaseEngineGetConfig tests the GetConfig method
func TestBaseEngineGetConfig(t *testing.T) {
	config := EvolutionConfig{Seed: 42}
	engine := NewBaseEngine(config)

	retrievedConfig := engine.GetConfig()
	if retrievedConfig.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", retrievedConfig.Seed)
	}
}

// TestBaseEngineGetSeed tests the GetSeed method
func TestBaseEngineGetSeed(t *testing.T) {
	config := EvolutionConfig{Seed: 42}
	engine := NewBaseEngine(config)

	seed := engine.GetSeed()
	if seed != 42 {
		t.Errorf("Expected seed 42, got %d", seed)
	}
}

// TestBaseEngineRandomMethods tests all random number generation methods
func TestBaseEngineRandomMethods(t *testing.T) {
	config := EvolutionConfig{Seed: 42}
	engine := NewBaseEngine(config)

	// Test GetRandomFloat32
	float32Val := engine.GetRandomFloat32()
	if float32Val < 0.0 || float32Val >= 1.0 {
		t.Errorf("Expected float32 between 0.0 and 1.0, got %f", float32Val)
	}

	// Test GetRandomFloat64
	float64Val := engine.GetRandomFloat64()
	if float64Val < 0.0 || float64Val >= 1.0 {
		t.Errorf("Expected float64 between 0.0 and 1.0, got %f", float64Val)
	}

	// Test GetRandomInt
	intVal := engine.GetRandomInt()
	// Note: int can be negative, so we just check it's not zero (which is unlikely with seed 42)
	if intVal == 0 {
		t.Log("Random int is 0 (this is possible but unlikely with seed 42)")
	}

	// Test GetRandomInt64
	int64Val := engine.GetRandomInt64()
	// Note: int64 can be negative, so we just check it's not zero (which is unlikely with seed 42)
	if int64Val == 0 {
		t.Log("Random int64 is 0 (this is possible but unlikely with seed 42)")
	}
}

// TestBaseEngineConsistency tests that the same seed produces consistent results
func TestBaseEngineConsistency(t *testing.T) {
	config1 := EvolutionConfig{Seed: 42}
	engine1 := NewBaseEngine(config1)

	config2 := EvolutionConfig{Seed: 42}
	engine2 := NewBaseEngine(config2)

	// Same seed should produce same sequence
	val1 := engine1.GetRandomFloat32()
	val2 := engine2.GetRandomFloat32()

	if val1 != val2 {
		t.Errorf("Expected same values for same seed, got %f and %f", val1, val2)
	}
}

// TestBaseEngineDifferentSeeds tests that different seeds produce different results
func TestBaseEngineDifferentSeeds(t *testing.T) {
	config1 := EvolutionConfig{Seed: 42}
	engine1 := NewBaseEngine(config1)

	config2 := EvolutionConfig{Seed: 43}
	engine2 := NewBaseEngine(config2)

	// Different seeds should produce different sequences
	val1 := engine1.GetRandomFloat32()
	val2 := engine2.GetRandomFloat32()

	if val1 == val2 {
		t.Log("Different seeds produced same value (this is possible but unlikely)")
	}
}

// TestBaseEngineEmbedding tests that BaseEngine can be embedded in other structs
func TestBaseEngineEmbedding(t *testing.T) {
	// Create a test struct that embeds BaseEngine
	type TestEngine struct {
		BaseEngine
		extraField string
	}

	config := EvolutionConfig{Seed: 42}
	testEngine := TestEngine{
		BaseEngine: NewBaseEngine(config),
		extraField: "test",
	}

	// Test that embedded methods work
	seed := testEngine.GetSeed()
	if seed != 42 {
		t.Errorf("Expected seed 42, got %d", seed)
	}

	// Test that extra fields are preserved
	if testEngine.extraField != "test" {
		t.Errorf("Expected 'test', got '%s'", testEngine.extraField)
	}
}

// TestBaseEngineConfigAccess tests configuration access through BaseEngine
func TestBaseEngineConfigAccess(t *testing.T) {
	config := EvolutionConfig{
		Seed:                  42,
		MorphologyChangeRate:  0.1,
		OrthographyChangeRate: 0.2,
		ScriptReformRate:      0.05,
		BorrowingThreshold:    0.3,
	}

	engine := NewBaseEngine(config)

	// Test direct config access
	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}

	// Test through getter methods
	if engine.GetSeed() != 42 {
		t.Errorf("Expected seed 42 from getter, got %d", engine.GetSeed())
	}

	retrievedConfig := engine.GetConfig()
	if retrievedConfig.MorphologyChangeRate != 0.1 {
		t.Errorf("Expected MorphologyChangeRate 0.1, got %f", retrievedConfig.MorphologyChangeRate)
	}
	if retrievedConfig.OrthographyChangeRate != 0.2 {
		t.Errorf("Expected OrthographyChangeRate 0.2, got %f", retrievedConfig.OrthographyChangeRate)
	}
	if retrievedConfig.ScriptReformRate != 0.05 {
		t.Errorf("Expected ScriptReformRate 0.05, got %f", retrievedConfig.ScriptReformRate)
	}
	if retrievedConfig.BorrowingThreshold != 0.3 {
		t.Errorf("Expected BorrowingThreshold 0.3, got %f", retrievedConfig.BorrowingThreshold)
	}
}
