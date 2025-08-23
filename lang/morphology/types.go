package morphology

import (
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// MorphemeType is an enumeration type that identifies the category of a morpheme,
// such as whether it's a root, prefix, suffix, or other affix type.
type MorphemeType byte

const (
	MorphemeTypeUnknown MorphemeType = iota
	MorphemeTypeRoot
	MorphemeTypePrefix
	MorphemeTypeSuffix
	MorphemeTypeInfix
	MorphemeTypeCircumfix
)

var morphemeTypeEnum = []string{
	"unknown",
	"root",
	"prefix",
	"suffix",
	"infix",
	"circumfix",
}

// String returns the string representation of the MorphemeType.
func (mt MorphemeType) String() string {
	if mt > MorphemeTypeCircumfix {
		return morphemeTypeEnum[0]
	}
	return morphemeTypeEnum[mt]
}

// MorphemeFrequency represents how common a morpheme is in a language.
type MorphemeFrequency byte

const (
	MorphemeFrequencyUnknown MorphemeFrequency = iota
	MorphemeFrequencyRare
	MorphemeFrequencyUncommon
	MorphemeFrequencyCommon
	MorphemeFrequencyVeryCommon
)

var morphemeFrequencyEnum = []string{
	"unknown",
	"rare",
	"uncommon",
	"common",
	"very_common",
}

// String returns the string representation of the MorphemeFrequency.
func (mf MorphemeFrequency) String() string {
	if mf > MorphemeFrequencyVeryCommon {
		return morphemeFrequencyEnum[0]
	}
	return morphemeFrequencyEnum[mf]
}

// Morpheme represents the smallest meaningful unit of language with its
// phonological form, meaning, and morphological properties.
type Morpheme struct {
	ID        string            `json:"id"`
	Type      MorphemeType      `json:"type"`
	Meaning   string            `json:"meaning"`
	Phonemes  []phoneme.Phoneme `json:"phonemes"`
	Syllables []string          `json:"syllables"` // Syllable strings from phonology
	Weight    float32           `json:"weight"`
	Culture   string            `json:"culture,omitempty"`
	Frequency MorphemeFrequency `json:"frequency"`
}

// MorphemeList is a collection of morphemes that supports both weighted and
// unweighted random selection.
type MorphemeList []Morpheme

// UnweightedChoice selects a random morpheme from the list with equal probability
// for each morpheme, ignoring their weight values.
func (ml MorphemeList) UnweightedChoice(rng *rand.Rand) *Morpheme {
	if len(ml) == 0 {
		return nil
	}
	return &ml[rng.IntN(len(ml))]
}

// WeightedChoice selects a random morpheme from the list using their weight
// values to determine selection probability. Higher weights increase the chance
// of selection.
func (ml MorphemeList) WeightedChoice(rng *rand.Rand) *Morpheme {
	if len(ml) == 0 {
		return nil
	}

	var total float32
	for _, m := range ml {
		total += m.Weight
	}

	if total == 0 {
		panic("MorphemeList with zero total weight cannot be sampled from")
	}

	target := rng.Float32() * total
	var cumulative float32
	for i := range ml {
		cumulative += ml[i].Weight
		if cumulative >= target {
			return &ml[i]
		}
	}

	// Edge case fallback — shouldn't normally happen
	return &ml[len(ml)-1]
}

// WordCategory represents the grammatical category of a word.
type WordCategory byte

const (
	WordCategoryUnknown WordCategory = iota
	WordCategoryNoun
	WordCategoryVerb
	WordCategoryAdjective
	WordCategoryAdverb
	WordCategoryPronoun
	WordCategoryPreposition
	WordCategoryConjunction
	WordCategoryInterjection
)

var wordCategoryEnum = []string{
	"unknown",
	"noun",
	"verb",
	"adjective",
	"adverb",
	"pronoun",
	"preposition",
	"conjunction",
	"interjection",
}

// String returns the string representation of the WordCategory.
func (wc WordCategory) String() string {
	if wc > WordCategoryInterjection {
		return wordCategoryEnum[0]
	}
	return wordCategoryEnum[wc]
}

// WordFrequency represents how common a word is in a language.
type WordFrequency byte

const (
	WordFrequencyUnknown WordFrequency = iota
	WordFrequencyRare
	WordFrequencyUncommon
	WordFrequencyCommon
	WordFrequencyVeryCommon
)

var wordFrequencyEnum = []string{
	"unknown",
	"rare",
	"uncommon",
	"common",
	"very_common",
}

// String returns the string representation of the WordFrequency.
func (wf WordFrequency) String() string {
	if wf > WordFrequencyVeryCommon {
		return wordFrequencyEnum[0]
	}
	return wordFrequencyEnum[wf]
}

// Word represents a complete word formed from one or more morphemes with
// its phonological form, meaning, and grammatical properties.
type Word struct {
	ID        string            `json:"id"`
	Morphemes []Morpheme        `json:"morphemes"`
	Meaning   string            `json:"meaning"`
	Phonemes  []phoneme.Phoneme `json:"phonemes"`
	Written   string            `json:"written"`
	Category  WordCategory      `json:"category"`
	Culture   string            `json:"culture,omitempty"`
	Frequency WordFrequency     `json:"frequency"`
	Weight    float32           `json:"weight"`
}

// WordList is a collection of words that supports both weighted and
// unweighted random selection.
type WordList []Word

// UnweightedChoice selects a random word from the list with equal probability
// for each word, ignoring their weight values.
func (wl WordList) UnweightedChoice(rng *rand.Rand) *Word {
	if len(wl) == 0 {
		return nil
	}
	return &wl[rng.IntN(len(wl))]
}

// WeightedChoice selects a random word from the list using their weight
// values to determine selection probability. Higher weights increase the chance
// of selection.
func (wl WordList) WeightedChoice(rng *rand.Rand) *Word {
	if len(wl) == 0 {
		return nil
	}

	var total float32
	for _, w := range wl {
		total += w.Weight
	}

	if total == 0 {
		panic("WordList with zero total weight cannot be sampled from")
	}

	target := rng.Float32() * total
	var cumulative float32
	for i := range wl {
		cumulative += wl[i].Weight
		if cumulative >= target {
			return &wl[i]
		}
	}

	// Edge case fallback — shouldn't normally happen
	return &wl[len(wl)-1]
}

// NameType represents the category of a generated name.
type NameType byte

const (
	NameTypeUnknown NameType = iota
	NameTypePersonal
	NameTypeFamily
	NameTypePlace
	NameTypeTitle
	NameTypeDeity
	NameTypeArtifact
)

var nameTypeEnum = []string{
	"unknown",
	"personal",
	"family",
	"place",
	"title",
	"deity",
	"artifact",
}

// String returns the string representation of the NameType.
func (nt NameType) String() string {
	if nt > NameTypeArtifact {
		return nameTypeEnum[0]
	}
	return nameTypeEnum[nt]
}

// Name represents a culturally appropriate name for people, places, or concepts.
type Name struct {
	ID        string   `json:"id"`
	Type      NameType `json:"type"`
	Value     string   `json:"value"`
	Meaning   string   `json:"meaning"`
	Culture   string   `json:"culture"`
	Gender    string   `json:"gender,omitempty"`
	Formality string   `json:"formality,omitempty"`
}
