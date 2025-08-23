package interlingua

import (
	"fmt"
	"strings"
)

// EnglishRealizer implements the Realizer interface for English.
type EnglishRealizer struct {
	seed     int64
	resolver ConceptResolver
}

// NewEnglishRealizer creates a new English realizer with the specified seed and concept resolver.
func NewEnglishRealizer(seed int64, resolver ConceptResolver) *EnglishRealizer {
	return &EnglishRealizer{
		seed:     seed,
		resolver: resolver,
	}
}

// LanguageCode returns the language code for English.
func (r *EnglishRealizer) LanguageCode() string {
	return "en"
}

// Realize converts an Interlingua Document to English tokens.
func (r *EnglishRealizer) Realize(doc Document) ([]string, Trace, error) {
	if err := Validate(doc); err != nil {
		return nil, Trace{}, fmt.Errorf("document validation failed: %w", err)
	}

	// Initialize trace
	trace := Trace{
		SourceLanguageID: doc.Trace.SourceLanguageID,
		ConstructionID:   doc.Trace.ConstructionID,
		TokenToNode:      make(map[int]string),
		Notes:            make([]Note, 0),
	}

	// Build entity lookup map
	entityMap := make(map[string]Entity)
	for _, entity := range doc.Entities {
		entityMap[entity.ID] = entity
	}

	var tokens []string
	tokenIdx := 0

	// Process each event
	for _, event := range doc.Events {
		eventTokens, eventTrace, err := r.realizeEvent(event, entityMap, tokenIdx)
		if err != nil {
			return nil, trace, fmt.Errorf("failed to realize event %s: %w", event.ID, err)
		}

		// Merge traces
		for k, v := range eventTrace.TokenToNode {
			trace.TokenToNode[k] = v
		}
		trace.Notes = append(trace.Notes, eventTrace.Notes...)

		tokens = append(tokens, eventTokens...)
		tokenIdx += len(eventTokens)
	}

	return tokens, trace, nil
}

