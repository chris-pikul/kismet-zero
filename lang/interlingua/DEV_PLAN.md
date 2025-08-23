## Interlingua Development Plan (Agent Guide)

This plan guides implementing the `lang/interlingua` subsystem as specified in `lang/interlingua/README.md`. It is task-oriented, with clear milestones, outputs, and tests. Follow Kismet workspace rules (lowerCamelCase JSON tags, no global state, deterministic behavior, table-driven tests, ≥70% coverage).

### Scope and Objectives
- Build a semantic pivot for bi-directional translation between generated languages and user languages (EN first).
- Emit Interlingua Documents during conlang generation; realize/analyze user languages via rule-based components.
- Maintain determinism, reversibility, and loss-aware diagnostics.

### Constraints and Standards
- Do not modify existing logic without explicit instruction; integrate via opt-in adapters/injection points.
- Use dependency injection; avoid global state.
- Follow `lang` conventions: lowerCamelCase JSON tags; explicit error handling; small, testable functions.
- Ensure deterministic behavior via `Language.Seed`.
- Provide GoDoc comments for exported symbols; add `package.go` with package-level docs.
- Write unit tests for all exported functions; table-driven; aim ≥70% coverage.

### Deliverables and Artifacts
- New package: `lang/interlingua` with documented types, validators, realizer/analyzer scaffolds.
- English realizer and analyzer v0 with round-trip tests.
- Emission helpers for conlang generation (hooks) without altering existing packages by default.
- Test fixtures: JSON Documents, expected surfaces, golden files.
- Diagnostics utilities and validators.

### Proposed Package Layout (create files)
- `lang/interlingua/package.go`
- `lang/interlingua/types.go`
- `lang/interlingua/trace.go`
- `lang/interlingua/validate.go`
- `lang/interlingua/concepts.go`
- `lang/interlingua/emit_conlang.go`
- `lang/interlingua/realize_en.go`
- `lang/interlingua/analyze_en.go`
- `lang/interlingua/testdata/` (fixtures)
- `lang/interlingua/types_test.go`, `realize_en_test.go`, `analyze_en_test.go`, `validate_test.go`

### Key Interfaces (implement in types.go)
```go
package interlingua

const SchemaVersion = "1.0"

type ConceptID string
type Role string

type TAM struct {
    Tense, Aspect, Mood, Polarity, Evidential, Voice, Deixis string
}

type NominalFeatures struct {
    Person int
    Number, Gender, Animacy, Definiteness, Case, Classifier, Deixis string
}

type Discourse struct {
    Topic bool
    Focus bool
    Givenness string
    Coref string
}

type Entity struct {
    ID        string
    Concept   ConceptID
    LemmaHint string
    Feats     NominalFeatures
    Discourse Discourse
}

type Event struct {
    ID        string
    Predicate ConceptID
    LemmaHint string
    TAM       TAM
    Roles     map[Role]string
    Adjuncts  map[Role]string
    Valency   int
}

type Note struct {
    Severity string // info|warn|loss|error
    Code     string
    Message  string
}

type Trace struct {
    SourceLanguageID string
    TokenToNode      map[int]string
    ConstructionID   string
    Notes            []Note
}

type Document struct {
    SchemaVersion string
    LanguageID    string
    Entities      []Entity
    Events        []Event
    Trace         Trace
}

// Realizer renders a Document to tokens for a target language.
type Realizer interface {
    LanguageCode() string
    Realize(doc Document) ([]string, Trace, error)
}

// Analyzer builds a Document from tokens for a source language.
type Analyzer interface {
    LanguageCode() string
    Analyze(tokens []string) (Document, error)
}
```

### Milestones

#### Milestone 0: Bootstrap and Types
- Goal: Establish package skeleton and core types with validators.
- Tasks:
  - Create files per layout; add `package.go` docs.
  - Implement constants for roles (agent, patient, recipient, …).
  - Implement `Validate(doc Document) error` in `validate.go`:
    - Required fields present.
    - Role references resolve to valid Entity IDs.
    - Valency consistent with roles.
  - Implement `Trace` helpers:
    - `AddNote(severity, code, msg string)`
    - `MapToken(tokenIdx int, nodeID string)`
- Outputs: Compilable package with tests.
- Tests:
  - Types marshaling/unmarshaling JSON.
  - Validator detects missing roles/IDs; passes valid example (from README).
  - Trace helpers append notes and mappings.

#### Milestone 1: Concept Inventory (Minimal)
- Goal: Map lexemes to ConceptIDs with sense priors.
- Tasks:
  - Define interfaces in `concepts.go`:
    - `type ConceptResolver interface { Resolve(lemma string, lang string) (ConceptID, bool) }`
    - In-memory resolver with seed data for core domains (transfer, motion, possession, creation, perception, speech, statives, time, location).
  - Provide `Register(lemma, lang string, id ConceptID)` for tests.
- Outputs: Resolver usable by realizer/analyzer.
- Tests: Resolution success, fallback behavior with boolean false.

