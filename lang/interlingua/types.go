package interlingua

// SchemaVersion allows evolution of the format over time.
const SchemaVersion = "1.0"

// ConceptID identifies a sense-level concept; can be internal or linked externally.
type ConceptID string

// Role labels are semantic.
type Role string

const (
	RoleAgent       Role = "agent"     // arg0
	RolePatient     Role = "patient"   // arg1
	RoleRecipient   Role = "recipient" // arg2
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

// TAM represents tense, aspect, mood, and related verbal features.
type TAM struct {
	Tense      string `json:"tense,omitempty"`      // past|pres|fut|remote|...
	Aspect     string `json:"aspect,omitempty"`     // prog|perf|hab|iter|comp|...
	Mood       string `json:"mood,omitempty"`       // ind|subj|imp|cond|opt|...
	Polarity   string `json:"polarity,omitempty"`   // pos|neg
	Evidential string `json:"evidential,omitempty"` // dir|inf|rep|vis|aud|...
	Voice      string `json:"voice,omitempty"`      // act|pass|mid|appl|caus|...
	Deixis     string `json:"deixis,omitempty"`     // proximal|distal (verbal)
}

// NominalFeatures represents grammatical features of nominal expressions.
type NominalFeatures struct {
	Person       int    `json:"person,omitempty"`       // 1|2|3 (allow 4 for obviative later)
	Number       string `json:"number,omitempty"`       // sg|pl|du|paucal|...
	Gender       string `json:"gender,omitempty"`       // masc|fem|neut|anim|inan|...
	Animacy      string `json:"animacy,omitempty"`      // anim|inan
	Definiteness string `json:"definiteness,omitempty"` // def|indef|specific|generic
	Case         string `json:"case,omitempty"`         // nom|acc|erg|abs|dat|gen|loc|ins|...
	Classifier   string `json:"classifier,omitempty"`   // shape|measure|human|...
	Deixis       string `json:"deixis,omitempty"`       // proximal|distal (nominal)
}

// Discourse represents information structure and pragmatic properties.
type Discourse struct {
	Topic     bool   `json:"topic,omitempty"`
	Focus     bool   `json:"focus,omitempty"`
	Givenness string `json:"givenness,omitempty"` // new|given|bridged
	Coref     string `json:"coref,omitempty"`     // coreference ID
}

// Entity represents a referent in the discourse.
type Entity struct {
	ID        string          `json:"id"`                  // stable within document
	Concept   ConceptID       `json:"concept"`             // sense
	LemmaHint string          `json:"lemmaHint,omitempty"` // for reversibility
	Feats     NominalFeatures `json:"feats,omitempty"`
	Discourse Discourse       `json:"discourse,omitempty"`
}

// Event represents a predicate with its arguments and adjuncts.
type Event struct {
	ID        string          `json:"id"`
	Predicate ConceptID       `json:"predicate"` // sense predicate (e.g., give-01)
	LemmaHint string          `json:"lemmaHint,omitempty"`
	TAM       TAM             `json:"tam"`
	Roles     map[Role]string `json:"roles,omitempty"`    // Role -> Entity.ID
	Adjuncts  map[Role]string `json:"adjuncts,omitempty"` // time/location/manner -> Entity.ID
	Valency   int             `json:"valency,omitempty"`
}

// Note represents a diagnostic or informational message.
type Note struct {
	Severity string `json:"severity"` // info|warn|loss|error
	Code     string `json:"code"`     // e.g., "feature.dropped.evidential"
	Message  string `json:"message"`
}

// Trace represents alignment metadata and diagnostic information.
type Trace struct {
	SourceLanguageID string         `json:"sourceLanguageId,omitempty"` // origin language
	TokenToNode      map[int]string `json:"tokenToNode,omitempty"`      // token idx -> node ID
	ConstructionID   string         `json:"constructionId,omitempty"`   // source template identifier
	Notes            []Note         `json:"notes,omitempty"`
}

// Document represents a complete semantic graph for a sentence or discourse unit.
type Document struct {
	SchemaVersion string   `json:"schemaVersion"`
	LanguageID    string   `json:"languageId,omitempty"` // language of surface form emitted/analyzed
	Entities      []Entity `json:"entities"`
	Events        []Event  `json:"events"`
	Trace         Trace    `json:"trace,omitempty"`
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
