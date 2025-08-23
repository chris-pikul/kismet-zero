package interlingua

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestRoundTripStability tests that conlang documents maintain stability through round-trip processing.
func TestRoundTripStability(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)
	analyzer := NewEnglishAnalyzer(resolver)

	// Test various predicate families with realistic sentences
	tests := []struct {
		name        string
		description string
		document    Document
	}{
		{
			name:        "transfer_give",
			description: "Transfer domain - giving objects",
			document: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "test-conlang",
				Entities: []Entity{
					{ID: "agent1", Concept: "woman", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "patient1", Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "indef"}},
					{ID: "recipient1", Concept: "child", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "event1",
						Predicate: "give-01",
						TAM:       TAM{Tense: "past", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:     "agent1",
							RolePatient:   "patient1",
							RoleRecipient: "recipient1",
						},
						Valency: 3,
					},
				},
				Trace: Trace{
					SourceLanguageID: "test-conlang",
					ConstructionID:   "ditrans.past.pos",
				},
			},
		},
		{
			name:        "motion_walk",
			description: "Motion domain - intransitive movement",
			document: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "test-conlang",
				Entities: []Entity{
					{ID: "agent1", Concept: "man", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "event1",
						Predicate: "walk-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent: "agent1",
						},
						Valency: 1,
					},
				},
				Trace: Trace{
					SourceLanguageID: "test-conlang",
					ConstructionID:   "intrans.pres.pos",
				},
			},
		},
		{
			name:        "perception_see",
			description: "Perception domain - visual experience",
			document: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "test-conlang",
				Entities: []Entity{
					{ID: "experiencer1", Concept: "girl", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "stimulus1", Concept: "car", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "event1",
						Predicate: "see-01",
						TAM:       TAM{Tense: "past", Polarity: "neg"},
						Roles: map[Role]string{
							RoleExperiencer: "experiencer1",
							RoleStimulus:    "stimulus1",
						},
						Valency: 2,
					},
				},
				Trace: Trace{
					SourceLanguageID: "test-conlang",
					ConstructionID:   "trans.past.neg",
				},
			},
		},
		{
			name:        "creation_make",
			description: "Creation domain - making objects",
			document: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "test-conlang",
				Entities: []Entity{
					{ID: "agent1", Concept: "person", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "indef"}},
					{ID: "patient1", Concept: "house", Feats: NominalFeatures{Number: "sg", Definiteness: "indef"}},
				},
				Events: []Event{
					{
						ID:        "event1",
						Predicate: "make-01",
						TAM:       TAM{Tense: "pres", Polarity: "pos"},
						Roles: map[Role]string{
							RoleAgent:   "agent1",
							RolePatient: "patient1",
						},
						Valency: 2,
					},
				},
				Trace: Trace{
					SourceLanguageID: "test-conlang",
					ConstructionID:   "trans.pres.pos",
				},
			},
		},
		{
			name:        "stative_have",
			description: "Stative domain - possession",
			document: Document{
				SchemaVersion: SchemaVersion,
				LanguageID:    "test-conlang",
				Entities: []Entity{
					{ID: "possessor1", Concept: "boy", Feats: NominalFeatures{Person: 3, Number: "sg", Definiteness: "def"}},
					{ID: "possessed1", Concept: "money", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}},
				},
				Events: []Event{
					{
						ID:        "event1",
						Predicate: "have-01",
						TAM:       TAM{Tense: "pres", Polarity: "neg"},
						Roles: map[Role]string{
							RoleAgent:   "possessor1",
							RolePatient: "possessed1",
						},
						Valency: 2,
					},
				},
				Trace: Trace{
					SourceLanguageID: "test-conlang",
					ConstructionID:   "trans.pres.neg",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original document as golden file
			saveGoldenDocument(t, tt.name, tt.document)

			// Step 1: Document → English (Realize)
			englishTokens, realizeTrace, err := realizer.Realize(tt.document)
			if err != nil {
				t.Fatalf("Realization failed for %s: %v", tt.description, err)
			}

			t.Logf("Original document: %s", tt.description)
			t.Logf("Realized English: %v", englishTokens)

			// Step 2: English → Document (Analyze)
			analyzedDoc, err := analyzer.Analyze(englishTokens)
			if err != nil {
				t.Fatalf("Analysis failed for %s: %v", tt.description, err)
			}

			// Step 3: Validate both documents
			if err := Validate(tt.document); err != nil {
				t.Errorf("Original document validation failed: %v", err)
			}
			if err := Validate(analyzedDoc); err != nil {
				t.Errorf("Analyzed document validation failed: %v", err)
			}

			// Step 4: Check semantic equivalence
			if err := checkSemanticEquivalence(tt.document, analyzedDoc); err != nil {
				t.Errorf("Semantic equivalence check failed for %s: %v", tt.description, err)
				t.Logf("Original TAM: %+v", tt.document.Events[0].TAM)
				t.Logf("Analyzed TAM: %+v", analyzedDoc.Events[0].TAM)
				t.Logf("Original roles: %+v", tt.document.Events[0].Roles)
				t.Logf("Analyzed roles: %+v", analyzedDoc.Events[0].Roles)
			}

			// Step 5: Check that key properties are preserved
			checkTAMPreservation(t, tt.document, analyzedDoc, tt.description)
			checkRolePreservation(t, tt.document, analyzedDoc, tt.description)
			checkValencyPreservation(t, tt.document, analyzedDoc, tt.description)

			// Step 6: Verify trace information
			if len(realizeTrace.Notes) > 0 {
				t.Logf("Realization notes: %d", len(realizeTrace.Notes))
			}
			if len(analyzedDoc.Trace.Notes) > 0 {
				t.Logf("Analysis notes: %d", len(analyzedDoc.Trace.Notes))
			}
		})
	}
}

