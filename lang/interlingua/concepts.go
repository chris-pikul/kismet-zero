package interlingua

// ConceptResolver provides access to concept inventory and sense resolution.
type ConceptResolver interface {
	// Resolve attempts to find a ConceptID for a given lemma in a language.
	// Returns the ConceptID and true if found, empty string and false otherwise.
	Resolve(lemma, lang string) (ConceptID, bool)
}

// InMemoryConceptResolver provides a simple in-memory concept inventory.
type InMemoryConceptResolver struct {
	concepts map[string]map[string]ConceptID // lang -> lemma -> ConceptID
}

// NewInMemoryConceptResolver creates a new resolver with seed data for core domains.
func NewInMemoryConceptResolver() *InMemoryConceptResolver {
	resolver := &InMemoryConceptResolver{
		concepts: make(map[string]map[string]ConceptID),
	}

	// Initialize with core concepts for English
	resolver.concepts["en"] = map[string]ConceptID{
		// Transfer domain
		"give": "give-01",
		"send": "send-01",
		"hand": "hand-01",
		"pass": "pass-01",
		"lend": "lend-01",
		"rent": "rent-01",
		"sell": "sell-01",
		"buy":  "buy-01",
		"owe":  "owe-01",

		// Motion domain
		"move":    "move-00",
		"go":      "go-01",
		"come":    "come-01",
		"walk":    "walk-01",
		"run":     "run-01",
		"fly":     "fly-01",
		"swim":    "swim-01",
		"jump":    "jump-01",
		"climb":   "climb-01",
		"carry":   "carry-01",
		"bring":   "bring-01",
		"take":    "take-01",
		"receive": "receive-01",
		"get":     "get-01",
		"put":     "put-01",

		// Possession domain
		"have":    "have-01",
		"own":     "own-01",
		"possess": "possess-01",
		"belong":  "belong-01",
		"keep":    "keep-01",
		"hold":    "hold-01",

		// Creation domain
		"make":   "make-01",
		"create": "create-01",
		"build":  "build-01",
		"write":  "write-01",
		"draw":   "draw-01",
		"paint":  "paint-01",
		"sing":   "sing-01",
		"play":   "play-01",

		// Perception domain
		"see":    "see-01",
		"look":   "look-01",
		"watch":  "watch-01",
		"hear":   "hear-01",
		"listen": "listen-01",
		"feel":   "feel-01",
		"touch":  "touch-01",
		"smell":  "smell-01",
		"taste":  "taste-01",

		// Speech domain
		"say":     "say-01",
		"tell":    "tell-01",
		"speak":   "speak-01",
		"talk":    "talk-01",
		"ask":     "ask-01",
		"answer":  "answer-01",
		"call":    "call-01",
		"name":    "name-01",
		"promise": "promise-01",
		"offer":   "offer-01",

		// Statives domain
		"be":     "be-01",
		"become": "become-01",
		"seem":   "seem-01",
		"appear": "appear-01",
		"remain": "remain-01",
		"stay":   "stay-01",
		"lie":    "lie-01",
		"stand":  "stand-01",
		"sit":    "sit-01",

		// Time domain
		"begin":    "begin-01",
		"start":    "start-01",
		"end":      "end-01",
		"finish":   "finish-01",
		"continue": "continue-01",
		"stop":     "stop-01",
		"wait":     "wait-01",
		"happen":   "happen-01",
		"occur":    "occur-01",

		// Location domain
		"live": "live-01",

		// Common entities
		"boy":       "boy-01",
		"girl":      "girl-01",
		"man":       "man-01",
		"woman":     "woman-01",
		"person":    "person-01",
		"child":     "child-01",
		"book":      "book-01",
		"house":     "house-01",
		"car":       "car-01",
		"tree":      "tree-01",
		"water":     "water-01",
		"food":      "food-01",
		"money":     "money-01",
		"time":      "time-01",
		"day":       "day-01",
		"night":     "night-01",
		"year":      "year-01",
		"yesterday": "yesterday-01",
		"today":     "today-01",
		"tomorrow":  "tomorrow-01",
	}

	return resolver
}

// Resolve attempts to find a ConceptID for a given lemma in a language.
func (r *InMemoryConceptResolver) Resolve(lemma, lang string) (ConceptID, bool) {
	langConcepts, exists := r.concepts[lang]
	if !exists {
		return "", false
	}

	conceptID, found := langConcepts[lemma]
	return conceptID, found
}

// Register adds a new lemma-concept mapping for a language.
// This is primarily for testing and extension purposes.
func (r *InMemoryConceptResolver) Register(lemma, lang string, conceptID ConceptID) {
	if r.concepts[lang] == nil {
		r.concepts[lang] = make(map[string]ConceptID)
	}
	r.concepts[lang][lemma] = conceptID
}

// GetSupportedLanguages returns a list of languages with concept inventories.
func (r *InMemoryConceptResolver) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(r.concepts))
	for lang := range r.concepts {
		languages = append(languages, lang)
	}
	return languages
}

// GetConceptCount returns the number of concepts available for a language.
func (r *InMemoryConceptResolver) GetConceptCount(lang string) int {
	langConcepts, exists := r.concepts[lang]
	if !exists {
		return 0
	}
	return len(langConcepts)
}
