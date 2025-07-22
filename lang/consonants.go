package lang

import "fmt"

// ConsonantManner is an enumerated type describing the "manner of articulation" for a
// phoneme. It is a general axis for consonants.
type ConsonantManner byte

const (
	// Error value for [ConsonantManner] indicating an invalid enum
	ConsonantMannerUnknown ConsonantManner = iota
	ConsonantMannerPlosive
	ConsonantMannerNasal
	ConsonantMannerFricative
	ConsonantMannerSibilants
	ConsonantMannerLateralFricative
	ConsonantMannerAffricate
	ConsonantMannerVibrant
	ConsonantMannerFlap
	ConsonantMannerTrill
	ConsonantMannerApproximant
	ConsonantMannerSemivowel
	ConsonantMannerDiphthong
	ConsonantMannerLateralApproximant
	ConsonantMannerEjective
	ConsonantMannerImplosive
	ConsonantMannerClick
	ConsonantMannerPercussive
)

var mannerEnums = []string{
	"unknown",
	"plosive",
	"nasal",
	"fricative",
	"sibilants",
	"lateral-fricative",
	"affricate",
	"vibrant",
	"flap",
	"trill",
	"approximant",
	"semivowel",
	"diphthong",
	"lateral-approximant",
	"ejective",
	"implosive",
	"click",
	"percussive",
}

func (m ConsonantManner) String() string {
	if m > ConsonantMannerPercussive {
		return mannerEnums[0]
	}
	return mannerEnums[m]
}

func (m *ConsonantManner) UnmarshalText(src []byte) error {
	str := string(src)
	*m = ParseConsonantManner(str)
	if *m == ConsonantMannerUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantManner '%s'", str)
	}
	return nil
}

func ParseConsonantManner(manner string) ConsonantManner {
	for i, m := range mannerEnums {
		if m == manner {
			return ConsonantManner(i)
		}
	}
	return ConsonantMannerUnknown
}

type ConsonantSpec struct {
	Manner ConsonantManner `json:"manner"`
}
