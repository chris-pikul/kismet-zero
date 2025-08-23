package orthography

import (
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// GraphemeType is an enumeration type that identifies the category of a grapheme,
// such as whether it represents a letter, syllable, or logogram.
type GraphemeType byte

const (
	GraphemeTypeUnknown GraphemeType = iota
	GraphemeTypeLetter
	GraphemeTypeSyllable
	GraphemeTypeLogogram
)

var graphemeTypeEnum = []string{
	"unknown",
	"letter",
	"syllable",
	"logogram",
}

// String returns the string representation of the GraphemeType.
func (gt GraphemeType) String() string {
	if gt > GraphemeTypeLogogram {
		return graphemeTypeEnum[0]
	}
	return graphemeTypeEnum[gt]
}

// Grapheme represents a single written character or symbol that maps to one or more phonemes.
type Grapheme struct {
	Symbol     string           `json:"symbol"`
	Type       GraphemeType     `json:"type"`
	PhonemeRef *phoneme.Phoneme `json:"phonemeRef,omitempty"`
	Weight     float32          `json:"weight"`
}

// WritingStyle is an enumeration type that identifies the style of writing system.
type WritingStyle byte

const (
	WritingStyleUnknown WritingStyle = iota
	WritingStyleAlphabetic
	WritingStyleSyllabic
	WritingStyleLogographic
)

var writingStyleEnum = []string{
	"unknown",
	"alphabetic",
	"syllabic",
	"logographic",
}

// String returns the string representation of the WritingStyle.
func (ws WritingStyle) String() string {
	if ws > WritingStyleLogographic {
		return writingStyleEnum[0]
	}
	return writingStyleEnum[ws]
}

// OrthographyMapping represents a mapping from a phoneme to one or more graphemes.
type OrthographyMapping struct {
	Phoneme   *phoneme.Phoneme `json:"phoneme"`
	Graphemes []Grapheme       `json:"graphemes"`
	Primary   bool             `json:"primary"`
}

// WritingSystem defines a complete writing system with style, graphemes, and mappings.
type WritingSystem struct {
	Style     WritingStyle         `json:"style"`
	Graphemes []Grapheme           `json:"graphemes"`
	Mappings  []OrthographyMapping `json:"mappings"`
	Name      string               `json:"name"`
	Culture   string               `json:"culture,omitempty"`
}

// GraphemeList is a collection of graphemes that supports both weighted and
// unweighted random selection.
type GraphemeList []Grapheme

// UnweightedChoice selects a random grapheme from the list with equal probability
// for each grapheme, ignoring their weight values.
func (gl GraphemeList) UnweightedChoice(rng *rand.Rand) *Grapheme {
	if len(gl) == 0 {
		return nil
	}
	return &gl[rng.IntN(len(gl))]
}

// WeightedChoice selects a random grapheme from the list using their weight
// values to determine selection probability. Higher weights increase the chance
// of selection.
func (gl GraphemeList) WeightedChoice(rng *rand.Rand) *Grapheme {
	if len(gl) == 0 {
		return nil
	}

	var total float32
	for _, g := range gl {
		total += g.Weight
	}

	if total == 0 {
		panic("GraphemeList with zero total weight cannot be sampled from")
	}

	target := rng.Float32() * total
	var cumulative float32
	for i := range gl {
		cumulative += gl[i].Weight
		if cumulative >= target {
			return &gl[i]
		}
	}

	// Edge case fallback — shouldn't normally happen
	return &gl[len(gl)-1]
}
