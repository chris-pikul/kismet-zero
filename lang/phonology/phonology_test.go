package phonology

import (
	"math/rand/v2"
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestNewSyllableTemplate(t *testing.T) {
	onset := SyllableSlot{Required: true, MaxSize: 2}
	nucleus := SyllableSlot{Required: true, MaxSize: 1}
	coda := SyllableSlot{Required: false, MaxSize: 0}
	weight := 5

	template := NewSyllableTemplate(onset, nucleus, coda, weight)

	if template.Onset != onset {
		t.Errorf("Expected onset %v, got %v", onset, template.Onset)
	}
	if template.Nucleus != nucleus {
		t.Errorf("Expected nucleus %v, got %v", nucleus, template.Nucleus)
	}
	if template.Coda != coda {
		t.Errorf("Expected coda %v, got %v", coda, template.Coda)
	}
	if template.Weight != weight {
		t.Errorf("Expected weight %d, got %d", weight, template.Weight)
	}
}

func TestCommonTemplates(t *testing.T) {
	testCases := []struct {
		name     string
		template SyllableTemplate
		expected struct {
			onsetRequired   bool
			onsetMaxSize    int
			nucleusRequired bool
			nucleusMaxSize  int
			codaRequired    bool
			codaMaxSize     int
			weight          int
		}
	}{
		{
			name:     "TemplateCV",
			template: TemplateCV,
			expected: struct {
				onsetRequired   bool
				onsetMaxSize    int
				nucleusRequired bool
				nucleusMaxSize  int
				codaRequired    bool
				codaMaxSize     int
				weight          int
			}{
				onsetRequired:   true,
				onsetMaxSize:    1,
				nucleusRequired: true,
				nucleusMaxSize:  1,
				codaRequired:    false,
				codaMaxSize:     0,
				weight:          10,
			},
		},
		{
			name:     "TemplateCVC",
			template: TemplateCVC,
			expected: struct {
				onsetRequired   bool
				onsetMaxSize    int
				nucleusRequired bool
				nucleusMaxSize  int
				codaRequired    bool
				codaMaxSize     int
				weight          int
			}{
				onsetRequired:   true,
				onsetMaxSize:    1,
				nucleusRequired: true,
				nucleusMaxSize:  1,
				codaRequired:    true,
				codaMaxSize:     1,
				weight:          8,
			},
		},
		{
			name:     "TemplateV",
			template: TemplateV,
			expected: struct {
				onsetRequired   bool
				onsetMaxSize    int
				nucleusRequired bool
				nucleusMaxSize  int
				codaRequired    bool
				codaMaxSize     int
				weight          int
			}{
				onsetRequired:   false,
				onsetMaxSize:    0,
				nucleusRequired: true,
				nucleusMaxSize:  1,
				codaRequired:    false,
				codaMaxSize:     0,
				weight:          3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.template.Onset.Required != tc.expected.onsetRequired {
				t.Errorf("Expected onset required %v, got %v", tc.expected.onsetRequired, tc.template.Onset.Required)
			}
			if tc.template.Onset.MaxSize != tc.expected.onsetMaxSize {
				t.Errorf("Expected onset max size %d, got %d", tc.expected.onsetMaxSize, tc.template.Onset.MaxSize)
			}
			if tc.template.Nucleus.Required != tc.expected.nucleusRequired {
				t.Errorf("Expected nucleus required %v, got %v", tc.expected.nucleusRequired, tc.template.Nucleus.Required)
			}
			if tc.template.Nucleus.MaxSize != tc.expected.nucleusMaxSize {
				t.Errorf("Expected nucleus max size %d, got %d", tc.expected.nucleusMaxSize, tc.template.Nucleus.MaxSize)
			}
			if tc.template.Coda.Required != tc.expected.codaRequired {
				t.Errorf("Expected coda required %v, got %v", tc.expected.codaRequired, tc.template.Coda.Required)
			}
			if tc.template.Coda.MaxSize != tc.expected.codaMaxSize {
				t.Errorf("Expected coda max size %d, got %d", tc.expected.codaMaxSize, tc.template.Coda.MaxSize)
			}
			if tc.template.Weight != tc.expected.weight {
				t.Errorf("Expected weight %d, got %d", tc.expected.weight, tc.template.Weight)
			}
		})
	}
}

func TestNewPhonology(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	if phonology.pool != pool {
		t.Errorf("Expected pool %v, got %v", pool, phonology.pool)
	}
	if len(phonology.templates) != 0 {
		t.Errorf("Expected empty templates, got %d", len(phonology.templates))
	}
	if len(phonology.rules) != 0 {
		t.Errorf("Expected empty rules, got %d", len(phonology.rules))
	}
}

func TestAddTemplate(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	template := SyllableTemplate{Weight: 5}
	phonology.AddTemplate(template, 10)

	if len(phonology.templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(phonology.templates))
	}

	addedTemplate := phonology.templates[0]
	if addedTemplate.Weight != 10 {
		t.Errorf("Expected weight 10, got %d", addedTemplate.Weight)
	}
}

