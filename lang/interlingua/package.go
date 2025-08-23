// Package interlingua provides a semantic pivot for bi-directional translation
// between procedurally generated languages and user-facing languages.
//
// The interlingua package serves as a semantic layer that enables:
//   - Emission of Interlingua Documents during conlang generation
//   - Realization of Documents to target languages (English, French, etc.)
//   - Analysis of source languages into Interlingua Documents
//   - Loss-aware translation with feature preservation and diagnostics
//
// Core Components:
//   - Document: Semantic graph representation of sentences
//   - Events: Predicates with TAM, roles, and adjuncts
//   - Entities: Referents with nominal features and discourse properties
//   - Trace: Alignment metadata and diagnostic notes
//   - Validators: Document coherence and role resolution checks
//
// The package integrates with existing lang components via opt-in hooks,
// maintaining determinism through Language.Seed and providing comprehensive
// diagnostics for translation quality assessment.
//
// Example usage:
//
//	doc := interlingua.NewDocument("en")
//	entityID := interlingua.AddEntity(doc, interlingua.Entity{
//		Concept: "boy",
//		Feats: interlingua.NominalFeatures{Number: "sg", Definiteness: "def"},
//	})
//	eventID := interlingua.AddEvent(doc, interlingua.Event{
//		Predicate: "give-01",
//		TAM: interlingua.TAM{Tense: "past", Polarity: "pos"},
//		Valency: 3,
//	})
//	interlingua.SetRole(doc.Events[0], interlingua.RoleAgent, entityID)
//
//	realizer := interlingua.NewEnglishRealizer(42, nil)
//	tokens, trace, err := realizer.Realize(*doc)
package interlingua
