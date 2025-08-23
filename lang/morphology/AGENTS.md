# Agent Guidance: lang/morphology

## Purpose
This package defines the **word-building layer** that constructs meaningful linguistic units from phonological building blocks.  
It acts as the bridge between syllable generation (`lang/phonology`) and sentence construction (`lang/grammar`), enabling the creation of morphemes, words, and lexicons with cultural context.

## Responsibilities
- Define **morpheme types** and structures (roots, affixes, infixes, circumfixes).
- Implement **word formation** using derivational and inflectional rules.
- Generate **lexicons** with culturally appropriate word inventories.
- Create **names** for people, places, and cultural concepts.
- Support **morphological analysis** and word structure parsing.
- Maintain **phonological well-formedness** when constructing words.
- Provide **weighted selection** for morpheme and word generation.

## Boundaries
- Do not define phonemes here. Import them from `lang/phoneme`.
- Do not implement phonotactics here. That belongs to `lang/phonology`.
- Do not create writing systems here. That belongs to `lang/orthography`.
- Do not construct sentences here. That belongs to `lang/grammar`.
- Do not model diachronic change here. That belongs to `lang/evolution`.

## Rules
- Represent morphological rules as **composable units** (`MorphologicalRule` interface).
- Keep word generation **deterministic** under seeded randomness.
- Maintain **cultural context** for appropriate semantic domains.
- Use **unit tests** to verify rule application and word formation.
- Support **different language types** (agglutinative, fusional, isolating).

## Testing Expectations
- Verify morphemes respect phonological constraints.
- Verify word formation follows defined rules.
- Verify lexicon generation maintains cultural consistency.
- Verify name generation produces appropriate patterns.
- Use table-driven tests for different morphological patterns.
- Test integration with phonology and orthography packages.

## Example Extension Points
- Adding new morpheme types or word formation patterns.
- Supporting language-family specific morphological systems.
- Expanding semantic domain coverage for cultural contexts.
- Future integration with grammar and evolution packages.
