package interlingua

import (
	"testing"
)

func TestNewEnglishRealizer(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	if realizer == nil {
		t.Fatal("Expected realizer to be created")
	}

	if realizer.LanguageCode() != "en" {
		t.Errorf("Expected language code 'en', got %s", realizer.LanguageCode())
	}
}

func TestEnglishRealizerBasicClauses(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	tests := []struct {
		name     string
		doc      Document
		expected []string
	}{
		{
			name: "simple transitive positive present",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
					},
				},
			},
			expected: []string{"the", "boy", "gives", "the", "book"},
		},
		{
			name: "simple transitive negative past",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "past", Polarity: "neg"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
					},
				},
			},
			expected: []string{"the", "boy", "didn't", "give", "the", "book"},
		},
		{
			name: "ditransitive positive past",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
					{ID: "e3", Concept: "girl", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "past", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:     "e1",
							RolePatient:   "e2",
							RoleRecipient: "e3",
						},
					},
				},
			},
			expected: []string{"the", "boy", "gave", "the", "book", "to", "the", "girl"},
		},
		{
			name: "intransitive positive present",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "walk-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent: "e1",
						},
					},
				},
			},
			expected: []string{"the", "boy", "walks"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, trace, err := realizer.Realize(tt.doc)
			if err != nil {
				t.Fatalf("Realization failed: %v", err)
			}

			// Check token count
			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got:      %v", tokens)
			}

			// Check token content
			for i, expected := range tt.expected {
				if i < len(tokens) && tokens[i] != expected {
					t.Errorf("Token %d: expected '%s', got '%s'", i, expected, tokens[i])
				}
			}

			// Check trace
			if trace.SourceLanguageID != tt.doc.Trace.SourceLanguageID {
				t.Errorf("Expected source language '%s', got '%s'", tt.doc.Trace.SourceLanguageID, trace.SourceLanguageID)
			}
		})
	}
}

func TestEnglishRealizerAdjuncts(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	tests := []struct {
		name     string
		doc      Document
		expected []string
	}{
		{
			name: "time adjunct",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
					{ID: "t1", Concept: "yesterday"},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "past", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
						Adjuncts: map[Role]string{
							RoleTime: "t1",
						},
					},
				},
			},
			expected: []string{"the", "boy", "gave", "the", "book", "yesterday"},
		},
		{
			name: "location adjunct",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
					{ID: "l1", Concept: "house", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
						Adjuncts: map[Role]string{
							RoleLocation: "l1",
						},
					},
				},
			},
			expected: []string{"the", "boy", "gives", "the", "book", "in", "the", "house"},
		},
		{
			name: "manner adjunct",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
					{ID: "m1", Concept: "care", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
						Adjuncts: map[Role]string{
							RoleManner: "m1",
						},
					},
				},
			},
			expected: []string{"the", "boy", "gives", "the", "book", "with", "the", "care"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, trace, err := realizer.Realize(tt.doc)
			if err != nil {
				t.Fatalf("Realization failed: %v", err)
			}

			// Check token count
			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got:      %v", tokens)
			}

			// Check token content
			for i, expected := range tt.expected {
				if i < len(tokens) && tokens[i] != expected {
					t.Errorf("Token %d: expected '%s', got '%s'", i, expected, tokens[i])
				}
			}

			// Check that adjuncts are properly mapped in trace
			if len(trace.TokenToNode) == 0 {
				t.Error("Expected token-to-node mapping in trace")
			}
		})
	}
}

func TestEnglishRealizerDeterminism(t *testing.T) {
	resolver := NewInMemoryConceptResolver()

	// Create a test document
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{ID: "e1", Concept: "boy", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
			{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "give-01",
				TAM:       TAM{Tense: "past", Polarity: "pos"},
				Roles: map[Role]string{
					RoleAgent:   "e1",
					RolePatient: "e2",
				},
			},
		},
	}

	// Test multiple runs with same seed
	realizer1 := NewEnglishRealizer(42, resolver)
	realizer2 := NewEnglishRealizer(42, resolver)

	tokens1, _, err := realizer1.Realize(doc)
	if err != nil {
		t.Fatalf("First realization failed: %v", err)
	}

	tokens2, _, err := realizer2.Realize(doc)
	if err != nil {
		t.Fatalf("Second realization failed: %v", err)
	}

	// Results should be identical
	if len(tokens1) != len(tokens2) {
		t.Errorf("Token count mismatch: %d vs %d", len(tokens1), len(tokens2))
	}

	for i, token1 := range tokens1 {
		if i < len(tokens2) && token1 != tokens2[i] {
			t.Errorf("Token %d mismatch: '%s' vs '%s'", i, token1, tokens2[i])
		}
	}
}

func TestEnglishRealizerFallbackBehavior(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	// Test with unknown predicate
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "unknown-predicate-01",
				TAM:       TAM{Tense: "pres", Polarity: "pos"},
				Roles: map[Role]string{
					RoleAgent: "e1",
				},
			},
		},
	}

	tokens, trace, err := realizer.Realize(doc)
	if err != nil {
		t.Fatalf("Realization failed: %v", err)
	}

	// Should have fallback behavior
	if len(tokens) == 0 {
		t.Error("Expected fallback tokens")
	}

	// Should have warning note about fallback
	warnings := trace.GetNotesByCode("verb.fallback")
	if len(warnings) == 0 {
		t.Error("Expected warning note about verb fallback")
	}
}

