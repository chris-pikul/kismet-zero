package lang

import "fmt"

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

func (vh VowelHeight) String() string {
	if vh > VowelHeightHigh {
		return vowelHeightEnum[0]
	}
	return vowelHeightEnum[vh]
}

func (vh *VowelHeight) UnmarshalText(src []byte) error {
	str := string(src)
	*vh = ParseVowelHeight(str)
	if *vh == VowelHeightUnknown && str != "unknown" {
		return fmt.Errorf("unknown VowelHeight '%s'", str)
	}
	return nil
}

func ParseVowelHeight(vowelHeight string) VowelHeight {
	for i, e := range vowelHeightEnum {
		if vowelHeight == e {
			return VowelHeight(i)
		}
	}
	return VowelHeightUnknown
}

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

func (vb VowelBackness) String() string {
	if vb > VowelBacknessBack {
		return vowelBacknessEnum[0]
	}
	return vowelBacknessEnum[vb]
}

func (vb *VowelBackness) UnmarshalText(src []byte) error {
	str := string(src)
	*vb = ParseVowelBackness(str)
	if *vb == VowelBacknessUnknown && str != "unknown" {
		return fmt.Errorf("unknown VowelBackness '%s'", str)
	}
	return nil
}

func ParseVowelBackness(vowelBackness string) VowelBackness {
	for i, e := range vowelBacknessEnum {
		if vowelBackness == e {
			return VowelBackness(i)
		}
	}
	return VowelBacknessUnknown
}

type VowelSpec struct {
	Height    VowelHeight   `json:"height"`
	Backness  VowelBackness `json:"backness"`
	Rounded   bool          `json:"rounded"`
	Diphthong bool          `json:"diphthong"`
}