// realizeEvent converts a single event to English tokens.
func (r *EnglishRealizer) realizeEvent(event Event, entityMap map[string]Entity, startTokenIdx int) ([]string, Trace, error) {
	trace := Trace{
		TokenToNode: make(map[int]string),
		Notes:       make([]Note, 0),
	}

	var tokens []string
	tokenIdx := startTokenIdx

	// Get the main verb
	verb := r.getVerb(event.Predicate, event.TAM)
	if verb == "" {
		// Fallback to lemma hint or predicate
		if event.LemmaHint != "" {
			verb = event.LemmaHint
		} else {
			verb = string(event.Predicate)
			// Remove sense number if present
			if idx := strings.Index(verb, "-"); idx != -1 {
				verb = verb[:idx]
			}
		}
		trace.AddNote("warn", "verb.fallback", fmt.Sprintf("Using fallback verb '%s' for predicate %s", verb, event.Predicate))
	}

	// Add subject (agent or experiencer) first for SVO order
	var subjectID string
	var hasSubject bool
	if agentID, hasAgent := event.Roles[RoleAgent]; hasAgent {
		subjectID = agentID
		hasSubject = true
	} else if experiencerID, hasExperiencer := event.Roles[RoleExperiencer]; hasExperiencer {
		subjectID = experiencerID
		hasSubject = true
	}

	if hasSubject {
		subject := entityMap[subjectID]
		subjectTokens := r.realizeEntity(subject, true) // true = is subject
		tokens = append(tokens, subjectTokens...)
		for i, _ := range subjectTokens {
			trace.MapToken(tokenIdx+i, subjectID)
		}
		tokenIdx += len(subjectTokens)
	}

	// Handle negation and verb
	if event.TAM.Polarity == "neg" {
		// Add negation marker based on tense
		if event.TAM.Tense == "past" {
			tokens = append(tokens, "didn't")
		} else {
			// Check if subject is third person singular
			if hasSubject {
				if subject, exists := entityMap[subjectID]; exists {
					if subject.Feats.Person == 3 && subject.Feats.Number == "sg" {
						tokens = append(tokens, "doesn't")
					} else {
						tokens = append(tokens, "don't")
					}
				} else {
					tokens = append(tokens, "don't") // Default
				}
			} else {
				tokens = append(tokens, "don't") // Default
			}
		}
		trace.MapToken(tokenIdx, event.ID)
		tokenIdx++

		// Add base form of verb for negative sentences
		tokens = append(tokens, verb)
		trace.MapToken(tokenIdx, event.ID)
		tokenIdx++
	} else {
		// For positive sentences, add inflected verb
		if event.TAM.Tense == "past" {
			verb = r.getPastTense(verb)
		} else if event.TAM.Tense == "pres" {
			// Check if subject is third person singular
			if hasSubject {
				if subject, exists := entityMap[subjectID]; exists {
					if subject.Feats.Person == 3 && subject.Feats.Number == "sg" {
						verb = r.getThirdPersonSingular(verb)
					}
				}
			}
		}
		tokens = append(tokens, verb)
		trace.MapToken(tokenIdx, event.ID)
		tokenIdx++
	}

	// Add direct object (patient or stimulus)
	var objectID string
	var hasObject bool
	if patientID, hasPatient := event.Roles[RolePatient]; hasPatient {
		objectID = patientID
		hasObject = true
	} else if stimulusID, hasStimulus := event.Roles[RoleStimulus]; hasStimulus {
		objectID = stimulusID
		hasObject = true
	}

	if hasObject {
		object := entityMap[objectID]
		objectTokens := r.realizeEntity(object, false) // false = not subject
		tokens = append(tokens, objectTokens...)
		for i, _ := range objectTokens {
			trace.MapToken(tokenIdx+i, objectID)
		}
		tokenIdx += len(objectTokens)
	}

	// Add indirect object (recipient) if present
	if recipientID, hasRecipient := event.Roles[RoleRecipient]; hasRecipient {
		recipient := entityMap[recipientID]
		tokens = append(tokens, "to")
		recipientTokens := r.realizeEntity(recipient, false)
		tokens = append(tokens, recipientTokens...)
		for i, _ := range recipientTokens {
			trace.MapToken(tokenIdx+i, recipientID)
		}
		tokenIdx += len(recipientTokens) + 1 // +1 for "to"
	}

	// Add adjuncts
	adjunctTokens, adjunctTrace := r.realizeAdjuncts(event, entityMap, tokenIdx)
	tokens = append(tokens, adjunctTokens...)

	// Merge adjunct trace
	for k, v := range adjunctTrace.TokenToNode {
		trace.TokenToNode[k] = v
	}
	trace.Notes = append(trace.Notes, adjunctTrace.Notes...)

	return tokens, trace, nil
}

// realizeEntity converts an entity to English tokens.
func (r *EnglishRealizer) realizeEntity(entity Entity, isSubject bool) []string {
	var tokens []string

	// Add determiner/article
	if entity.Feats.Definiteness == "def" {
		tokens = append(tokens, "the")
	} else if entity.Feats.Definiteness == "indef" {
		// Check if we need "a" or "an"
		if entity.LemmaHint != "" {
			firstChar := strings.ToLower(entity.LemmaHint[0:1])
			if strings.Contains("aeiou", firstChar) {
				tokens = append(tokens, "an")
			} else {
				tokens = append(tokens, "a")
			}
		} else {
			tokens = append(tokens, "a")
		}
	}

	// Add the noun
	noun := r.getNoun(entity.Concept, entity.LemmaHint)
	if noun == "" {
		// Fallback to lemma hint or concept
		if entity.LemmaHint != "" {
			noun = entity.LemmaHint
		} else {
			noun = string(entity.Concept)
			// Remove sense number if present
			if idx := strings.Index(noun, "-"); idx != -1 {
				noun = noun[:idx]
			}
		}
	}

	// Handle pluralization
	if entity.Feats.Number == "pl" {
		noun = r.pluralize(noun)
	}

	tokens = append(tokens, noun)

	return tokens
}