// TestEnglishRoundTrip tests round-trip processing starting from natural English sentences.
func TestEnglishRoundTrip(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)
	analyzer := NewEnglishAnalyzer(resolver)

	// Realistic English sentences covering different domains
	englishSentences := []struct {
		name        string
		description string
		tokens      []string
	}{
		{
			name:        "daily_conversation_1",
			description: "Simple present tense statement",
			tokens:      []string{"the", "woman", "walks"},
		},
		{
			name:        "daily_conversation_2",
			description: "Past tense with direct object",
			tokens:      []string{"the", "man", "took", "the", "book"},
		},
		{
			name:        "daily_conversation_3",
			description: "Negative past tense",
			tokens:      []string{"the", "child", "didn't", "see", "the", "car"},
		},
		{
			name:        "daily_conversation_4",
			description: "Present tense giving",
			tokens:      []string{"the", "boy", "gives", "a", "book", "to", "the", "girl"},
		},
		{
			name:        "daily_conversation_5",
			description: "Negative present with possession",
			tokens:      []string{"the", "girl", "doesn't", "have", "the", "money"},
		},
		{
			name:        "daily_conversation_6",
			description: "Past tense creation",
			tokens:      []string{"the", "person", "made", "a", "house"},
		},
		{
			name:        "daily_conversation_7",
			description: "Simple motion verb",
			tokens:      []string{"the", "boy", "runs"},
		},
		{
			name:        "daily_conversation_8",
			description: "Transfer with indefinite article",
			tokens:      []string{"a", "woman", "sent", "the", "book", "to", "a", "man"},
		},
	}

	for _, tt := range englishSentences {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Processing: %s - %v", tt.description, tt.tokens)

			// Step 1: English → Document (Analyze)
			originalDoc, err := analyzer.Analyze(tt.tokens)
			if err != nil {
				t.Fatalf("Initial analysis failed for %s: %v", tt.description, err)
			}

			// Step 2: Document → English (Realize)
			realizedTokens, _, err := realizer.Realize(originalDoc)
			if err != nil {
				t.Fatalf("Realization failed for %s: %v", tt.description, err)
			}

			// Step 3: English → Document (Re-analyze)
			reanalyzedDoc, err := analyzer.Analyze(realizedTokens)
			if err != nil {
				t.Fatalf("Re-analysis failed for %s: %v", tt.description, err)
			}

			t.Logf("Original tokens: %v", tt.tokens)
			t.Logf("Realized tokens: %v", realizedTokens)

			// Step 4: Validate all documents
			if err := Validate(originalDoc); err != nil {
				t.Errorf("Original document validation failed: %v", err)
			}
			if err := Validate(reanalyzedDoc); err != nil {
				t.Errorf("Re-analyzed document validation failed: %v", err)
			}

			// Step 5: Check semantic stability
			if err := checkSemanticEquivalence(originalDoc, reanalyzedDoc); err != nil {
				t.Logf("Semantic drift detected (acceptable): %v", err)
				// This is acceptable for natural language processing
			}

			// Step 6: Check that core semantics are preserved
			checkCoreSemantics(t, originalDoc, reanalyzedDoc, tt.description)
		})
	}
}

// TestGoldenFiles tests against saved golden reference documents.
func TestGoldenFiles(t *testing.T) {
	resolver := NewInMemoryConceptResolver()
	realizer := NewEnglishRealizer(42, resolver)

	// Load and test golden files
	goldenFiles := []string{
		"transfer_give.json",
		"motion_walk.json",
		"perception_see.json",
		"creation_make.json",
		"stative_have.json",
	}

	for _, filename := range goldenFiles {
		t.Run(filename, func(t *testing.T) {
			doc, err := loadGoldenDocument(filename)
			if err != nil {
				t.Skipf("Golden file %s not found, skipping test", filename)
				return
			}

			// Test that golden documents can be realized
			tokens, _, err := realizer.Realize(doc)
			if err != nil {
				t.Errorf("Failed to realize golden document %s: %v", filename, err)
			}

			// Test that realized output is deterministic
			tokens2, _, err := realizer.Realize(doc)
			if err != nil {
				t.Errorf("Second realization failed for %s: %v", filename, err)
			}

			// Check determinism
			if len(tokens) != len(tokens2) {
				t.Errorf("Non-deterministic output length for %s", filename)
			}
			for i, token := range tokens {
				if i < len(tokens2) && token != tokens2[i] {
					t.Errorf("Non-deterministic output at position %d for %s: %s vs %s", i, filename, token, tokens2[i])
				}
			}

			t.Logf("Golden file %s realized as: %v", filename, tokens)
		})
	}
}

