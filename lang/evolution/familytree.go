package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// LanguageFamilyTree manages the genealogical relationships between languages.
type LanguageFamilyTree struct {
	root   *LanguageNode
	nodes  map[string]*LanguageNode // LanguageID -> LanguageNode
	config EvolutionConfig
}

// NewLanguageFamilyTree creates a new language family tree.
func NewLanguageFamilyTree(config EvolutionConfig) *LanguageFamilyTree {
	return &LanguageFamilyTree{
		nodes:  make(map[string]*LanguageNode),
		config: config,
	}
}

// AddLanguage adds a language to the family tree.
func (lft *LanguageFamilyTree) AddLanguage(language *lang.Language, parent *lang.Language) error {
	if language == nil {
		return fmt.Errorf("cannot add nil language")
	}

	// Create the language node
	node := &LanguageNode{
		Language:       language,
		Evolution:      make([]EvolutionEvent, 0),
		ContactHistory: make([]ContactEvent, 0),
	}

	// If this is the first language, make it the root
	if lft.root == nil {
		lft.root = node
	} else if parent != nil {
		// Find the parent node
		parentNode, exists := lft.nodes[parent.ID.String()]
		if !exists {
			return fmt.Errorf("parent language %s not found in family tree", parent.ID.String())
		}

		// Set up the parent-child relationship
		node.Parent = parentNode
		parentNode.Children = append(parentNode.Children, node)
		node.DivergenceDate = time.Now()
	}

	// Add to the nodes map
	lft.nodes[language.ID.String()] = node

	return nil
}

// GetLanguageNode retrieves a language node by its ID.
func (lft *LanguageFamilyTree) GetLanguageNode(languageID string) (*LanguageNode, bool) {
	node, exists := lft.nodes[languageID]
	return node, exists
}

// GetAncestors returns all ancestors of a language up to the root.
func (lft *LanguageFamilyTree) GetAncestors(languageID string) []*LanguageNode {
	var ancestors []*LanguageNode

	node, exists := lft.nodes[languageID]
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

// GetDescendants returns all descendants of a language.
func (lft *LanguageFamilyTree) GetDescendants(languageID string) []*LanguageNode {
	var descendants []*LanguageNode

	node, exists := lft.nodes[languageID]
	if !exists {
		return descendants
	}

	// Recursively collect all descendants
	lft.collectDescendants(node, &descendants)

	return descendants
}

// collectDescendants recursively collects all descendants of a node.
func (lft *LanguageFamilyTree) collectDescendants(node *LanguageNode, descendants *[]*LanguageNode) {
	for _, child := range node.Children {
		*descendants = append(*descendants, child)
		lft.collectDescendants(child, descendants)
	}
}

// GetSiblings returns all languages that share the same parent.
func (lft *LanguageFamilyTree) GetSiblings(languageID string) []*LanguageNode {
	var siblings []*LanguageNode

	node, exists := lft.nodes[languageID]
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

// GetCommonAncestor finds the most recent common ancestor of two languages.
func (lft *LanguageFamilyTree) GetCommonAncestor(lang1ID, lang2ID string) (*LanguageNode, error) {
	_, exists1 := lft.nodes[lang1ID]
	_, exists2 := lft.nodes[lang2ID]

	if !exists1 || !exists2 {
		return nil, fmt.Errorf("one or both languages not found in family tree")
	}

	// Get ancestors of both languages
	ancestors1 := lft.GetAncestors(lang1ID)
	ancestors2 := lft.GetAncestors(lang2ID)

	// Find the most recent common ancestor
	for _, ancestor1 := range ancestors1 {
		for _, ancestor2 := range ancestors2 {
			if ancestor1 == ancestor2 {
				return ancestor1, nil
			}
		}
	}

	// If no common ancestor found, return root
	return lft.root, nil
}

// CalculateDivergenceTime estimates when two languages diverged from their common ancestor.
func (lft *LanguageFamilyTree) CalculateDivergenceTime(lang1ID, lang2ID string) (time.Time, error) {
	commonAncestor, err := lft.GetCommonAncestor(lang1ID, lang2ID)
	if err != nil {
		return time.Time{}, err
	}

	// For now, use a simple heuristic based on the number of generations
	// In practice, this would use more sophisticated linguistic dating methods

	generations1 := lft.countGenerations(commonAncestor, lft.nodes[lang1ID])
	generations2 := lft.countGenerations(commonAncestor, lft.nodes[lang2ID])

	// Assume each generation represents roughly 500 years of linguistic change
	divergenceYears := (generations1 + generations2) * 500

	return time.Now().AddDate(-int(divergenceYears), 0, 0), nil
}

// countGenerations counts the number of generations between two nodes.
func (lft *LanguageFamilyTree) countGenerations(from, to *LanguageNode) int {
	if from == to {
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

// AddEvolutionEvent adds an evolution event to a language's history.
func (lft *LanguageFamilyTree) AddEvolutionEvent(languageID string, event EvolutionEvent) error {
	node, exists := lft.nodes[languageID]
	if !exists {
		return fmt.Errorf("language %s not found in family tree", languageID)
	}

	node.Evolution = append(node.Evolution, event)
	return nil
}

// AddContactEvent adds a contact event to a language's history.
func (lft *LanguageFamilyTree) AddContactEvent(languageID string, event ContactEvent) error {
	node, exists := lft.nodes[languageID]
	if !exists {
		return fmt.Errorf("language %s not found in family tree", languageID)
	}

	node.ContactHistory = append(node.ContactHistory, event)
	return nil
}

// GetEvolutionHistory returns the complete evolution history of a language.
func (lft *LanguageFamilyTree) GetEvolutionHistory(languageID string) ([]EvolutionEvent, error) {
	node, exists := lft.nodes[languageID]
	if !exists {
		return nil, fmt.Errorf("language %s not found in family tree", languageID)
	}

	return node.Evolution, nil
}

// GetContactHistory returns the contact history of a language.
func (lft *LanguageFamilyTree) GetContactHistory(languageID string) ([]ContactEvent, error) {
	node, exists := lft.nodes[languageID]
	if !exists {
		return nil, fmt.Errorf("language %s not found in family tree", languageID)
	}

	return node.ContactHistory, nil
}

// PrintFamilyTree prints a text representation of the family tree.
func (lft *LanguageFamilyTree) PrintFamilyTree() string {
	if lft.root == nil {
		return "Empty family tree"
	}

	return lft.printNode(lft.root, 0)
}

// printNode recursively prints a node and its children with proper indentation.
func (lft *LanguageFamilyTree) printNode(node *LanguageNode, depth int) string {
	if node == nil {
		return ""
	}

	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	result := fmt.Sprintf("%s%s (%s)\n", indent, node.Language.Name, node.Language.ID.String())

	for _, child := range node.Children {
		result += lft.printNode(child, depth+1)
	}

	return result
}
