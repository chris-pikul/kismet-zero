# lang/phonology

## Purpose
This package defines the **rules for how phonemes combine into syllables and sound patterns**.  
It acts as the bridge between raw phoneme inventories (`lang/phoneme`) and word-building (`lang/morphology`).

## Responsibilities
- Define syllable templates (CV, CVC, CCV, etc.).
- Enforce phonotactic rules on phoneme combinations.
- Provide generators for valid syllables from a phoneme pool.
- Support weighted probabilities for syllable structures and clusters.
- Allow customization for different linguistic "dialects" or cultures.

## Boundaries
- Do not define phonemes here. Import them from `lang/phoneme`.
- Do not construct morphemes or words. That is for `lang/morphology`.
- Do not apply diachronic change. That is for `lang/evolution`.

## Rules
- Represent phonotactic rules as **composable units** (`PhonotacticRule` interface).
- Keep syllable generation **deterministic** under seeded randomness.
- Maintain **low-level focus**: only syllables, not larger linguistic units.
- Use **unit tests** to verify rule enforcement (e.g., disallow forbidden clusters).

## Testing Expectations
- Verify syllables respect defined templates.
- Verify invalid combinations are filtered out.
- Verify probabilities are respected when selecting structures.
- Use table-driven tests for syllable generation.