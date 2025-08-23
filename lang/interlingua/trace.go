package interlingua

// AddNote adds a diagnostic note to the trace.
func (t *Trace) AddNote(severity, code, message string) {
	t.Notes = append(t.Notes, Note{
		Severity: severity,
		Code:     code,
		Message:  message,
	})
}

// MapToken maps a token index to a node ID for alignment tracking.
func (t *Trace) MapToken(tokenIdx int, nodeID string) {
	if t.TokenToNode == nil {
		t.TokenToNode = make(map[int]string)
	}
	t.TokenToNode[tokenIdx] = nodeID
}

// GetNodeForToken retrieves the node ID mapped to a token index.
func (t *Trace) GetNodeForToken(tokenIdx int) (string, bool) {
	if t.TokenToNode == nil {
		return "", false
	}
	nodeID, exists := t.TokenToNode[tokenIdx]
	return nodeID, exists
}

// SetSourceLanguage sets the source language ID for the trace.
func (t *Trace) SetSourceLanguage(langID string) {
	t.SourceLanguageID = langID
}

// SetConstructionID sets the construction identifier for the trace.
func (t *Trace) SetConstructionID(constructionID string) {
	t.ConstructionID = constructionID
}

// HasNotes returns true if the trace contains any notes.
func (t *Trace) HasNotes() bool {
	return len(t.Notes) > 0
}

// GetNotesBySeverity returns all notes with the specified severity level.
func (t *Trace) GetNotesBySeverity(severity string) []Note {
	var filtered []Note
	for _, note := range t.Notes {
		if note.Severity == severity {
			filtered = append(filtered, note)
		}
	}
	return filtered
}

// GetNotesByCode returns all notes with the specified code.
func (t *Trace) GetNotesByCode(code string) []Note {
	var filtered []Note
	for _, note := range t.Notes {
		if note.Code == code {
			filtered = append(filtered, note)
		}
	}
	return filtered
}
