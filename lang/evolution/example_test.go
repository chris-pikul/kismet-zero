package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
	"github.com/chris-pikul/kismet-zero/lang/orthography"
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// Evolution demonstrates how to use the evolution package to evolve languages.
func Evolution() {
	// Create an evolution engine with default configuration
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create a base language (e.g., Proto-Indo-European)
	protoLang := createProtoLanguage()

	// Evolve the language over time
	evolvedLang, evolutionEvent, err := engine.EvolveLanguage(
		protoLang,
		"early_evolution",
		time.Hour*24*365*100, // 100 years
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to evolve language: %v", err))
	}

	fmt.Printf("Language evolved: %s\n", evolvedLang.Name)
	fmt.Printf("Evolution event: %s\n", evolutionEvent.Name)
	fmt.Printf("Number of changes: %d\n", len(evolutionEvent.Changes))

	// Create a child language (e.g., Germanic)
	germanicID := lang.LanguageID{
		Family:   "indoeuropean",
		Branch:   "germanic",
		Language: "proto_germanic",
	}

	childLang, _, err := engine.CreateChildLanguage(
		evolvedLang,
		germanicID,
		"Proto-Germanic",
		"middle_evolution",
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create child language: %v", err))
	}

	fmt.Printf("Child language created: %s\n", childLang.Name)
	fmt.Printf("Parent: %s\n", childLang.ParentID.String())

	// Simulate contact with another language
	contactLang := createContactLanguage()

	contactEvolved, evolutionEvent, err := engine.EvolveLanguageFromContact(
		childLang,
		contactLang,
		ContactTypeTrade,
		0.6,                  // Medium intensity
		time.Hour*24*365*100, // 100 years
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to evolve through contact: %v", err))
	}

	fmt.Printf("Language evolved through contact: %s\n", contactEvolved.Name)
	fmt.Printf("Evolution trigger: %s\n", evolutionEvent.TriggerType)

	// Print the family tree
	familyTree := engine.GetFamilyTree()
	fmt.Println("\nLanguage Family Tree:")
	fmt.Println(familyTree.PrintFamilyTree())

	// Get evolution history
	history, err := engine.GetEvolutionHistory(protoLang.ID.String())
	if err != nil {
		panic(fmt.Sprintf("Failed to get evolution history: %v", err))
	}

	fmt.Printf("\nEvolution history for %s:\n", protoLang.Name)
	for _, event := range history {
		fmt.Printf("- %s: %s\n", event.Timestamp.Format("2006-01-02"), event.Description)
	}
}

// createProtoLanguage creates a sample proto-language for demonstration.
func createProtoLanguage() *lang.Language {
	// Create a basic phoneme pool
	pool := phoneme.NewPool()

	// Create a simple phonology
	phonology := phonology.NewPhonology(pool)

	// Create the language
	langID := lang.LanguageID{
		Family:   "indoeuropean",
		Branch:   "proto",
		Language: "proto_indoeuropean",
	}

	language := lang.NewLanguage(langID, "Proto-Indo-European", lang.LanguageTypeNatural, 42)
	language.SetPhonology(phonology)
	language.Culture = "proto_indoeuropean"
	language.Description = "Reconstructed ancestor of Indo-European languages"

	return language
}

// createContactLanguage creates a sample language for contact simulation.
func createContactLanguage() *lang.Language {
	// Create a basic phoneme pool
	pool := phoneme.NewPool()

	// Create a simple phonology
	phonology := phonology.NewPhonology(pool)

	// Create the language
	langID := lang.LanguageID{
		Family:   "uralic",
		Branch:   "finnic",
		Language: "proto_finnic",
	}

	language := lang.NewLanguage(langID, "Proto-Finnic", lang.LanguageTypeNatural, 43)
	language.SetPhonology(phonology)
	language.Culture = "proto_finnic"
	language.Description = "Reconstructed ancestor of Finnic languages"

	return language
}

