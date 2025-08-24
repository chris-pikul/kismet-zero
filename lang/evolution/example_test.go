package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
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

	protoLang := lang.NewLanguage(langID, "Proto-Indo-European", lang.LanguageTypeNatural, 42)
	protoLang.SetPhonology(phonology)
	protoLang.Culture = "Proto-Indo-European"

	return protoLang
}

// createContactLanguage creates a sample contact language for demonstration.
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

	contactLang := lang.NewLanguage(langID, "Proto-Finnic", lang.LanguageTypeNatural, 42)
	contactLang.SetPhonology(phonology)
	contactLang.Culture = "Uralic"

	return contactLang
}

// CulturalInfluence demonstrates cultural influence modeling.
func CulturalInfluence() {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create source and target languages
	sourceLang := createProtoLanguage()
	targetLang := createContactLanguage()

	// Define cultural identities
	sourceIdentity := CulturalIdentity{
		CultureName:           "Source Culture",
		WritingSystemPrestige: 0.8,
		InnovationTendency:    0.7,
		PreservationInstinct:  0.6,
		MagicalTradition:      0.3,
		ReligiousInfluence:    0.4,
		EconomicPower:         0.6,
		MilitaryPower:         0.5,
		CulturalConfidence:    0.7,
	}

	targetIdentity := CulturalIdentity{
		CultureName:           "Target Culture",
		WritingSystemPrestige: 0.4,
		InnovationTendency:    0.3,
		PreservationInstinct:  0.8,
		MagicalTradition:      0.7,
		ReligiousInfluence:    0.6,
		EconomicPower:         0.4,
		MilitaryPower:         0.3,
		CulturalConfidence:    0.5,
	}

	// Simulate cultural influence
	changes, event, err := engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeTrade,
		time.Hour*24*365*50, // 50 years
		sourceIdentity,
		targetIdentity,
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to simulate cultural influence: %v", err))
	}

	fmt.Printf("Cultural influence event: %s\n", event.Description)
	fmt.Printf("Number of linguistic changes: %d\n", len(changes))
	fmt.Printf("Influence intensity: %.2f\n", event.Intensity)

	// Print details of changes
	for i, change := range changes {
		fmt.Printf("Change %d: %s (%s)\n", i+1, change.Description, change.Type.String())
	}

	fmt.Println("This demonstrates how different cultural contact types")
	fmt.Println("result in varying levels of linguistic and cultural influence,")
	fmt.Println("with sophisticated modeling of cultural compatibility factors.")
}

// DialectFormation demonstrates how to create and evolve dialects.
func DialectFormation() {
	// Create an evolution engine
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create a parent language
	parentLangID := lang.LanguageID{
		Family:   "hum",
		Branch:   "germanic",
		Language: "english",
	}
	parentLang := lang.NewLanguage(parentLangID, "English", lang.LanguageTypeNatural, 42)
	parentLang.Culture = "British"

	// Create a geographic dialect (e.g., Appalachian English)
	appalachianRegion := GeographicRegion{
		ID:           "appalachia",
		Name:         "Appalachian Mountains",
		Latitude:     35.0,
		Longitude:    -82.0,
		Climate:      "temperate",
		Terrain:      "mountain",
		Population:   3000,
		Urbanization: 0.2,
	}

	appalachianDialect, formationEvent, err := engine.CreateGeographicDialect(
		parentLang,
		"appalachian",
		"Appalachian English",
		appalachianRegion,
		"colonial",
	)
	if err != nil {
		fmt.Printf("Error creating dialect: %v\n", err)
		return
	}

	fmt.Printf("Created dialect: %s\n", appalachianDialect.Name)
	fmt.Printf("Dialect type: %s\n", appalachianDialect.Type.String())
	fmt.Printf("Formation factors: %v\n", formationEvent.GeographicFactors)
	fmt.Printf("Intelligibility with parent: %.2f\n", appalachianDialect.Features.IntelligibilityScore)

	// Create a social dialect (e.g., Cockney English)
	cockneyDialect, cockneyEvent, err := engine.CreateSocialDialect(
		parentLang,
		"cockney",
		"Cockney English",
		"working",
		0.9, // High urbanization
		"industrial",
	)
	if err != nil {
		fmt.Printf("Error creating social dialect: %v\n", err)
		return
	}

	fmt.Printf("\nCreated social dialect: %s\n", cockneyDialect.Name)
	fmt.Printf("Dialect type: %s\n", cockneyDialect.Type.String())
	fmt.Printf("Social factors: %v\n", cockneyEvent.SocialFactors)

	// Evolve the Appalachian dialect over time
	evolvedDialect, evolutionEvent, err := engine.EvolveDialect(
		appalachianDialect,
		"modern",
		time.Hour*24*365*200, // 200 years
	)
	if err != nil {
		fmt.Printf("Error evolving dialect: %v\n", err)
		return
	}

	fmt.Printf("\nEvolved dialect: %s\n", evolvedDialect.Name)
	fmt.Printf("Evolution changes: %d\n", len(evolutionEvent.Changes))
	fmt.Printf("New intelligibility: %.2f\n", evolvedDialect.Features.IntelligibilityScore)

	// Get all dialects of the parent language
	dialects := engine.GetFamilyTree().GetDialectsByParent(parentLang.ID.String())
	fmt.Printf("\nTotal dialects of %s: %d\n", parentLang.Name, len(dialects))
}

