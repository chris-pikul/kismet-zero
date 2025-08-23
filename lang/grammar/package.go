// Package grammar defines the sentence-level construction layer that organizes
// morphological building blocks into structured, grammatically correct phrases and sentences.
//
// This package builds on lang/phoneme, lang/phonology, and lang/morphology by
// providing syntactic organization, agreement systems, and sentence generation.
// It enables the creation of culturally appropriate sentence structures while
// maintaining grammatical consistency and supporting different language types.
//
// The grammar package serves as the bridge between word-level morphology and
// higher-level text generation in the lang system:
//   - phoneme: provides sound units for phonological validation
//   - phonology: ensures syllable and word boundary well-formedness
//   - morphology: supplies words with categories and inflections
//   - orthography: enables written representation of sentences
//   - evolution: applies diachronic grammatical change over time
//
// Responsibilities:
//   - Define syntactic structures and word order patterns (SVO, SOV, VSO, etc.).
//   - Implement agreement systems for case, number, gender, tense, and aspect.
//   - Generate sentences and phrases using morphological building blocks.
//   - Support different language types through flexible rule-based systems.
//   - Maintain phonological and grammatical well-formedness.
//   - Provide weighted selection for grammar rule application.
//
// Boundaries:
//   - Does not define phonemes (delegates to lang/phoneme).
//   - Does not implement phonotactics (delegates to lang/phonology).
//   - Does not construct morphemes or words (delegates to lang/morphology).
//   - Does not create writing systems (delegates to lang/orthography).
//   - Does not model diachronic change (handled in lang/evolution).
//
// This package must remain focused on sentence-level construction and syntactic
// organization. It should provide clean abstractions for different grammatical
// systems while maintaining deterministic behavior through seeded randomness.
package grammar
