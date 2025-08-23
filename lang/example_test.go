package lang_test

import (
	"fmt"
	"log"

	"github.com/chris-pikul/kismet-zero/lang"
)

// ExampleLanguageID demonstrates how to create and use LanguageID structs.
func ExampleLanguageID() {
	// Create a simple language ID
	simpleID := lang.LanguageID{
		Family: "hum",
	}
	fmt.Printf("Simple ID: %s\n", simpleID.String())

	// Create a more complex language ID
	complexID := lang.LanguageID{
		Family:   "elv",
		Branch:   "wood",
		Language: "sindarin",
		Script:   "tengwar",
	}
	fmt.Printf("Complex ID: %s\n", complexID.String())

	// Parse a language ID from a string
	parsedID, err := lang.ParseLanguageID("orc-black-speech")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed ID: %s\n", parsedID.String())

	// Output:
	// Simple ID: hum
	// Complex ID: elv-wood-sindarin-tengwar
	// Parsed ID: orc-black-speech
}

// ExampleLanguage demonstrates how to create and configure a Language.
func ExampleLanguage() {
	// Create a new language with a unique ID and seed
	id := lang.LanguageID{
		Family:   "elv",
		Branch:   "high",
		Language: "quenya",
	}

	language := lang.NewLanguage(id, "Quenya", lang.LanguageTypeConstructed, 12345)

	// Configure the language
	language.SetComplexity(lang.LanguageComplexityComplex)
	language.SetCulture("elvish")
	language.SetDescription("The ancient language of the High Elves")

	// Check if the language is complete
	if !language.IsComplete() {
		missing := language.GetMissingComponents()
		fmt.Printf("Language is missing: %v\n", missing)
	}

	fmt.Printf("Created language: %s\n", language.String())
	fmt.Printf("Type: %s, Complexity: %s\n", language.Type.String(), language.Complexity.String())

	// Output:
	// Language is missing: [phonology orthography morphology grammar]
	// Created language: Language(elv-high-quenya: Quenya)
	// Type: constructed, Complexity: complex
}

// ExampleLanguage_clone demonstrates how to clone and evolve languages.
func ExampleLanguage_clone() {
	// Create an original language
	originalID := lang.LanguageID{Family: "hum", Language: "english"}
	original := lang.NewLanguage(originalID, "English", lang.LanguageTypeNatural, 1000)
	original.SetComplexity(lang.LanguageComplexityModerate)

	// Clone it to create a dialect
	dialectID := lang.LanguageID{Family: "hum", Language: "english", Dialect: "american"}
	dialect := original.Clone(dialectID, 2000)

	fmt.Printf("Original: %s\n", original.String())
	fmt.Printf("Dialect: %s\n", dialect.String())
	fmt.Printf("Dialect parent: %s\n", dialect.ParentID.String())

	// Output:
	// Original: Language(hum-english: English)
	// Dialect: Language(hum-english-american: English (clone))
	// Dialect parent: hum-english
}
