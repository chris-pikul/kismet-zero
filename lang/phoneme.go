package lang

type Phoneme struct {
	Symbol string      `json:"symbol"`
	Type   PhonemeType `json:"type"`
	Weight float32     `json:"weight"`

	Consonant *ConsonantSpec `json:"consonant,omitempty"`
	Vowel     *VowelSpec     `json:"vowel,omitempty"`
}
