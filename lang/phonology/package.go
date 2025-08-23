// Package phonology defines rules for how phonemes combine into syllables and sound patterns.
//
// This package builds on lang/phoneme by providing phonotactic rules, syllable
// templates, and generators for valid syllables. It enforces constraints on how
// consonants and vowels may combine, allowing probabilistic weighting and
// customization for different linguistic systems.
//
// Responsibilities:
//   - Define and apply syllable templates (CV, CVC, CCV, etc.).
//   - Provide a Phonology type encapsulating inventory, rules, and syllable generator.
//   - Enforce phonotactic constraints through composable rules.
//   - Support probabilistic weighting of syllable structures and clusters.
//
// Boundaries:
//   - Does not define phonemes (delegates to lang/phoneme).
//   - Does not construct morphemes or words (handled in lang/morphology).
//   - Does not model diachronic change (handled in lang/evolution).
//
// This package must remain focused on syllable-level constraints and generation.
package phonology