// ExampleDialectFeatures demonstrates the distinctive features of dialects.
func ExampleDialectFeatures() {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create a coastal dialect
	coastalRegion := GeographicRegion{
		ID:           "coastal",
		Name:         "Coastal Region",
		Climate:      "tropical",
		Terrain:      "coastal",
		Population:   12000,
		Urbanization: 0.7,
	}

	parentLang := lang.NewLanguage(
		lang.LanguageID{Family: "test", Branch: "test", Language: "testlang"},
		"Test Language",
		lang.LanguageTypeNatural,
		42,
	)

	dialect, _, err := engine.CreateGeographicDialect(
		parentLang,
		"coastal_dialect",
		"Coastal Dialect",
		coastalRegion,
		"modern",
	)
	if err != nil {
		fmt.Printf("Error creating dialect: %v\n", err)
		return
	}

	fmt.Printf("Dialect: %s\n", dialect.Name)
	fmt.Printf("Phonological features: %v\n", dialect.Features.PhonologicalFeatures)
	fmt.Printf("Lexical features: %v\n", dialect.Features.LexicalFeatures)
	fmt.Printf("Grammatical features: %v\n", dialect.Features.GrammaticalFeatures)
	fmt.Printf("Intelligibility: %.2f\n", dialect.Features.IntelligibilityScore)
}

// ExampleLinguisticAdaptation demonstrates the enhanced linguistic adaptation system.
func ExampleLinguisticAdaptation() {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create a test language
	testLang := lang.NewLanguage(
		lang.LanguageID{Family: "test", Branch: "test", Language: "testlang"},
		"Test Language",
		lang.LanguageTypeNatural,
		42,
	)

	fmt.Printf("Language: %s\n", testLang.Name)

	// Generate a random phonological adaptation
	adaptation, err := engine.GenerateRandomAdaptation(testLang, "phonological")
	if err != nil {
		fmt.Printf("Error generating adaptation: %v\n", err)
		return
	}

	fmt.Printf("Generated adaptation: %s\n", adaptation.Description)
	fmt.Printf("Adaptation type: %s\n", adaptation.Type)
	fmt.Printf("Complexity change: %.2f\n", adaptation.ComplexityChange)

	// Apply the phonological adaptation
	change, err := engine.ApplyLinguisticAdaptation(testLang, adaptation)
	if err != nil {
		fmt.Printf("Error applying adaptation: %v\n", err)
		return
	}

	fmt.Printf("Applied change: %s\n", change.Description)
	fmt.Printf("Change type: %s\n", change.Type.String())
	fmt.Printf("Change direction: %s\n", change.Direction.String())
	fmt.Printf("Intensity: %.2f\n", change.Intensity)

	// Test borrowing pattern application
	borrowingPattern := &BorrowingPattern{
		SelectiveAdoption:     0.7,
		AdaptationStrength:    0.8,
		IntegrationDepth:      0.6,
		ResistanceLevel:       0.3,
		PrestigeSensitivity:   0.5,
		HybridizationTendency: 0.4,
	}

	borrowingChange, err := engine.ApplyBorrowingPattern(testLang, borrowingPattern, "Source Language", "cultural")
	if err != nil {
		fmt.Printf("Failed to apply borrowing pattern: %v\n", err)
	} else {
		fmt.Printf("Applied borrowing pattern: %s\n", borrowingChange.Description)
	}
}
