// Package lang provides procedural linguistic generation and evolution for
// Kismet, supporting the creation of languages tied to cultures, with mechanisms
// for diachronic change and cross-cultural blending.
//
// The lang package serves as the orchestration layer that coordinates between
// specialized linguistic modules:
//   - phoneme: fundamental sound units and pools
//   - phonology: syllable structure and phonotactic rules
//   - orthography: writing system mappings and styles
//   - morphology: word formation and lexicon generation
//   - grammar: sentence structure and agreement systems
//   - evolution: language change and diachronic evolution
//   - cultural: cultural influence and contact modeling
//   - variation: dialect formation and language variation
//   - interlingua: semantic analysis and generation
//
// Core Responsibilities:
//   - Define the Language struct that integrates all linguistic components
//   - Provide language identification and classification systems
//   - Enable language creation, evolution, and cultural linking
//   - Support high-level API for name and text generation
//   - Maintain consistency across all linguistic subsystems
//   - Enable deterministic but diverse language generation
//
// The package follows Go best practices with explicit error handling,
// dependency injection, and comprehensive testing. All exported types
// and functions are documented with GoDoc comments.
package lang