// SoundChanges demonstrates how sound changes work in language evolution.
func SoundChanges() {
	config := DefaultEvolutionConfig(123)
	soundEngine := NewSoundChangeEngine(config)

	// Create a sample phonology
	pool := phoneme.NewPool()
	phonology := phonology.NewPhonology(pool)

	// Apply sound changes for a specific era
	changes, modifiedPhonology := soundEngine.ApplySoundChanges(phonology, "early_evolution")

	fmt.Printf("Applied %d sound changes\n", len(changes))
	fmt.Printf("Phonology modified: %v\n", modifiedPhonology != nil)

	// Print the changes
	for _, change := range changes {
		fmt.Printf("- %s: %s\n", change.Type.String(), change.Description)
	}
}

// ContactEvolution demonstrates how contact between languages works.
func ContactEvolution() {
	config := DefaultEvolutionConfig(456)
	contactEngine := NewContactEvolutionEngine(config)

	// Create two sample languages
	lang1 := createProtoLanguage()
	lang2 := createContactLanguage()

	// Simulate trade contact
	changes, contactEvent := contactEngine.SimulateContact(
		lang1, // source language
		lang2, // target language
		ContactTypeTrade,
		0.7,                  // High intensity
		time.Hour*24*365*200, // 200 years
	)

	fmt.Printf("Contact event: %s\n", contactEvent.Description)
	fmt.Printf("Applied %d changes\n", len(changes))
	fmt.Printf("Lexical borrowing: %v\n", contactEvent.LexicalBorrowing)
	fmt.Printf("Phonological borrowing: %v\n", contactEvent.PhonologicalBorrowing)
	fmt.Printf("Grammatical influence: %v\n", contactEvent.GrammaticalInfluence)

	// Calculate contact influence
	influence := contactEngine.CalculateContactInfluence(
		lang1,
		lang2,
		ContactTypeTrade,
		time.Hour*24*365*200,
	)

	fmt.Printf("Calculated influence: %.2f\n", influence)
}

