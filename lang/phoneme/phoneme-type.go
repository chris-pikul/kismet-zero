package phoneme

import "fmt"

// PhonemeType is an enumeration type that identifies the category of a phoneme,
// such as whether it is a "consonant" or a "vowel".
type PhonemeType byte

const (
	PhonemeTypeUnknown PhonemeType = iota
	PhonemeTypeConsonant
	PhonemeTypeVowel
)

var phonemeTypeEnum = []string{
	"unknown",
	"consonant",
	"vowel",
}

// String returns the string representation of the PhonemeType.
func (pt PhonemeType) String() string {
	if pt > PhonemeTypeVowel {
		return phonemeTypeEnum[0]
	}
	return phonemeTypeEnum[pt]
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for JSON
// deserialization of PhonemeType values.
func (pt *PhonemeType) UnmarshalText(src []byte) error {
	str := string(src)
	*pt = ParsePhonemeType(str)
	if *pt == PhonemeTypeUnknown && str != "unknown" {
		return fmt.Errorf("invalid PhonemeType '%s'", str)
	}
	return nil
}

// ParsePhonemeType converts a string representation to a PhonemeType enum value.
// Returns PhonemeTypeUnknown if the string is not recognized.
func ParsePhonemeType(phonemeType string) PhonemeType {
	for i, e := range phonemeTypeEnum {
		if e == phonemeType {
			return PhonemeType(i)
		}
	}
	return PhonemeTypeUnknown
}

// MarshalJSON implements the json.Marshaler interface.
func (pt PhonemeType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pt.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (pt *PhonemeType) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed := ParsePhonemeType(str)
	if parsed == PhonemeTypeUnknown && str != "unknown" {
		return fmt.Errorf("invalid PhonemeType '%s'", str)
	}

	*pt = parsed
	return nil
}
