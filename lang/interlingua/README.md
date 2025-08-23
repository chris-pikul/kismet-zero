# Interlingua (Semantic Pivot) — Kismet `lang/interlingua`

This document specifies the Interlingua subsystem: a semantic pivot for bi‑directional translation between procedurally generated languages and user-facing languages (e.g., English, French). It serves as the authoritative knowledge base for architecture, data model, and algorithms, guiding later implementation.

## Motivation

- **Cross-language translation**: Enable translation between any generated language and user languages without N×M transfer rules.
- **Determinism and reversibility**: Preserve meaning and construction choices so conlang → interlingua → conlang is stable.
- **Scalability**: Add new target languages by writing only a realizer/analyzer pair against the pivot.
- **Tight integration with `lang`**: Leverage existing phonology, morphology, grammar, orthography models; avoid brittle parsing by instrumenting generators.

## Goals

- **Bi-directional pivot**: Conlang ↔ Interlingua ↔ UserLang.
- **Loss-aware**: Preserve features even when targets cannot realize them; record degradations.
- **Deterministic**: Seeded, reproducible outputs consistent with `lang.Language`.
- **Composable**: Clean contracts with `phonology`, `orthography`, `morphology`, `grammar`.
- **Extensible**: Feature inventory covers present and future grammar (evidentiality, applicatives, classifiers, switch-reference).

## Non-goals (v0)

- Full free-text MT for arbitrary user text beyond supported analyzers.
- Statistical/neural modeling; rely on rule-based deterministic realization.
- Discourse planning, pragmatics beyond basic topic/focus flags.

## High-level Architecture

- **Interlingua Document**: A semantic graph per sentence (or small discourse unit) consisting of:
  - Events (predicates) with TAM/polarity/voice/etc.
  - Entities with nominal features (person/number/gender/case/definiteness/classifier).
  - Role edges (agent, patient, recipient, etc.) and adjuncts (time, location, manner).
  - Trace metadata for alignment and reversibility (construction IDs, token ↔ node alignments, notes).
- **Pipelines**:
  - Conlang generation → emit Interlingua alongside surface text using known templates and morpheme feature rules.
  - Interlingua → UserLang (realization): map roles to syntax; linearize with language-specific templates.
  - UserLang → Interlingua (analysis): heuristic parse → role mapping; use lexicon with ConceptIDs; mark uncertainties.

Data flows:
- Conlang → Interlingua → English/French
- English/French → Interlingua → Conlang

Interlingua is the stable, loss-aware source of truth.

## Core Data Model

A compact Go-oriented schema for the pivot. Enumerations are strings for extensibility; JSON tags are lowerCamelCase per repo rules.