// realizeAdjuncts converts adjuncts to English tokens.
func (r *EnglishRealizer) realizeAdjuncts(event Event, entityMap map[string]Entity, startTokenIdx int) ([]string, Trace) {
	trace := Trace{
		TokenToNode: make(map[int]string),
		Notes:       make([]Note, 0),
	}

	var tokens []string
	tokenIdx := startTokenIdx

	// Time adjuncts
	if timeID, hasTime := event.Adjuncts[RoleTime]; hasTime {
		timeEntity := entityMap[timeID]
		timeTokens := r.realizeTimeAdjunct(timeEntity)
		tokens = append(tokens, timeTokens...)
		for i, _ := range timeTokens {
			trace.MapToken(tokenIdx+i, timeID)
		}
		tokenIdx += len(timeTokens)
	}

	// Location adjuncts
	if locationID, hasLocation := event.Adjuncts[RoleLocation]; hasLocation {
		locationEntity := entityMap[locationID]
		locationTokens := r.realizeLocationAdjunct(locationEntity)
		tokens = append(tokens, locationTokens...)
		for i, _ := range locationTokens {
			trace.MapToken(tokenIdx+i, locationID)
		}
		tokenIdx += len(locationTokens)
	}

	// Manner adjuncts
	if mannerID, hasManner := event.Adjuncts[RoleManner]; hasManner {
		mannerEntity := entityMap[mannerID]
		mannerTokens := r.realizeMannerAdjunct(mannerEntity)
		tokens = append(tokens, mannerTokens...)
		for i, _ := range mannerTokens {
			trace.MapToken(tokenIdx+i, mannerID)
		}
		tokenIdx += len(mannerTokens)
	}

	return tokens, trace
}

// realizeTimeAdjunct converts a time adjunct to English tokens.
func (r *EnglishRealizer) realizeTimeAdjunct(timeEntity Entity) []string {
	// Handle common time expressions
	timeWord := r.getTimeWord(timeEntity.Concept, timeEntity.LemmaHint)
	if timeWord != "" {
		return []string{timeWord}
	}

	// Fallback to basic realization
	return r.realizeEntity(timeEntity, false)
}

// realizeLocationAdjunct converts a location adjunct to English tokens.
func (r *EnglishRealizer) realizeLocationAdjunct(locationEntity Entity) []string {
	// Add preposition for location
	preposition := "in"
	if locationEntity.Concept == "house" || locationEntity.Concept == "building" {
		preposition = "in"
	} else if locationEntity.Concept == "tree" || locationEntity.Concept == "mountain" {
		preposition = "on"
	}

	tokens := []string{preposition}
	entityTokens := r.realizeEntity(locationEntity, false)
	tokens = append(tokens, entityTokens...)
	return tokens
}

// realizeMannerAdjunct converts a manner adjunct to English tokens.
func (r *EnglishRealizer) realizeMannerAdjunct(mannerEntity Entity) []string {
	// Add "with" for manner
	tokens := []string{"with"}
	entityTokens := r.realizeEntity(mannerEntity, false)
	tokens = append(tokens, entityTokens...)
	return tokens
}

