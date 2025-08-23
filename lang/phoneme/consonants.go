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

// ConsonantPlace is an enumerated type describing the "place of articulation"
// for a consonant phoneme. It represents where in the vocal tract the
// primary constriction occurs during sound production.
type ConsonantPlace byte

const (
	// Error value for ConsonantPlace indicating an invalid enum
	ConsonantPlaceUnknown ConsonantPlace = iota
	ConsonantPlaceBilabial
	ConsonantPlaceLabiodental
	ConsonantPlaceDental
	ConsonantPlaceAlveolar
	ConsonantPlacePostalveolar
	ConsonantPlaceRetroflex
	ConsonantPlacePalatal
	ConsonantPlaceVelar
	ConsonantPlaceUvular
	ConsonantPlacePharyngeal
	ConsonantPlaceGlottal
)

var placeEnums = []string{
	"unknown",
	"bilabial",
	"labiodental",
	"dental",
	"alveolar",
	"postalveolar",
	"retroflex",
	"palatal",
	"velar",
	"uvular",
	"pharyngeal",
	"glottal",
}

// String returns the string representation of the ConsonantPlace.
func (p ConsonantPlace) String() string {
	if p > ConsonantPlaceGlottal {
		return placeEnums[0]
	}
	return placeEnums[p]
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for JSON
// deserialization of ConsonantPlace values.
func (p *ConsonantPlace) UnmarshalText(src []byte) error {
	str := string(src)
	*p = ParseConsonantPlace(str)
	if *p == ConsonantPlaceUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantPlace '%s'", str)
	}
	return nil
}

// ParseConsonantPlace converts a string representation to a ConsonantPlace
// enum value. Returns ConsonantPlaceUnknown if the string is not recognized.
func ParseConsonantPlace(place string) ConsonantPlace {
	for i, p := range placeEnums {
		if p == place {
			return ConsonantPlace(i)
		}
	}
	return ConsonantPlaceUnknown
}

// MarshalJSON implements the json.Marshaler interface.
func (cp ConsonantPlace) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cp.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (cp *ConsonantPlace) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed := ParseConsonantPlace(str)
	if parsed == ConsonantPlaceUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantPlace '%s'", str)
	}

	*cp = parsed
	return nil
}

// ConsonantVoicing is an enumerated type describing the voicing state
// of a consonant phoneme. It represents whether the vocal cords vibrate
// during sound production.
type ConsonantVoicing byte

const (
	// Error value for ConsonantVoicing indicating an invalid enum
	ConsonantVoicingUnknown ConsonantVoicing = iota
	ConsonantVoicingVoiced
	ConsonantVoicingUnvoiced
)

var voicingEnums = []string{
	"unknown",
	"voiced",
	"unvoiced",
}

// String returns the string representation of the ConsonantVoicing.
func (v ConsonantVoicing) String() string {
	if v > ConsonantVoicingUnvoiced {
		return voicingEnums[0]
	}
	return voicingEnums[v]
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for JSON
// deserialization of ConsonantVoicing values.
func (v *ConsonantVoicing) UnmarshalText(src []byte) error {
	str := string(src)
	*v = ParseConsonantVoicing(str)
	if *v == ConsonantVoicingUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantVoicing '%s'", str)
	}
	return nil
}

// ParseConsonantVoicing converts a string representation to a ConsonantVoicing
// enum value. Returns ConsonantVoicingUnknown if the string is not recognized.
func ParseConsonantVoicing(voicing string) ConsonantVoicing {
	for i, v := range voicingEnums {
		if v == voicing {
			return ConsonantVoicing(i)
		}
	}
	return ConsonantVoicingUnknown
}

// MarshalJSON implements the json.Marshaler interface.
func (cv ConsonantVoicing) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cv.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (cv *ConsonantVoicing) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed := ParseConsonantVoicing(str)
	if parsed == ConsonantVoicingUnknown && str != "unknown" {
		return fmt.Errorf("invalid ConsonantVoicing '%s'", str)
	}

	*cv = parsed
	return nil
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
	Manner  ConsonantManner  `json:"manner"`
	Place   ConsonantPlace   `json:"place"`
	Voicing ConsonantVoicing `json:"voicing"`
}
