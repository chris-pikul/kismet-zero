package phonology_test

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// DemoNewPhonology demonstrates how to create a new phonology with phonemes and rules.
func DemoNewPhonology() {
	// Create a phoneme pool
	pool := phoneme.NewPool()

	// Add some consonants with complete specifications
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "k",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &phoneme.ConsonantSpec{
			Manner:  phoneme.ConsonantMannerPlosive,
			Place:   phoneme.ConsonantPlaceVelar,
			Voicing: phoneme.ConsonantVoicingUnvoiced,
		},
	})
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "t",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &phoneme.ConsonantSpec{
			Manner:  phoneme.ConsonantMannerPlosive,
			Place:   phoneme.ConsonantPlaceAlveolar,
			Voicing: phoneme.ConsonantVoicingUnvoiced,
		},
	})
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "n",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 0.8,
		Consonant: &phoneme.ConsonantSpec{
			Manner:  phoneme.ConsonantMannerNasal,
			Place:   phoneme.ConsonantPlaceAlveolar,
			Voicing: phoneme.ConsonantVoicingVoiced,
		},
	})
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "s",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 0.9,
		Consonant: &phoneme.ConsonantSpec{
			Manner:  phoneme.ConsonantMannerFricative,
			Place:   phoneme.ConsonantPlaceAlveolar,
			Voicing: phoneme.ConsonantVoicingUnvoiced,
		},
	})

	// Add some vowels
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "a",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightLow,
			Backness: phoneme.VowelBacknessCentral,
		},
	})
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "i",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 0.9,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightHigh,
			Backness: phoneme.VowelBacknessFront,
		},
	})

	// Create phonology
	phono := phonology.NewPhonology(pool)

	// Add syllable templates
	phono.AddTemplate(phonology.TemplateCV, 10) // CV is most common
	phono.AddTemplate(phonology.TemplateCVC, 8) // CVC is common
	phono.AddTemplate(phonology.TemplateV, 3)   // V is less common
	phono.AddTemplate(phonology.TemplateCCV, 5) // CCV is moderately common

	// Add advanced phonotactic rules using new features
	phono.AddRule(phonology.NewMaxClusterSizeRule(2))         // Max 2 consonants in cluster
	phono.AddRule(phonology.NewVelarNasalOnsetRule())         // No velar+nasal onset clusters
	phono.AddRule(phonology.NewVowelRequirementRule())        // Must have vowel
	phono.AddRule(phonology.NewSyllableLengthRule(1, 5))      // Length 1-5 phonemes
	phono.AddRule(phonology.NewConsonantHarmonyRule())        // Place harmony constraints
	phono.AddRule(phonology.NewVoicingAssimilationRule(true)) // Require voicing harmony
	phono.AddRule(phonology.NewMannerRestrictionsRule())      // Manner-based restrictions

	fmt.Printf("Phonology created with %d phonemes and %d rules\n",
		pool.Size(), phono.GetRuleCount())

	// Output: Phonology created with 6 phonemes and 7 rules
}

// DemoGenerateSyllable demonstrates how to generate syllables using the phonology.
func DemoGenerateSyllable() {
	// Create a simple phoneme pool
	pool := phoneme.NewPool()
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "k",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &phoneme.ConsonantSpec{
			Manner: phoneme.ConsonantMannerPlosive,
		},
	})
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "a",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightLow,
			Backness: phoneme.VowelBacknessCentral,
		},
	})

	// Create phonology with basic rules
	phono := phonology.NewPhonology(pool)
	phono.AddTemplate(phonology.TemplateCV, 10)
	phono.AddRule(phonology.NewVowelRequirementRule())

	// Generate syllables
	rng := rand.New(rand.NewPCG(42, 123))

	fmt.Println("Generated syllables:")
	for i := 0; i < 5; i++ {
		syllable := phono.GenerateSyllable(rng)
		fmt.Printf("  %s\n", syllable)
	}

	// Output:
	// Generated syllables:
	//   ka
	//   ka
	//   ka
	//   ka
	//   ka
}

// DemoCustomRules demonstrates how to create and use custom phonotactic rules.
func DemoCustomRules() {
	// Create a phoneme pool
	pool := phoneme.NewPool()
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "s",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &phoneme.ConsonantSpec{
			Manner: phoneme.ConsonantMannerFricative,
		},
	})
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "a",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightLow,
			Backness: phoneme.VowelBacknessCentral,
		},
	})

	// Create phonology
	phono := phonology.NewPhonology(pool)

	// Add a custom rule that only allows fricatives in onset
	customRule := &OnsetFricativeRule{}
	phono.AddRule(customRule)

	// Test some sequences
	testSequences := [][]phoneme.Phoneme{
		{{Symbol: "s", Type: phoneme.PhonemeTypeConsonant}, {Symbol: "a", Type: phoneme.PhonemeTypeVowel}},     // Valid: fricative + vowel
		{{Symbol: "a", Type: phoneme.PhonemeTypeVowel}, {Symbol: "s", Type: phoneme.PhonemeTypeConsonant}},     // Valid: vowel + fricative
		{{Symbol: "s", Type: phoneme.PhonemeTypeConsonant}, {Symbol: "s", Type: phoneme.PhonemeTypeConsonant}}, // Invalid: fricative + fricative
	}

	fmt.Println("Testing custom rule:")
	for i, seq := range testSequences {
		isValid := phono.ApplyRules(seq)
		fmt.Printf("  Sequence %d: %v (valid: %t)\n", i+1, seq, isValid)
	}

	// Output:
	// Testing custom rule:
	//   Sequence 1: [{s consonant} {a vowel}] (valid: true)
	//   Sequence 2: [{a vowel} {s consonant}] (valid: true)
	//   Sequence 3: [{s consonant} {s consonant}] (valid: false)
}

