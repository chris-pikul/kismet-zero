package lang

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// SoundChange represents a phonological change that can occur in a language.
type SoundChange struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Probability float32 `json:"probability"` // 0.0 to 1.0
	Era         string  `json:"era"`         // When this change typically happens
}

// CommonSoundChanges provides a library of historically attested sound changes.
var CommonSoundChanges = []SoundChange{
	{
		ID:          "lenition_intervocalic",
		Name:        "Intervocalic Lenition",
		Description: "Stops become fricatives between vowels",
		Probability: 0.6,
		Era:         "early_evolution",
	},
	{
		ID:          "palatalization_before_front_vowels",
		Name:        "Palatalization",
		Description: "Consonants become palatalized before front vowels",
		Probability: 0.5,
		Era:         "middle_evolution",
	},
	{
		ID:          "final_consonant_deletion",
		Name:        "Final Consonant Deletion",
		Description: "Word-final consonants are lost",
		Probability: 0.4,
		Era:         "late_evolution",
	},
	{
		ID:          "vowel_merger",
		Name:        "Vowel Merger",
		Description: "Similar vowels merge into one",
		Probability: 0.3,
		Era:         "late_evolution",
	},
}

// applySoundChanges applies realistic sound changes to a language's phonology.
func applySoundChanges(language *Language, timePeriod string, seed int64) []string {
	if language.Phonology == nil {
		return nil
	}

	rng := rand.New(rand.NewPCG(uint64(seed), 0))
	var changes []string

	// Determine which era we're in based on time period
	var era string
	switch timePeriod {
	case "100_years", "short":
		era = "early_evolution"
	case "500_years", "medium":
		era = "middle_evolution"
	case "1000_years", "long":
		era = "late_evolution"
	case "2000_years", "very_long":
		era = "late_evolution"
	default:
		era = "evolution"
	}

	// Apply applicable sound changes
	for _, change := range CommonSoundChanges {
		if change.Era == era && rng.Float32() < change.Probability {
			changes = append(changes, change.Name)

			// Apply the change to the phonology
			applySpecificSoundChange(language.Phonology, change)
		}
	}

	return changes
}

// applySpecificSoundChange applies a specific sound change to the phonology.
func applySpecificSoundChange(ph *phonology.Phonology, change SoundChange) {
	// This is a simplified implementation that simulates the effect of sound changes
	// In a full implementation, this would modify the actual phoneme inventory

	switch change.ID {
	case "lenition_intervocalic":
		// Simulate lenition by adding some fricatives to the inventory
		// This is a placeholder - real implementation would be more sophisticated
	case "palatalization_before_front_vowels":
		// Simulate palatalization
		// This is a placeholder - real implementation would be more sophisticated
	case "final_consonant_deletion":
		// Simulate final consonant deletion
		// This would affect syllable templates
		// This is a placeholder - real implementation would be more sophisticated
	case "vowel_merger":
		// Simulate vowel merger
		// This is a placeholder - real implementation would be more sophisticated
	}
}

// generateRandomSoundChange creates a new sound change based on the current phonology.
func generateRandomSoundChange(language *Language, era string, seed int64) SoundChange {
	// Analyze the current phonology to generate plausible changes
	var change SoundChange

	// Generate a generic change for now
	change = SoundChange{
		ID:          fmt.Sprintf("generated_change_%d", time.Now().Unix()),
		Name:        "Phonological Drift",
		Description: "Natural phonological change over time",
		Probability: 0.3,
		Era:         era,
	}

	return change
}
