package interlingua

import (
	"testing"
)

func TestNewEnglishAnalyzer(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created")
	}

	if analyzer.LanguageCode() != "en" {
		t.Errorf("Expected language code 'en', got %s", analyzer.LanguageCode())
	}
}

func TestEnglishAnalyzerBasicClauses(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	tests := []struct {
		name     string
		tokens   []string
		expected Document
	}{
		{
			name:   "simple transitive positive present",
			tokens: []string{"the", "boy", "gives", "the", "book"},
			expected: Document{
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
						Valency: 2,
					},
				},
			},
		},
		{
			name:   "simple transitive negative past",
			tokens: []string{"didn't", "give", "the", "boy", "the", "book"},
			expected: Document{
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
						Valency: 2,
					},
				},
			},
		},
		{
			name:   "ditransitive positive past",
			tokens: []string{"gave", "the", "boy", "the", "book", "to", "the", "girl"},
			expected: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "en",
				Entities: []Entity{
					{ID: "e1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "e2", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
					{ID: "e3", Concept: "girl", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
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
						Valency: 3,
					},
				},
			},
		},
		{
			name:   "intransitive positive present",
			tokens: []string{"the", "boy", "walks"},
			expected: Document{
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
						Valency: 1,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := analyzer.Analyze(tt.tokens)
			if err != nil {
				t.Fatalf("Analysis failed: %v", err)
			}

			// Check basic structure
			if doc.SchemaVersion != tt.expected.SchemaVersion {
				t.Errorf("Expected schema version %s, got %s", tt.expected.SchemaVersion, doc.SchemaVersion)
			}
			if doc.LanguageID != tt.expected.LanguageID {
				t.Errorf("Expected language ID %s, got %s", tt.expected.LanguageID, doc.LanguageID)
			}

			// Check entities
			if len(doc.Entities) != len(tt.expected.Entities) {
				t.Errorf("Expected %d entities, got %d", len(tt.expected.Entities), len(doc.Entities))
			}

			// Check events
			if len(doc.Events) != len(tt.expected.Events) {
				t.Errorf("Expected %d events, got %d", len(tt.expected.Events), len(doc.Events))
			}

			if len(doc.Events) > 0 && len(tt.expected.Events) > 0 {
				event := doc.Events[0]
				expectedEvent := tt.expected.Events[0]

				// Check predicate
				if event.Predicate != expectedEvent.Predicate {
					t.Errorf("Expected predicate %s, got %s", expectedEvent.Predicate, event.Predicate)
				}

				// Check TAM
				if event.TAM.Tense != expectedEvent.TAM.Tense {
					t.Errorf("Expected tense %s, got %s", expectedEvent.TAM.Tense, event.TAM.Tense)
				}
				if event.TAM.Polarity != expectedEvent.TAM.Polarity {
					t.Errorf("Expected polarity %s, got %s", expectedEvent.TAM.Polarity, event.TAM.Polarity)
				}

				// Check valency
				if event.Valency != expectedEvent.Valency {
					t.Errorf("Expected valency %d, got %d", expectedEvent.Valency, event.Valency)
				}

				// Check roles
				if len(event.Roles) != len(expectedEvent.Roles) {
					t.Errorf("Expected %d roles, got %d", len(expectedEvent.Roles), len(event.Roles))
				}
			}

			// Check trace
			if doc.Trace.SourceLanguageID != "en" {
				t.Errorf("Expected source language 'en', got %s", doc.Trace.SourceLanguageID)
			}
		})
	}
}

func TestEnglishAnalyzerTAMDetection(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	tests := []struct {
		name             string
		tokens           []string
		expectedTense    string
		expectedPolarity string
	}{
		{
			name:             "present positive",
			tokens:           []string{"gives", "the", "boy", "the", "book"},
			expectedTense:    "pres",
			expectedPolarity: "pos",
		},
		{
			name:             "past positive",
			tokens:           []string{"gave", "the", "boy", "the", "book"},
			expectedTense:    "past",
			expectedPolarity: "pos",
		},
		{
			name:             "present negative",
			tokens:           []string{"doesn't", "the", "boy", "the", "book"},
			expectedTense:    "pres",
			expectedPolarity: "neg",
		},
		{
			name:             "past negative",
			tokens:           []string{"didn't", "the", "boy", "the", "book"},
			expectedTense:    "past",
			expectedPolarity: "neg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := analyzer.Analyze(tt.tokens)
			if err != nil {
				t.Fatalf("Analysis failed: %v", err)
			}

			if len(doc.Events) == 0 {
				t.Fatal("No events found in document")
			}

			event := doc.Events[0]
			if event.TAM.Tense != tt.expectedTense {
				t.Errorf("Expected tense %s, got %s", tt.expectedTense, event.TAM.Tense)
			}
			if event.TAM.Polarity != tt.expectedPolarity {
				t.Errorf("Expected polarity %s, got %s", tt.expectedPolarity, event.TAM.Polarity)
			}
		})
	}
}

