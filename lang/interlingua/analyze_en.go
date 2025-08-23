package interlingua

import (
	"fmt"
	"strings"
)

// EnglishAnalyzer implements the Analyzer interface for English.
type EnglishAnalyzer struct {
	resolver ConceptResolver
}

// NewEnglishAnalyzer creates a new English analyzer with the specified concept resolver.
func NewEnglishAnalyzer(resolver ConceptResolver) *EnglishAnalyzer {
	return &EnglishAnalyzer{
		resolver: resolver,
	}
}

// LanguageCode returns the language code for English.
func (a *EnglishAnalyzer) LanguageCode() string {
	return "en"
}

// Analyze converts English tokens to an Interlingua Document.
func (a *EnglishAnalyzer) Analyze(tokens []string) (Document, error) {
	if len(tokens) == 0 {
		return Document{}, fmt.Errorf("no tokens to analyze")
	}

	doc := NewDocument("en")
	doc.Trace.SetSourceLanguage("en")

	// Simple tokenization and pattern matching
	analysis, err := a.parseSimplePattern(tokens)
	if err != nil {
		return Document{}, fmt.Errorf("failed to parse pattern: %w", err)
	}

	// Build entities
	entityMap := make(map[Role]Entity)
	for role, entity := range analysis.entities {
		entityID := fmt.Sprintf("e%d", len(doc.Entities)+1)
		entity.ID = entityID
		AddEntity(&doc, entity)
		entityMap[role] = entity
	}

	// Build event
	event := Event{
		ID:        fmt.Sprintf("v%d", len(doc.Events)+1),
		Predicate: analysis.predicate,
		TAM:       analysis.tam,
		Roles:     make(map[Role]string),
		Valency:   len(analysis.entities),
	}

	// Assign roles
	for role, entity := range entityMap {
		event.Roles[role] = entity.ID
	}

	// Add event to document
	AddEvent(&doc, event)

	// Add diagnostic notes
	if analysis.notes != nil {
		doc.Trace.Notes = append(doc.Trace.Notes, analysis.notes...)
	}

	return doc, nil
}

// parseResult holds the parsed information from pattern matching.
type parseResult struct {
	predicate ConceptID
	tam       TAM
	entities  map[Role]Entity
	notes     []Note
}

// parseSimplePattern performs lightweight parsing using simple patterns.
func (a *EnglishAnalyzer) parseSimplePattern(tokens []string) (*parseResult, error) {
	result := &parseResult{
		entities: make(map[Role]Entity),
		notes:    make([]Note, 0),
		tam:      TAM{Polarity: "pos"}, // Default to positive
	}

	// Simple pattern matching for basic clause types
	if len(tokens) < 3 {
		return nil, fmt.Errorf("insufficient tokens for analysis")
	}

	// Check for negation
	negIndex := -1
	for i, token := range tokens {
		if token == "didn't" || token == "don't" || token == "doesn't" {
			negIndex = i
			result.tam.Polarity = "neg"
			break
		}
	}

	// Find verb (usually after negation or at beginning)
	verbIndex := 0
	if negIndex != -1 {
		verbIndex = negIndex + 1
	}

	// Find the actual verb by looking for a word that could be a verb
	for i := verbIndex; i < len(tokens); i++ {
		if a.isVerb(tokens[i]) {
			verbIndex = i
			break
		}
	}

	if verbIndex >= len(tokens) {
		return nil, fmt.Errorf("no verb found")
	}

	verb := tokens[verbIndex]
	result.predicate = a.resolveVerb(verb)

	// Determine word order pattern
	// Check if this is VSO (verb-subject-object) or SVO (subject-verb-object)
	// For now, assume VSO if verb is at position 0 or 1, SVO otherwise
	isVSO := verbIndex <= 1
	result.notes = append(result.notes, Note{
		Severity: "info",
		Code:     "word.order.detected",
		Message:  fmt.Sprintf("Detected %s word order", map[bool]string{true: "VSO", false: "SVO"}[isVSO]),
	})

	// Determine tense
	if negIndex != -1 {
		// Check the specific negation marker
		negMarker := tokens[negIndex]
		if negMarker == "didn't" {
			result.tam.Tense = "past"
		} else if negMarker == "don't" || negMarker == "doesn't" {
			result.tam.Tense = "pres"
		} else {
			// Default to present for other negation markers
			result.tam.Tense = "pres"
		}
	} else {
		result.tam.Tense = a.detectTense(verb)
	}

	// Parse arguments based on position
	a.parseArguments(tokens, verbIndex, result)

	// Add diagnostic notes for assumptions
	if result.tam.Tense == "" {
		result.tam.Tense = "pres" // Default to present
		result.notes = append(result.notes, Note{
			Severity: "warn",
			Code:     "tense.defaulted",
			Message:  "Tense defaulted to present",
		})
	}

	return result, nil
}

