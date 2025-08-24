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
