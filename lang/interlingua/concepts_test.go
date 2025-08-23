package interlingua

import (
	"testing"
)

func TestNewInMemoryConceptResolver(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	if resolver == nil {
		t.Fatal("Expected resolver to be created")
	}

	// Check that English concepts are loaded
	enCount := resolver.GetConceptCount("en")
	if enCount == 0 {
		t.Error("Expected English concepts to be loaded")
	}

	// Check that other languages are not loaded initially
	frCount := resolver.GetConceptCount("fr")
	if frCount != 0 {
		t.Errorf("Expected 0 French concepts, got %d", frCount)
	}
}

func TestConceptResolverResolve(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Test successful resolution
	conceptID, found := resolver.Resolve("give", "en")
	if !found {
		t.Error("Expected to find concept for 'give'")
	}
	if conceptID != "give-01" {
		t.Errorf("Expected concept ID 'give-01', got %s", conceptID)
	}

	// Test successful resolution for different concept types
	conceptID, found = resolver.Resolve("boy", "en")
	if !found {
		t.Error("Expected to find concept for 'boy'")
	}
	if conceptID != "boy-01" {
		t.Errorf("Expected concept ID 'boy-01', got %s", conceptID)
	}

	// Test resolution for motion concepts
	conceptID, found = resolver.Resolve("move", "en")
	if !found {
		t.Error("Expected to find concept for 'move'")
	}
	if conceptID != "move-00" {
		t.Errorf("Expected concept ID 'move-00', got %s", conceptID)
	}
}

func TestConceptResolverResolveNotFound(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Test non-existent lemma
	conceptID, found := resolver.Resolve("nonexistent", "en")
	if found {
		t.Error("Expected not to find concept for 'nonexistent'")
	}
	if conceptID != "" {
		t.Errorf("Expected empty concept ID, got %s", conceptID)
	}

	// Test non-existent language
	conceptID, found = resolver.Resolve("give", "fr")
	if found {
		t.Error("Expected not to find concept for 'give' in French")
	}
	if conceptID != "" {
		t.Errorf("Expected empty concept ID, got %s", conceptID)
	}
}

func TestConceptResolverRegister(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Register a new concept
	resolver.Register("bonjour", "fr", "hello-01")

	// Verify it can be resolved
	conceptID, found := resolver.Resolve("bonjour", "fr")
	if !found {
		t.Error("Expected to find registered concept 'bonjour'")
	}
	if conceptID != "hello-01" {
		t.Errorf("Expected concept ID 'hello-01', got %s", conceptID)
	}

	// Verify English concepts are still available
	conceptID, found = resolver.Resolve("give", "en")
	if !found {
		t.Error("Expected English concepts to still be available")
	}

	// Verify French count increased
	frCount := resolver.GetConceptCount("fr")
	if frCount != 1 {
		t.Errorf("Expected 1 French concept, got %d", frCount)
	}
}

func TestConceptResolverGetSupportedLanguages(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	languages := resolver.GetSupportedLanguages()

	// Should have at least English
	if len(languages) < 1 {
		t.Error("Expected at least one supported language")
	}

	// Check that English is included
	hasEnglish := false
	for _, lang := range languages {
		if lang == "en" {
			hasEnglish = true
			break
		}
	}
	if !hasEnglish {
		t.Error("Expected English to be in supported languages")
	}
}

func TestConceptResolverGetConceptCount(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Test English count
	enCount := resolver.GetConceptCount("en")
	if enCount < 50 { // Should have many concepts
		t.Errorf("Expected many English concepts, got %d", enCount)
	}

	// Test non-existent language
	frCount := resolver.GetConceptCount("fr")
	if frCount != 0 {
		t.Errorf("Expected 0 French concepts, got %d", frCount)
	}

	// Test after registration
	resolver.Register("bonjour", "fr", "hello-01")
	frCount = resolver.GetConceptCount("fr")
	if frCount != 1 {
		t.Errorf("Expected 1 French concept after registration, got %d", frCount)
	}
}

func TestConceptResolverCoreDomains(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Test transfer domain
	transferConcepts := []string{"give", "send", "hand", "pass", "lend", "rent", "sell", "buy", "owe"}
	for _, concept := range transferConcepts {
		conceptID, found := resolver.Resolve(concept, "en")
		if !found {
			t.Errorf("Expected to find transfer concept '%s'", concept)
		}
		if conceptID == "" {
			t.Errorf("Expected non-empty concept ID for '%s'", concept)
		}
	}

	// Test motion domain
	motionConcepts := []string{"move", "go", "come", "walk", "run", "fly", "swim", "jump", "climb"}
	for _, concept := range motionConcepts {
		conceptID, found := resolver.Resolve(concept, "en")
		if !found {
			t.Errorf("Expected to find motion concept '%s'", concept)
		}
		if conceptID == "" {
			t.Errorf("Expected non-empty concept ID for '%s'", concept)
		}
	}

	// Test perception domain
	perceptionConcepts := []string{"see", "look", "watch", "hear", "listen", "feel", "touch", "smell", "taste"}
	for _, concept := range perceptionConcepts {
		conceptID, found := resolver.Resolve(concept, "en")
		if !found {
			t.Errorf("Expected to find perception concept '%s'", concept)
		}
		if conceptID == "" {
			t.Errorf("Expected non-empty concept ID for '%s'", concept)
		}
	}
}

func TestConceptResolverEntityConcepts(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Test common entity concepts
	entityConcepts := []string{"boy", "girl", "man", "woman", "person", "child", "book", "house", "car", "tree"}
	for _, concept := range entityConcepts {
		conceptID, found := resolver.Resolve(concept, "en")
		if !found {
			t.Errorf("Expected to find entity concept '%s'", concept)
		}
		if conceptID == "" {
			t.Errorf("Expected non-empty concept ID for '%s'", concept)
		}
	}
}

func TestConceptResolverTimeConcepts(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Test time-related concepts
	timeConcepts := []string{"yesterday", "today", "tomorrow", "time", "day", "night", "year"}
	for _, concept := range timeConcepts {
		conceptID, found := resolver.Resolve(concept, "en")
		if !found {
			t.Errorf("Expected to find time concept '%s'", concept)
		}
		if conceptID == "" {
			t.Errorf("Expected non-empty concept ID for '%s'", concept)
		}
	}
}
