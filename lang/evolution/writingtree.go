package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/orthography"
)

// WritingSystemNode represents a node in the writing system family tree.
type WritingSystemNode struct {
	WritingSystem *orthography.WritingSystem `json:"writingSystem"`
	Parent        *WritingSystemNode         `json:"parent,omitempty"`
	Children      []*WritingSystemNode       `json:"children,omitempty"`

	// Evolution history
	EvolutionHistory []OrthographicChange `json:"evolutionHistory"`
	BorrowingHistory []BorrowingEvent     `json:"borrowingHistory"`

	// Relationship metadata
	DivergenceDate time.Time `json:"divergenceDate,omitempty"`
	DivergenceType string    `json:"divergenceType,omitempty"` // "natural", "reform", "borrowing", "contact"

	// Writing system genealogy
	AncestralFeatures []string `json:"ancestralFeatures,omitempty"` // Features inherited from ancestors
	Innovations       []string `json:"innovations,omitempty"`       // New features developed
	LostFeatures      []string `json:"lostFeatures,omitempty"`      // Features lost over time
}

// BorrowingEvent represents borrowing of writing system features from another system.
type BorrowingEvent struct {
	ID               string    `json:"id"`
	Timestamp        time.Time `json:"timestamp"`
	SourceID         string    `json:"sourceID"`         // ID of source writing system
	TargetID         string    `json:"targetID"`         // ID of target writing system
	BorrowedFeatures []string  `json:"borrowedFeatures"` // What was borrowed
	BorrowingType    string    `json:"borrowingType"`    // "graphemes", "mappings", "style"
	Intensity        float32   `json:"intensity"`        // 0.0 to 1.0, strength of borrowing
	Description      string    `json:"description,omitempty"`
}

// WritingSystemFamilyTree manages genealogical relationships between writing systems.
type WritingSystemFamilyTree struct {
	root   *WritingSystemNode
	nodes  map[string]*WritingSystemNode // WritingSystem.Name -> WritingSystemNode
	config EvolutionConfig
}

// NewWritingSystemFamilyTree creates a new writing system family tree.
func NewWritingSystemFamilyTree(config EvolutionConfig) *WritingSystemFamilyTree {
	return &WritingSystemFamilyTree{
		nodes:  make(map[string]*WritingSystemNode),
		config: config,
	}
}

// AddWritingSystem adds a writing system to the family tree.
// If parent is provided, establishes a parent-child relationship.
func (wsft *WritingSystemFamilyTree) AddWritingSystem(
	writingSystem *orthography.WritingSystem,
	parent *orthography.WritingSystem,
	divergenceType string,
) error {
	if writingSystem == nil {
		return fmt.Errorf("cannot add nil writing system")
	}

	// Create the writing system node
	node := &WritingSystemNode{
		WritingSystem:    writingSystem,
		EvolutionHistory: make([]OrthographicChange, 0),
		BorrowingHistory: make([]BorrowingEvent, 0),
		DivergenceType:   divergenceType,
	}

	// If this is the first writing system, make it the root
	if wsft.root == nil {
		wsft.root = node
		node.DivergenceDate = time.Now()
	} else if parent != nil {
		// Find parent node
		parentNode, exists := wsft.nodes[parent.Name]
		if !exists {
			return fmt.Errorf("parent writing system %s not found in family tree", parent.Name)
		}

		// Set up the parent-child relationship
		node.Parent = parentNode
		parentNode.Children = append(parentNode.Children, node)
		node.DivergenceDate = time.Now()
	}

	// Add to the nodes map
	wsft.nodes[writingSystem.Name] = node

	return nil
}

// GetWritingSystemNode retrieves a writing system node by name.
func (wsft *WritingSystemFamilyTree) GetWritingSystemNode(name string) (*WritingSystemNode, bool) {
	node, exists := wsft.nodes[name]
	return node, exists
}