```go
package interlingua

// SchemaVersion allows evolution of the format over time.
const SchemaVersion = "1.0"

// ConceptID identifies a sense-level concept; can be internal or linked externally.
type ConceptID string

// Role labels are semantic.
type Role string

const (
    RoleAgent       Role = "agent"       // arg0
    RolePatient     Role = "patient"     // arg1
    RoleRecipient   Role = "recipient"   // arg2
    RoleExperiencer Role = "experiencer"
    RoleStimulus    Role = "stimulus"
    RoleInstrument  Role = "instrument"
    RoleLocation    Role = "location"
    RoleTime        Role = "time"
    RoleManner      Role = "manner"
    RoleCause       Role = "cause"
    RoleBeneficiary Role = "beneficiary"
    RoleSource      Role = "source"
    RoleGoal        Role = "goal"
)

type TAM struct {
    Tense      string `json:"tense,omitempty"`      // past|pres|fut|remote|...
    Aspect     string `json:"aspect,omitempty"`     // prog|perf|hab|iter|comp|...
    Mood       string `json:"mood,omitempty"`       // ind|subj|imp|cond|opt|...
    Polarity   string `json:"polarity,omitempty"`   // pos|neg
    Evidential string `json:"evidential,omitempty"` // dir|inf|rep|vis|aud|...
    Voice      string `json:"voice,omitempty"`      // act|pass|mid|appl|caus|...
    Deixis     string `json:"deixis,omitempty"`     // proximal|distal (verbal)
}

type NominalFeatures struct {
    Person       int    `json:"person,omitempty"`        // 1|2|3 (allow 4 for obviative later)
    Number       string `json:"number,omitempty"`        // sg|pl|du|paucal|...
    Gender       string `json:"gender,omitempty"`        // masc|fem|neut|anim|inan|...
    Animacy      string `json:"animacy,omitempty"`       // anim|inan
    Definiteness string `json:"definiteness,omitempty"`  // def|indef|specific|generic
    Case         string `json:"case,omitempty"`          // nom|acc|erg|abs|dat|gen|loc|ins|...
    Classifier   string `json:"classifier,omitempty"`    // shape|measure|human|...
    Deixis       string `json:"deixis,omitempty"`        // proximal|distal (nominal)
}

type Discourse struct {
    Topic     bool   `json:"topic,omitempty"`
    Focus     bool   `json:"focus,omitempty"`
    Givenness string `json:"givenness,omitempty"`   // new|given|bridged
    Coref     string `json:"coref,omitempty"`       // coreference ID
}

type Entity struct {
    ID        string          `json:"id"`                    // stable within document
    Concept   ConceptID       `json:"concept"`               // sense
    LemmaHint string          `json:"lemmaHint,omitempty"`   // for reversibility
    Feats     NominalFeatures `json:"feats,omitempty"`
    Discourse Discourse       `json:"discourse,omitempty"`
}

type Event struct {
    ID        string            `json:"id"`
    Predicate ConceptID         `json:"predicate"`           // sense predicate (e.g., give-01)
    LemmaHint string            `json:"lemmaHint,omitempty"`
    TAM       TAM               `json:"tam"`
    Roles     map[Role]string   `json:"roles,omitempty"`     // Role -> Entity.ID
    Adjuncts  map[Role]string   `json:"adjuncts,omitempty"`  // time/location/manner → Entity.ID
    Valency   int               `json:"valency,omitempty"`
}

type Note struct {
    Severity string `json:"severity"` // info|warn|loss|error
    Code     string `json:"code"`     // e.g., "feature.dropped.evidential"
    Message  string `json:"message"`
}

type Trace struct {
    SourceLanguageID string         `json:"sourceLanguageId,omitempty"` // origin language
    TokenToNode      map[int]string `json:"tokenToNode,omitempty"`       // token idx → node ID
    ConstructionID   string         `json:"constructionId,omitempty"`    // source template identifier
    Notes            []Note         `json:"notes,omitempty"`
}

type Document struct {
    SchemaVersion string   `json:"schemaVersion"`
    LanguageID    string   `json:"languageId,omitempty"` // language of surface form emitted/analyzed
    Entities      []Entity `json:"entities"`
    Events        []Event  `json:"events"`
    Trace         Trace    `json:"trace,omitempty"`
}
```

### Concept Inventory

- A global inventory maps lemmas to `ConceptID`s and sense metadata:
  - Generated languages: assign `ConceptID`s during lexicon generation; store priors.
  - User languages: maintain dictionaries for high-frequency lemmas; back off to generic concepts (e.g., `move-00`) for OOV.
- Multiword expressions: treat as single concepts with realization hints (e.g., `make_up_mind-01`).

### Construction IDs and Lemma Hints

- Each grammar template in a conlang has a `constructionID` (e.g., `decl.transitive.svo.neg.pst`).
- `lemmaHint` carries the base lexeme for faithful reverse realization.

## Algorithms and Pipelines

### A. Conlang Generation → Interlingua Emission

Instead of parsing, instrument existing generators.

1. **Morphology hook**:
   - Maintain a morpheme→feature table (e.g., `-ed` → `{tense:past}`).
   - Emit accumulated features per word to support event/entity feature assembly.

2. **Grammar hook**:
   - When selecting a sentence template, assign a `constructionID`.
   - Map syntactic roles in the template to semantic roles (agent/patient/etc.).
   - Construct `Event` with `TAM`, attach `Entity` nodes with `NominalFeatures`.
   - Fill `Trace.TokenToNode` as tokens are linearized.

3. **Output**:
   - Surface string (existing), plus `Document` (semantic) and `Trace` (alignment).
   - Deterministic given `Language.Seed`.

