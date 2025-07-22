package lang

type PhonemePool struct {
	Consonants PhonemeList `json:"consonants"`
	Vowels     PhonemeList `json:"vowels"`
}