func TestEnglishRealizerPluralization(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	tests := []struct {
		name     string
		doc      Document
		expected []string
	}{
		{
			name: "singular definite",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "walk-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent: "e1",
						},
					},
				},
			},
			expected: []string{"the", "boy", "walks"},
		},
		{
			name: "plural definite",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "pl", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "walk-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent: "e1",
						},
					},
				},
			},
			expected: []string{"the", "boys", "walk"},
		},
		{
			name: "singular indefinite",
			doc: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "indef"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "walk-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent: "e1",
						},
					},
				},
			},
			expected: []string{"a", "boy", "walks"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _, err := realizer.Realize(tt.doc)
			if err != nil {
				t.Fatalf("Realization failed: %v", err)
			}

			// Check token count
			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got:      %v", tokens)
			}

			// Check token content
			for i, expected := range tt.expected {
				if i < len(tokens) && tokens[i] != expected {
					t.Errorf("Token %d: expected '%s', got '%s'", i, expected, tokens[i])
				}
			}
		})
	}
}

func TestEnglishRealizerREADMEExample(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	// Recreate the README example exactly
	// Surface (conlang, glossed): `neg-give.pst boy-nom book-acc girl-dat yesterday`
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "elv-wood-sindarin",
		Entities: []Entity{
			{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
			{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
			{ID: "e3", Concept: "girl", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
			{ID: "t1", Concept: "yesterday"},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "give-01",
				TAM:       TAM{Tense: "past", Polarity: "neg"},
				Roles: map[Role]string{
					RoleAgent:     "e1",
					RolePatient:   "e2",
					RoleRecipient: "e3",
				},
				Adjuncts: map[Role]string{
					RoleTime: "t1",
				},
			},
		},
		Trace: Trace{
			SourceLanguageID: "elv-wood-sindarin",
			ConstructionID:   "decl.transitive.svo.neg.pst",
		},
	}

	tokens, trace, err := realizer.Realize(doc)
	if err != nil {
		t.Fatalf("README example realization failed: %v", err)
	}

	// Expected: "the boy didn't give the book to the girl yesterday"
	expected := []string{"the", "boy", "didn't", "give", "the", "book", "to", "the", "girl", "yesterday"}

	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
		t.Errorf("Expected: %v", expected)
		t.Errorf("Got:      %v", tokens)
	}

	// Check token content
	for i, exp := range expected {
		if i < len(tokens) && tokens[i] != exp {
			t.Errorf("Token %d: expected '%s', got '%s'", i, exp, tokens[i])
		}
	}

	// Check trace
	if trace.SourceLanguageID != "elv-wood-sindarin" {
		t.Errorf("Expected source language 'elv-wood-sindarin', got '%s'", trace.SourceLanguageID)
	}
	if trace.ConstructionID != "decl.transitive.svo.neg.pst" {
		t.Errorf("Expected construction ID 'decl.transitive.svo.neg.pst', got '%s'", trace.ConstructionID)
	}

	// Check token mapping
	if len(trace.TokenToNode) == 0 {
		t.Error("Expected token-to-node mapping in trace")
	}
}

func TestEnglishRealizerVerbTenses(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	tests := []struct {
		name     string
		tense    string
		expected string
	}{
		{"present", "pres", "gives"},
		{"past", "past", "gave"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: tt.tense, Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
					},
				},
			}

			tokens, _, err := realizer.Realize(doc)
			if err != nil {
				t.Fatalf("Realization failed: %v", err)
			}

			// Find the verb token
			var verb string
			for _, token := range tokens {
				if token == "gives" || token == "gave" {
					verb = token
					break
				}
			}

			if verb != tt.expected {
				t.Errorf("Expected verb '%s', got '%s'", tt.expected, verb)
			}
		})
	}
}

func TestEnglishRealizerArticleSelection(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	tests := []struct {
		name     string
		lemma    string
		expected string
	}{
		{"consonant start", "book", "a"},
		{"vowel start", "apple", "an"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: ConceptID(tt.lemma), Feats: NominalFeatures{Number: "sg", Definiteness: "indef"}, LemmaHint: tt.lemma},
				},
				Events: []Event{
					{
						ID:        "v1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "e1",
							RolePatient: "e2",
						},
					},
				},
			}

			tokens, _, err := realizer.Realize(doc)
			if err != nil {
				t.Fatalf("Realization failed: %v", err)
			}

			// Find the article token
			var article string
			for i, token := range tokens {
				if token == "a" || token == "an" {
					article = token
					// Check that it's followed by the noun
					if i+1 < len(tokens) && tokens[i+1] == tt.lemma {
						break
					}
				}
			}

			if article != tt.expected {
				t.Errorf("Expected article '%s', got '%s'", tt.expected, article)
			}
		})
	}
}
