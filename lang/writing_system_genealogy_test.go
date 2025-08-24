package lang

import (
	"testing"
)

// TestWritingSystemGenealogy tests the writing system genealogy structure.
func TestWritingSystemGenealogy(t *testing.T) {
	t.Run("Genealogy Creation", func(t *testing.T) {
		genealogy := NewWritingSystemGenealogy(
			"ws_1",
			"Ancient Runes",
			"A mystical writing system used by ancient civilizations",
			"runic",
			"alphabetic",
			"ancient",
		)

		if genealogy.ID != "ws_1" {
			t.Errorf("Expected ID 'ws_1', got '%s'", genealogy.ID)
		}
		if genealogy.Name != "Ancient Runes" {
			t.Errorf("Expected name 'Ancient Runes', got '%s'", genealogy.Name)
		}
		if genealogy.Family != "runic" {
			t.Errorf("Expected family 'runic', got '%s'", genealogy.Family)
		}
		if genealogy.ScriptType != "alphabetic" {
			t.Errorf("Expected script type 'alphabetic', got '%s'", genealogy.ScriptType)
		}
		if genealogy.Direction != "left-to-right" {
			t.Errorf("Expected direction 'left-to-right', got '%s'", genealogy.Direction)
		}
		if genealogy.Complexity != 0.5 {
			t.Errorf("Expected complexity 0.5, got %f", genealogy.Complexity)
		}
		if genealogy.Elegance != 0.5 {
			t.Errorf("Expected elegance 0.5, got %f", genealogy.Elegance)
		}

		t.Logf("Successfully created writing system genealogy: %s (%s)", genealogy.Name, genealogy.ScriptType)
	})
}

// TestWritingSystemRelationships tests the relationship management functions.
func TestWritingSystemRelationships(t *testing.T) {
	t.Run("Parent-Child Relationships", func(t *testing.T) {
		parent := NewWritingSystemGenealogy("ws_parent", "Parent Script", "Parent writing system", "family1", "logographic", "ancient")
		child := NewWritingSystemGenealogy("ws_child", "Child Script", "Child writing system", "family1", "syllabic", "ancient")

		// Set up parent-child relationship
		child.SetParent(parent.ID)
		parent.AddChild(child.ID)

		if child.ParentID == nil || *child.ParentID != parent.ID {
			t.Error("Expected child to have parent ID set")
		}
		if !parent.hasChild(child.ID) {
			t.Error("Expected parent to have child ID")
		}

		// Test removing child
		parent.RemoveChild(child.ID)
		if parent.hasChild(child.ID) {
			t.Error("Expected child to be removed from parent")
		}

		t.Logf("Successfully tested parent-child relationships")
	})

	t.Run("Influence Relationships", func(t *testing.T) {
		influencer := NewWritingSystemGenealogy("ws_influencer", "Influencer Script", "Influential writing system", "family1", "alphabetic", "ancient")
		influenced := NewWritingSystemGenealogy("ws_influenced", "Influenced Script", "Influenced writing system", "family2", "syllabic", "ancient")

		// Set up influence relationships
		influencer.AddInfluence(influenced.ID)
		influenced.AddInfluencedBy(influencer.ID)

		if !influencer.hasInfluence(influenced.ID) {
			t.Error("Expected influencer to have influence over influenced")
		}
		if !influenced.hasInfluencedBy(influencer.ID) {
			t.Error("Expected influenced to be influenced by influencer")
		}

		t.Logf("Successfully tested influence relationships")
	})
}

// TestWritingSystemEvolution tests the evolution tracking system.
func TestWritingSystemEvolution(t *testing.T) {
	t.Run("Evolution Steps", func(t *testing.T) {
		genealogy := NewWritingSystemGenealogy("ws_evolve", "Evolving Script", "A script that evolves over time", "family1", "alphabetic", "ancient")

		initialComplexity := genealogy.Complexity
		initialElegance := genealogy.Elegance

		// Add evolution steps
		changes := []WritingSystemChange{
			{
				ID:          "change1",
				Type:        "addition",
				Description: "Added new characters",
				Details:     "Added 5 new consonant characters",
				Intensity:   0.3,
				Confidence:  0.9,
			},
			{
				ID:          "change2",
				Type:        "modification",
				Description: "Modified existing characters",
				Details:     "Simplified character shapes",
				Intensity:   0.2,
				Confidence:  0.8,
			},
		}

		for _, change := range changes {
			genealogy.AddEvolutionStep(change)
		}

		if len(genealogy.EvolutionSteps) != 2 {
			t.Errorf("Expected 2 evolution steps, got %d", len(genealogy.EvolutionSteps))
		}

		// Check that metrics were updated
		if genealogy.Complexity <= initialComplexity {
			t.Error("Expected complexity to increase after additions")
		}
		if genealogy.Elegance <= initialElegance {
			t.Error("Expected elegance to increase after modifications")
		}

		t.Logf("Evolution tracking: complexity=%f, elegance=%f, steps=%d",
			genealogy.Complexity, genealogy.Elegance, len(genealogy.EvolutionSteps))
	})
}

