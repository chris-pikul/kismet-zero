// Package morphology defines the word-building layer that constructs meaningful
// linguistic units from phonological building blocks.
//
// This package builds on lang/phoneme, lang/phonology, and lang/orthography by
// providing morpheme construction, word formation rules, and lexicon generation.
// It enables the creation of culturally appropriate words and names while
// maintaining phonological well-formedness and supporting different language types.
//
// Responsibilities:
//   - Define morpheme types and structures (roots, affixes, infixes, circumfixes).
//   - Implement word formation using derivational and inflectional rules.
//   - Generate lexicons with culturally appropriate word inventories.
//   - Create names for people, places, and cultural concepts.
//   - Support morphological analysis and word structure parsing.
//   - Maintain phonological well-formedness when constructing words.
//   - Provide weighted selection for morpheme and word generation.
//
// Boundaries:
//   - Does not define phonemes (delegates to lang/phoneme).
//   - Does not implement phonotactics (delegates to lang/phonology).
//   - Does not create writing systems (delegates to lang/orthography).
//   - Does not construct sentences (handled in lang/grammar).
//   - Does not model diachronic change (handled in lang/evolution).
//
// This package must remain focused on word-level construction and morphological
// rules. It should provide clean abstractions for different morphological
// systems while maintaining deterministic behavior through seeded randomness.
package morphology