// parseArguments parses the arguments based on position and patterns.
func (a *EnglishAnalyzer) parseArguments(tokens []string, verbIndex int, result *parseResult) {
	// Simple heuristics: subject→agent, dobj→patient, PP "to"→recipient

	// Determine word order based on verb position
	isVSO := verbIndex <= 1

	if isVSO {
		// VSO order: verb-subject-object
		// Subject is after verb
		subjectIndex := a.findSubjectIndexVSO(tokens, verbIndex)
		if subjectIndex != -1 {
			subject := a.parseNounPhrase(tokens, subjectIndex)
			// Only set defaults if not already set
			if subject.Feats.Person == 0 {
				subject.Feats.Person = 3 // Default to third person
			}
			if subject.Feats.Number == "" {
				subject.Feats.Number = "sg" // Default to singular
			}
			result.entities[RoleAgent] = subject

			// Direct object is after subject
			dobjIndex := a.findDirectObjectIndexVSO(tokens, verbIndex, subjectIndex)
			if dobjIndex != -1 {
				patient := a.parseNounPhrase(tokens, dobjIndex)
				result.entities[RolePatient] = patient
			}
		}
	} else {
		// SVO order: subject-verb-object
		// Subject is before verb
		subjectIndex := a.findSubjectIndexSVO(tokens, verbIndex)
		if subjectIndex != -1 {
			subject := a.parseNounPhrase(tokens, subjectIndex)
			// Only set defaults if not already set
			if subject.Feats.Person == 0 {
				subject.Feats.Person = 3 // Default to third person
			}
			if subject.Feats.Number == "" {
				subject.Feats.Number = "sg" // Default to singular
			}
			result.entities[RoleAgent] = subject
		}

		// Direct object is after verb
		dobjIndex := a.findDirectObjectIndexSVO(tokens, verbIndex)
		if dobjIndex != -1 {
			patient := a.parseNounPhrase(tokens, dobjIndex)
			result.entities[RolePatient] = patient
		}
	}

	// Find indirect object (usually "to" + NP) - same for both orders
	recipientIndex := a.findIndirectObjectIndex(tokens, verbIndex)
	if recipientIndex != -1 {
		recipient := a.parseNounPhrase(tokens, recipientIndex)
		result.entities[RoleRecipient] = recipient
	}

	// Find time adjuncts
	timeIndex := a.findTimeAdjunctIndex(tokens)
	if timeIndex != -1 {
		_ = a.parseTimeAdjunct(tokens, timeIndex)
		// Note: We'll add this as an adjunct later
		result.notes = append(result.notes, Note{
			Severity: "info",
			Code:     "time.adjunct.found",
			Message:  fmt.Sprintf("Time adjunct found: %s", tokens[timeIndex]),
		})
	}

	// Find location adjuncts
	locationIndex := a.findLocationAdjunctIndex(tokens)
	if locationIndex != -1 {
		_ = a.parseLocationAdjunct(tokens, locationIndex)
		// Note: We'll add this as an adjunct later
		result.notes = append(result.notes, Note{
			Severity: "info",
			Code:     "location.adjunct.found",
			Message:  fmt.Sprintf("Location adjunct found: %s", tokens[locationIndex]),
		})
	}
}

// findSubjectIndexVSO finds the subject noun phrase index for VSO order.
func (a *EnglishAnalyzer) findSubjectIndexVSO(tokens []string, verbIndex int) int {
	// Subject is after verb in VSO order
	for i := verbIndex + 1; i < len(tokens); i++ {
		if a.isNounPhrase(tokens, i) {
			return i
		}
	}
	return -1
}

