package lang

import (
	"fmt"
	"log"
)

// ExampleHighLevelAPI demonstrates how to use the high-level language generation API.
func ExampleHighLevelAPI() {
	fmt.Println("=== Kismet Language Generation API Example ===")

	// 1. Create a random language with cultural preferences
	fmt.Println("1. Creating an Elvish language...")
	elvishLang, err := CreateRandomLanguage("elvish", "high_elf", 42)
	if err != nil {
		log.Fatalf("Failed to create elvish language: %v", err)
	}

	fmt.Printf("   Created: %s\n", elvishLang.Name)
	fmt.Printf("   Culture: %s\n", elvishLang.Culture)
	fmt.Printf("   Complexity: %s\n", elvishLang.Complexity.String())
	fmt.Printf("   Description: %s\n\n", elvishLang.Description)

	// 2. Generate names for elvish characters
	fmt.Println("2. Generating Elvish names...")
	elvishNames, err := GenerateNames(elvishLang, "person", 5, "noble", 42)
	if err != nil {
		log.Fatalf("Failed to generate names: %v", err)
	}

	fmt.Printf("   Generated names: %v\n\n", elvishNames)

	// 3. Generate sample text in the language
	fmt.Println("3. Generating sample text...")
	elvishText, err := GenerateText(elvishLang, "poetry", "short", "nature", 42)
	if err != nil {
		log.Fatalf("Failed to generate text: %v", err)
	}

	fmt.Printf("   Sample poetry:\n%s\n\n", elvishText)

	// 4. Create a language family
	fmt.Println("4. Creating a language family...")
	protoLang, err := CreateRandomLanguage("proto", "ancient", 100)
	if err != nil {
		log.Fatalf("Failed to create proto language: %v", err)
	}

	family, err := CreateLanguageFamily(protoLang, 3, "1000_years", 100)
	if err != nil {
		log.Fatalf("Failed to create language family: %v", err)
	}

	fmt.Printf("   Proto language: %s\n", protoLang.Name)
	fmt.Printf("   Family members: %d\n", len(family))
	for i, lang := range family {
		if i == 0 {
			fmt.Printf("   - %s (proto)\n", lang.Name)
		} else {
			fmt.Printf("   - %s (branch)\n", lang.Name)
		}
	}
	fmt.Println()

	// 5. Simulate language contact
	fmt.Println("5. Simulating language contact...")
	humanLang, err := CreateRandomLanguage("human", "northern_human", 200)
	if err != nil {
		log.Fatalf("Failed to create human language: %v", err)
	}

	contactLang1, contactLang2, err := SimulateLanguageContact(
		elvishLang, humanLang, "trade", 0.7, "500_years", 300)
	if err != nil {
		log.Fatalf("Failed to simulate contact: %v", err)
	}

	fmt.Printf("   Elvish after contact: %s\n", contactLang1.Description)
	fmt.Printf("   Human after contact: %s\n\n", contactLang2.Description)

	// 6. Get language information
	fmt.Println("6. Language information...")
	info := GetLanguageInfo(elvishLang)
	fmt.Printf("   Phoneme count: %d\n", info.BasicStats.PhonemeCount)
	fmt.Printf("   Consonant count: %d\n", info.BasicStats.ConsonantCount)
	fmt.Printf("   Vowel count: %d\n", info.BasicStats.VowelCount)
	fmt.Printf("   Cultural features: %s\n", info.CulturalFeatures.TypicalDomains)
	fmt.Println()

	// 7. Validate language
	fmt.Println("7. Validating language...")
	validation := ValidateLanguage(elvishLang)
	if validation.IsValid {
		fmt.Println("   ✓ Language is valid and complete")
	} else {
		fmt.Printf("   ✗ Language has issues: %v\n", validation.Issues)
	}

	if len(validation.Warnings) > 0 {
		fmt.Printf("   ⚠ Warnings: %v\n", validation.Warnings)
	}

	if len(validation.Suggestions) > 0 {
		fmt.Printf("   💡 Suggestions: %v\n", validation.Suggestions)
	}
	fmt.Println()

	// 8. Create a cultural language set
	fmt.Println("8. Creating a cultural language ecosystem...")
	culturalLangs, err := CreateCulturalLanguageSet("dwarvish", 4, "high", 400)
	if err != nil {
		log.Fatalf("Failed to create cultural language set: %v", err)
	}

	fmt.Printf("   Created %d dwarf languages:\n", len(culturalLangs))
	for i, lang := range culturalLangs {
		if i == 0 {
			fmt.Printf("   - %s (main)\n", lang.Name)
		} else {
			fmt.Printf("   - %s (variant)\n", lang.Name)
		}
	}
	fmt.Println()

	// 9. Export language data
	fmt.Println("9. Exporting language data...")
	exportData, err := ExportLanguage(elvishLang, "json")
	if err != nil {
		log.Fatalf("Failed to export language: %v", err)
	}

	fmt.Printf("   Exported data length: %d bytes\n", len(exportData))
	fmt.Printf("   Export preview: %s...\n\n", string(exportData[:100]))

	fmt.Println("=== Example completed successfully! ===")
}

// ExampleQuickStart demonstrates the simplest way to use the API.
func ExampleQuickStart() {
	fmt.Println("=== Quick Start Example ===")

	// Create a language in one line
	lang, err := CreateRandomLanguage("orcish", "black_orc", 123)
	if err != nil {
		log.Fatalf("Failed to create language: %v", err)
	}

	// Generate some names
	names, err := GenerateNames(lang, "person", 3, "warrior", 123)
	if err != nil {
		log.Fatalf("Failed to generate names: %v", err)
	}

	// Generate some text
	text, err := GenerateText(lang, "dialogue", "short", "battle", 123)
	if err != nil {
		log.Fatalf("Failed to generate text: %v", err)
	}

	fmt.Printf("Created %s language\n", lang.Name)
	fmt.Printf("Generated names: %v\n", names)
	fmt.Printf("Generated text: %s\n", text)
}

// ExampleCulturalVariations shows how different cultures create different languages.
func ExampleCulturalVariations() {
	fmt.Println("=== Cultural Variations Example ===")

	cultures := []string{"elvish", "dwarvish", "orcish", "human", "ancient"}

	for _, culture := range cultures {
		lang, err := CreateRandomLanguage(culture, culture, int64(len(culture)))
		if err != nil {
			log.Printf("Failed to create %s language: %v", culture, err)
			continue
		}

		fmt.Printf("%s:\n", culture)
		fmt.Printf("  Name: %s\n", lang.Name)
		fmt.Printf("  Complexity: %s\n", lang.Complexity.String())
		fmt.Printf("  Description: %s\n\n", lang.Description)
	}
}