// GetAncestors returns all ancestors of a writing system.
func (wsft *WritingSystemFamilyTree) GetAncestors(name string) []*WritingSystemNode {
	var ancestors []*WritingSystemNode

	node, exists := wsft.nodes[name]
	if !exists {
		return ancestors
	}

	current := node.Parent
	for current != nil {
		ancestors = append(ancestors, current)
		current = current.Parent
	}

	return ancestors
}

// GetDescendants returns all descendants of a writing system.
func (wsft *WritingSystemFamilyTree) GetDescendants(name string) []*WritingSystemNode {
	var descendants []*WritingSystemNode

	node, exists := wsft.nodes[name]
	if !exists {
		return descendants
	}

	// Recursively collect all descendants
	wsft.collectDescendants(node, &descendants)

	return descendants
}

// collectDescendants recursively collects all descendants of a node.
func (wsft *WritingSystemFamilyTree) collectDescendants(node *WritingSystemNode, descendants *[]*WritingSystemNode) {
	if node == nil {
		return
	}

	for _, child := range node.Children {
		*descendants = append(*descendants, child)
		wsft.collectDescendants(child, descendants)
	}
}

// GetSiblings returns all siblings of a writing system.
func (wsft *WritingSystemFamilyTree) GetSiblings(name string) []*WritingSystemNode {
	var siblings []*WritingSystemNode

	node, exists := wsft.nodes[name]
	if !exists || node.Parent == nil {
		return siblings
	}

	// Add all children of the parent except the current node
	for _, child := range node.Parent.Children {
		if child != node {
			siblings = append(siblings, child)
		}
	}

	return siblings
}

// GetCommonAncestor finds the most recent common ancestor of two writing systems.
func (wsft *WritingSystemFamilyTree) GetCommonAncestor(name1, name2 string) (*WritingSystemNode, error) {
	_, exists1 := wsft.nodes[name1]
	_, exists2 := wsft.nodes[name2]

	if !exists1 || !exists2 {
		return nil, fmt.Errorf("one or both writing systems not found in family tree")
	}

	// Get ancestors of both writing systems
	ancestors1 := wsft.GetAncestors(name1)
	ancestors2 := wsft.GetAncestors(name2)

	// Find the most recent common ancestor
	for _, ancestor1 := range ancestors1 {
		for _, ancestor2 := range ancestors2 {
			if ancestor1 == ancestor2 {
				return ancestor1, nil
			}
		}
	}

	// If no common ancestor found, return root
	return wsft.root, nil
}

// EstimateDivergenceTime estimates when two writing systems diverged.
func (wsft *WritingSystemFamilyTree) EstimateDivergenceTime(name1, name2 string) (time.Time, error) {
	commonAncestor, err := wsft.GetCommonAncestor(name1, name2)
	if err != nil {
		return time.Time{}, err
	}

	// For now, use a simple heuristic based on the number of generations
	// In practice, this would use more sophisticated dating methods
	generations1 := wsft.countGenerations(commonAncestor, wsft.nodes[name1])
	generations2 := wsft.countGenerations(commonAncestor, wsft.nodes[name2])

	// Assume each generation represents roughly 200 years of writing system change
	// (Writing systems change more slowly than spoken languages)
	divergenceYears := (generations1 + generations2) * 200

	return time.Now().AddDate(-int(divergenceYears), 0, 0), nil
}

// countGenerations counts the number of generations between two nodes.
func (wsft *WritingSystemFamilyTree) countGenerations(from, to *WritingSystemNode) int {
	if from == nil || to == nil {
		return 0
	}

	count := 0
	current := to

	for current != nil && current != from {
		count++
		current = current.Parent
	}

	return count
}

// AddEvolutionEvent adds an evolution event to a writing system's history.
func (wsft *WritingSystemFamilyTree) AddEvolutionEvent(name string, change OrthographicChange) error {
	node, exists := wsft.nodes[name]
	if !exists {
		return fmt.Errorf("writing system %s not found in family tree", name)
	}

	node.EvolutionHistory = append(node.EvolutionHistory, change)
	return nil
}

