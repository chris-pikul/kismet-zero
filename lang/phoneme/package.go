// Package phoneme defines the fundamental sound units of language generation.
//
// This package provides strongly typed representations of consonants and vowels,
// including their articulatory features (manner, place, voicing, height, backness, rounding).
// Phonemes are grouped into pools that allow weighted random selection while
// maintaining deduplication guarantees.
//
// The phoneme package serves as the foundation for higher-level linguistic
// modules in the lang system:
//   - phonology: builds rules for valid phoneme combinations
//   - morphology: constructs words from morphemes
//   - grammar: organizes phrases and sentences
//   - evolution: applies diachronic change and cultural blending
//
// Responsibilities:
//   - Define typed enums for consonant and vowel features.
//   - Provide the Phoneme struct with feature metadata.
//   - Expose utilities for managing phoneme pools.
//   - Ensure deterministic behavior through seeded randomness in tests.
//
// This package must remain low-level and composable. It should not implement
// phonotactics, morphology, or diachronic change. Those are the responsibility
// of higher-level packages.
package phoneme
