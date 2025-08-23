# Agent Guidance: lang/phoneme

## Purpose
This package defines the **phonemic building blocks** of a language.  
It provides types, constants, and utilities for working with **consonants, vowels, and phonological features** that higher-level modules (phonology, morphology, grammar) depend on.

## Responsibilities
- Define the **Phoneme** struct and associated specifications:
  - Consonant manner, place, voicing
  - Vowel height, backness, rounding
- Provide **phoneme pools** for consonants and vowels.
- Support **weighted randomness** for selection of phonemes.
- Ensure **deduplication** of phonemes within a pool.
- Act as the **foundation layer** for phonology generation.

## Boundaries
- Do not implement **phonotactics** (syllable rules) here. That belongs to `lang/phonology`.
- Do not build words or morphemes here. That belongs to `lang/morphology`.
- Do not model diachronic change here. That belongs to `lang/evolution`.

## Rules
- Maintain **clear and strongly typed enums** for phoneme categories.
- For probabilistic logic, use **fixed seeds** for reproducibility.

## Example Extension Points
- Adding new consonant manners or vowel features.
- Expanding phoneme rarity weighting.
- Future integration with orthography mappings.