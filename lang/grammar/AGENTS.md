# Grammar Package Overview

The grammar package provides sentence-level construction and syntactic organization for
Kismet's language generation system, building on the morphological foundation to create
structured, grammatically correct phrases and sentences.

## Core Goals
• Define syntactic structures (SVO, SOV, VSO, etc.) with configurable patterns
  and weights for different language types.
• Implement comprehensive agreement systems (case, number, gender, tense, aspect)
  that maintain grammatical consistency across sentence elements.
• Generate sentences and phrases using morphological building blocks while
  preserving phonological well-formedness.
• Support diverse language types (agglutinative, fusional, isolating) through
  flexible rule-based systems.
• Maintain cultural and linguistic consistency by integrating with phonology,
  morphology, and orthography systems.

## Key Features
1. Sentence Structure:
  • Configurable word order patterns (SVO, SOV, VSO, VOS, OVS, OSV).
  • Template-based sentence generation for common types (declarative, interrogative, imperative).
  • Support for complex sentences with subordinate clauses and coordination.
2. Agreement Systems:
  • Case marking (nominative, accusative, genitive, dative, etc.).
  • Number agreement (singular, dual, plural).
  • Gender agreement (masculine, feminine, neuter, animate, inanimate).
  • Tense and aspect marking (past, present, future, perfective, imperfective).
3. Phrase Construction:
  • Noun phrases with proper article, adjective, and modifier ordering.
  • Verb phrases with tense, aspect, and mood marking.
  • Prepositional and adverbial phrase structures.
4. Syntactic Rules:
  • Composable rule system for different grammatical phenomena.
  • Weighted rule application for probabilistic grammar generation.
  • Validation of sentence structure and grammatical well-formedness.
5. Integration:
  • Receives Word objects from morphology package.
  • Ensures phonological consistency with phonology package.
  • Provides complete sentences for higher-level text generation.
  • Maintains cultural context and language-specific patterns.

## Package Structure
grammar/
  • package.go — Package documentation and overview
  • types.go — Core grammar types, enums, and interfaces
  • syntax.go — Sentence structure patterns and templates
  • agreement.go — Agreement rules and systems
  • generator.go — Sentence and phrase generation logic

## Boundaries
• Does not define phonemes (delegates to lang/phoneme).
• Does not implement phonotactics (delegates to lang/phonology).
• Does not construct morphemes or words (delegates to lang/morphology).
• Does not create writing systems (delegates to lang/orthography).
• Does not model diachronic change (handled in lang/evolution).

## Design Principles
• Use composition over inheritance with interfaces for abstractions.
• Maintain deterministic behavior through seeded randomness.
• Support weighted selection for grammar rule application.
• Provide clean abstractions for different grammatical systems.
• Ensure cultural and linguistic consistency across all generated content.