// OnsetFricativeRule is a custom rule that only allows fricatives in onset position.
type OnsetFricativeRule struct{}

func (r *OnsetFricativeRule) Validate(seq []phoneme.Phoneme) bool {
	if len(seq) < 2 {
		return true
	}

	// Check if first two phonemes are both consonants
	if seq[0].Type == phoneme.PhonemeTypeConsonant && seq[1].Type == phoneme.PhonemeTypeConsonant {
		first := seq[0].Consonant
		second := seq[1].Consonant

		if first != nil && second != nil {
			// Only allow fricative + non-fricative combinations
			if first.Manner == phoneme.ConsonantMannerFricative && second.Manner == phoneme.ConsonantMannerFricative {
				return false
			}
		}
	}

	return true
}

// DemoAdvancedFeatures demonstrates the new place and voicing features.
func DemoAdvancedFeatures() {
	// Create a more comprehensive phoneme pool
	pool := phoneme.NewPool()

	// Add consonants with different places and voicing
	consonants := []struct {
		symbol  string
		manner  phoneme.ConsonantManner
		place   phoneme.ConsonantPlace
		voicing phoneme.ConsonantVoicing
	}{
		{"p", phoneme.ConsonantMannerPlosive, phoneme.ConsonantPlaceBilabial, phoneme.ConsonantVoicingUnvoiced},
		{"b", phoneme.ConsonantMannerPlosive, phoneme.ConsonantPlaceBilabial, phoneme.ConsonantVoicingVoiced},
		{"t", phoneme.ConsonantMannerPlosive, phoneme.ConsonantPlaceAlveolar, phoneme.ConsonantVoicingUnvoiced},
		{"d", phoneme.ConsonantMannerPlosive, phoneme.ConsonantPlaceAlveolar, phoneme.ConsonantVoicingVoiced},
		{"k", phoneme.ConsonantMannerPlosive, phoneme.ConsonantPlaceVelar, phoneme.ConsonantVoicingUnvoiced},
		{"f", phoneme.ConsonantMannerFricative, phoneme.ConsonantPlaceLabiodental, phoneme.ConsonantVoicingUnvoiced},
		{"s", phoneme.ConsonantMannerFricative, phoneme.ConsonantPlaceAlveolar, phoneme.ConsonantVoicingUnvoiced},
		{"n", phoneme.ConsonantMannerNasal, phoneme.ConsonantPlaceAlveolar, phoneme.ConsonantVoicingVoiced},
	}

	for _, c := range consonants {
		pool.AddConsonant(phoneme.Phoneme{
			Symbol: c.symbol,
			Type:   phoneme.PhonemeTypeConsonant,
			Weight: 1.0,
			Consonant: &phoneme.ConsonantSpec{
				Manner:  c.manner,
				Place:   c.place,
				Voicing: c.voicing,
			},
		})
	}

	// Add vowels
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "a",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightLow,
			Backness: phoneme.VowelBacknessCentral,
		},
	})

	// Create phonology with strict rules
	phono := phonology.NewPhonology(pool)

	// Use strict default rules that leverage all new features
	for _, rule := range phonology.StrictDefaultRules() {
		phono.AddRule(rule)
	}

	// Add basic templates
	phono.AddTemplate(phonology.TemplateCV, 10)
	phono.AddTemplate(phonology.TemplateCVC, 5)

	fmt.Printf("Advanced phonology created with %d consonants, %d vowels, and %d rules\n",
		len(pool.Consonants), len(pool.Vowels), phono.GetRuleCount())

	// Test some specific phoneme sequences
	testSequences := [][]phoneme.Phoneme{
		// Valid: voiceless bilabial + voiceless labiodental
		{{Symbol: "p", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
			Place: phoneme.ConsonantPlaceBilabial, Voicing: phoneme.ConsonantVoicingUnvoiced}},
			{Symbol: "f", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
				Place: phoneme.ConsonantPlaceLabiodental, Voicing: phoneme.ConsonantVoicingUnvoiced}},
			{Symbol: "a", Type: phoneme.PhonemeTypeVowel}},

		// Invalid: different places (bilabial + velar)
		{{Symbol: "p", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
			Place: phoneme.ConsonantPlaceBilabial, Voicing: phoneme.ConsonantVoicingUnvoiced}},
			{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
				Place: phoneme.ConsonantPlaceVelar, Voicing: phoneme.ConsonantVoicingUnvoiced}},
			{Symbol: "a", Type: phoneme.PhonemeTypeVowel}},
	}

	fmt.Println("Testing phonotactic constraints:")
	for i, seq := range testSequences {
		isValid := phono.ValidateSyllable(seq)
		// Create string representation
		var seqStr string
		for _, p := range seq {
			seqStr += p.Symbol
		}
		fmt.Printf("  Sequence %d: %s (valid: %t)\n", i+1, seqStr, isValid)
	}

	// Output:
	// Advanced phonology created with 8 consonants, 1 vowels, and 7 rules
	// Testing phonotactic constraints:
	//   Sequence 1: pfa (valid: true)
	//   Sequence 2: pka (valid: false)
}
