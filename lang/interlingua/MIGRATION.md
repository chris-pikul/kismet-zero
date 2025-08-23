# Interlingua Integration Migration (After Plan)

Audience: maintainers integrating `lang/interlingua` across `morphology`, `grammar`, and optionally `lang` root after the pivot package is implemented.

This plan is additive and non-breaking. It promotes first-class concept IDs, feature annotations, and construction metadata to enable deterministic, reversible interlingua emission and realization. No existing logic needs to change unless you opt to wire automatic emission.

## Prerequisites
- `lang/interlingua` package implemented per `README.md`:
  - Core types (`Document`, `Event`, `Entity`, `Role`, `ConceptID`, `TAM`, `Trace`, etc.).
  - Validators and trace utilities.
  - English realizer/analyzer v0 (optional at this stage).
- Tests and fixtures in place.

## Goals
- **Expose semantic hooks** in `morphology` and `grammar`:
  - Words carry one or more `ConceptID`s.
  - Morphemes optionally declare feature deltas (e.g., `{tense:past}`, `{polarity:neg}`).
  - Grammar templates expose stable `constructionID` and a semantic role plan.
  - Optional emission observer to collect token↔source alignment.
- **Type-safety**: prefer type aliases to avoid stringly-typed fields where feasible.
- **No behavior changes** unless you wire emission; fields are optional.

---

## Package: lang/morphology

### 1) Add ConceptIDs to `Word`
- **Why**: Deterministic sense mapping conlang ↔ interlingua without external sidecars.
- **Change**: Add a new optional field; no behavior change.

```go
// types.go
type Word struct {
    ID        string            `json:"id"`
    Morphemes []Morpheme        `json:"morphemes"`
    Meaning   string            `json:"meaning"`
    Phonemes  []phoneme.Phoneme `json:"phonemes"`
    Written   string            `json:"written"`
    Category  WordCategory      `json:"category"`
    Culture   string            `json:"culture,omitempty"`
    Frequency WordFrequency     `json:"frequency"`
    Weight    float32           `json:"weight"`
    Agreement AgreementFeatures `json:"agreement,omitempty"`

    // NEW (optional, additive)
    ConceptIDs []interlingua.ConceptID `json:"conceptIds,omitempty"`
}
```

- **Type-safety hint**: If you prefer to avoid importing `interlingua` here, use a type alias that can be swapped later:
```go
// concepts_alias.go (in morphology)
// Option A: direct alias (requires import; best type-safety)
type ConceptID = interlingua.ConceptID

// Option B: local alias now, upgrade later
type ConceptID string
```

### 2) Optional: Morpheme feature annotations
- **Why**: Emit `TAM`, case, polarity, etc. without parsing.
- **Change**: Add a compact, typed feature delta container.

```go
// types.go
type FeatureKey string
type FeatureVal string

const (
    FeatureTense     FeatureKey = "tense"
    FeatureAspect    FeatureKey = "aspect"
    FeatureMood      FeatureKey = "mood"
    FeaturePolarity  FeatureKey = "polarity"
    FeatureCase      FeatureKey = "case"
    FeatureVoice     FeatureKey = "voice"
    FeatureEvidential FeatureKey = "evidential"
)

type Morpheme struct {
    ID        string            `json:"id"`
    Type      MorphemeType      `json:"type"`
    Meaning   string            `json:"meaning"`
    Phonemes  []phoneme.Phoneme `json:"phonemes"`
    Syllables []string          `json:"syllables"`
    Weight    float32           `json:"weight"`
    Culture   string            `json:"culture,omitempty"`
    Frequency MorphemeFrequency `json:"frequency"`

    // NEW (optional, additive)
    Features map[FeatureKey]FeatureVal `json:"features,omitempty"`
}
```

- **Type-safety hint**: You can align `FeatureKey`/`FeatureVal` with interlingua’s enumerations later, or alias them:
```go
// align with interlingua (optional)
type Role = interlingua.Role
```

### 3) Lexicon export/import (optional)
- Extend `ExportLexicon` and related APIs to include `ConceptIDs` so external tools remain consistent.

---

## Package: lang/grammar

### 1) Tag sentence templates with stable IDs and role plan
- **Why**: Provide `Trace.ConstructionID` and deterministic semantic role mapping.

```go
// syntax.go
type SentenceTemplate struct {
    // existing fields...

    // NEW (optional, additive)
    ID    string            // e.g., "decl.transitive.svo.neg.pst"
    Roles map[string]string // e.g., {"subject":"agent","object":"patient","iobj":"recipient"}
}
```

- **Type-safety hint**: Replace `map[string]string` with typed roles:
```go
// Option A: import-alias interlingua.Role (best type-safety)
type Role = interlingua.Role
Roles map[string]Role

// Option B: define local Role string and alias later
type Role string
```

