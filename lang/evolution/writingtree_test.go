package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/orthography"
)

func TestNewWritingSystemFamilyTree(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	if tree == nil {
		t.Fatal("Writing system family tree should not be nil")
	}

	if tree.root != nil {
		t.Error("New tree should have nil root")
	}

	if len(tree.nodes) != 0 {
		t.Error("New tree should have no nodes")
	}
}

func TestWritingSystemFamilyTree_AddWritingSystem(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create test writing systems
	ws1 := createCustomTestWritingSystem("Proto-Writing", orthography.WritingStyleAlphabetic)
	ws2 := createCustomTestWritingSystem("Child-Writing", orthography.WritingStyleAlphabetic)

	// Test adding root writing system
	err := tree.AddWritingSystem(ws1, nil, "origin")
	if err != nil {
		t.Errorf("Failed to add root writing system: %v", err)
	}

	if tree.root == nil {
		t.Error("Root should be set after adding first writing system")
	}

	if tree.root.WritingSystem != ws1 {
		t.Error("Root should contain the first writing system")
	}

	// Test adding child writing system
	err = tree.AddWritingSystem(ws2, ws1, "natural_evolution")
	if err != nil {
		t.Errorf("Failed to add child writing system: %v", err)
	}

	// Verify parent-child relationship
	childNode, exists := tree.GetWritingSystemNode(ws2.Name)
	if !exists {
		t.Error("Child writing system should be in tree")
	}

	if childNode.Parent != tree.root {
		t.Error("Child should have correct parent")
	}

	if len(tree.root.Children) != 1 || tree.root.Children[0] != childNode {
		t.Error("Parent should have correct children")
	}

	// Test adding with non-existent parent
	ws3 := createCustomTestWritingSystem("Orphan-Writing", orthography.WritingStyleAlphabetic)
	nonExistentParent := createCustomTestWritingSystem("NonExistent", orthography.WritingStyleAlphabetic)
	err = tree.AddWritingSystem(ws3, nonExistentParent, "evolution")
	if err == nil {
		t.Error("Should fail when adding with non-existent parent")
	}
}

func TestWritingSystemFamilyTree_GetWritingSystemNode(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Test getting from empty tree
	node, exists := tree.GetWritingSystemNode("nonexistent")
	if exists {
		t.Error("Should not find node in empty tree")
	}
	if node != nil {
		t.Error("Should return nil node when not found")
	}

	// Add a writing system and test retrieval
	ws := createCustomTestWritingSystem("Test-Writing", orthography.WritingStyleAlphabetic)
	tree.AddWritingSystem(ws, nil, "origin")

	node, exists = tree.GetWritingSystemNode(ws.Name)
	if !exists {
		t.Error("Should find added writing system")
	}
	if node == nil {
		t.Error("Should return non-nil node when found")
	}
	if node.WritingSystem != ws {
		t.Error("Returned node should contain correct writing system")
	}
}

func TestWritingSystemFamilyTree_GetAncestors(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create a simple tree: root -> child -> grandchild
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child := createCustomTestWritingSystem("Child", orthography.WritingStyleAlphabetic)
	grandchild := createCustomTestWritingSystem("Grandchild", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child, root, "evolution")
	tree.AddWritingSystem(grandchild, child, "evolution")

	// Test ancestors of grandchild
	ancestors := tree.GetAncestors(grandchild.Name)
	if len(ancestors) != 2 {
		t.Errorf("Expected 2 ancestors, got %d", len(ancestors))
	}

	// First ancestor should be child, second should be root
	if ancestors[0].WritingSystem != child {
		t.Error("First ancestor should be child")
	}
	if ancestors[1].WritingSystem != root {
		t.Error("Second ancestor should be root")
	}

	// Test ancestors of root
	rootAncestors := tree.GetAncestors(root.Name)
	if len(rootAncestors) != 0 {
		t.Errorf("Root should have no ancestors, got %d", len(rootAncestors))
	}
}