#### Milestone 2: Emission Hooks for Conlangs
- Goal: Enable generating a Document alongside surface forms without changing existing packages by default.
- Tasks:
  - In `emit_conlang.go`, provide builder helpers:
    - `NewDocument(langID string) Document`
    - `AddEntity(doc *Document, e Entity) string` (returns ID)
    - `AddEvent(doc *Document, ev Event) string`
    - `SetRole(ev *Event, role Role, entityID string)`
    - `AttachAdjunct(ev *Event, role Role, entityID string)`
  - Provide a pure function to build a sentence Document from known template info:
    - `BuildFromTemplate(langID, constructionID string, predicate ConceptID, tam TAM, args map[Role]Entity, adjuncts map[Role]Entity) (Document, error)`
  - Keep this opt-in; do not alter `grammar` package.
- Outputs: Emission API callable from generators later.
- Tests: Construct the README example via `BuildFromTemplate`; validate JSON.

#### Milestone 3: English Realizer v0
- Goal: Interlingua → English tokens for simple clauses.
- Tasks:
  - Implement `enRealizer` that satisfies `Realizer`.
  - Support:
    - Declarative active clauses with roles agent/patient/recipient.
    - Tense: past/present; Polarity: pos/neg.
    - Basic time/location adjuncts via PPs or adverbs (“yesterday”).
    - Article selection for definiteness; basic plural -s.
  - Determinism: tie-breaks use a seeded RNG injected or fixed ordering.
  - Loss handling: evidentiality becomes adverb (“apparently”) or note with severity `loss`.
- Outputs: `realize_en.go`.
- Tests:
  - Table-driven fixtures: K examples per predicate type.
  - Golden outputs; assert deterministic results across runs.
  - Ensure `Trace.Notes` appended on degraded features.

#### Milestone 4: English Analyzer v0
- Goal: English → Interlingua for simple clauses.
- Tasks:
  - Implement `enAnalyzer` satisfying `Analyzer`.
  - Lightweight parsing:
    - Tokenize; simple patterns; optional UD-like dependency via a small third-party or internal heuristic (if external, gate behind build tag).
  - Heuristics:
    - subject→agent (unless passive), dobj→patient, PP “to”→recipient.
    - Past via “did/didn’t” or -ed; negation via “not/never”.
    - Pronouns for person/number; nouns default sg; definiteness from articles.
  - Map lemmas to ConceptIDs via `ConceptResolver`; fallback with warn note.
- Outputs: `analyze_en.go`.
- Tests:
  - Surface → Document equals expected semantics for the fixtures used in Milestone 3 (semantic equivalence, not token-for-token).
  - Add diagnostics for defaults (warn) and losses (loss).

#### Milestone 5: Round-trip and Invariance Tests
- Goal: Guarantee stability for conlangs and semantic equivalence for English.
- Tasks:
  - Conlang simulated via `BuildFromTemplate`:
    - Document → EN → Analyze → Document’ ≈ Document (roles/TAM/polarity equal).
  - Golden tests placed in `testdata/`.
- Outputs: `roundtrip_test.go`.
- Tests: Table-driven across predicate families.

#### Milestone 6: Integration Points (Opt-in)
- Goal: Provide clean hooks for `grammar`/`morphology` without altering logic.
- Tasks:
  - Define small adapter interfaces in `emit_conlang.go`:
    - `type FeatureEmitter interface { FeaturesForMorpheme(m string) (map[string]string, bool) }`
    - `type ConstructionTagger interface { CurrentConstructionID() string }`
  - Provide sample adapters and documentation on how generators can call `BuildFromTemplate`.
- Outputs: Adapters + docs comments.
- Tests: Mock adapters in tests to simulate morphology/grammar signals.

#### Milestone 7: Documentation and Examples
- Goal: Developer UX.
- Tasks:
  - Add usage examples in `package.go`.
  - Add a minimal example program under `lang/interlingua/examples/` that builds the README example and prints EN.
- Outputs: Docs compile; examples run.
- Tests: Example test ensures output matches golden string.

### Testing Matrix (core)
- Clause types: intransitive, transitive, ditransitive (give).
- TAM: present/past; negation.
- NPs: sg/pl; def/indef; pronouns.
- Adjuncts: time (“yesterday”), location (“in the house”).
- Passive voice (planned extension).
- Features with loss: evidentiality, classifiers (record loss notes).

### Diagnostics and Validation
- `Validate(doc)` returns descriptive errors with codes.
- `Trace.AddNote` used consistently for assumptions and losses.
- A `CheckSemanticEquivalence(a, b Document) error` helper for tests (same roles/TAM/polarity ignoring IDs/order).

### Performance and Determinism
- Keep realization O(n) over tokens.
- Cache lemma→ConceptID lookups in a small LRU (package-scope injected via constructor; avoid globals).
- Deterministic tie-break ordering; seed injected via constructor for realizer if needed:
  - `NewEnglishRealizer(seed int64, resolver ConceptResolver) Realizer`

### Definition of Done (project-level)
- All milestones complete; all tests passing with ≥70% coverage.
- Round-trip tests stable; golden files approved.
- Package exported docs in place; examples run.
- No changes to existing `lang` subpackages required for basic functionality.

### Risks and Mitigations
- Ambiguity in analysis: mitigate with notes and defaults; keep analyzer conservative.
- Category mismatch (evidentiality, classifiers): represent in Document; realize periphrastically or note loss.
- Scope creep: stick to Milestone coverage; add features iteratively.

### How to Run
- `go test ./...` (ensure golden fixtures under `lang/interlingua/testdata`).
- Example: `go test -run Example ./lang/interlingua/...`