Round-trip invariant (conlangs): surface → Document → surface’ == surface (modulo orthographic trivialities).

### B. Interlingua → User Language Realization

Rule-based realization per target language.

1. **Planning**:
   - Choose clause type and word order based on roles and `TAM`.
   - Decide voice (active/passive) if marked; otherwise target default.

2. **Syntactic mapping**:
   - Map roles to subject/object/oblique/prepositional phrases (language-specific).
   - Insert auxiliaries/particles for `TAM`, `Polarity`, `Mood`.

3. **Morphology/lexicalization**:
   - Select lemma from target lexicon by `ConceptID` (+ `LemmaHint` for tie-breaking).
   - Inflect per target rules; handle agreement.

4. **Linearization**:
   - Apply target templates; produce tokens and optional updated `Trace`.

5. **Loss handling**:
   - If feature cannot be realized (e.g., evidentiality in English), add `Note{Severity:"loss"}` and optionally periphrastic fallback (“apparently”).

### C. User Language → Interlingua Analysis (v0: English-first)

1. **Parsing**:
   - Lightweight UD-like dependency parse.
   - Identify main predicates, arguments, and adjuncts.

2. **Role assignment**:
   - Heuristics: subject→agent (default), dobj→patient, iobj/prep “to”→recipient, etc.
   - Passive detection to flip roles.

3. **Feature extraction**:
   - Tense/aspect/mood/polarity from auxiliaries/inflections.
   - Number/person/gender where recoverable; default otherwise.

4. **Sense mapping**:
   - Lemma→`ConceptID` via lexicon; fallback to generic concepts with `Note{Severity:"warn"}`.

5. **Document and Trace**:
   - Build `Document`; record defaults/uncertainties in `Trace.Notes`.

## Integration with Existing Packages

- `phoneme`/`phonology`: No direct dependency; operate at semantic layer. Phonology influences morphology/orthography only.
- `orthography`: Used for surface rendering only; Interlingua is orthography-agnostic.
- `morphology`:
  - Provide morpheme→feature mapping.
  - Expose functions to assemble nominal/verbal feature bundles deterministically.
  - Export lexicon entries with `ConceptID`s and sense priors.
- `grammar`:
  - Tag templates with `constructionID`.
  - Provide role mapping tables per template (syntactic→semantic).
  - Emit token-level alignment during linearization.
- `lang` root:
  - `Language` gains no new hard dependency; Interlingua operates alongside as an optional layer.
  - Use `Language.ID` to fill `Document.LanguageID` and `Trace.SourceLanguageID`.
  - Determinism via `Language.RNG`.

## Proposed Package Layout

```
lang/interlingua/
  README.md                // this document
  package.go               // GoDoc: package overview
  types.go                 // Document, Event, Entity, features, notes, roles
  concepts.go              // Concept inventory interfaces + linking (WordNet/Wikidata optional)
  trace.go                 // Trace helpers, diagnostics, loss tracking
  emit_conlang.go          // Hooks/types used by morphology/grammar to emit Document
  realize_en.go            // Interlingua → English realizer (v0)
  analyze_en.go            // English → Interlingua analyzer (v0)
  realize_fr.go            // Interlingua → French (planned)
  analyze_fr.go            // French → Interlingua (planned)
  testdata/                // Round-trip fixtures
```

Interfaces (sketch):

```go
// Realizer renders a Document to tokens for a target language.
type Realizer interface {
    LanguageCode() string              // "en", "fr", ...
    Realize(doc Document) ([]string, Trace, error)
}

// Analyzer builds a Document from tokens for a source language.
type Analyzer interface {
    LanguageCode() string
    Analyze(tokens []string) (Document, error)
}
```

## Versioning and Persistence

- `Document.SchemaVersion` allows migrations.
- JSON is the canonical wire/storage format; MsgPack optional.
- Store alongside generated text for reproducibility and evolution audits.

## Error Handling and Diagnostics

- `Trace.Notes` captures:
  - **info**: realizational choices (e.g., passive realized as active).
  - **warn**: defaults/assumptions (e.g., number=sg assumed).
  - **loss**: features dropped or approximated (e.g., evidentiality → adverb).
  - **error**: inconsistent roles, missing arguments when required by valency.