// getVerb returns the appropriate English verb form for a concept and TAM.
func (r *EnglishRealizer) getVerb(concept ConceptID, tam TAM) string {
	// Try to resolve the concept
	if r.resolver != nil {
		if lemma, found := r.reverseResolve(concept, "en"); found {
			return lemma
		}
	}

	// Fallback to concept-based mapping
	conceptStr := string(concept)
	if idx := strings.Index(conceptStr, "-"); idx != -1 {
		conceptStr = conceptStr[:idx]
	}

	// Common verb mappings
	verbMap := map[string]string{
		"give":    "give",
		"send":    "send",
		"move":    "move",
		"see":     "see",
		"hear":    "hear",
		"make":    "make",
		"take":    "take",
		"bring":   "bring",
		"carry":   "carry",
		"walk":    "walk",
		"run":     "run",
		"fly":     "fly",
		"swim":    "swim",
		"jump":    "jump",
		"climb":   "climb",
		"write":   "write",
		"draw":    "draw",
		"paint":   "paint",
		"sing":    "sing",
		"play":    "play",
		"say":     "say",
		"tell":    "tell",
		"speak":   "speak",
		"talk":    "talk",
		"ask":     "ask",
		"answer":  "answer",
		"call":    "call",
		"name":    "name",
		"promise": "promise",
		"offer":   "offer",
		"have":    "have",
		"own":     "own",
		"keep":    "keep",
		"hold":    "hold",
		"begin":   "begin",
		"start":   "start",
		"end":     "end",
		"stop":    "stop",
		"wait":    "wait",
		"live":    "live",
		"stay":    "stay",
		"remain":  "remain",
		"lie":     "lie",
		"stand":   "stand",
		"sit":     "sit",
	}

	if verb, found := verbMap[conceptStr]; found {
		return verb
	}

	return ""
}

// getNoun returns the appropriate English noun for a concept.
func (r *EnglishRealizer) getNoun(concept ConceptID, lemmaHint string) string {
	// Try to resolve the concept
	if r.resolver != nil {
		if lemma, found := r.reverseResolve(concept, "en"); found {
			return lemma
		}
	}

	// Fallback to concept-based mapping
	conceptStr := string(concept)
	if idx := strings.Index(conceptStr, "-"); idx != -1 {
		conceptStr = conceptStr[:idx]
	}

	// Common noun mappings
	nounMap := map[string]string{
		"boy":    "boy",
		"girl":   "girl",
		"man":    "man",
		"woman":  "woman",
		"person": "person",
		"child":  "child",
		"book":   "book",
		"house":  "house",
		"car":    "car",
		"tree":   "tree",
		"water":  "water",
		"food":   "food",
		"money":  "money",
		"time":   "time",
		"day":    "day",
		"night":  "night",
		"year":   "year",
	}

	if noun, found := nounMap[conceptStr]; found {
		return noun
	}

	return ""
}

// getTimeWord returns the appropriate English time expression.
func (r *EnglishRealizer) getTimeWord(concept ConceptID, lemmaHint string) string {
	// Try to resolve the concept
	if r.resolver != nil {
		if lemma, found := r.reverseResolve(concept, "en"); found {
			return lemma
		}
	}

	// Common time word mappings
	timeMap := map[string]string{
		"yesterday": "yesterday",
		"today":     "today",
		"tomorrow":  "tomorrow",
		"now":       "now",
		"then":      "then",
		"soon":      "soon",
		"later":     "later",
		"earlier":   "earlier",
	}

	conceptStr := string(concept)
	if idx := strings.Index(conceptStr, "-"); idx != -1 {
		conceptStr = conceptStr[:idx]
	}

	if timeWord, found := timeMap[conceptStr]; found {
		return timeWord
	}

	return ""
}

// getThirdPersonSingular returns the third person singular present form of a verb.
func (r *EnglishRealizer) getThirdPersonSingular(verb string) string {
	// Irregular verbs
	irregularMap := map[string]string{
		"have": "has",
		"do":   "does",
		"go":   "goes",
	}

	if thirdForm, found := irregularMap[verb]; found {
		return thirdForm
	}

	// Regular verb rules
	if strings.HasSuffix(verb, "y") && !strings.Contains("aeiou", strings.ToLower(verb[len(verb)-2:len(verb)-1])) {
		return verb[:len(verb)-1] + "ies"
	}
	if strings.HasSuffix(verb, "s") || strings.HasSuffix(verb, "sh") || strings.HasSuffix(verb, "ch") || strings.HasSuffix(verb, "x") || strings.HasSuffix(verb, "z") {
		return verb + "es"
	}

	// Default: add -s
	return verb + "s"
}

