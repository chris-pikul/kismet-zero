package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// SoundChangeRule defines a systematic phonological change that can occur in a language.
type SoundChangeRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Change parameters
	FromPhonemes []phoneme.Phoneme `json:"fromPhonemes"` // Phonemes that change
	ToPhonemes   []phoneme.Phoneme `json:"toPhonemes"`   // What they become
	Environment  string            `json:"environment"`  // Context where change occurs

	// Probability and timing
	Probability float32 `json:"probability"`   // 0.0 to 1.0, chance of occurring
	Era         string  `json:"era,omitempty"` // When this change typically happens

	// Examples from historical linguistics
	HistoricalExample string `json:"historicalExample,omitempty"`
}

// CommonSoundChangeRules provides a library of historically attested sound changes.
var CommonSoundChangeRules = []SoundChangeRule{
	{
		ID:          "grimm_law_stops",
		Name:        "Grimm's Law - Stop Consonant Shift",
		Description: "Voiced stops become voiceless, voiceless stops become fricatives",
		FromPhonemes: []phoneme.Phoneme{
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceBilabial,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceAlveolar,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceVelar,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
		},
		ToPhonemes: []phoneme.Phoneme{
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceBilabial,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceAlveolar,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceVelar,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
		},
		Environment:       "word-initial or stressed syllable",
		Probability:       0.6,
		Era:               "early_evolution",
		HistoricalExample: "PIE *bʰ > Germanic *b, PIE *dʰ > Germanic *d",
	},
	{
		ID:          "lenition_intervocalic",
		Name:        "Intervocalic Lenition",
		Description: "Stops become fricatives between vowels",
		FromPhonemes: []phoneme.Phoneme{
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
		},
		ToPhonemes: []phoneme.Phoneme{
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerFricative,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerFricative,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
		},
		Environment:       "between vowels",
		Probability:       0.7,
		Era:               "middle_evolution",
		HistoricalExample: "Latin 'civitas' > Spanish 'ciudad'",
	},
	{
		ID:          "palatalization_front_vowels",
		Name:        "Palatalization Before Front Vowels",
		Description: "Velar consonants become palatal before front vowels",
		FromPhonemes: []phoneme.Phoneme{
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceVelar,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlaceVelar,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
		},
		ToPhonemes: []phoneme.Phoneme{
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlacePalatal,
					Voicing: phoneme.ConsonantVoicingUnvoiced,
				},
			},
			{
				Type: phoneme.PhonemeTypeConsonant,
				Consonant: &phoneme.ConsonantSpec{
					Manner:  phoneme.ConsonantMannerPlosive,
					Place:   phoneme.ConsonantPlacePalatal,
					Voicing: phoneme.ConsonantVoicingVoiced,
				},
			},
		},
		Environment:       "before front vowels",
		Probability:       0.5,
		Era:               "middle_evolution",
		HistoricalExample: "Latin 'centum' > Italian 'cento'",
	},
}

// SoundChangeEngine manages the application of sound change rules to a phonology.
type SoundChangeEngine struct {
	BaseEngine
	rules []SoundChangeRule
}

// NewSoundChangeEngine creates a new sound change engine with the given configuration.
func NewSoundChangeEngine(config EvolutionConfig) *SoundChangeEngine {
	return &SoundChangeEngine{
		BaseEngine: NewBaseEngine(config),
		rules:      CommonSoundChangeRules,
	}
}

// AddRule adds a custom sound change rule to the engine.
func (sce *SoundChangeEngine) AddRule(rule SoundChangeRule) {
	sce.rules = append(sce.rules, rule)
}

// ApplySoundChanges applies applicable sound change rules to a phonology.
// Returns a list of changes that were applied and the modified phonology.
func (sce *SoundChangeEngine) ApplySoundChanges(ph *phonology.Phonology, era string) ([]LinguisticChange, *phonology.Phonology) {
	var changes []LinguisticChange
	modifiedPhonology := sce.clonePhonology(ph)

	for _, rule := range sce.rules {
		if rule.Era != "" && rule.Era != era {
			continue // Skip rules not applicable to this era
		}

		if sce.GetRandomFloat32() < rule.Probability {
			change := sce.applyRule(rule, modifiedPhonology, era)
			if change != nil {
				changes = append(changes, *change)
			}
		}
	}

	return changes, modifiedPhonology
}

// applyRule applies a single sound change rule to the phonology.
func (sce *SoundChangeEngine) applyRule(rule SoundChangeRule, ph *phonology.Phonology, era string) *LinguisticChange {
	// This is a simplified implementation - in practice, you'd need to:
	// 1. Find all instances of the target phonemes in the language's lexicon
	// 2. Apply the change according to the environment rules
	// 3. Update the phoneme inventory if new phonemes are introduced

	change := &LinguisticChange{
		ID:          fmt.Sprintf("sound_change_%s_%d", rule.ID, time.Now().Unix()),
		Type:        ChangeTypeSoundShift,
		Direction:   ChangeDirectionModifying,
		Description: fmt.Sprintf("Applied %s: %s", rule.Name, rule.Description),
		Details:     rule.HistoricalExample,
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     "natural_evolution",
		Intensity:   0.5,
	}

	// Mark affected phonemes
	for _, phoneme := range rule.FromPhonemes {
		change.AffectedPhonemes = append(change.AffectedPhonemes, phoneme.Symbol)
	}

	return change
}

// clonePhonology creates a deep copy of the phonology for modification.
func (sce *SoundChangeEngine) clonePhonology(ph *phonology.Phonology) *phonology.Phonology {
	// This is a simplified clone - in practice, you'd need to deep copy
	// all phonology components including the phoneme pool, templates, and rules
	return ph
}