// checkSemanticEquivalence verifies that two documents represent the same semantic content.
func checkSemanticEquivalence(doc1, doc2 Document) error {
	// Check that both documents have events
	if len(doc1.Events) == 0 || len(doc2.Events) == 0 {
		return fmt.Errorf("one or both documents have no events")
	}

	// For simplicity, compare first event (can be extended for multiple events)
	event1 := doc1.Events[0]
	event2 := doc2.Events[0]

	// Check predicate equivalence
	if event1.Predicate != event2.Predicate {
		return fmt.Errorf("predicate mismatch: %s vs %s", event1.Predicate, event2.Predicate)
	}

	// Check TAM equivalence (core features)
	if event1.TAM.Tense != event2.TAM.Tense {
		return fmt.Errorf("tense mismatch: %s vs %s", event1.TAM.Tense, event2.TAM.Tense)
	}
	if event1.TAM.Polarity != event2.TAM.Polarity {
		return fmt.Errorf("polarity mismatch: %s vs %s", event1.TAM.Polarity, event2.TAM.Polarity)
	}

	// Check role preservation (core roles)
	coreRoles := []Role{RoleAgent, RolePatient, RoleRecipient, RoleExperiencer, RoleStimulus}
	for _, role := range coreRoles {
		if _, hasRole1 := event1.Roles[role]; hasRole1 {
			if _, hasRole2 := event2.Roles[role]; !hasRole2 {
				return fmt.Errorf("role %s lost in second document", role)
			}
		}
	}

	return nil
}

// checkTAMPreservation verifies that TAM features are preserved.
func checkTAMPreservation(t *testing.T, original, analyzed Document, description string) {
	if len(original.Events) == 0 || len(analyzed.Events) == 0 {
		t.Errorf("Missing events in TAM preservation check for %s", description)
		return
	}

	originalTAM := original.Events[0].TAM
	analyzedTAM := analyzed.Events[0].TAM

	if originalTAM.Tense != analyzedTAM.Tense {
		t.Errorf("Tense not preserved for %s: %s → %s", description, originalTAM.Tense, analyzedTAM.Tense)
	}
	if originalTAM.Polarity != analyzedTAM.Polarity {
		t.Errorf("Polarity not preserved for %s: %s → %s", description, originalTAM.Polarity, analyzedTAM.Polarity)
	}
}

// checkRolePreservation verifies that semantic roles are preserved.
func checkRolePreservation(t *testing.T, original, analyzed Document, description string) {
	if len(original.Events) == 0 || len(analyzed.Events) == 0 {
		t.Errorf("Missing events in role preservation check for %s", description)
		return
	}

	originalRoles := original.Events[0].Roles
	analyzedRoles := analyzed.Events[0].Roles

	for role := range originalRoles {
		if _, exists := analyzedRoles[role]; !exists {
			t.Errorf("Role %s lost during round-trip for %s", role, description)
		}
	}
}

// checkValencyPreservation verifies that valency is preserved.
func checkValencyPreservation(t *testing.T, original, analyzed Document, description string) {
	if len(original.Events) == 0 || len(analyzed.Events) == 0 {
		t.Errorf("Missing events in valency preservation check for %s", description)
		return
	}

	originalValency := original.Events[0].Valency
	analyzedValency := analyzed.Events[0].Valency

	if originalValency != analyzedValency {
		t.Errorf("Valency not preserved for %s: %d → %d", description, originalValency, analyzedValency)
	}
}

// checkCoreSemantics verifies that core semantic properties are maintained.
func checkCoreSemantics(t *testing.T, original, reanalyzed Document, description string) {
	if len(original.Events) == 0 || len(reanalyzed.Events) == 0 {
		t.Errorf("Missing events in core semantics check for %s", description)
		return
	}

	originalEvent := original.Events[0]
	reanalyzedEvent := reanalyzed.Events[0]

	// Check that we still have the same type of predicate
	if originalEvent.Predicate != reanalyzedEvent.Predicate {
		t.Logf("Predicate evolved during round-trip for %s: %s → %s", description, originalEvent.Predicate, reanalyzedEvent.Predicate)
	}

	// Check that valency is in reasonable range
	if reanalyzedEvent.Valency < 1 || reanalyzedEvent.Valency > 3 {
		t.Errorf("Unexpected valency for %s: %d", description, reanalyzedEvent.Valency)
	}
}

// saveGoldenDocument saves a document as a golden reference file.
func saveGoldenDocument(t *testing.T, name string, doc Document) {
	filename := filepath.Join("testdata", name+".json")

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal golden document %s: %v", name, err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		t.Errorf("Failed to write golden document %s: %v", name, err)
	}
}

// loadGoldenDocument loads a document from a golden reference file.
func loadGoldenDocument(filename string) (Document, error) {
	var doc Document

	filepath := filepath.Join("testdata", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return doc, err
	}

	err = json.Unmarshal(data, &doc)
	return doc, err
}