func TestWritingSystemFamilyTree_GetDescendants(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create a simple tree: root -> child1, child2 -> grandchild
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child1 := createCustomTestWritingSystem("Child1", orthography.WritingStyleAlphabetic)
	child2 := createCustomTestWritingSystem("Child2", orthography.WritingStyleAlphabetic)
	grandchild := createCustomTestWritingSystem("Grandchild", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child1, root, "evolution")
	tree.AddWritingSystem(child2, root, "evolution")
	tree.AddWritingSystem(grandchild, child2, "evolution")

	// Test descendants of root
	descendants := tree.GetDescendants(root.Name)
	if len(descendants) != 3 {
		t.Errorf("Expected 3 descendants, got %d", len(descendants))
	}

	// Test descendants of child1
	child1Descendants := tree.GetDescendants(child1.Name)
	if len(child1Descendants) != 0 {
		t.Errorf("Child1 should have no descendants, got %d", len(child1Descendants))
	}

	// Test descendants of child2
	child2Descendants := tree.GetDescendants(child2.Name)
	if len(child2Descendants) != 1 {
		t.Errorf("Child2 should have 1 descendant, got %d", len(child2Descendants))
	}
}

func TestWritingSystemFamilyTree_GetSiblings(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create a simple tree: root -> child1, child2
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child1 := createCustomTestWritingSystem("Child1", orthography.WritingStyleAlphabetic)
	child2 := createCustomTestWritingSystem("Child2", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child1, root, "evolution")
	tree.AddWritingSystem(child2, root, "evolution")

	// Test siblings of child1
	siblings := tree.GetSiblings(child1.Name)
	if len(siblings) != 1 {
		t.Errorf("Expected 1 sibling, got %d", len(siblings))
	}
	if siblings[0].WritingSystem != child2 {
		t.Error("Child1 should have Child2 as sibling")
	}

	// Test siblings of root
	rootSiblings := tree.GetSiblings(root.Name)
	if len(rootSiblings) != 0 {
		t.Errorf("Root should have no siblings, got %d", len(rootSiblings))
	}
}

func TestWritingSystemFamilyTree_GetCommonAncestor(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create a simple tree: root -> child1 -> grandchild1, child2 -> grandchild2
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child1 := createCustomTestWritingSystem("Child1", orthography.WritingStyleAlphabetic)
	child2 := createCustomTestWritingSystem("Child2", orthography.WritingStyleAlphabetic)
	grandchild1 := createCustomTestWritingSystem("Grandchild1", orthography.WritingStyleAlphabetic)
	grandchild2 := createCustomTestWritingSystem("Grandchild2", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child1, root, "evolution")
	tree.AddWritingSystem(child2, root, "evolution")
	tree.AddWritingSystem(grandchild1, child1, "evolution")
	tree.AddWritingSystem(grandchild2, child2, "evolution")

	// Test common ancestor of siblings
	commonAncestor, err := tree.GetCommonAncestor(child1.Name, child2.Name)
	if err != nil {
		t.Errorf("Failed to find common ancestor: %v", err)
	}
	if commonAncestor == nil {
		t.Error("Common ancestor should not be nil")
	}
	if commonAncestor.WritingSystem != root {
		t.Error("Common ancestor should be root")
	}

	// Test common ancestor of cousins
	commonAncestor, err = tree.GetCommonAncestor(grandchild1.Name, grandchild2.Name)
	if err != nil {
		t.Errorf("Failed to find common ancestor: %v", err)
	}
	if commonAncestor == nil {
		t.Error("Common ancestor should not be nil")
	}
	if commonAncestor.WritingSystem != root {
		t.Error("Common ancestor should be root")
	}

	// Test with non-existent writing system
	_, err = tree.GetCommonAncestor("nonexistent", child1.Name)
	if err == nil {
		t.Error("Should fail when one writing system doesn't exist")
	}
}

