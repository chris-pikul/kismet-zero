// Package orthography defines writing system mappings that convert phonemes to graphemes.
//
// This package builds on lang/phoneme and lang/phonology by providing the bridge
// between spoken language and written representation. It enables consistent text
// generation for names, words, and phrases by mapping phonological elements to
// appropriate writing system characters.
//
// Responsibilities:
//   - Define grapheme-phoneme mappings for consistent written representation.
//   - Support multiple writing system styles (alphabetic, syllabic, logographic).
//   - Provide orthography generators that create writing systems for specific cultures.
//   - Ensure consistency between spoken and written forms of generated languages.
//   - Support customization of writing systems based on cultural preferences.
//
// Boundaries:
//   - Does not define phonemes (delegates to lang/phoneme).
//   - Does not implement phonotactics (delegates to lang/phonology).
//   - Does not construct morphemes or words (handled in lang/morphology).
//   - Does not model diachronic change (handled in lang/evolution).
//
// This package must remain focused on writing system representation and mapping.
// It should provide clean abstractions for different orthographic styles while
// maintaining deterministic behavior through seeded randomness.
package orthography
