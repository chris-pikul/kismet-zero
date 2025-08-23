package phoneme

// PhonemePool represents a collection of phonemes organized by type,
// providing the foundation for language-specific phonology generation.
type PhonemePool struct {
	Consonants PhonemeList `json:"consonants"`
	Vowels     PhonemeList `json:"vowels"`
}

// NewPool creates a new empty PhonemePool.
func NewPool() *PhonemePool {
	return &PhonemePool{
		Consonants: make(PhonemeList, 0),
		Vowels:     make(PhonemeList, 0),
	}
}

// AddConsonant adds a consonant phoneme to the pool, ensuring no duplicates
// based on symbol.
func (p *PhonemePool) AddConsonant(phoneme Phoneme) {
	if phoneme.Type != PhonemeTypeConsonant {
		return
	}

	// Check for duplicates
	for _, existing := range p.Consonants {
		if existing.Symbol == phoneme.Symbol {
			return
		}
	}

	p.Consonants = append(p.Consonants, phoneme)
}

// AddVowel adds a vowel phoneme to the pool, ensuring no duplicates
// based on symbol.
func (p *PhonemePool) AddVowel(phoneme Phoneme) {
	if phoneme.Type != PhonemeTypeVowel {
		return
	}

	// Check for duplicates
	for _, existing := range p.Vowels {
		if existing.Symbol == phoneme.Symbol {
			return
		}
	}

	p.Vowels = append(p.Vowels, phoneme)
}

// GetConsonantBySymbol retrieves a consonant phoneme by its symbol.
// Returns nil if not found.
func (p *PhonemePool) GetConsonantBySymbol(symbol string) *Phoneme {
	for _, consonant := range p.Consonants {
		if consonant.Symbol == symbol {
			return &consonant
		}
	}
	return nil
}

// GetVowelBySymbol retrieves a vowel phoneme by its symbol.
// Returns nil if not found.
func (p *PhonemePool) GetVowelBySymbol(symbol string) *Phoneme {
	for _, vowel := range p.Vowels {
		if vowel.Symbol == symbol {
			return &vowel
		}
	}
	return nil
}

// Size returns the total number of phonemes in the pool.
func (p *PhonemePool) Size() int {
	return len(p.Consonants) + len(p.Vowels)
}