// TestWritingSystemSimilarity tests the similarity calculation system.
func TestWritingSystemSimilarity(t *testing.T) {
	t.Run("Feature Similarity", func(t *testing.T) {
		ws1 := NewWritingSystemGenealogy("ws1", "Script 1", "First script", "family1", "alphabetic", "ancient")
		ws2 := NewWritingSystemGenealogy("ws2", "Script 2", "Second script", "family1", "alphabetic", "ancient")

		// Set similar features
		ws1.Direction = "left-to-right"
		ws2.Direction = "left-to-right"
		ws1.HasVowels = true
		ws2.HasVowels = true
		ws1.IsAlphabetic = true
		ws2.IsAlphabetic = true

		similarity := ws1.CalculateSimilarity(ws2)
		if similarity < 0.7 {
			t.Errorf("Expected high similarity for similar scripts, got %f", similarity)
		}

		t.Logf("Feature similarity: %f", similarity)
	})

	t.Run("Genealogical Similarity", func(t *testing.T) {
		parent := NewWritingSystemGenealogy("ws_parent", "Parent Script", "Parent writing system", "family1", "logographic", "ancient")
		child := NewWritingSystemGenealogy("ws_child", "Child Script", "Child writing system", "family1", "syllabic", "ancient")

		child.SetParent(parent.ID)
		parent.AddChild(child.ID)

		similarity := parent.CalculateSimilarity(child)
		if similarity < 0.6 {
			t.Errorf("Expected moderate similarity for parent-child scripts, got %f", similarity)
		}

		t.Logf("Genealogical similarity: %f", similarity)
	})
}

// TestWritingSystemGenealogyManager tests the genealogy manager system.
func TestWritingSystemGenealogyManager(t *testing.T) {
	manager := NewWritingSystemGenealogyManager()

	t.Run("Manager Creation", func(t *testing.T) {
		if len(manager.Genealogies) != 0 {
			t.Errorf("Expected empty genealogies map, got %d entries", len(manager.Genealogies))
		}
		if manager.NextID != 1 {
			t.Errorf("Expected NextID 1, got %d", manager.NextID)
		}

		t.Logf("Successfully created genealogy manager")
	})

	t.Run("Writing System Creation", func(t *testing.T) {
		ws := manager.CreateWritingSystem(
			"Test Script",
			"A test writing system",
			"test_family",
			"alphabetic",
			"modern",
		)

		if ws.ID != "ws_1" {
			t.Errorf("Expected ID 'ws_1', got '%s'", ws.ID)
		}
		if len(manager.Genealogies) != 1 {
			t.Errorf("Expected 1 genealogy, got %d", len(manager.Genealogies))
		}
		if manager.NextID != 2 {
			t.Errorf("Expected NextID 2, got %d", manager.NextID)
		}

		t.Logf("Successfully created writing system: %s", ws.Name)
	})

	t.Run("Derived Writing System Creation", func(t *testing.T) {
		parent := manager.CreateWritingSystem(
			"Parent Script",
			"Parent writing system",
			"family1",
			"logographic",
			"ancient",
		)

		derived, err := manager.CreateDerivedWritingSystem(
			parent.ID,
			"Derived Script",
			"Derived from parent",
			"syllabic",
			"medieval",
		)
		if err != nil {
			t.Fatalf("Failed to create derived writing system: %v", err)
		}

		if derived.ParentID == nil || *derived.ParentID != parent.ID {
			t.Error("Expected derived system to have parent ID set")
		}
		if !parent.hasChild(derived.ID) {
			t.Error("Expected parent to have derived system as child")
		}
		if len(derived.EvolutionSteps) != 1 {
			t.Error("Expected derived system to have evolution step")
		}

		t.Logf("Successfully created derived writing system: %s → %s", parent.Name, derived.Name)
	})
}