### 2) Surface metadata from the generator
- **Why**: Build token↔node alignments without intrusive changes.

```go
// generator.go
type RealizationMeta struct {
    TemplateID   string
    TokenToSource map[int]string // e.g., "verb","subj","obj","iobj","adv-time"
}

// where sentence is realized, return or optionally fill this meta
```

- If changing return signatures is undesirable, add an observer interface:

```go
// generator.go
type Observer interface {
    OnTemplateChosen(id string)
    OnToken(index int, source string) // "verb","subj","obj","iobj","adv-time"
}

type Options struct {
    Observer Observer
}

func (g *Generator) GenerateSentence(..., opts *Options) (/* existing return types */) {
    // if opts != nil && opts.Observer != nil { notify observer }
}
```

- **Type-safety hint**: Define `const` source tags and/or use a typed enum to avoid string drift.

---

## Package: lang (root)

### 1) Optional interlingua services accessor
- Keep `Language` unchanged. Add a small helper to attach interlingua services per language if desired.

```go
// interlingua_adapter.go (new, optional)
type InterlinguaServices struct {
    Realizers map[string]interlingua.Realizer // "en","fr",...
    Analyzers map[string]interlingua.Analyzer
}

func (l *Language) Interlingua() *InterlinguaServices {
    // lazy init or dependency-injected externally
    return nil
}
```

- This avoids global state and keeps interlingua opt-in.

---

## Wiring Emission (optional, post-migration)

If you choose to auto-emit `Document` during conlang generation:

- In grammar generation:
  - Set `template.ID` on selection.
  - Map syntactic slots to `Roles`.
  - While linearizing tokens, call observer hooks to fill `RealizationMeta` and later `Trace.TokenToNode`.

- In morphology:
  - Accumulate morpheme `Features` to compute `TAM`, `Case`, `Polarity`.
  - Use `Word.ConceptIDs` to select predicate and entity concepts.

- In interlingua:
  - Use `emit_conlang.go` helpers to assemble `Document` + `Trace`.
  - Validate with `Validate(doc)`.

This remains optional; you can construct `Document` in a separate pass, using the metadata fields as inputs.

---

## Type-Safe Alias Strategy

- **Preferred**: Directly use `interlingua` types where referenced:
  - `morphology.Word.ConceptIDs []interlingua.ConceptID`
  - `grammar.SentenceTemplate.Roles map[string]interlingua.Role`
- **Alternative**: Local aliases to decouple now and tighten later:
  - `type ConceptID string`, `type Role string`
  - Later, migrate to `type ConceptID = interlingua.ConceptID` and `type Role = interlingua.Role` once the package stabilizes.

This preserves compile-time safety while allowing staged adoption.

---

## Backwards Compatibility

- All changes are additive:
  - New fields are `omitempty`; zero-cost for JSON/backward consumers.
  - Observer/Options are nil-safe.
- No existing behavior changes unless you explicitly wire interlingua emission.

---

## Migration Steps

1. **Add fields and interfaces**:
   - `morphology/types.go`: `ConceptIDs`, `Features`.
   - `grammar/syntax.go`: `SentenceTemplate.ID`, `SentenceTemplate.Roles`.
   - `grammar/generator.go`: `Observer` and/or `RealizationMeta` (optional).

2. **Compile and run tests**:
   - Ensure no behavior changes.
   - Add tests ensuring new fields round-trip in JSON when set.

3. **Plumb IDs in generators** (optional):
   - Assign `template.ID` when selecting templates.
   - Populate `Roles` maps in templates where known.

4. **Adopt type aliases** (optional):
   - Switch `ConceptID`/`Role` to alias interlingua types once package is locked.

5. **Enable interlingua emission** (optional):
   - Implement a small adapter that listens to `Observer` and uses morphology features to build `Document` via `emit_conlang.go`.

---

## Testing Checklist

- **Morphology**:
  - Words with `ConceptIDs` serialize/deserialize.
  - Morphemes with `Features` combine into expected `TAM`/Case/Polarity in a unit test (pure function).

- **Grammar**:
  - Templates carry stable `ID` strings.
  - Observer receives `OnTemplateChosen` and `OnToken` in expected order.

- **Interlingua integration**:
  - Build README example via real generator + emission hooks; validate `Document`.
  - Round-trip tests remain deterministic across seeds.

---

## Rollback plan

- Since changes are additive:
  - Remove references to new fields in your code; leave struct fields (harmless).
  - Disable observers/options by passing `nil`.
  - Interlingua remains a separate package; no cycles introduced.

---

## Notes

- Keep enumerations future-proof by using strings in JSON.
- Prefer constructor injection for any interlingua services; avoid global singletons.
- Maintain stable `constructionID` naming; treat them as part of the public surface for reproducibility.
- Track `SchemaVersion` in stored `Document`s for forward migration tools.