- Validation utilities:
  - Validate role coherence (e.g., `recipient` requires ditransitive predicate).
  - Validate feature compatibility (e.g., case systems vs. prepositional encoding).

## Determinism, Caching, and Performance

- Deterministic realization from `Language.Seed` for conlangs (choice ties broken by seeded RNG).
- Cache lexeme→ConceptID resolutions and common constructions.
- Memoize morphological inflection templates per target language.

## Testing Strategy

- **Round-trip tests (conlangs)**: surface → Document → surface’ equals input.
- **Semantic invariance tests**: Document → target → analysis → Document’ semantically equivalent (role/TAM/polarity).
- **Table-driven**: Cover core predicates (motion, transfer, possession, creation, perception, speech, statives).
- **Golden fixtures**: Store `Document` JSON next to expected surface strings.

## Evolution and Language Families

- `ConceptID`s remain stable across evolution; sense splits/merges tracked with a mapping log:
  - split: `X → {X.1, X.2}` with migration notes.
  - merge: `{Y.1, Y.2} → Y`.
- Children dialects inherit concept inventories and feature systems; Document emission remains consistent.

## Extensibility

- Add roles (e.g., `comitative`, `ablative`) without breaking old docs (string enums).
- Extend `TAM`, `NominalFeatures` gradually; consumers should ignore unknown fields.
- Support richer information structure (topic/focus subtypes) in future versions.

## Example

Surface (conlang, glossed): `neg-give.pst boy-nom book-acc girl-dat yesterday`

Document (abridged):

```json
{
  "schemaVersion": "1.0",
  "languageId": "elv-wood-sindarin",
  "entities": [
    {"id":"e1","concept":"boy","feats":{"number":"sg","definiteness":"def"}},
    {"id":"e2","concept":"book","feats":{"number":"sg","definiteness":"def"}},
    {"id":"e3","concept":"girl","feats":{"number":"sg","definiteness":"def"}},
    {"id":"t1","concept":"yesterday"}
  ],
  "events": [
    {
      "id":"v1",
      "predicate":"give-01",
      "tam":{"tense":"past","polarity":"neg"},
      "roles":{"agent":"e1","patient":"e2","recipient":"e3"},
      "adjuncts":{"time":"t1"},
      "valency":3
    }
  ],
  "trace":{
    "sourceLanguageId":"elv-wood-sindarin",
    "constructionId":"decl.transitive.svo.neg.pst",
    "tokenToNode":{"0":"v1","1":"e1","2":"e2","3":"e3","4":"t1"},
    "notes":[]
  }
}
```

English realization: “The boy didn’t give the book to the girl yesterday.”

Loss examples:
- If target cannot realize evidentiality: add `{"severity":"loss","code":"feature.dropped.evidential","message":"English lacks evidential morphology"}`.

## Implementation Roadmap (for later planning)

1. `types.go`, `trace.go`, validators.
2. Morphology/grammar emission hooks for conlangs.
3. English realizer with core clause types; unit tests + golden fixtures.
4. English analyzer (limited UD heuristics); semantic equivalence tests.
5. Concept inventory (bootstrap with core concepts + mapping tools).
6. French realizer/analyzer; extend features as needed.
7. Performance pass: caching, memoization, serialization.

## Appendix A: Role Mapping Guidelines

- Transitive: agent→subject, patient→dobj (case/preposition depends on target).
- Ditransitive: recipient→IO (dative or PP “to”).
- Passive: agent demoted/omitted; patient→subject.
- Motion: source/goal roles map to “from/to” PPs or case.
- Cause: periphrastic “because (of)” or causative morphology if target supports.

## Appendix B: Feature Defaults (analysis)

- Person: default 3 unless pronoun indicates otherwise.
- Number: default sg unless plural markers present.
- Gender: unknown unless lexical or pronominally marked.
- Case: infer from syntax/prepositions; otherwise unknown.
- Definiteness: infer from articles/determiners; default indef.

---

This README defines the Interlingua as a semantic, loss-aware, bi-directional pivot integrated with Kismet’s `lang` package. It is designed to be deterministic, extensible, and practical for procedurally generated languages and user-facing translation.