// OrthographicEvolution demonstrates orthographic evolution features.
func OrthographicEvolution() {
	config := DefaultEvolutionConfig(789)
	orthoEngine := NewOrthographicEvolutionEngine(config)

	// Create a sample writing system
	writingSystem := &orthography.WritingSystem{
		Style: orthography.WritingStyleAlphabetic,
		Graphemes: []orthography.Grapheme{
			{Symbol: "a", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "b", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "c", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
		},
		Name:    "Test Writing System",
		Culture: "test_culture",
	}

	// Apply orthographic changes
	changes := orthoEngine.ApplyOrthographicChanges(writingSystem, "middle_evolution")

	fmt.Printf("Applied %d orthographic changes\n", len(changes))

	// Print the changes
	for _, change := range changes {
		fmt.Printf("Change: %s (%s)\n", change.Description, change.Type.String())
		fmt.Printf("  Details: %s\n", change.Details)
		fmt.Printf("  Complexity change: %.2f\n", change.ComplexityChange)
	}

	// Generate a script reform
	reform := orthoEngine.GenerateOrthographicReform(writingSystem, "spelling", "reform_era")
	fmt.Printf("\nScript reform: %s\n", reform.Description)
	fmt.Printf("  Type: %s\n", reform.ScriptReformType)
	fmt.Printf("  Impact: %.2f\n", reform.ComplexityChange)

	// Calculate complexity
	complexity := orthoEngine.CalculateOrthographicComplexity(writingSystem)
	fmt.Printf("Writing system complexity: %.2f\n", complexity)
}

// WritingSystemFamilyTreeExample demonstrates the writing system family tree functionality.
func WritingSystemFamilyTreeExample() {
	config := DefaultEvolutionConfig(999)
	tree := NewWritingSystemFamilyTree(config)

	// Create a proto-writing system (root)
	protoWriting := &orthography.WritingSystem{
		Style: orthography.WritingStyleAlphabetic,
		Graphemes: []orthography.Grapheme{
			{Symbol: "a", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "b", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "c", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
		},
		Name:    "Proto-Alphabetic",
		Culture: "proto_culture",
	}

	// Add the proto-writing system as root
	err := tree.AddWritingSystem(protoWriting, nil, "origin")
	if err != nil {
		panic(fmt.Sprintf("Failed to add proto-writing system: %v", err))
	}

	// Create a child writing system
	childWriting := &orthography.WritingSystem{
		Style: orthography.WritingStyleAlphabetic,
		Graphemes: []orthography.Grapheme{
			{Symbol: "a", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "b", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "c", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "d", Type: orthography.GraphemeTypeLetter, Weight: 1.0}, // New grapheme
		},
		Name:    "Child-Alphabetic",
		Culture: "child_culture",
	}

	// Add the child writing system
	err = tree.AddWritingSystem(childWriting, protoWriting, "natural_evolution")
	if err != nil {
		panic(fmt.Sprintf("Failed to add child writing system: %v", err))
	}

	// Create a sibling writing system
	siblingWriting := &orthography.WritingSystem{
		Style: orthography.WritingStyleSyllabic, // Different style
		Graphemes: []orthography.Grapheme{
			{Symbol: "ba", Type: orthography.GraphemeTypeSyllable, Weight: 1.0},
			{Symbol: "ca", Type: orthography.GraphemeTypeSyllable, Weight: 1.0},
		},
		Name:    "Sibling-Syllabic",
		Culture: "sibling_culture",
	}

	// Add the sibling writing system
	err = tree.AddWritingSystem(siblingWriting, protoWriting, "style_evolution")
	if err != nil {
		panic(fmt.Sprintf("Failed to add sibling writing system: %v", err))
	}

	// Print the family tree
	fmt.Println("Writing System Family Tree:")
	fmt.Println(tree.PrintFamilyTree())

	// Get ancestors of the child
	ancestors := tree.GetAncestors(childWriting.Name)
	fmt.Printf("\nAncestors of %s: %d\n", childWriting.Name, len(ancestors))

	// Get descendants of the proto system
	descendants := tree.GetDescendants(protoWriting.Name)
	fmt.Printf("Descendants of %s: %d\n", protoWriting.Name, len(descendants))

	// Get siblings of the child
	siblings := tree.GetSiblings(childWriting.Name)
	fmt.Printf("Siblings of %s: %d\n", childWriting.Name, len(siblings))

	// Calculate similarity between siblings
	similarity, err := tree.CalculateWritingSystemSimilarity(childWriting.Name, siblingWriting.Name)
	if err != nil {
		panic(fmt.Sprintf("Failed to calculate similarity: %v", err))
	}
	fmt.Printf("Similarity between siblings: %.2f\n", similarity)

	// Estimate divergence time
	divergenceTime, err := tree.EstimateDivergenceTime(childWriting.Name, siblingWriting.Name)
	if err != nil {
		panic(fmt.Sprintf("Failed to estimate divergence time: %v", err))
	}
	fmt.Printf("Estimated divergence time: %v\n", divergenceTime.Format("2006-01-02"))

	// Add some evolution events
	evolutionChange := OrthographicChange{
		ID:          "test_evolution",
		Type:        OrthographicChangeTypeInnovation,
		Description: "Added new grapheme",
		Timestamp:   time.Now(),
	}

	err = tree.AddEvolutionEvent(childWriting.Name, evolutionChange)
	if err != nil {
		panic(fmt.Sprintf("Failed to add evolution event: %v", err))
	}

	// Add a borrowing event
	borrowingEvent := BorrowingEvent{
		ID:               "test_borrowing",
		Timestamp:        time.Now(),
		SourceID:         siblingWriting.Name,
		TargetID:         childWriting.Name,
		BorrowedFeatures: []string{"syllable_structure"},
		BorrowingType:    "style",
		Intensity:        0.6,
	}

	err = tree.AddBorrowingEvent(childWriting.Name, borrowingEvent)
	if err != nil {
		panic(fmt.Sprintf("Failed to add borrowing event: %v", err))
	}

	// Get evolution history
	evolutionHistory, err := tree.GetEvolutionHistory(childWriting.Name)
	if err != nil {
		panic(fmt.Sprintf("Failed to get evolution history: %v", err))
	}
	fmt.Printf("\nEvolution history for %s: %d events\n", childWriting.Name, len(evolutionHistory))

	// Get borrowing history
	borrowingHistory, err := tree.GetBorrowingHistory(childWriting.Name)
	if err != nil {
		panic(fmt.Sprintf("Failed to get borrowing history: %v", err))
	}
	fmt.Printf("Borrowing history for %s: %d events\n", childWriting.Name, len(borrowingHistory))

	// Test lineage generation
	evolutionEngine := NewOrthographicEvolutionEngine(config)
	newLineage, err := tree.GenerateWritingSystemLineage(
		childWriting,
		"Grandchild",
		orthography.WritingStyleLogographic,
		evolutionEngine,
		"modern_era",
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to generate lineage: %v", err))
	}

	fmt.Printf("\nGenerated new lineage: %s (%s)\n", newLineage.Name, newLineage.Style.String())

	// Print updated family tree
	fmt.Println("\nUpdated Writing System Family Tree:")
	fmt.Println(tree.PrintFamilyTree())
}

// CulturalInfluenceExample demonstrates the sophisticated cultural influence modeling.
func CulturalInfluenceExample() {
	config := DefaultEvolutionConfig(777)
	engine := NewCulturalInfluenceEngine(config)

	// Create source language (dominant empire)
	sourceLang := &lang.Language{
		ID:      lang.LanguageID{Family: "imperial", Branch: "central", Language: "imperial_common"},
		Name:    "Imperial Common",
		Culture: "imperial_empire",
		Seed:    123,
	}

	// Create target language (frontier colony)
	targetLang := &lang.Language{
		ID:      lang.LanguageID{Family: "frontier", Branch: "western", Language: "frontier_tongue"},
		Name:    "Frontier Tongue",
		Culture: "frontier_colony",
		Seed:    456,
	}

	// Define imperial cultural identity (powerful, prestigious, magical)
	imperialIdentity := CulturalIdentity{
		CultureName:           "Imperial Empire",
		WritingSystemPrestige: 0.9, // Very prestigious writing system
		InnovationTendency:    0.6, // Moderate innovation
		PreservationInstinct:  0.7, // High preservation of traditions
		MagicalTradition:      0.8, // Strong magical practices
		ReligiousInfluence:    0.9, // Dominant religious system
		EconomicPower:         0.9, // Economic powerhouse
		MilitaryPower:         0.9, // Military superpower
		CulturalConfidence:    0.9, // Very confident in superiority
	}

	// Define frontier cultural identity (adaptable, innovative, less powerful)
	frontierIdentity := CulturalIdentity{
		CultureName:           "Frontier Colony",
		WritingSystemPrestige: 0.3, // Less prestigious writing
		InnovationTendency:    0.8, // Very innovative
		PreservationInstinct:  0.2, // Low preservation instinct
		MagicalTradition:      0.4, // Moderate magical practices
		ReligiousInfluence:    0.3, // Less religious influence
		EconomicPower:         0.4, // Moderate economy
		MilitaryPower:         0.3, // Weak military
		CulturalConfidence:    0.4, // Moderate cultural confidence
	}

	fmt.Println("=== Cultural Influence Simulation ===")
	fmt.Printf("Source: %s (%s)\n", sourceLang.Name, sourceLang.Culture)
	fmt.Printf("Target: %s (%s)\n", targetLang.Name, targetLang.Culture)
	fmt.Println()

	// Simulate different types of cultural contact

	// 1. Conquest (high intensity)
	fmt.Println("1. CONQUEST SCENARIO")
	changes, event := engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeConquest,
		time.Hour*24*365*100, // 100 years
		imperialIdentity,
		frontierIdentity,
	)

	fmt.Printf("  Contact Type: %s\n", event.ContactType.String())
	fmt.Printf("  Intensity: %.2f\n", event.Intensity)
	fmt.Printf("  Orthographic Influence: %t\n", event.OrthographicInfluence)
	fmt.Printf("  Lexical Influence: %t\n", event.LexicalInfluence)
	fmt.Printf("  Phonological Influence: %t\n", event.PhonologicalInfluence)
	fmt.Printf("  Grammatical Influence: %t\n", event.GrammaticalInfluence)
	fmt.Printf("  Magical Influence: %t\n", event.MagicalInfluence)
	fmt.Printf("  Religious Influence: %t\n", event.ReligiousInfluence)
	fmt.Printf("  Changes Generated: %d\n", len(changes))
	fmt.Println()

	// 2. Trade (moderate intensity)
	fmt.Println("2. TRADE SCENARIO")
	changes, event = engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeTrade,
		time.Hour*24*365*50, // 50 years
		imperialIdentity,
		frontierIdentity,
	)

	fmt.Printf("  Contact Type: %s\n", event.ContactType.String())
	fmt.Printf("  Intensity: %.2f\n", event.Intensity)
	fmt.Printf("  Orthographic Influence: %t\n", event.OrthographicInfluence)
	fmt.Printf("  Lexical Influence: %t\n", event.LexicalInfluence)
	fmt.Printf("  Changes Generated: %d\n", len(changes))
	fmt.Println()

	// 3. Religious influence (moderate-high intensity)
	fmt.Println("3. RELIGIOUS INFLUENCE SCENARIO")
	changes, event = engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeReligious,
		time.Hour*24*365*75, // 75 years
		imperialIdentity,
		frontierIdentity,
	)

	fmt.Printf("  Contact Type: %s\n", event.ContactType.String())
	fmt.Printf("  Intensity: %.2f\n", event.Intensity)
	fmt.Printf("  Religious Influence: %t\n", event.ReligiousInfluence)
	fmt.Printf("  Changes Generated: %d\n", len(changes))
	fmt.Println()

	// 4. Ancient influence (high intensity, mysterious)
	fmt.Println("4. ANCIENT INFLUENCE SCENARIO")
	changes, event = engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeAncient,
		time.Hour*24*365*200, // 200 years
		imperialIdentity,
		frontierIdentity,
	)

	fmt.Printf("  Contact Type: %s\n", event.ContactType.String())
	fmt.Printf("  Intensity: %.2f\n", event.Intensity)
	fmt.Printf("  Ancient Influence: %t\n", event.AncientInfluence)
	fmt.Printf("  Changes Generated: %d\n", len(changes))
	fmt.Println()

	// 5. Magical influence (special case)
	fmt.Println("5. MAGICAL INFLUENCE SCENARIO")
	changes, event = engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeMagical,
		time.Hour*24*365*60, // 60 years
		imperialIdentity,
		frontierIdentity,
	)

	fmt.Printf("  Contact Type: %s\n", event.ContactType.String())
	fmt.Printf("  Intensity: %.2f\n", event.Intensity)
	fmt.Printf("  Magical Influence: %t\n", event.MagicalInfluence)
	fmt.Printf("  Changes Generated: %d\n", len(changes))
	fmt.Println()

	// Show cultural compatibility analysis
	fmt.Println("=== CULTURAL COMPATIBILITY ANALYSIS ===")
	compatibility := engine.calculateCulturalCompatibility(sourceLang, targetLang, imperialIdentity, frontierIdentity)

	fmt.Printf("Linguistic Similarity: %.2f\n", compatibility.LinguisticSimilarity)
	fmt.Printf("Historical Relationship: %.2f\n", compatibility.HistoricalRelationship)
	fmt.Printf("Power Balance: %.2f (negative = source dominant)\n", compatibility.PowerBalance)
	fmt.Printf("Geographic Proximity: %.2f\n", compatibility.GeographicProximity)
	fmt.Printf("Cultural Values: %.2f\n", compatibility.CulturalValues)
	fmt.Printf("Magical Affinity: %.2f\n", compatibility.MagicalAffinity)
	fmt.Printf("Religious Compatibility: %.2f\n", compatibility.ReligiousCompatibility)
	fmt.Printf("Economic Interdependence: %.2f\n", compatibility.EconomicInterdependence)
	fmt.Println()

	// Show resistance and success factors
	fmt.Println("=== INFLUENCE FACTORS ===")
	resistanceFactors := engine.identifyResistanceFactors(frontierIdentity, compatibility)
	successFactors := engine.identifySuccessFactors(imperialIdentity, compatibility)

	fmt.Printf("Resistance Factors: %v\n", resistanceFactors)
	fmt.Printf("Success Factors: %v\n", successFactors)
	fmt.Println()

	// Show adaptation notes
	if len(changes) > 0 {
		adaptationNotes := engine.generateAdaptationNotes(changes, compatibility)
		fmt.Printf("Adaptation Notes: %s\n", adaptationNotes)
	}

	fmt.Println("=== SCENARIO COMPLETE ===")
	fmt.Println("This demonstrates how different cultural contact types")
	fmt.Println("result in varying levels of linguistic and cultural influence,")
	fmt.Println("with sophisticated modeling of cultural compatibility factors.")
}
