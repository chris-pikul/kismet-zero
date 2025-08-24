package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

func TestNewEvolutionEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	if engine == nil {
		t.Fatal("Evolution engine should not be nil")
	}

	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}

	if engine.familyTree == nil {
		t.Error("Family tree should be initialized")
	}
}

func TestDefaultEvolutionConfig(t *testing.T) {
	config := DefaultEvolutionConfig(123)

	if config.Seed != 123 {
		t.Errorf("Expected seed 123, got %d", config.Seed)
	}

	if config.NaturalChangeRate <= 0 || config.NaturalChangeRate > 1 {
		t.Errorf("NaturalChangeRate should be between 0 and 1, got %f", config.NaturalChangeRate)
	}

	if config.SoundShiftProbability <= 0 || config.SoundShiftProbability > 1 {
		t.Errorf("SoundShiftProbability should be between 0 and 1, got %f", config.SoundShiftProbability)
	}
}

func TestSoundChangeEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewSoundChangeEngine(config)

	if engine == nil {
		t.Fatal("Sound change engine should not be nil")
	}

	if len(engine.rules) == 0 {
		t.Error("Should have default sound change rules")
	}

	// Test adding a custom rule
	customRule := SoundChangeRule{
		ID:          "test_rule",
		Name:        "Test Rule",
		Description: "A test sound change rule",
		Probability: 0.5,
	}

	engine.AddRule(customRule)

	if len(engine.rules) != len(CommonSoundChangeRules)+1 {
		t.Error("Custom rule should be added")
	}
}

func TestMorphologicalEvolutionEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewMorphologicalEvolutionEngine(config)

	if engine == nil {
		t.Fatal("Morphological evolution engine should not be nil")
	}

	// Test with nil components (should not panic)
	changes := engine.ApplyMorphologicalChanges(nil, nil, "test_era")

	if changes == nil {
		t.Error("Should return empty slice, not nil")
	}
}

func TestContactEvolutionEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewContactEvolutionEngine(config)

	if engine == nil {
		t.Fatal("Contact evolution engine should not be nil")
	}

	// Test contact influence calculation
	lang1 := createTestLanguage("lang1", "culture1")
	lang2 := createTestLanguage("lang2", "culture2")

	influence := engine.CalculateContactInfluence(
		lang1,
		lang2,
		ContactTypeTrade,
		time.Hour*24*365*100,
	)

	if influence < 0 || influence > 1 {
		t.Errorf("Influence should be between 0 and 1, got %f", influence)
	}
}

func TestLanguageFamilyTree(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewLanguageFamilyTree(config)

	if tree == nil {
		t.Fatal("Family tree should not be nil")
	}

	// Test adding languages
	lang1 := createTestLanguage("lang1", "culture1")
	lang2 := createTestLanguage("lang2", "culture2")

	err := tree.AddLanguage(lang1, nil)
	if err != nil {
		t.Errorf("Failed to add root language: %v", err)
	}

	err = tree.AddLanguage(lang2, lang1)
	if err != nil {
		t.Errorf("Failed to add child language: %v", err)
	}

	// Test retrieving nodes
	node1, exists := tree.GetLanguageNode(lang1.ID.String())
	if !exists {
		t.Error("Should find language 1")
	}

	node2, exists := tree.GetLanguageNode(lang2.ID.String())
	if !exists {
		t.Error("Should find language 2")
	}

	// Test parent-child relationship
	if node2.Parent != node1 {
		t.Error("Child should have correct parent")
	}

	if len(node1.Children) != 1 || node1.Children[0] != node2 {
		t.Error("Parent should have correct children")
	}
}

func TestEvolutionEngine_EvolveLanguage(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	lang1 := createTestLanguage("test_lang", "test_culture")

	evolved, event, err := engine.EvolveLanguage(lang1, "test_era", time.Hour*24*365*100)
	if err != nil {
		t.Errorf("Failed to evolve language: %v", err)
	}

	if evolved == nil {
		t.Fatal("Evolved language should not be nil")
	}

	if event.ID == "" {
		t.Error("Evolution event should have an ID")
	}

	// Note: Changes might be empty if random evolution doesn't trigger
	// This is acceptable behavior for the evolution system
	if event.Changes == nil {
		t.Error("Changes should not be nil")
	}

	// Test that evolution timestamp is set
	if evolved.EvolvedAt.IsZero() {
		t.Error("EvolvedAt should be set")
	}
}

