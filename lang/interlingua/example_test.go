package interlingua

import (
	"fmt"
)

// ExampleDocument demonstrates creating a new interlingua document.
func ExampleDocument() {
	// Create a new document
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{
				ID:      "e1",
				Concept: "boy",
				Feats: NominalFeatures{
					Number:       "sg",
					Definiteness: "def",
				},
			},
			{
				ID:      "e2",
				Concept: "book",
				Feats: NominalFeatures{
					Number:       "sg",
					Definiteness: "def",
				},
			},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "give-01",
				TAM: TAM{
					Tense:    "past",
					Polarity: "pos",
				},
				Roles: map[Role]string{
					RoleAgent:   "e1",
					RolePatient: "e2",
				},
				Valency: 2,
			},
		},
		Trace: Trace{
			SourceLanguageID: "en",
			ConstructionID:   "decl.transitive.svo.pos.pst",
			TokenToNode: map[int]string{
				0: "v1",
				1: "e1",
				2: "e2",
			},
		},
	}

	// Validate the document
	if err := Validate(doc); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	// Validate the document (JSON marshaling would go here in real usage)
	fmt.Printf("Document created successfully\n")
	// Output: Document created successfully
}

// ExampleConceptResolver demonstrates using the concept resolver.
func ExampleConceptResolver() {
	resolver := NewInMemoryConceptResolver()

	// Resolve some concepts
	concepts := []string{"give", "boy", "book", "move"}
	for _, lemma := range concepts {
		if conceptID, found := resolver.Resolve(lemma, "en"); found {
			fmt.Printf("'%s' -> %s\n", lemma, conceptID)
		} else {
			fmt.Printf("'%s' -> not found\n", lemma)
		}
	}

	// Register a new concept
	resolver.Register("bonjour", "fr", "hello-01")
	if conceptID, found := resolver.Resolve("bonjour", "fr"); found {
		fmt.Printf("'bonjour' -> %s\n", conceptID)
	}

	// Output:
	// 'give' -> give-01
	// 'boy' -> boy-01
	// 'book' -> book-01
	// 'move' -> move-00
	// 'bonjour' -> hello-01
}

// ExampleTrace demonstrates using trace helpers for diagnostics.
func ExampleTrace() {
	trace := &Trace{}

	// Set metadata
	trace.SetSourceLanguage("en")
	trace.SetConstructionID("decl.transitive.svo.neg.pst")

	// Add diagnostic notes
	trace.AddNote("info", "construction.selected", "Selected SVO transitive construction")
	trace.AddNote("warn", "feature.defaulted", "Number defaulted to singular")
	trace.AddNote("loss", "feature.dropped.evidential", "Evidentiality not supported in target language")

	// Map tokens to nodes
	trace.MapToken(0, "v1")
	trace.MapToken(1, "e1")
	trace.MapToken(2, "e2")

	// Check for specific types of notes
	if trace.HasNotes() {
		fmt.Printf("Trace has %d notes\n", len(trace.Notes))

		warnings := trace.GetNotesBySeverity("warn")
		fmt.Printf("Found %d warnings\n", len(warnings))

		losses := trace.GetNotesBySeverity("loss")
		fmt.Printf("Found %d feature losses\n", len(losses))
	}

	// Output:
	// Trace has 3 notes
	// Found 1 warnings
	// Found 1 feature losses
}