// AddBorrowingEvent adds a borrowing event to a writing system's history.
func (wsft *WritingSystemFamilyTree) AddBorrowingEvent(name string, event BorrowingEvent) error {
	node, exists := wsft.nodes[name]
	if !exists {
		return fmt.Errorf("writing system %s not found in family tree", name)
	}

	node.BorrowingHistory = append(node.BorrowingHistory, event)
	return nil
}

// GetEvolutionHistory returns the evolution history of a writing system.
func (wsft *WritingSystemFamilyTree) GetEvolutionHistory(name string) ([]OrthographicChange, error) {
	node, exists := wsft.nodes[name]
	if !exists {
		return nil, fmt.Errorf("writing system %s not found in family tree", name)
	}

	return node.EvolutionHistory, nil
}

// GetBorrowingHistory returns the borrowing history of a writing system.
func (wsft *WritingSystemFamilyTree) GetBorrowingHistory(name string) ([]BorrowingEvent, error) {
	node, exists := wsft.nodes[name]
	if !exists {
		return nil, fmt.Errorf("writing system %s not found in family tree", name)
	}

	return node.BorrowingHistory, nil
}

// PrintFamilyTree returns a string representation of the writing system family tree.
func (wsft *WritingSystemFamilyTree) PrintFamilyTree() string {
	if wsft.root == nil {
		return "Empty writing system family tree"
	}

	return wsft.printNode(wsft.root, 0)
}

// printNode recursively prints a node and its children.
func (wsft *WritingSystemFamilyTree) printNode(node *WritingSystemNode, depth int) string {
	if node == nil {
		return ""
	}

	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	result := fmt.Sprintf("%s%s (%s)\n", indent, node.WritingSystem.Name, node.WritingSystem.Style.String())

	for _, child := range node.Children {
		result += wsft.printNode(child, depth+1)
	}

	return result
}

// CalculateWritingSystemSimilarity calculates how similar two writing systems are.
func (wsft *WritingSystemFamilyTree) CalculateWritingSystemSimilarity(name1, name2 string) (float32, error) {
	node1, exists1 := wsft.nodes[name1]
	node2, exists2 := wsft.nodes[name2]

	if !exists1 || !exists2 {
		return 0.0, fmt.Errorf("one or both writing systems not found in family tree")
	}

	// Base similarity from common ancestor
	commonAncestor, err := wsft.GetCommonAncestor(name1, name2)
	if err != nil {
		return 0.0, err
	}

	// Calculate generational distance
	generations1 := wsft.countGenerations(commonAncestor, node1)
	generations2 := wsft.countGenerations(commonAncestor, node2)
	totalGenerations := generations1 + generations2

	// Similarity decreases with generational distance
	// Base similarity starts at 1.0 and decreases by 0.1 per generation
	baseSimilarity := 1.0 - float64(totalGenerations)*0.1
	if baseSimilarity < 0.0 {
		baseSimilarity = 0.0
	}

	// Additional similarity from shared features
	sharedFeatures := wsft.countSharedFeatures(node1, node2)
	featureSimilarity := float64(sharedFeatures) * 0.05

	totalSimilarity := baseSimilarity + featureSimilarity
	if totalSimilarity > 1.0 {
		totalSimilarity = 1.0
	}

	return float32(totalSimilarity), nil
}

// countSharedFeatures counts the number of shared features between two writing systems.
func (wsft *WritingSystemFamilyTree) countSharedFeatures(node1, node2 *WritingSystemNode) int {
	shared := 0

	// Count shared grapheme types
	graphemeTypes1 := make(map[orthography.GraphemeType]bool)
	for _, g := range node1.WritingSystem.Graphemes {
		graphemeTypes1[g.Type] = true
	}

	for _, g := range node2.WritingSystem.Graphemes {
		if graphemeTypes1[g.Type] {
			shared++
		}
	}

	// Count shared writing style
	if node1.WritingSystem.Style == node2.WritingSystem.Style {
		shared += 2
	}

	return shared
}

