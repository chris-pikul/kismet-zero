package interlingua

import (
	"testing"
)

func TestTraceAddNote(t *testing.T) {
	trace := &Trace{}

	// Add notes with different severities
	trace.AddNote("info", "test.info", "Information message")
	trace.AddNote("warn", "test.warn", "Warning message")
	trace.AddNote("loss", "test.loss", "Loss message")
	trace.AddNote("error", "test.error", "Error message")

	if len(trace.Notes) != 4 {
		t.Errorf("Expected 4 notes, got %d", len(trace.Notes))
	}

	// Check first note
	if trace.Notes[0].Severity != "info" {
		t.Errorf("Expected severity 'info', got %s", trace.Notes[0].Severity)
	}
	if trace.Notes[0].Code != "test.info" {
		t.Errorf("Expected code 'test.info', got %s", trace.Notes[0].Code)
	}
	if trace.Notes[0].Message != "Information message" {
		t.Errorf("Expected message 'Information message', got %s", trace.Notes[0].Message)
	}

	// Check last note
	if trace.Notes[3].Severity != "error" {
		t.Errorf("Expected severity 'error', got %s", trace.Notes[3].Severity)
	}
}

func TestTraceMapToken(t *testing.T) {
	trace := &Trace{}

	// Map tokens to nodes
	trace.MapToken(0, "v1")
	trace.MapToken(1, "e1")
	trace.MapToken(2, "e2")

	if len(trace.TokenToNode) != 3 {
		t.Errorf("Expected 3 token mappings, got %d", len(trace.TokenToNode))
	}

	// Check mappings
	if trace.TokenToNode[0] != "v1" {
		t.Errorf("Expected token 0 to map to 'v1', got %s", trace.TokenToNode[0])
	}
	if trace.TokenToNode[1] != "e1" {
		t.Errorf("Expected token 1 to map to 'e1', got %s", trace.TokenToNode[1])
	}
	if trace.TokenToNode[2] != "e2" {
		t.Errorf("Expected token 2 to map to 'e2', got %s", trace.TokenToNode[2])
	}
}

func TestTraceGetNodeForToken(t *testing.T) {
	trace := &Trace{}

	// Test with empty trace
	nodeID, exists := trace.GetNodeForToken(0)
	if exists {
		t.Error("Expected no mapping for empty trace")
	}
	if nodeID != "" {
		t.Errorf("Expected empty node ID, got %s", nodeID)
	}

	// Add mapping and test
	trace.MapToken(0, "v1")
	nodeID, exists = trace.GetNodeForToken(0)
	if !exists {
		t.Error("Expected mapping to exist")
	}
	if nodeID != "v1" {
		t.Errorf("Expected node ID 'v1', got %s", nodeID)
	}

	// Test non-existent token
	nodeID, exists = trace.GetNodeForToken(1)
	if exists {
		t.Error("Expected no mapping for non-existent token")
	}
	if nodeID != "" {
		t.Errorf("Expected empty node ID, got %s", nodeID)
	}
}

func TestTraceSetSourceLanguage(t *testing.T) {
	trace := &Trace{}

	trace.SetSourceLanguage("en")
	if trace.SourceLanguageID != "en" {
		t.Errorf("Expected source language 'en', got %s", trace.SourceLanguageID)
	}

	trace.SetSourceLanguage("fr")
	if trace.SourceLanguageID != "fr" {
		t.Errorf("Expected source language 'fr', got %s", trace.SourceLanguageID)
	}
}

func TestTraceSetConstructionID(t *testing.T) {
	trace := &Trace{}

	trace.SetConstructionID("decl.transitive.svo.pos.pst")
	if trace.ConstructionID != "decl.transitive.svo.pos.pst" {
		t.Errorf("Expected construction ID 'decl.transitive.svo.pos.pst', got %s", trace.ConstructionID)
	}

	trace.SetConstructionID("decl.intransitive.sv.neg.pres")
	if trace.ConstructionID != "decl.intransitive.sv.neg.pres" {
		t.Errorf("Expected construction ID 'decl.intransitive.sv.neg.pres', got %s", trace.ConstructionID)
	}
}

func TestTraceHasNotes(t *testing.T) {
	trace := &Trace{}

	// Initially no notes
	if trace.HasNotes() {
		t.Error("Expected no notes initially")
	}

	// Add a note
	trace.AddNote("info", "test", "test message")
	if !trace.HasNotes() {
		t.Error("Expected notes after adding one")
	}
}

func TestTraceGetNotesBySeverity(t *testing.T) {
	trace := &Trace{}

	// Add notes with different severities
	trace.AddNote("info", "info1", "Info 1")
	trace.AddNote("info", "info2", "Info 2")
	trace.AddNote("warn", "warn1", "Warning 1")
	trace.AddNote("error", "error1", "Error 1")

	// Get info notes
	infoNotes := trace.GetNotesBySeverity("info")
	if len(infoNotes) != 2 {
		t.Errorf("Expected 2 info notes, got %d", len(infoNotes))
	}

	// Get warn notes
	warnNotes := trace.GetNotesBySeverity("warn")
	if len(warnNotes) != 1 {
		t.Errorf("Expected 1 warn note, got %d", len(warnNotes))
	}

	// Get non-existent severity
	nonexistentNotes := trace.GetNotesBySeverity("nonexistent")
	if len(nonexistentNotes) != 0 {
		t.Errorf("Expected 0 notes for non-existent severity, got %d", len(nonexistentNotes))
	}
}

func TestTraceGetNotesByCode(t *testing.T) {
	trace := &Trace{}

	// Add notes with different codes
	trace.AddNote("info", "test.code1", "Message 1")
	trace.AddNote("warn", "test.code2", "Message 2")
	trace.AddNote("info", "test.code1", "Message 3") // Same code, different severity

	// Get notes by code
	code1Notes := trace.GetNotesByCode("test.code1")
	if len(code1Notes) != 2 {
		t.Errorf("Expected 2 notes for code 'test.code1', got %d", len(code1Notes))
	}

	code2Notes := trace.GetNotesByCode("test.code2")
	if len(code2Notes) != 1 {
		t.Errorf("Expected 1 note for code 'test.code2', got %d", len(code2Notes))
	}

	// Get non-existent code
	nonexistentNotes := trace.GetNotesByCode("nonexistent")
	if len(nonexistentNotes) != 0 {
		t.Errorf("Expected 0 notes for non-existent code, got %d", len(nonexistentNotes))
	}
}

func TestTraceInitialization(t *testing.T) {
	trace := &Trace{}

	// Test that maps are properly initialized when first used
	trace.MapToken(0, "v1")
	if trace.TokenToNode == nil {
		t.Error("Expected TokenToNode map to be initialized")
	}

	// Test that notes slice is properly initialized
	trace.AddNote("info", "test", "test message")
	if trace.Notes == nil {
		t.Error("Expected Notes slice to be initialized")
	}
}
