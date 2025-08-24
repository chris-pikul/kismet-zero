package evolution

import (
	"strings"
	"testing"
)

// TestCreateTestLanguage tests the CreateTestLanguage function
func TestCreateTestLanguage(t *testing.T) {
	lang := CreateTestLanguage("TestLang", "TestCulture")

	if lang == nil {
		t.Fatal("Expected language to be created")
	}

	if lang.Name != "TestLang" {
		t.Errorf("Expected name 'TestLang', got '%s'", lang.Name)
	}

	if lang.Culture != "TestCulture" {
		t.Errorf("Expected culture 'TestCulture', got '%s'", lang.Culture)
	}

	if lang.Description != "Test language for evolution testing" {
		t.Errorf("Expected description 'Test language for evolution testing', got '%s'", lang.Description)
	}

	if lang.ID.Family != "test" {
		t.Errorf("Expected family 'test', got '%s'", lang.ID.Family)
	}

	if lang.ID.Branch != "test" {
		t.Errorf("Expected branch 'test', got '%s'", lang.ID.Branch)
	}

	if lang.ID.Language != "TestLang" {
		t.Errorf("Expected language 'TestLang', got '%s'", lang.ID.Language)
	}

	if lang.Phonology == nil {
		t.Error("Expected phonology to be set")
	}
}

// TestCreateTestLanguageForFamilyTree tests the CreateTestLanguageForFamilyTree function
func TestCreateTestLanguageForFamilyTree(t *testing.T) {
	lang := CreateTestLanguageForFamilyTree("FamilyLang", "FamilyCulture")

	if lang == nil {
		t.Fatal("Expected language to be created")
	}

	if lang.Name != "FamilyLang" {
		t.Errorf("Expected name 'FamilyLang', got '%s'", lang.Name)
	}

	if lang.Culture != "FamilyCulture" {
		t.Errorf("Expected culture 'FamilyCulture', got '%s'", lang.Culture)
	}

	// This should be the same as CreateTestLanguage
	expectedLang := CreateTestLanguage("FamilyLang", "FamilyCulture")
	if lang.Name != expectedLang.Name {
		t.Errorf("Expected same name as CreateTestLanguage, got '%s' vs '%s'", lang.Name, expectedLang.Name)
	}
}

// TestCreateTestWritingSystem tests the CreateTestWritingSystem function
func TestCreateTestWritingSystem(t *testing.T) {
	ws := CreateTestWritingSystem()

	if ws == nil {
		t.Fatal("Expected writing system to be created")
	}

	// The function currently returns a basic struct, so we just verify it's not nil
	// In the future, this could be expanded with more detailed testing
}

// TestContains tests the Contains utility function
func TestContains(t *testing.T) {
	// Test basic containment
	if !Contains("hello world", "world") {
		t.Error("Expected 'hello world' to contain 'world'")
	}

	if !Contains("hello world", "hello") {
		t.Error("Expected 'hello world' to contain 'hello'")
	}

	if !Contains("hello world", "o w") {
		t.Error("Expected 'hello world' to contain 'o w'")
	}

	// Test edge cases
	if Contains("hello world", "xyz") {
		t.Error("Expected 'hello world' not to contain 'xyz'")
	}

	if Contains("hello world", "") {
		t.Error("Expected 'hello world' not to contain empty string")
	}

	// Test exact matches
	if !Contains("hello", "hello") {
		t.Error("Expected 'hello' to contain 'hello'")
	}

	// Test partial matches at boundaries
	if !Contains("hello world", "hello") {
		t.Error("Expected 'hello world' to contain 'hello'")
	}

	if !Contains("hello world", "world") {
		t.Error("Expected 'hello world' to contain 'world'")
	}
}

// TestSplitLines tests the SplitLines utility function
func TestSplitLines(t *testing.T) {
	// Test basic line splitting
	input := "line1\nline2\nline3"
	lines := SplitLines(input)

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "line1" {
		t.Errorf("Expected first line 'line1', got '%s'", lines[0])
	}

	if lines[1] != "line2" {
		t.Errorf("Expected second line 'line2', got '%s'", lines[1])
	}

	if lines[2] != "line3" {
		t.Errorf("Expected third line 'line3', got '%s'", lines[2])
	}

	// Test with no newlines
	input = "single line"
	lines = SplitLines(input)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}

	if lines[0] != "single line" {
		t.Errorf("Expected 'single line', got '%s'", lines[0])
	}

	// Test with empty string
	input = ""
	lines = SplitLines(input)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line for empty string, got %d", len(lines))
	}

	if lines[0] != "" {
		t.Errorf("Expected empty line, got '%s'", lines[0])
	}

	// Test with trailing newline
	input = "line1\nline2\n"
	lines = SplitLines(input)

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if lines[0] != "line1" {
		t.Errorf("Expected first line 'line1', got '%s'", lines[0])
	}

	if lines[1] != "line2" {
		t.Errorf("Expected second line 'line2', got '%s'", lines[1])
	}
}

// TestStartsWith tests the StartsWith utility function
func TestStartsWith(t *testing.T) {
	// Test basic prefix matching
	if !StartsWith("hello world", "hello") {
		t.Error("Expected 'hello world' to start with 'hello'")
	}

	if !StartsWith("hello world", "h") {
		t.Error("Expected 'hello world' to start with 'h'")
	}

	if !StartsWith("hello world", "hello world") {
		t.Error("Expected 'hello world' to start with 'hello world'")
	}

	// Test edge cases
	if StartsWith("hello world", "world") {
		t.Error("Expected 'hello world' not to start with 'world'")
	}

	if StartsWith("hello world", "xyz") {
		t.Error("Expected 'hello world' not to start with 'xyz'")
	}

	// Empty string is technically a prefix of any string
	if !StartsWith("hello world", "") {
		t.Error("Expected 'hello world' to start with empty string (empty string is a valid prefix)")
	}

	// Test exact matches
	if !StartsWith("hello", "hello") {
		t.Error("Expected 'hello' to start with 'hello'")
	}

	// Test with short strings
	if !StartsWith("h", "h") {
		t.Error("Expected 'h' to start with 'h'")
	}

	if StartsWith("h", "hello") {
		t.Error("Expected 'h' not to start with 'hello'")
	}
}

// TestUtilityFunctionsEdgeCases tests edge cases for utility functions
func TestUtilityFunctionsEdgeCases(t *testing.T) {
	// Test Contains with very long strings
	longString := strings.Repeat("a", 1000) + "b" + strings.Repeat("c", 1000)
	if !Contains(longString, "b") {
		t.Error("Expected long string to contain 'b'")
	}

	// Test SplitLines with mixed line endings
	mixedInput := "line1\r\nline2\nline3\r"
	lines := SplitLines(mixedInput)
	// Note: This only handles \n, not \r or \r\n
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Test StartsWith with unicode
	unicodeString := "café"
	if !StartsWith(unicodeString, "caf") {
		t.Error("Expected 'café' to start with 'caf'")
	}
}

// TestUtilityFunctionsPerformance tests basic performance characteristics
func TestUtilityFunctionsPerformance(t *testing.T) {
	// Test that functions don't hang on large inputs
	largeString := strings.Repeat("a", 10000)

	// Test Contains performance
	start := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Contains(largeString, "b")
		}
	})
	t.Logf("Contains performance: %s", start)

	// Test SplitLines performance
	start = testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SplitLines(largeString)
		}
	})
	t.Logf("SplitLines performance: %s", start)

	// Test StartsWith performance
	start = testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StartsWith(largeString, "a")
		}
	})
	t.Logf("StartsWith performance: %s", start)
}