// GenerateWritingSystemLineage creates a new writing system that evolved from a parent.
func (wsft *WritingSystemFamilyTree) GenerateWritingSystemLineage(
	parent *orthography.WritingSystem,
	childName string,
	childStyle orthography.WritingStyle,
	evolutionEngine *OrthographicEvolutionEngine,
	era string,
) (*orthography.WritingSystem, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent writing system cannot be nil")
	}

	// Create a new writing system based on the parent
	childWritingSystem := &orthography.WritingSystem{
		Style:     childStyle,
		Graphemes: make([]orthography.Grapheme, 0),
		Mappings:  make([]orthography.OrthographyMapping, 0),
		Name:      childName,
		Culture:   parent.Culture,
	}

	// Copy and modify graphemes from parent
	for _, parentGrapheme := range parent.Graphemes {
		// Apply some evolution to the grapheme
		evolvedGrapheme := wsft.evolveGrapheme(parentGrapheme, evolutionEngine, era)
		childWritingSystem.Graphemes = append(childWritingSystem.Graphemes, evolvedGrapheme)
	}

	// Apply orthographic changes to create divergence
	changes := evolutionEngine.ApplyOrthographicChanges(childWritingSystem, era)

	// Always add the child to the family tree
	err := wsft.AddWritingSystem(childWritingSystem, parent, "natural_evolution")
	if err != nil {
		return nil, fmt.Errorf("failed to add child writing system to family tree: %w", err)
	}

	// Record the evolution events if there are any
	if len(changes) > 0 {
		for _, change := range changes {
			err := wsft.AddEvolutionEvent(childName, change)
			if err != nil {
				return nil, fmt.Errorf("failed to record evolution event: %w", err)
			}
		}
	}

	return childWritingSystem, nil
}

// evolveGrapheme applies evolutionary changes to a grapheme.
func (wsft *WritingSystemFamilyTree) evolveGrapheme(
	grapheme orthography.Grapheme,
	engine *OrthographicEvolutionEngine,
	era string,
) orthography.Grapheme {
	// Simple evolution: randomly modify the grapheme
	evolved := grapheme

	// Use a simple deterministic evolution based on the symbol
	// In practice, this would use proper RNG from the engine
	symbolHash := 0
	for _, r := range evolved.Symbol {
		symbolHash += int(r)
	}

	// 10% chance to change the symbol slightly (deterministic based on hash)
	if (symbolHash % 10) == 0 {
		// Simple symbol evolution (in practice, this would be more sophisticated)
		switch evolved.Symbol {
		case "a":
			evolved.Symbol = "ɑ"
		case "e":
			evolved.Symbol = "ɛ"
		case "i":
			evolved.Symbol = "ɪ"
		case "o":
			evolved.Symbol = "ɔ"
		case "u":
			evolved.Symbol = "ʊ"
		}
	}

	// 5% chance to change the weight (deterministic based on hash)
	if (symbolHash % 20) == 0 {
		weightChange := float32((symbolHash%4)-2) * 0.1
		evolved.Weight += weightChange
		if evolved.Weight < 0.1 {
			evolved.Weight = 0.1
		}
		if evolved.Weight > 2.0 {
			evolved.Weight = 2.0
		}
	}

	return evolved
}

// GetRoot returns the root of the writing system family tree.
func (wsft *WritingSystemFamilyTree) GetRoot() *WritingSystemNode {
	return wsft.root
}

// GetWritingSystemCount returns the total number of writing systems in the tree.
func (wsft *WritingSystemFamilyTree) GetWritingSystemCount() int {
	return len(wsft.nodes)
}

// GetWritingSystemNames returns all writing system names in the tree.
func (wsft *WritingSystemFamilyTree) GetWritingSystemNames() []string {
	names := make([]string, 0, len(wsft.nodes))
	for name := range wsft.nodes {
		names = append(names, name)
	}
	return names
}
