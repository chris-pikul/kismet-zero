package phoneme

import "fmt"

// VowelHeight describes the vertical position of the tongue during vowel
// articulation, ranging from high (close) to low (open).
type VowelHeight byte

const (
	VowelHeightUnknown VowelHeight = iota
	VowelHeightLow
	VowelHeightMid
	VowelHeightHigh
)

var vowelHeightEnum = []string{
	"unknown",
	"low",
	"mid",
	"high",
}

// String returns the string representation of the VowelHeight.
func (vh VowelHeight) String() string {
	if vh > VowelHeightHigh {
		return vowelHeightEnum[0]
	}
	return vowelHeightEnum[vh]
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for JSON
// deserialization of VowelHeight values.
func (vh *VowelHeight) UnmarshalText(src []byte) error {
	str := string(src)
	*vh = ParseVowelHeight(str)
	if *vh == VowelHeightUnknown && str != "unknown" {
		return fmt.Errorf("unknown VowelHeight '%s'", str)
	}
	return nil
}

// ParseVowelHeight converts a string representation to a VowelHeight enum value.
// Returns VowelHeightUnknown if the string is not recognized.
func ParseVowelHeight(vowelHeight string) VowelHeight {
	for i, e := range vowelHeightEnum {
		if e == vowelHeight {
			return VowelHeight(i)
		}
	}
	return VowelHeightUnknown
}

// MarshalJSON implements the json.Marshaler interface.
func (vh VowelHeight) MarshalJSON() ([]byte, error) {
	return []byte(`"` + vh.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (vh *VowelHeight) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed := ParseVowelHeight(str)
	if parsed == VowelHeightUnknown && str != "unknown" {
		return fmt.Errorf("invalid VowelHeight '%s'", str)
	}

	*vh = parsed
	return nil
}

// VowelBackness describes the horizontal position of the tongue during vowel
// articulation, ranging from front to back.
type VowelBackness byte

const (
	VowelBacknessUnknown VowelBackness = iota
	VowelBacknessFront
	VowelBacknessCentral
	VowelBacknessBack
)

var vowelBacknessEnum = []string{
	"unknown",
	"front",
	"central",
	"back",
}

// String returns the string representation of the VowelBackness.
func (vb VowelBackness) String() string {
	if vb > VowelBacknessBack {
		return vowelBacknessEnum[0]
	}
	return vowelBacknessEnum[vb]
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for JSON
// deserialization of VowelBackness values.
func (vb *VowelBackness) UnmarshalText(src []byte) error {
	str := string(src)
	*vb = ParseVowelBackness(str)
	if *vb == VowelBacknessUnknown && str != "unknown" {
		return fmt.Errorf("unknown VowelBackness '%s'", str)
	}
	return nil
}

// ParseVowelBackness converts a string representation to a VowelBackness
// enum value. Returns VowelBacknessUnknown if the string is not recognized.
func ParseVowelBackness(vowelBackness string) VowelBackness {
	for i, e := range vowelBacknessEnum {
		if e == vowelBackness {
			return VowelBackness(i)
		}
	}
	return VowelBacknessUnknown
}

// MarshalJSON implements the json.Marshaler interface.
func (vb VowelBackness) MarshalJSON() ([]byte, error) {
	return []byte(`"` + vb.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (vb *VowelBackness) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed := ParseVowelBackness(str)
	if parsed == VowelBacknessUnknown && str != "unknown" {
		return fmt.Errorf("invalid VowelBackness '%s'", str)
	}

	*vb = parsed
	return nil
}

// VowelSpec describes the articulatory features of a vowel phoneme.
type VowelSpec struct {
	Height    VowelHeight   `json:"height"`
	Backness  VowelBackness `json:"backness"`
	Rounded   bool          `json:"rounded"`
	Diphthong bool          `json:"diphthong"`
}