func TestPickTemplate(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	// Test with no templates - should return default
	rng := rand.New(rand.NewPCG(1, 2))
	template := phonology.PickTemplate(rng)
	if template != TemplateCV {
		t.Errorf("Expected default TemplateCV, got %v", template)
	}

	// Add templates and test weighted selection
	phonology.AddTemplate(TemplateCV, 10)
	phonology.AddTemplate(TemplateCVC, 5)

	// Test multiple picks to ensure randomness
	templates := make(map[SyllableTemplate]int)
	for i := 0; i < 1000; i++ {
		picked := phonology.PickTemplate(rng)
		templates[picked]++
	}

	// Should have picked both templates
	if len(templates) < 2 {
		t.Errorf("Expected at least 2 different templates, got %d", len(templates))
	}
}

func TestAddRule(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	rule := NewMaxClusterSizeRule(3)
	phonology.AddRule(rule)

	if len(phonology.rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(phonology.rules))
	}
}

func TestApplyRules(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	// Add a simple rule
	rule := NewMaxClusterSizeRule(2)
	phonology.AddRule(rule)

	// Test valid sequence
	validSeq := []phoneme.Phoneme{
		{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
		{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
	}
	if !phonology.ApplyRules(validSeq) {
		t.Error("Expected valid sequence to pass rules")
	}

	// Test invalid sequence (too many consonants)
	invalidSeq := []phoneme.Phoneme{
		{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
		{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
		{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
		{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
	}
	if phonology.ApplyRules(invalidSeq) {
		t.Error("Expected invalid sequence to fail rules")
	}
}

func TestValidateSyllable(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name:    "Empty sequence",
			seq:     []phoneme.Phoneme{},
			isValid: false,
		},
		{
			name: "Single vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Single consonant",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: false,
		},
		{
			name: "CV syllable",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "CVC syllable",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := phonology.ValidateSyllable(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestGenerateSyllable(t *testing.T) {
	pool := phoneme.NewPool()

	// Add some test phonemes
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "k",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &phoneme.ConsonantSpec{
			Manner: phoneme.ConsonantMannerPlosive,
		},
	})
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "a",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightLow,
			Backness: phoneme.VowelBacknessCentral,
		},
	})

	phonology := NewPhonology(pool)
	phonology.AddTemplate(TemplateCV, 10)

	rng := rand.New(rand.NewPCG(1, 2))
	syllable := phonology.GenerateSyllable(rng)

	if syllable == "" {
		t.Error("Expected non-empty syllable")
	}

	// Should contain both consonant and vowel
	if len(syllable) < 2 {
		t.Errorf("Expected syllable length >= 2, got %d", len(syllable))
	}
}

func TestFillTemplate(t *testing.T) {
	pool := phoneme.NewPool()

	// Add test phonemes
	pool.AddConsonant(phoneme.Phoneme{
		Symbol: "k",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &phoneme.ConsonantSpec{
			Manner: phoneme.ConsonantMannerPlosive,
		},
	})
	pool.AddVowel(phoneme.Phoneme{
		Symbol: "a",
		Type:   phoneme.PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &phoneme.VowelSpec{
			Height:   phoneme.VowelHeightLow,
			Backness: phoneme.VowelBacknessCentral,
		},
	})

	phonology := NewPhonology(pool)
	rng := rand.New(rand.NewPCG(1, 2))

	template := TemplateCV
	result := phonology.fillTemplate(template, rng)

	if len(result) != 2 {
		t.Errorf("Expected 2 phonemes, got %d", len(result))
	}

	if result[0].Type != phoneme.PhonemeTypeConsonant {
		t.Errorf("Expected first phoneme to be consonant, got %v", result[0].Type)
	}

	if result[1].Type != phoneme.PhonemeTypeVowel {
		t.Errorf("Expected second phoneme to be vowel, got %v", result[1].Type)
	}
}

func TestPhonemesToString(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	phonemes := []phoneme.Phoneme{
		{Symbol: "k"},
		{Symbol: "a"},
		{Symbol: "t"},
	}

	result := phonology.phonemesToString(phonemes)
	expected := "kat"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGetPool(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	retrievedPool := phonology.GetPool()
	if retrievedPool != pool {
		t.Errorf("Expected pool %v, got %v", pool, retrievedPool)
	}
}

func TestGetTemplates(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	phonology.AddTemplate(TemplateCV, 10)
	templates := phonology.GetTemplates()

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	// Test that it's a copy, not a reference
	templates[0].Weight = 999
	if phonology.templates[0].Weight == 999 {
		t.Error("Expected template copy, got reference")
	}
}

func TestGetRuleCount(t *testing.T) {
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	if phonology.GetRuleCount() != 0 {
		t.Errorf("Expected 0 rules, got %d", phonology.GetRuleCount())
	}

	phonology.AddRule(NewMaxClusterSizeRule(3))
	if phonology.GetRuleCount() != 1 {
		t.Errorf("Expected 1 rule, got %d", phonology.GetRuleCount())
	}
}