func TestEnglishAnalyzerEntityFeatures(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	tests := []struct {
		name             string
		tokens           []string
		expectedFeatures map[string]NominalFeatures
	}{
		{
			name:   "definite singular",
			tokens: []string{"the", "boy", "walks"},
			expectedFeatures: map[string]NominalFeatures{
				"boy": {Person: 3, Number: "sg", Definiteness: "def"},
			},
		},
		{
			name:   "indefinite singular",
			tokens: []string{"a", "boy", "walks"},
			expectedFeatures: map[string]NominalFeatures{
				"boy": {Person: 3, Number: "sg", Definiteness: "indef"},
			},
		},
		{
			name:   "definite plural",
			tokens: []string{"the", "boys", "walk"},
			expectedFeatures: map[string]NominalFeatures{
				"boys": {Person: 3, Number: "pl", Definiteness: "def"},
			},
		},
		{
			name:   "pronoun first person",
			tokens: []string{"I", "walk"},
			expectedFeatures: map[string]NominalFeatures{
				"I": {Person: 1, Number: "sg"},
			},
		},
		{
			name:   "pronoun second person",
			tokens: []string{"you", "walk"},
			expectedFeatures: map[string]NominalFeatures{
				"you": {Person: 2, Number: "sg"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := analyzer.Analyze(tt.tokens)
			if err != nil {
				t.Fatalf("Analysis failed: %v", err)
			}

			// Check entity features
			for expectedNoun, expectedFeats := range tt.expectedFeatures {
				found := false
				for _, entity := range doc.Entities {
					if entity.LemmaHint == expectedNoun {
						found = true
						if entity.Feats.Person != expectedFeats.Person {
							t.Errorf("Expected person %d for %s, got %d", expectedFeats.Person, expectedNoun, entity.Feats.Person)
						}
						if entity.Feats.Number != expectedFeats.Number {
							t.Errorf("Expected number %s for %s, got %s", expectedFeats.Number, expectedNoun, entity.Feats.Number)
						}
						if entity.Feats.Definiteness != expectedFeats.Definiteness {
							t.Errorf("Expected definiteness %s for %s, got %s", expectedFeats.Definiteness, expectedNoun, entity.Feats.Definiteness)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected entity with noun %s not found", expectedNoun)
				}
			}
		})
	}
}

func TestEnglishAnalyzerAdjuncts(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	tests := []struct {
		name          string
		tokens        []string
		expectedNotes []string
	}{
		{
			name:          "time adjunct",
			tokens:        []string{"gave", "the", "boy", "the", "book", "yesterday"},
			expectedNotes: []string{"time.adjunct.found"},
		},
		{
			name:          "location adjunct",
			tokens:        []string{"gave", "the", "boy", "the", "book", "in", "the", "house"},
			expectedNotes: []string{"location.adjunct.found"},
		},
		{
			name:          "both adjuncts",
			tokens:        []string{"gave", "the", "boy", "the", "book", "yesterday", "in", "the", "house"},
			expectedNotes: []string{"time.adjunct.found", "location.adjunct.found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := analyzer.Analyze(tt.tokens)
			if err != nil {
				t.Fatalf("Analysis failed: %v", err)
			}

			// Check for expected notes
			for _, expectedCode := range tt.expectedNotes {
				found := false
				for _, note := range doc.Trace.Notes {
					if note.Code == expectedCode {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected note with code %s not found", expectedCode)
				}
			}
		})
	}
}

func TestEnglishAnalyzerErrorHandling(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	tests := []struct {
		name        string
		tokens      []string
		expectError bool
	}{
		{
			name:        "empty tokens",
			tokens:      []string{},
			expectError: true,
		},
		{
			name:        "insufficient tokens",
			tokens:      []string{"the", "boy"},
			expectError: true,
		},
		{
			name:        "valid tokens",
			tokens:      []string{"the", "boy", "walks"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := analyzer.Analyze(tt.tokens)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestEnglishAnalyzerConceptResolution(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)

	tests := []struct {
		name              string
		tokens            []string
		expectedPredicate ConceptID
	}{
		{
			name:              "known verb",
			tokens:            []string{"gave", "the", "boy", "the", "book"},
			expectedPredicate: "give-01",
		},
		{
			name:              "known noun",
			tokens:            []string{"the", "boy", "walks"},
			expectedPredicate: "walk-01",
		},
		{
			name:              "unknown verb",
			tokens:            []string{"xyzzed", "the", "boy"},
			expectedPredicate: "xyzzed-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := analyzer.Analyze(tt.tokens)
			if err != nil {
				t.Fatalf("Analysis failed: %v", err)
			}

			if len(doc.Events) == 0 {
				t.Fatal("No events found in document")
			}

			event := doc.Events[0]
			if event.Predicate != tt.expectedPredicate {
				t.Errorf("Expected predicate %s, got %s", tt.expectedPredicate, event.Predicate)
			}
		})
	}
}

func TestEnglishAnalyzerRoundTrip(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	analyzer := NewEnglishAnalyzer(resolver)
	realizer := NewEnglishRealizer(42, resolver)

	// Test round-trip: English → Interlingua → English
	originalTokens := []string{"gave", "the", "boy", "the", "book", "to", "the", "girl"}

	// Analyze: English → Interlingua
	doc, err := analyzer.Analyze(originalTokens)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Realize: Interlingua → English
	realizedTokens, _, err := realizer.Realize(doc)
	if err != nil {
		t.Fatalf("Realization failed: %v", err)
	}

	// Check that we got some tokens back
	if len(realizedTokens) == 0 {
		t.Error("Realization returned no tokens")
	}

	// Check that the document is valid
	if err := Validate(doc); err != nil {
		t.Errorf("Generated document validation failed: %v", err)
	}

	// Check that we have the expected structure
	if len(doc.Entities) < 3 {
		t.Errorf("Expected at least 3 entities, got %d", len(doc.Entities))
	}
	if len(doc.Events) < 1 {
		t.Errorf("Expected at least 1 event, got %d", len(doc.Events))
	}

	event := doc.Events[0]
	if event.Predicate != "give-01" {
		t.Errorf("Expected predicate 'give-01', got %s", event.Predicate)
	}
	if event.TAM.Tense != "past" {
		t.Errorf("Expected tense 'past', got %s", event.TAM.Tense)
	}
	if event.TAM.Polarity != "pos" {
		t.Errorf("Expected polarity 'pos', got %s", event.TAM.Polarity)
	}
	if event.Valency != 3 {
		t.Errorf("Expected valency 3, got %d", event.Valency)
	}
}
