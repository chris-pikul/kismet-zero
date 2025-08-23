# Agent Guidance: lang/orthography

## Purpose
This package defines the **writing system mappings** that convert phonemes to graphemes (letters, symbols, or characters).  
It provides the bridge between spoken language (phonology) and written representation, enabling consistent text generation for names, words, and phrases.

## Responsibilities
- Define **grapheme-phoneme mappings** that represent how sounds are written.
- Support multiple **writing system styles** (alphabetic, syllabic, logographic).
- Provide **orthography generators** that create writing systems for specific cultures.
- Ensure **consistency** between spoken and written forms of generated languages.
- Support **customization** of writing systems based on cultural preferences.

## Boundaries
- Do not define phonemes here. Import them from `lang/phoneme`.
- Do not implement phonotactics here. That belongs to `lang/phonology`.
- Do not build words or morphemes here. That belongs to `lang/morphology`.
- Do not model diachronic change here. That belongs to `lang/evolution`.

## Rules
- Maintain **clear mapping relationships** between phonemes and graphemes.
- Support **configurable writing system styles** with sensible defaults.
- Ensure **deterministic behavior** under seeded randomness for reproducibility.
- Use **unit tests** to verify mapping consistency and style adherence.

## Example Extension Points
- Adding new writing system styles (abugida, abjad, etc.).
- Supporting cultural variations in grapheme selection.
- Future integration with font/rendering systems.
- Historical orthography evolution patterns.

## Testing Expectations
- Verify phoneme-to-grapheme mappings are consistent.
- Verify writing system styles generate appropriate character sets.
- Verify orthography generation respects cultural parameters.
- Use table-driven tests for different writing system configurations.