// findDirectObjectIndexVSO finds the direct object noun phrase index for VSO order.
func (a *EnglishAnalyzer) findDirectObjectIndexVSO(tokens []string, verbIndex int, subjectIndex int) int {
	// Direct object is after subject in VSO order
	for i := subjectIndex + 1; i < len(tokens); i++ {
		if a.isNounPhrase(tokens, i) {
			return i
		}
	}
	return -1
}

// findSubjectIndexSVO finds the subject noun phrase index for SVO order.
func (a *EnglishAnalyzer) findSubjectIndexSVO(tokens []string, verbIndex int) int {
	// Subject is before verb in SVO order
	for i := verbIndex - 1; i >= 0; i-- {
		if a.isNounPhrase(tokens, i) {
			return i
		}
	}
	return -1
}

// findDirectObjectIndexSVO finds the direct object noun phrase index for SVO order.
func (a *EnglishAnalyzer) findDirectObjectIndexSVO(tokens []string, verbIndex int) int {
	// Direct object is after verb in SVO order
	for i := verbIndex + 1; i < len(tokens); i++ {
		if a.isNounPhrase(tokens, i) {
			return i
		}
	}
	return -1
}

// findSubjectIndex finds the subject noun phrase index (legacy method).
func (a *EnglishAnalyzer) findSubjectIndex(tokens []string, verbIndex int) int {
	// Subject is usually the first noun phrase before the verb
	for i := verbIndex - 1; i >= 0; i-- {
		if a.isNounPhrase(tokens, i) {
			return i
		}
	}
	return -1
}

// findDirectObjectIndex finds the direct object noun phrase index.
func (a *EnglishAnalyzer) findDirectObjectIndex(tokens []string, verbIndex int) int {
	// Direct object is usually the first noun phrase after the verb
	for i := verbIndex + 1; i < len(tokens); i++ {
		if a.isNounPhrase(tokens, i) {
			return i
		}
	}
	return -1
}

// findIndirectObjectIndex finds the indirect object noun phrase index.
func (a *EnglishAnalyzer) findIndirectObjectIndex(tokens []string, verbIndex int) int {
	// Look for "to" + NP pattern
	for i := verbIndex + 1; i < len(tokens)-1; i++ {
		if tokens[i] == "to" && a.isNounPhrase(tokens, i+1) {
			return i + 1
		}
	}
	return -1
}

// findTimeAdjunctIndex finds time adjunct tokens.
func (a *EnglishAnalyzer) findTimeAdjunctIndex(tokens []string) int {
	timeWords := []string{"yesterday", "today", "tomorrow", "now", "then", "soon", "later", "earlier"}
	for i, token := range tokens {
		for _, timeWord := range timeWords {
			if token == timeWord {
				return i
			}
		}
	}
	return -1
}

// findLocationAdjunctIndex finds location adjunct patterns.
func (a *EnglishAnalyzer) findLocationAdjunctIndex(tokens []string) int {
	// Look for preposition + NP patterns
	prepositions := []string{"in", "on", "at", "by", "near", "under", "over"}
	for i := 0; i < len(tokens)-1; i++ {
		for _, prep := range prepositions {
			if tokens[i] == prep && a.isNounPhrase(tokens, i+1) {
				return i
			}
		}
	}
	return -1
}

// isNounPhrase checks if a token position starts a noun phrase.
func (a *EnglishAnalyzer) isNounPhrase(tokens []string, index int) bool {
	if index < 0 || index >= len(tokens) {
		return false
	}

	token := tokens[index]

	// Check for articles - only return true if there's a noun following
	if token == "the" || token == "a" || token == "an" {
		// Check if there's a noun following
		if index+1 < len(tokens) {
			nextToken := tokens[index+1]
			// Check if next token is a noun or pronoun
			pronouns := []string{"he", "she", "it", "they", "we", "you", "I"}
			commonNouns := []string{"boy", "girl", "man", "woman", "book", "house", "car", "tree", "boys", "girls", "men", "women", "books", "houses", "cars", "trees", "child", "children", "person", "people", "money", "water", "food", "time", "day", "night", "year"}

			for _, pronoun := range pronouns {
				if nextToken == pronoun {
					return true
				}
			}
			for _, noun := range commonNouns {
				if nextToken == noun {
					return true
				}
			}
		}
		return false
	}

	// Check for pronouns
	pronouns := []string{"he", "she", "it", "they", "we", "you", "I"}
	for _, pronoun := range pronouns {
		if token == pronoun {
			return true
		}
	}

	// Check for common nouns
	commonNouns := []string{"boy", "girl", "man", "woman", "book", "house", "car", "tree", "boys", "girls", "men", "women", "books", "houses", "cars", "trees", "child", "children", "person", "people", "money", "water", "food", "time", "day", "night", "year"}
	for _, noun := range commonNouns {
		if token == noun {
			return true
		}
	}

	return false
}