func TestEvolutionEngine_CreateChildLanguage(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	parentLang := createTestLanguage("parent", "culture1")

	// Add parent to family tree first
	err := engine.familyTree.AddLanguage(parentLang, nil)
	if err != nil {
		t.Fatalf("Failed to add parent to family tree: %v", err)
	}

	childID := lang.LanguageID{
		Family:   "test",
		Branch:   "child",
		Language: "child_lang",
	}

	child, event, err := engine.CreateChildLanguage(parentLang, childID, "Child Language", "test_era")
	if err != nil {
		t.Errorf("Failed to create child language: %v", err)
	}

	if child == nil {
		t.Fatal("Child language should not be nil")
	}

	if child.ParentID == nil || *child.ParentID != parentLang.ID {
		t.Error("Child should have correct parent ID")
	}

	if event.ID == "" {
		t.Error("Evolution event should have an ID")
	}

	// Test family tree relationship
	familyTree := engine.GetFamilyTree()
	childNode, exists := familyTree.GetLanguageNode(child.ID.String())
	if !exists {
		t.Error("Child should be in family tree")
	}

	if childNode.Parent == nil {
		t.Error("Child node should have parent")
	}
}

func TestEvolutionEngine_ContactEvolution(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	targetLang := createTestLanguage("target", "culture1")
	sourceLang := createTestLanguage("source", "culture2")

	evolved, event, err := engine.EvolveLanguageFromContact(
		targetLang,
		sourceLang,
		ContactTypeTrade,
		0.6,
		time.Hour*24*365*100,
	)

	if err != nil {
		t.Errorf("Failed to evolve through contact: %v", err)
	}

	if evolved == nil {
		t.Fatal("Contact-evolved language should not be nil")
	}

	if event.TriggerType != "contact" {
		t.Errorf("Expected trigger type 'contact', got '%s'", event.TriggerType)
	}

	if event.TriggerCulture != sourceLang.Culture {
		t.Errorf("Expected trigger culture '%s', got '%s'", sourceLang.Culture, event.TriggerCulture)
	}
}

func TestChangeTypes(t *testing.T) {
	// Test ChangeType String method
	changeType := ChangeTypeSoundShift
	if changeType.String() != "sound_shift" {
		t.Errorf("Expected 'sound_shift', got '%s'", changeType.String())
	}

	// Test ChangeDirection String method
	changeDir := ChangeDirectionAdditive
	if changeDir.String() != "additive" {
		t.Errorf("Expected 'additive', got '%s'", changeDir.String())
	}

	// Test MorphologicalChangeType String method
	morphChangeType := MorphologicalChangeTypeSimplification
	if morphChangeType.String() != "simplification" {
		t.Errorf("Expected 'simplification', got '%s'", morphChangeType.String())
	}
}

func TestContactTypes(t *testing.T) {
	// Test ContactType String method
	contactType := ContactTypeTrade
	if contactType.String() != "trade" {
		t.Errorf("Expected 'trade', got '%s'", contactType.String())
	}

	// Test BorrowingType String method
	borrowingType := BorrowingTypeLexical
	if borrowingType.String() != "lexical" {
		t.Errorf("Expected 'lexical', got '%s'", borrowingType.String())
	}
}

func TestFamilyTreeTraversal(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewLanguageFamilyTree(config)

	// Create a simple tree: root -> child1 -> grandchild
	root := createTestLanguage("root", "culture1")
	child1 := createTestLanguage("child1", "culture1")
	grandchild := createTestLanguage("grandchild", "culture1")

	// Add languages
	tree.AddLanguage(root, nil)
	tree.AddLanguage(child1, root)
	tree.AddLanguage(grandchild, child1)

	// Test ancestors
	ancestors := tree.GetAncestors(grandchild.ID.String())
	if len(ancestors) != 2 {
		t.Errorf("Expected 2 ancestors, got %d", len(ancestors))
	}

	// Test descendants
	descendants := tree.GetDescendants(root.ID.String())
	if len(descendants) != 2 {
		t.Errorf("Expected 2 descendants, got %d", len(descendants))
	}

	// Test siblings
	siblings := tree.GetSiblings(child1.ID.String())
	if len(siblings) != 0 {
		t.Errorf("Expected 0 siblings, got %d", len(siblings))
	}

	// Test common ancestor
	commonAncestor, err := tree.GetCommonAncestor(child1.ID.String(), grandchild.ID.String())
	if err != nil {
		t.Errorf("Failed to find common ancestor: %v", err)
	}

	if commonAncestor == nil {
		t.Error("Common ancestor should not be nil")
	}
}

// Helper function to create test languages
func createTestLanguage(name, culture string) *lang.Language {
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