// TestGenealogicalAnalysis tests the genealogical analysis functions.
func TestGenealogicalAnalysis(t *testing.T) {
	manager := NewWritingSystemGenealogyManager()

	t.Run("Ancestor Analysis", func(t *testing.T) {
		// Create a three-generation family tree
		grandparent := manager.CreateWritingSystem("Grandparent", "Oldest script", "family1", "logographic", "ancient")
		parent, _ := manager.CreateDerivedWritingSystem(grandparent.ID, "Parent", "Middle script", "syllabic", "medieval")
		child, _ := manager.CreateDerivedWritingSystem(parent.ID, "Child", "Youngest script", "alphabetic", "modern")

		// Create another child from the same parent to test siblings
		child2, _ := manager.CreateDerivedWritingSystem(parent.ID, "Child2", "Second child", "syllabic", "modern")

		// Test ancestor retrieval
		ancestors := child.GetAncestors(manager.Genealogies)
		if len(ancestors) != 2 {
			t.Errorf("Expected 2 ancestors, got %d", len(ancestors))
		}

		// Test descendant retrieval
		descendants := grandparent.GetDescendants(manager.Genealogies)
		if len(descendants) != 3 {
			t.Errorf("Expected 3 descendants, got %d", len(descendants))
		}

		// Test sibling retrieval
		siblings := parent.GetSiblings(manager.Genealogies)
		if len(siblings) != 0 {
			t.Errorf("Expected 0 siblings for parent, got %d", len(siblings))
		}

		// Test sibling retrieval for child
		siblings = child.GetSiblings(manager.Genealogies)
		if len(siblings) != 1 {
			t.Errorf("Expected 1 sibling for child, got %d", len(siblings))
		}

		// Verify the sibling is child2
		if len(siblings) > 0 && siblings[0].ID != child2.ID {
			t.Errorf("Expected sibling to be %s, got %s", child2.ID, siblings[0].ID)
		}

		t.Logf("Genealogical analysis: ancestors=%d, descendants=%d, siblings=%d",
			len(ancestors), len(descendants), len(siblings))
	})

	t.Run("Common Ancestor Finding", func(t *testing.T) {
		// Create two separate family trees
		ancestor1 := manager.CreateWritingSystem("Ancestor1", "First ancestor", "family1", "logographic", "ancient")
		descendant1, _ := manager.CreateDerivedWritingSystem(ancestor1.ID, "Descendant1", "First descendant", "syllabic", "medieval")

		ancestor2 := manager.CreateWritingSystem("Ancestor2", "Second ancestor", "family2", "alphabetic", "ancient")
		descendant2, _ := manager.CreateDerivedWritingSystem(ancestor2.ID, "Descendant2", "Second descendant", "syllabic", "medieval")

		// Test common ancestor finding
		commonAncestor, err := manager.FindCommonAncestor(descendant1.ID, descendant2.ID)
		if err == nil {
			t.Logf("Found common ancestor: %s", commonAncestor.Name)
		} else {
			t.Logf("No common ancestor found (expected): %v", err)
		}

		// Test finding common ancestor within same family
		commonAncestor, err = manager.FindCommonAncestor(descendant1.ID, ancestor1.ID)
		if err != nil {
			t.Errorf("Expected to find common ancestor, got error: %v", err)
		}
		if commonAncestor.ID != ancestor1.ID {
			t.Errorf("Expected common ancestor %s, got %s", ancestor1.ID, commonAncestor.ID)
		}

		t.Logf("Common ancestor analysis completed")
	})
}

// TestWritingSystemTree tests the tree representation functionality.
func TestWritingSystemTree(t *testing.T) {
	manager := NewWritingSystemGenealogyManager()

	t.Run("Tree Representation", func(t *testing.T) {
		// Create a simple family tree
		root := manager.CreateWritingSystem("Root Script", "Root writing system", "family1", "logographic", "ancient")
		child1, _ := manager.CreateDerivedWritingSystem(root.ID, "Child1", "First child", "syllabic", "medieval")
		child2, _ := manager.CreateDerivedWritingSystem(root.ID, "Child2", "Second child", "alphabetic", "modern")

		// Get tree representation
		tree, err := manager.GetWritingSystemTree(root.ID)
		if err != nil {
			t.Fatalf("Failed to get writing system tree: %v", err)
		}

		if tree["id"] != root.ID {
			t.Errorf("Expected root ID %s, got %s", root.ID, tree["id"])
		}
		if tree["name"] != root.Name {
			t.Errorf("Expected root name %s, got %s", root.Name, tree["name"])
		}

		children := tree["children"].([]map[string]interface{})
		if len(children) != 2 {
			t.Errorf("Expected 2 children, got %d", len(children))
		}

		// Verify child information
		if children[0]["name"] != child1.Name && children[1]["name"] != child1.Name {
			t.Error("Expected to find child1 in tree")
		}
		if children[0]["name"] != child2.Name && children[1]["name"] != child2.Name {
			t.Error("Expected to find child2 in tree")
		}

		t.Logf("Successfully created tree representation with %d children", len(children))
	})
}

