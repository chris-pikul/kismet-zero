package lang

import "fmt"

// PhonemeType is a enumeration type currently to support the type of a phoneme,
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

func (pt PhonemeType) String() string {
	if pt > PhonemeTypeVowel {
		return phonemeTypeEnum[0]
	}
	return phonemeTypeEnum[pt]
}

func (pt *PhonemeType) UnmarshalText(src []byte) error {
	str := string(src)
	*pt = ParsePhonemeType(str)
	if *pt == PhonemeTypeUnknown && str != "unknown" {
		return fmt.Errorf("invalid PhonemeType '%s'", str)
	}
	return nil
}

func ParsePhonemeType(phonemeType string) PhonemeType {
	for i, e := range phonemeTypeEnum {
		if e == phonemeType {
			return PhonemeType(i)
		}
	}
	return PhonemeTypeUnknown
}