// parseNounPhrase parses a noun phrase starting at the given index.
func (a *EnglishAnalyzer) parseNounPhrase(tokens []string, startIndex int) Entity {
	entity := Entity{
		Feats: NominalFeatures{},
	}

	// Parse determiner
	if startIndex < len(tokens) {
		det := tokens[startIndex]
		if det == "the" {
			entity.Feats.Definiteness = "def"
			startIndex++
		} else if det == "a" || det == "an" {
			entity.Feats.Definiteness = "indef"
			startIndex++
		}
	}

	// Parse noun
	if startIndex < len(tokens) {
		noun := tokens[startIndex]
		entity.Concept = a.resolveNoun(noun)
		entity.LemmaHint = noun

		// Detect number
		if strings.HasSuffix(noun, "s") {
			entity.Feats.Number = "pl"
		} else {
			entity.Feats.Number = "sg"
		}

		// Detect person for pronouns
		if noun == "I" || noun == "we" {
			entity.Feats.Person = 1
		} else if noun == "you" {
			entity.Feats.Person = 2
		} else {
			entity.Feats.Person = 3
		}
	}

	return entity
}

// parseTimeAdjunct parses a time adjunct.
func (a *EnglishAnalyzer) parseTimeAdjunct(tokens []string, index int) Entity {
	timeWord := tokens[index]
	return Entity{
		Concept:   a.resolveTimeWord(timeWord),
		LemmaHint: timeWord,
		Feats:     NominalFeatures{},
	}
}

// parseLocationAdjunct parses a location adjunct.
func (a *EnglishAnalyzer) parseLocationAdjunct(tokens []string, index int) Entity {
	// Skip preposition, parse the noun phrase
	if index > 0 && a.isPreposition(tokens[index-1]) {
		return a.parseNounPhrase(tokens, index)
	}
	return Entity{
		Concept:   "location",
		LemmaHint: "location",
		Feats:     NominalFeatures{},
	}
}

// isVerb checks if a token could be a verb.
func (a *EnglishAnalyzer) isVerb(token string) bool {
	// Check for common verbs
	commonVerbs := []string{"give", "gives", "gave", "send", "sends", "sent", "move", "moves", "moved", "see", "sees", "saw", "walk", "walks", "walked", "run", "runs", "ran", "make", "makes", "made", "take", "takes", "took", "bring", "brings", "brought", "carry", "carries", "carried", "have", "has", "had"}
	for _, verb := range commonVerbs {
		if token == verb {
			return true
		}
	}

	// Check for verb-like patterns, but be more restrictive
	if strings.HasSuffix(token, "s") || strings.HasSuffix(token, "ed") {
		// Exclude articles and common non-verbs
		excluded := []string{"the", "this", "that", "these", "those", "his", "her", "its", "our", "your", "their"}
		for _, exclude := range excluded {
			if token == exclude {
				return false
			}
		}
		return true
	}

	return false
}

// isPreposition checks if a token is a preposition.
func (a *EnglishAnalyzer) isPreposition(token string) bool {
	prepositions := []string{"in", "on", "at", "by", "near", "under", "over", "to", "from", "with"}
	for _, prep := range prepositions {
		if token == prep {
			return true
		}
	}
	return false
}

