package interlingua

import (
	"fmt"
)

// NewDocument creates a new empty Document with the specified language ID.
func NewDocument(langID string) Document {
	return Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    langID,
		Entities:      make([]Entity, 0),
		Events:        make([]Event, 0),
		Trace: Trace{
			TokenToNode: make(map[int]string),
			Notes:       make([]Note, 0),
		},
	}
}

// AddEntity adds an entity to a document and returns its assigned ID.
// If the entity has no ID, one will be generated automatically.
func AddEntity(doc *Document, e Entity) string {
	if e.ID == "" {
		// Generate a unique ID if none provided
		e.ID = fmt.Sprintf("e%d", len(doc.Entities)+1)
	}

	doc.Entities = append(doc.Entities, e)
	return e.ID
}

// AddEvent adds an event to a document and returns its assigned ID.
// If the event has no ID, one will be generated automatically.
func AddEvent(doc *Document, ev Event) string {
	if ev.ID == "" {
		// Generate a unique ID if none provided
		ev.ID = fmt.Sprintf("v%d", len(doc.Events)+1)
	}

	doc.Events = append(doc.Events, ev)
	return ev.ID
}

// SetRole assigns a semantic role to an event, linking it to an entity.
func SetRole(ev *Event, role Role, entityID string) error {
	if ev.Roles == nil {
		ev.Roles = make(map[Role]string)
	}

	// Validate role compatibility
	if err := ValidateRoleCompatibility(ev.Predicate, role); err != nil {
		return fmt.Errorf("invalid role assignment: %w", err)
	}

	ev.Roles[role] = entityID
	return nil
}

// AttachAdjunct attaches an adjunct (time, location, manner, etc.) to an event.
func AttachAdjunct(ev *Event, role Role, entityID string) {
	if ev.Adjuncts == nil {
		ev.Adjuncts = make(map[Role]string)
	}

	ev.Adjuncts[role] = entityID
}

// BuildFromTemplate builds a complete sentence Document from known template information.
// This is a pure function that constructs the semantic representation without
// requiring access to existing packages.
func BuildFromTemplate(langID, constructionID string, predicate ConceptID, tam TAM, args map[Role]Entity, adjuncts map[Role]Entity) (Document, error) {
	doc := NewDocument(langID)
	doc.Trace.SetConstructionID(constructionID)

	// Add entities for arguments
	entityIDs := make(map[Role]string)
	for role, entity := range args {
		entityID := AddEntity(&doc, entity)
		entityIDs[role] = entityID
	}

	// Add entities for adjuncts
	adjunctIDs := make(map[Role]string)
	for role, entity := range adjuncts {
		entityID := AddEntity(&doc, entity)
		adjunctIDs[role] = entityID
	}

	// Create the main event
	event := Event{
		Predicate: predicate,
		TAM:       tam,
		Roles:     make(map[Role]string),
		Adjuncts:  make(map[Role]string),
		Valency:   len(args),
	}

	// Assign roles
	for role, entityID := range entityIDs {
		if err := SetRole(&event, role, entityID); err != nil {
			return Document{}, fmt.Errorf("failed to set role %s: %w", role, err)
		}
	}

	// Attach adjuncts
	for role, entityID := range adjunctIDs {
		AttachAdjunct(&event, role, entityID)
	}

	// Add the event to the document
	AddEvent(&doc, event)

	return doc, nil
}

// BuildSimpleTransitive builds a simple transitive sentence (agent + patient).
func BuildSimpleTransitive(langID, constructionID string, predicate ConceptID, tam TAM, agent, patient Entity) (Document, error) {
	args := map[Role]Entity{
		RoleAgent:   agent,
		RolePatient: patient,
	}

	return BuildFromTemplate(langID, constructionID, predicate, tam, args, nil)
}

// BuildDitransitive builds a ditransitive sentence (agent + patient + recipient).
func BuildDitransitive(langID, constructionID string, predicate ConceptID, tam TAM, agent, patient, recipient Entity) (Document, error) {
	args := map[Role]Entity{
		RoleAgent:     agent,
		RolePatient:   patient,
		RoleRecipient: recipient,
	}

	return BuildFromTemplate(langID, constructionID, predicate, tam, args, nil)
}

// BuildWithAdjuncts builds a sentence with optional adjuncts.
func BuildWithAdjuncts(langID, constructionID string, predicate ConceptID, tam TAM, args map[Role]Entity, adjuncts map[Role]Entity) (Document, error) {
	return BuildFromTemplate(langID, constructionID, predicate, tam, args, adjuncts)
}

// AddTimeAdjunct adds a time adjunct to an existing event.
func AddTimeAdjunct(doc *Document, eventIdx int, timeEntity Entity) error {
	if eventIdx >= len(doc.Events) {
		return fmt.Errorf("event index %d out of range", eventIdx)
	}

	timeID := AddEntity(doc, timeEntity)
	AttachAdjunct(&doc.Events[eventIdx], RoleTime, timeID)

	// Map the time token if we have token mapping
	if doc.Trace.TokenToNode != nil {
		tokenIdx := len(doc.Trace.TokenToNode)
		doc.Trace.MapToken(tokenIdx, timeID)
	}

	return nil
}

// AddLocationAdjunct adds a location adjunct to an existing event.
func AddLocationAdjunct(doc *Document, eventIdx int, locationEntity Entity) error {
	if eventIdx >= len(doc.Events) {
		return fmt.Errorf("event index %d out of range", eventIdx)
	}

	locationID := AddEntity(doc, locationEntity)
	AttachAdjunct(&doc.Events[eventIdx], RoleLocation, locationID)

	// Map the location token if we have token mapping
	if doc.Trace.TokenToNode != nil {
		tokenIdx := len(doc.Trace.TokenToNode)
		doc.Trace.MapToken(tokenIdx, locationID)
	}

	return nil
}

// MapTokensToNodes maps a sequence of tokens to their corresponding node IDs.
// This is useful for maintaining alignment between surface forms and semantic nodes.
func MapTokensToNodes(doc *Document, tokens []string, nodeIDs []string) {
	if len(tokens) != len(nodeIDs) {
		doc.Trace.AddNote("warn", "token.mapping.mismatch",
			fmt.Sprintf("Token count (%d) doesn't match node count (%d)", len(tokens), len(nodeIDs)))
	}

	// Map tokens to nodes, using the shorter length to avoid index out of bounds
	mappingLength := len(tokens)
	if len(nodeIDs) < mappingLength {
		mappingLength = len(nodeIDs)
	}

	for i := 0; i < mappingLength; i++ {
		doc.Trace.MapToken(i, nodeIDs[i])
	}
}

// SetSourceLanguage sets the source language for the document trace.
func SetSourceLanguage(doc *Document, sourceLangID string) {
	doc.Trace.SetSourceLanguage(sourceLangID)
}

// AddConstructionNote adds a note about the construction used.
func AddConstructionNote(doc *Document, note string) {
	doc.Trace.AddNote("info", "construction.info", note)
}

// AddFeatureNote adds a note about a specific linguistic feature.
func AddFeatureNote(doc *Document, severity, feature, message string) {
	code := fmt.Sprintf("feature.%s.%s", severity, feature)
	doc.Trace.AddNote(severity, code, message)
}