// getPastTense returns the past tense form of a verb.
func (r *EnglishRealizer) getPastTense(verb string) string {
	// Irregular verbs
	irregularMap := map[string]string{
		"give":    "gave",
		"send":    "sent",
		"see":     "saw",
		"hear":    "heard",
		"make":    "made",
		"take":    "took",
		"bring":   "brought",
		"carry":   "carried",
		"write":   "wrote",
		"draw":    "drew",
		"paint":   "painted",
		"sing":    "sang",
		"play":    "played",
		"say":     "said",
		"tell":    "told",
		"speak":   "spoke",
		"talk":    "talked",
		"ask":     "asked",
		"answer":  "answered",
		"call":    "called",
		"name":    "named",
		"promise": "promised",
		"offer":   "offered",
		"have":    "had",
		"own":     "owned",
		"keep":    "kept",
		"hold":    "held",
		"begin":   "began",
		"start":   "started",
		"end":     "ended",
		"stop":    "stopped",
		"wait":    "waited",
		"live":    "lived",
		"stay":    "stayed",
		"remain":  "remained",
		"lie":     "lay",
		"stand":   "stood",
		"sit":     "sat",
	}

	if pastForm, found := irregularMap[verb]; found {
		return pastForm
	}

	// Regular verb rules
	if strings.HasSuffix(verb, "e") {
		return verb + "d"
	}
	if strings.HasSuffix(verb, "y") && !strings.Contains("aeiou", strings.ToLower(verb[len(verb)-2:len(verb)-1])) {
		return verb[:len(verb)-1] + "ied"
	}
	if strings.HasSuffix(verb, "c") {
		return verb + "ked"
	}

	// Default: add -ed
	return verb + "ed"
}

// pluralize returns the plural form of a noun.
func (r *EnglishRealizer) pluralize(noun string) string {
	// Irregular plurals
	irregularMap := map[string]string{
		"child":  "children",
		"man":    "men",
		"woman":  "women",
		"person": "people",
		"foot":   "feet",
		"tooth":  "teeth",
		"goose":  "geese",
		"mouse":  "mice",
		"louse":  "lice",
		"ox":     "oxen",
	}

	if plural, found := irregularMap[noun]; found {
		return plural
	}

	// Regular plural rules
	if strings.HasSuffix(noun, "y") && !strings.Contains("aeiou", strings.ToLower(noun[len(noun)-2:len(noun)-1])) {
		return noun[:len(noun)-1] + "ies"
	}
	if strings.HasSuffix(noun, "f") || strings.HasSuffix(noun, "fe") {
		if strings.HasSuffix(noun, "fe") {
			return noun[:len(noun)-2] + "ves"
		}
		return noun[:len(noun)-1] + "ves"
	}
	if strings.HasSuffix(noun, "s") || strings.HasSuffix(noun, "sh") || strings.HasSuffix(noun, "ch") || strings.HasSuffix(noun, "x") || strings.HasSuffix(noun, "z") {
		return noun + "es"
	}

	// Default: add -s
	return noun + "s"
}

// reverseResolve attempts to find a lemma for a given ConceptID.
// This is a simple reverse lookup for the in-memory resolver.
func (r *EnglishRealizer) reverseResolve(conceptID ConceptID, lang string) (string, bool) {
	if r.resolver == nil {
		return "", false
	}

	// For the in-memory resolver, we can do a simple reverse lookup
	// This is not efficient for large concept inventories, but works for testing
	if inMemResolver, ok := r.resolver.(*InMemoryConceptResolver); ok {
		langConcepts := inMemResolver.concepts[lang]
		if langConcepts == nil {
			return "", false
		}

		for lemma, concept := range langConcepts {
			if concept == conceptID {
				return lemma, true
			}
		}
	}

	return "", false
}