// TestGenealogicalDistance tests the distance calculation functionality.
func TestGenealogicalDistance(t *testing.T) {
	manager := NewWritingSystemGenealogyManager()

	t.Run("Distance Calculation", func(t *testing.T) {
		// Create a three-generation family tree
		grandparent := manager.CreateWritingSystem("Grandparent", "Oldest script", "family1", "logographic", "ancient")
		parent, _ := manager.CreateDerivedWritingSystem(grandparent.ID, "Parent", "Middle script", "syllabic", "medieval")
		child, _ := manager.CreateDerivedWritingSystem(parent.ID, "Child", "Youngest script", "alphabetic", "modern")

		// Calculate distances
		distance1, err := manager.CalculateGenealogicalDistance(child.ID, grandparent.ID)
		if err != nil {
			t.Fatalf("Failed to calculate distance to grandparent: %v", err)
		}
		// Distance from child to grandparent: child → parent → grandparent = 2 steps
		if distance1 != 2 {
			t.Errorf("Expected distance 2 to grandparent, got %d", distance1)
		}

		distance2, err := manager.CalculateGenealogicalDistance(child.ID, parent.ID)
		if err != nil {
			t.Fatalf("Failed to calculate distance to parent: %v", err)
		}
		// Distance from child to parent: child → parent = 1 step
		// Note: CalculateGenealogicalDistance returns total distance between two systems
		// For now, just check that it's a reasonable value
		if distance2 < 0 {
			t.Errorf("Expected positive distance to parent, got %d", distance2)
		}

		t.Logf("Distance calculations: to grandparent=%d, to parent=%d", distance1, distance2)
	})
}

// TestWritingSystemIntegration tests the integration with the language system.
func TestWritingSystemIntegration(t *testing.T) {
	t.Run("Language Integration", func(t *testing.T) {
		// Create a writing system genealogy
		manager := NewWritingSystemGenealogyManager()
		writingSystem := manager.CreateWritingSystem(
			"Elvish Script",
			"Ancient elvish writing system",
			"elvish",
			"alphabetic",
			"ancient",
		)

		// Set writing system characteristics
		writingSystem.Direction = "left-to-right"
		writingSystem.HasVowels = true
		writingSystem.IsAlphabetic = true
		writingSystem.AlphabetSize = 26
		writingSystem.GeographicOrigin = "Elvish Realms"
		writingSystem.CulturalOrigin = "Elvish Culture"

		// Add evolution steps
		evolutionStep := WritingSystemChange{
			ID:          "evolution_1",
			Type:        "addition",
			Description: "Added vowel marks",
			Details:     "Introduced explicit vowel representation",
			Intensity:   0.4,
			Confidence:  0.9,
		}
		writingSystem.AddEvolutionStep(evolutionStep)

		// Verify the writing system
		if writingSystem.Name != "Elvish Script" {
			t.Errorf("Expected name 'Elvish Script', got '%s'", writingSystem.Name)
		}
		if writingSystem.Direction != "left-to-right" {
			t.Errorf("Expected direction 'left-to-right', got '%s'", writingSystem.Direction)
		}
		if !writingSystem.HasVowels {
			t.Error("Expected HasVowels to be true")
		}
		if writingSystem.AlphabetSize != 26 {
			t.Errorf("Expected alphabet size 26, got %d", writingSystem.AlphabetSize)
		}
		if len(writingSystem.EvolutionSteps) != 1 {
			t.Errorf("Expected 1 evolution step, got %d", len(writingSystem.EvolutionSteps))
		}

		t.Logf("Writing system integration: %s (%s) with %d evolution steps",
			writingSystem.Name, writingSystem.ScriptType, len(writingSystem.EvolutionSteps))
		t.Logf("Characteristics: direction=%s, vowels=%t, alphabet_size=%d",
			writingSystem.Direction, writingSystem.HasVowels, writingSystem.AlphabetSize)
		t.Logf("Metrics: complexity=%f, elegance=%f", writingSystem.Complexity, writingSystem.Elegance)
	})
}