func TestWritingSystemFamilyTree_EstimateDivergenceTime(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create a simple tree: root -> child
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child := createCustomTestWritingSystem("Child", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child, root, "evolution")

	// Test divergence time estimation
	divergenceTime, err := tree.EstimateDivergenceTime(root.Name, child.Name)
	if err != nil {
		t.Errorf("Failed to estimate divergence time: %v", err)
	}

	// Should be in the past
	if divergenceTime.After(time.Now()) {
		t.Error("Divergence time should be in the past")
	}

	// Should be roughly 200 years ago (1 generation * 200 years)
	expectedTime := time.Now().AddDate(-200, 0, 0)
	tolerance := time.Hour * 24 * 365 * 50 // 50 years tolerance

	if divergenceTime.Before(expectedTime.Add(-tolerance)) || divergenceTime.After(expectedTime.Add(tolerance)) {
		t.Errorf("Divergence time %v should be roughly %v", divergenceTime, expectedTime)
	}
}

func TestWritingSystemFamilyTree_AddEvolutionEvent(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Add a writing system
	ws := createCustomTestWritingSystem("Test-Writing", orthography.WritingStyleAlphabetic)
	tree.AddWritingSystem(ws, nil, "origin")

	// Create an evolution event
	change := OrthographicChange{
		ID:          "test_change",
		Type:        OrthographicChangeTypeSimplification,
		Description: "Test evolution",
		Timestamp:   time.Now(),
	}

	// Add the evolution event
	err := tree.AddEvolutionEvent(ws.Name, change)
	if err != nil {
		t.Errorf("Failed to add evolution event: %v", err)
	}

	// Verify the event was added
	history, err := tree.GetEvolutionHistory(ws.Name)
	if err != nil {
		t.Errorf("Failed to get evolution history: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("Expected 1 evolution event, got %d", len(history))
	}

	if history[0].ID != change.ID {
		t.Error("Evolution event ID should match")
	}
}

func TestWritingSystemFamilyTree_AddBorrowingEvent(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Add a writing system
	ws := createCustomTestWritingSystem("Test-Writing", orthography.WritingStyleAlphabetic)
	tree.AddWritingSystem(ws, nil, "origin")

	// Create a borrowing event
	borrowingEvent := BorrowingEvent{
		ID:               "test_borrowing",
		Timestamp:        time.Now(),
		SourceID:         "source_ws",
		TargetID:         ws.Name,
		BorrowedFeatures: []string{"graphemes"},
		BorrowingType:    "graphemes",
		Intensity:        0.8,
	}

	// Add the borrowing event
	err := tree.AddBorrowingEvent(ws.Name, borrowingEvent)
	if err != nil {
		t.Errorf("Failed to add borrowing event: %v", err)
	}

	// Verify the event was added
	history, err := tree.GetBorrowingHistory(ws.Name)
	if err != nil {
		t.Errorf("Failed to get borrowing history: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("Expected 1 borrowing event, got %d", len(history))
	}

	if history[0].ID != borrowingEvent.ID {
		t.Error("Borrowing event ID should match")
	}
}

func TestWritingSystemFamilyTree_CalculateWritingSystemSimilarity(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Create a simple tree: root -> child1, child2
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child1 := createCustomTestWritingSystem("Child1", orthography.WritingStyleAlphabetic)
	child2 := createCustomTestWritingSystem("Child2", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child1, root, "evolution")
	tree.AddWritingSystem(child2, root, "evolution")

	// Test similarity between siblings
	similarity, err := tree.CalculateWritingSystemSimilarity(child1.Name, child2.Name)
	if err != nil {
		t.Errorf("Failed to calculate similarity: %v", err)
	}

	// Siblings should have high similarity (same parent, 1 generation apart)
	if similarity < 0.8 {
		t.Errorf("Expected high similarity between siblings, got %f", similarity)
	}

	// Test similarity between parent and child
	parentChildSimilarity, err := tree.CalculateWritingSystemSimilarity(root.Name, child1.Name)
	if err != nil {
		t.Errorf("Failed to calculate parent-child similarity: %v", err)
	}

	// Parent-child should have very high similarity
	if parentChildSimilarity < 0.9 {
		t.Errorf("Expected very high similarity between parent and child, got %f", parentChildSimilarity)
	}
}

func TestWritingSystemFamilyTree_PrintFamilyTree(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Test empty tree
	emptyTree := tree.PrintFamilyTree()
	if emptyTree != "Empty writing system family tree" {
		t.Errorf("Expected empty tree message, got: %s", emptyTree)
	}

	// Create a simple tree: root -> child
	root := createCustomTestWritingSystem("Root", orthography.WritingStyleAlphabetic)
	child := createCustomTestWritingSystem("Child", orthography.WritingStyleAlphabetic)

	tree.AddWritingSystem(root, nil, "origin")
	tree.AddWritingSystem(child, root, "evolution")

	// Test populated tree
	treeString := tree.PrintFamilyTree()
	if treeString == "" {
		t.Error("Tree string should not be empty")
	}

	// Should contain both writing system names
	if len(treeString) < 20 {
		t.Error("Tree string should be reasonably long")
	}
}

func TestWritingSystemFamilyTree_GenerateWritingSystemLineage(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)
	engine := NewOrthographicEvolutionEngine(config)

	// Create a parent writing system
	parent := createCustomTestWritingSystem("Parent", orthography.WritingStyleAlphabetic)

	// Add the parent to the tree first
	err := tree.AddWritingSystem(parent, nil, "origin")
	if err != nil {
		t.Fatalf("Failed to add parent to tree: %v", err)
	}

	// Generate a child lineage
	var child *orthography.WritingSystem
	child, err = tree.GenerateWritingSystemLineage(
		parent,
		"Child",
		orthography.WritingStyleSyllabic,
		engine,
		"test_era",
	)

	if err != nil {
		t.Errorf("Failed to generate writing system lineage: %v", err)
	}

	if child == nil {
		t.Fatal("Generated child should not be nil")
	}

	if child.Name != "Child" {
		t.Errorf("Expected child name 'Child', got '%s'", child.Name)
	}

	if child.Style != orthography.WritingStyleSyllabic {
		t.Errorf("Expected child style 'syllabic', got '%s'", child.Style.String())
	}

	if child.Culture != parent.Culture {
		t.Error("Child should inherit culture from parent")
	}

	// Verify the child was added to the family tree
	childNode, exists := tree.GetWritingSystemNode(child.Name)
	if !exists {
		t.Error("Generated child should be in family tree")
	}

	if childNode.Parent == nil {
		t.Error("Generated child should have parent")
	}
}

func TestWritingSystemFamilyTree_UtilityMethods(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	tree := NewWritingSystemFamilyTree(config)

	// Test initial state
	if tree.GetRoot() != nil {
		t.Error("New tree should have nil root")
	}

	if tree.GetWritingSystemCount() != 0 {
		t.Error("New tree should have 0 writing systems")
	}

	if len(tree.GetWritingSystemNames()) != 0 {
		t.Error("New tree should have no writing system names")
	}

	// Add a writing system and test updated state
	ws := createCustomTestWritingSystem("Test-Writing", orthography.WritingStyleAlphabetic)
	tree.AddWritingSystem(ws, nil, "origin")

	if tree.GetRoot() == nil {
		t.Error("Tree should have root after adding writing system")
	}

	if tree.GetWritingSystemCount() != 1 {
		t.Error("Tree should have 1 writing system")
	}

	names := tree.GetWritingSystemNames()
	if len(names) != 1 {
		t.Error("Tree should have 1 writing system name")
	}

	if names[0] != ws.Name {
		t.Error("Writing system name should match")
	}
}

// Helper function to create a test writing system with custom name and style
func createCustomTestWritingSystem(name string, style orthography.WritingStyle) *orthography.WritingSystem {
	// Create some test graphemes
	graphemes := []orthography.Grapheme{
		{Symbol: "a", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
		{Symbol: "b", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
		{Symbol: "c", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
	}

	// Create some test mappings without depending on phoneme pool
	mappings := []orthography.OrthographyMapping{
		{
			Phoneme:   nil, // We don't need actual phonemes for this test
			Graphemes: []orthography.Grapheme{graphemes[0]},
			Primary:   true,
		},
		{
			Phoneme:   nil,                                                // We don't need actual phonemes for this test
			Graphemes: []orthography.Grapheme{graphemes[1], graphemes[2]}, // Irregular mapping
			Primary:   true,
		},
	}

	return &orthography.WritingSystem{
		Style:     style,
		Graphemes: graphemes,
		Mappings:  mappings,
		Name:      name,
		Culture:   "test_culture",
	}
}