// resolveVerb resolves a verb token to a ConceptID.
func (a *EnglishAnalyzer) resolveVerb(verb string) ConceptID {
	// Try to resolve via ConceptResolver
	if a.resolver != nil {
		if conceptID, found := a.resolver.Resolve(verb, "en"); found {
			return conceptID
		}
	}

	// Fallback to common verb mappings
	verbMap := map[string]ConceptID{
		"give":    "give-01",
		"gives":   "give-01",
		"gave":    "give-01",
		"send":    "send-01",
		"sends":   "send-01",
		"sent":    "send-01",
		"move":    "move-00",
		"moves":   "move-00",
		"moved":   "move-00",
		"see":     "see-01",
		"sees":    "see-01",
		"saw":     "see-01",
		"walk":    "walk-01",
		"walks":   "walk-01",
		"walked":  "walk-01",
		"run":     "run-01",
		"runs":    "run-01",
		"ran":     "run-01",
		"make":    "make-01",
		"makes":   "make-01",
		"made":    "make-01",
		"take":    "take-01",
		"takes":   "take-01",
		"took":    "take-01",
		"bring":   "bring-01",
		"brings":  "bring-01",
		"brought": "bring-01",
		"carry":   "carry-01",
		"carries": "carry-01",
		"carried": "carry-01",
		"have":    "have-01",
		"has":     "have-01",
		"had":     "have-01",
	}

	if conceptID, found := verbMap[verb]; found {
		return conceptID
	}

	// Default fallback
	return ConceptID(verb + "-01")
}

// resolveNoun resolves a noun token to a ConceptID.
func (a *EnglishAnalyzer) resolveNoun(noun string) ConceptID {
	// Try to resolve via ConceptResolver
	if a.resolver != nil {
		if conceptID, found := a.resolver.Resolve(noun, "en"); found {
			return conceptID
		}
	}

	// Fallback to common noun mappings
	nounMap := map[string]ConceptID{
		"boy":      "boy",
		"boys":     "boy",
		"girl":     "girl",
		"girls":    "girl",
		"man":      "man",
		"men":      "man",
		"woman":    "woman",
		"women":    "woman",
		"person":   "person",
		"people":   "person",
		"child":    "child",
		"children": "child",
		"book":     "book",
		"books":    "book",
		"house":    "house",
		"houses":   "house",
		"car":      "car",
		"cars":     "car",
		"tree":     "tree",
		"trees":    "tree",
		"water":    "water",
		"food":     "food",
		"money":    "money",
		"time":     "time",
		"day":      "day",
		"night":    "night",
		"year":     "year",
	}

	if conceptID, found := nounMap[noun]; found {
		return conceptID
	}

	// Default fallback
	return ConceptID(noun)
}

// resolveTimeWord resolves a time word to a ConceptID.
func (a *EnglishAnalyzer) resolveTimeWord(timeWord string) ConceptID {
	// Try to resolve via ConceptResolver
	if a.resolver != nil {
		if conceptID, found := a.resolver.Resolve(timeWord, "en"); found {
			return conceptID
		}
	}

	// Fallback to common time word mappings
	timeMap := map[string]ConceptID{
		"yesterday": "yesterday",
		"today":     "today",
		"tomorrow":  "tomorrow",
		"now":       "now",
		"then":      "then",
		"soon":      "soon",
		"later":     "later",
		"earlier":   "earlier",
	}

	if conceptID, found := timeMap[timeWord]; found {
		return conceptID
	}

	// Default fallback
	return ConceptID(timeWord)
}

// detectTense attempts to detect the tense of a verb.
func (a *EnglishAnalyzer) detectTense(verb string) string {
	// Check for past tense markers
	if strings.HasSuffix(verb, "ed") {
		return "past"
	}

	// Check for irregular past forms
	irregularPast := []string{"gave", "sent", "saw", "heard", "made", "took", "brought", "carried", "wrote", "drew", "sang", "said", "told", "spoke", "asked", "answered", "called", "named", "promised", "offered", "had", "owned", "kept", "held", "began", "started", "ended", "stopped", "waited", "lived", "stayed", "remained", "lay", "stood", "sat"}
	for _, pastForm := range irregularPast {
		if verb == pastForm {
			return "past"
		}
	}

	// Check for third person singular present
	if strings.HasSuffix(verb, "s") {
		return "pres"
	}

	// Default to present
	return "pres"
}
