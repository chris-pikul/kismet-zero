package morphology

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// Example demonstrates basic usage of the morphology package.
func Example() {
	// Create a phoneme pool
	pool := phoneme.NewPool()

	// Create a phonology with the pool
	phon := phonology.NewPhonology(pool)

	// Add some syllable templates
	phon.AddTemplate(phonology.TemplateCV, 10)
	phon.AddTemplate(phonology.TemplateCVC, 8)

	// Create a morphology generator
	generator := NewMorphologyGenerator(phon)

	// Set the culture
	generator.SetCulture("elvish")

	// Create a random number generator
	rng := rand.New(rand.NewPCG(1, 2))

	// Generate a root morpheme
	root, err := generator.GenerateRootMorpheme("walk", rng)
	if err != nil {
		fmt.Printf("Error generating root morpheme: %v\n", err)
		return
	}

	fmt.Printf("Generated root morpheme: %s (meaning: %s)\n", root.ID, root.Meaning)

	// Generate a word from the root
	word, err := generator.GenerateWord(root, WordCategoryVerb, rng)
	if err != nil {
		fmt.Printf("Error generating word: %v\n", err)
		return
	}

	fmt.Printf("Generated word: %s (meaning: %s, category: %s)\n",
		word.Written, word.Meaning, word.Category.String())

	// Generate a personal name
	name, err := generator.GeneratePersonalName("male", "formal", rng)
	if err != nil {
		fmt.Printf("Error generating name: %v\n", err)
		return
	}

	fmt.Printf("Generated name: %s (meaning: %s)\n", name.Value, name.Meaning)

	// Generate a basic lexicon
	err = generator.GenerateLexicon(20, rng)
	if err != nil {
		fmt.Printf("Error generating lexicon: %v\n", err)
		return
	}

	// Get statistics
	stats := generator.GetStatistics()
	fmt.Printf("Generated lexicon with %d total words\n", stats["total_words"])
}

// ExampleMorphemeBuilder demonstrates using the MorphemeBuilder directly.
func ExampleMorphemeBuilder() {
	// Create a phoneme pool
	pool := phoneme.NewPool()

	// Create a phonology
	phon := phonology.NewPhonology(pool)
	phon.AddTemplate(phonology.TemplateCV, 10)

	// Create a morpheme builder
	builder := NewMorphemeBuilder(phon)

	// Create a random number generator
	rng := rand.New(rand.NewPCG(1, 2))

	// Build a root morpheme
	root, err := builder.BuildRootMorpheme("water", "nordic", rng)
	if err != nil {
		fmt.Printf("Error building root morpheme: %v\n", err)
		return
	}

	fmt.Printf("Built root morpheme: %s\n", root.ID)

	// Build an affix morpheme
	suffix, err := builder.BuildAffixMorpheme(MorphemeTypeSuffix, "small", "nordic", rng)
	if err != nil {
		fmt.Printf("Error building suffix morpheme: %v\n", err)
		return
	}

	fmt.Printf("Built suffix morpheme: %s\n", suffix.ID)
}

// ExampleRuleSet demonstrates using morphological rules.
func ExampleRuleSet() {
	// Create a rule set
	ruleSet := NewRuleSet()

	// Add a derivational rule
	agentRule := NewDerivationalRule(
		"agent_noun",
		MorphemeTypeRoot,
		MorphemeTypeSuffix,
		"one who does",
		1.0,
	)
	ruleSet.AddDerivationalRule(agentRule)

	// Add an inflectional rule
	pluralRule := NewInflectionalRule(
		"plural",
		WordCategoryNoun,
		[]string{"plural"},
		"many",
		1.0,
	)
	ruleSet.AddInflectionalRule(pluralRule)

	fmt.Printf("Created rule set with %d total rules\n", ruleSet.GetRuleCount())
	fmt.Printf("Derivational rules: %d\n", ruleSet.GetDerivationalRuleCount())
	fmt.Printf("Inflectional rules: %d\n", ruleSet.GetInflectionalRuleCount())
}

// ExampleLexicon demonstrates lexicon management.
func ExampleLexicon() {
	// Create a lexicon
	lexicon := NewLexicon("elvish")

	// Create some words
	word1 := &Word{
		ID:       "word1",
		Meaning:  "tree",
		Category: WordCategoryNoun,
		Culture:  "elvish",
		Weight:   1.0,
	}

	word2 := &Word{
		ID:       "word2",
		Meaning:  "walk",
		Category: WordCategoryVerb,
		Culture:  "elvish",
		Weight:   1.0,
	}

	// Add words to lexicon
	lexicon.AddWord(word1)
	lexicon.AddWord(word2)

	fmt.Printf("Lexicon contains %d words\n", lexicon.GetWordCount())
	fmt.Printf("Nouns: %d\n", lexicon.GetWordCountByCategory(WordCategoryNoun))
	fmt.Printf("Verbs: %d\n", lexicon.GetWordCountByCategory(WordCategoryVerb))

	// Search for words
	results := lexicon.SearchWords("tree")
	fmt.Printf("Found %d words matching 'tree'\n", len(results))
}

// ExampleNameGenerator demonstrates name generation.
func ExampleNameGenerator() {
	// Create a phoneme pool and phonology
	pool := phoneme.NewPool()
	phon := phonology.NewPhonology(pool)
	phon.AddTemplate(phonology.TemplateCV, 10)

	// Create a lexicon
	lexicon := NewLexicon("dwarven")

	// Create a name generator
	nameGen := NewNameGenerator(phon, lexicon)

	// Create a random number generator
	rng := rand.New(rand.NewPCG(1, 2))

	// Generate different types of names
	personalName, err := nameGen.GeneratePersonalName("dwarven", "male", "formal", rng)
	if err == nil {
		fmt.Printf("Personal name: %s (%s)\n", personalName.Value, personalName.Meaning)
	}

	placeName, err := nameGen.GeneratePlaceName("dwarven", "mountain", rng)
	if err == nil {
		fmt.Printf("Place name: %s (%s)\n", placeName.Value, placeName.Meaning)
	}

	title, err := nameGen.GenerateTitle("dwarven", "king", rng)
	if err == nil {
		fmt.Printf("Title: %s (%s)\n", title.Value, title.Meaning)
	}
}
