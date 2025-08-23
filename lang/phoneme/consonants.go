package phoneme

import "fmt"

// ConsonantManner is an enumerated type describing the "manner of articulation"
// for a consonant phoneme. It represents how the airflow is modified during
// sound production.
type ConsonantManner byte

const (
	// Error value for ConsonantManner indicating an invalid enum
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

// String returns the string representation of the ConsonantManner.
func (m ConsonantManner) String() string {
	if m > ConsonantMannerPercussive {
		return mannerEnums[0]
	}
	return mannerEnums[m]
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for JSON
// deserialization of ConsonantManner values.
func (m *ConsonantManner) UnmarshalText(src []byte) error {
	str := string(src)
	*m = ParseConsonantManner(str)
	if *m == ConsonantMannerUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantManner '%s'", str)
	}
	return nil
}

// ParseConsonantManner converts a string representation to a ConsonantManner
// enum value. Returns ConsonantMannerUnknown if the string is not recognized.
func ParseConsonantManner(manner string) ConsonantManner {
	for i, m := range mannerEnums {
		if m == manner {
			return ConsonantManner(i)
		}
	}
	return ConsonantMannerUnknown
}

// MarshalJSON implements the json.Marshaler interface.
func (cm ConsonantManner) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cm.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (cm *ConsonantManner) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed := ParseConsonantManner(str)
	if parsed == ConsonantMannerUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantManner '%s'", str)
	}

	*cm = parsed
	return nil
}

// ConsonantSpec describes the articulatory features of a consonant phoneme.
type ConsonantSpec struct {
	Manner ConsonantManner `json:"manner"`
}
