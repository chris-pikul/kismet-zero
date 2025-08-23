# Language Package Overview

The lang package provides procedural linguistic generation and evolution for
Kismet, supporting the creation of languages tied to cultures, with mechanisms
for diachronic change and cross-cultural blending.

## Core Goals
• Define phonological systems (phoneme pools, phonology subsets,
  weighting, and distribution) to establish unique sound inventories.
• Generate orthography mappings for consistent written representation.
• Build morphology and word-formation rules to create lexicons (names,
  words, phrases) with internal coherence.
• Support grammar modeling (syntax, inflection, agreement systems) for
  structured phrase/sentence generation.
• Enable diachronic language evolution, allowing languages to shift over
  time, diverge, or converge.
• Provide mechanisms for linguistic contact, enabling morphological and
  phonological influence between neighboring cultures.
• Allow procedural but rule-driven randomness, so generated languages are
  plausible yet diverse.

## Key Features
1. Phoneme System:
  • Global phoneme pool with IPA-derived categories and rarity weights.
  • Phonology type: curated, weighted subset for each language.
  • Tools for random/unweighted phoneme selection.
2. Orthography:
  • Mapping from phonemes to graphemes (letters or symbols).
  • Configurable styles (alphabetic, syllabic, logographic).
3.	Morphology & Lexicon:
  • Rules for constructing morphemes and assembling them into words.
  • Name and word generators tied to cultural context.
  • Derivational and inflectional morphology systems.
4.	Grammar & Syntax:
  • Configurable sentence structures (SVO, SOV, etc.).
  • Agreement rules (gender, number, case).
  • Phrase and sentence-level generators.
5.	Language Evolution:
  • Diachronic changes (sound shifts, morphological simplifications).
  • Contact-induced change (loanwords, blended phonologies).
  • Divergence and convergence modeling to track family trees of languages.
6.	Integration:
  • Languages are first-class entities tied to cultures.
  • Provides downstream generators (names, lore, text fragments).
  • Supports consistency across world history and cultural interactions.

## Package Structure (Suggested)
phoneme/
  • phoneme.go — Phoneme definitions, types, rarity weights
  • pool.go — Global pool and utilities
  • phonology.go — Phonology subset selection and weighting
  • choice.go — Random and weighted choice methods

phonology/
  • phonology.go — Core phonology types and syllable generation
  • rules.go — Phonotactic rules and constraints
  • syllable templates and validation systems

orthography/
  • mapping.go — Grapheme-phoneme mapping
  • styles.go — Writing system styles (alphabetic, syllabic, logographic)
  • generator.go — Generate writing systems for cultures

morphology/
  • morpheme.go — Smallest meaningful units
  • word.go — Word construction and lexicon generation
  • rules.go — Derivational and inflectional rules

grammar/
  • syntax.go — Sentence patterns (SVO, SOV, VSO, etc.)
  • agreement.go — Case, number, gender agreement rules
  • generator.go — Sentence and phrase construction

evolution/
  • soundchange.go — Systematic sound shifts
  • morphology.go — Evolution of forms over time
  • contact.go — Blending, borrowing, and loanwords
  • familytree.go — Track genealogical relationships between languages

root files
  • lang.go — Core language types and orchestration
  • culturelink.go — Linking languages to cultures
  • generator.go — High-level API: new language, evolve, merge
