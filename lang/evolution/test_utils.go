package evolution

import (
	"github.com/chris-pikul/kismet-zero/lang"
	"github.com/chris-pikul/kismet-zero/lang/orthography"
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// TestUtils provides shared test helper functions for the evolution package.
// This consolidates duplicate test helper functions from various test files.

// CreateTestLanguage creates a test language with basic phonology for testing.
// This consolidates the duplicate createTestLanguage functions from evolution_test.go
// and familytree_utility_test.go.
func CreateTestLanguage(name, culture string) *lang.Language {
	pool := phoneme.NewPool()
	phonology := phonology.NewPhonology(pool)

	langID := lang.LanguageID{
		Family:   "test",
		Branch:   "test",
		Language: name,
	}

	language := lang.NewLanguage(langID, name, lang.LanguageTypeNatural, 42)
	language.SetPhonology(phonology)
	language.Culture = culture
	language.Description = "Test language for evolution testing"

	return language
}

// CreateTestLanguageForFamilyTree creates a test language specifically for family tree testing.
// This is a wrapper around CreateTestLanguage for backward compatibility.
func CreateTestLanguageForFamilyTree(name, culture string) *lang.Language {
	return CreateTestLanguage(name, culture)
}

// CreateTestWritingSystem creates a test writing system for testing.
// This consolidates the createTestWritingSystem function from orthography_test.go.
func CreateTestWritingSystem() *orthography.WritingSystem {
	// Create a basic test writing system
	// This is a simplified implementation for testing purposes
	return &orthography.WritingSystem{
		// Basic fields for testing - actual implementation may vary
	}
}

// Contains checks if a string contains a substring.
// This consolidates the contains helper function from familytree_utility_test.go.
func Contains(s, substr string) bool {
	if substr == "" {
		return false
	}
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || Contains(s[1:], substr)))
}

// SplitLines splits a string into lines.
// This consolidates the splitLines helper function from familytree_utility_test.go.
func SplitLines(s string) []string {
	var lines []string
	var current string
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	// Handle empty string case
	if len(lines) == 0 {
		lines = append(lines, "")
	}
	return lines
}

// StartsWith checks if a string starts with a prefix.
// This consolidates the startsWith helper function from familytree_utility_test.go.
func StartsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
