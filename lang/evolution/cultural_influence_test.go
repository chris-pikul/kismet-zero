package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

func TestNewCulturalInfluenceEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	if engine == nil {
		t.Fatal("Cultural influence engine should not be nil")
	}

	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}
}

func TestCulturalContactType_String(t *testing.T) {
	testCases := []struct {
		contactType CulturalContactType
		expected    string
	}{
		{CulturalContactTypeUnknown, "unknown"},
		{CulturalContactTypeTrade, "trade"},
		{CulturalContactTypeConquest, "conquest"},
		{CulturalContactTypeMigration, "migration"},
		{CulturalContactTypeReligious, "religious"},
		{CulturalContactTypeDiplomatic, "diplomatic"},
		{CulturalContactTypeAlliance, "alliance"},
		{CulturalContactTypeColonization, "colonization"},
		{CulturalContactTypeCulturalExchange, "cultural_exchange"},
		{CulturalContactTypeMagical, "magical"},
		{CulturalContactTypeAncient, "ancient"},
	}

	for _, tc := range testCases {
		result := tc.contactType.String()
		if result != tc.expected {
			t.Errorf("Expected '%s' for %v, got '%s'", tc.expected, tc.contactType, result)
		}
	}

	// Test out of range value
	outOfRange := CulturalContactType(99)
	if outOfRange.String() != "unknown" {
		t.Errorf("Expected 'unknown' for out of range value, got '%s'", outOfRange.String())
	}
}

func TestCulturalInfluenceEngine_SimulateCulturalInfluence(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	// Create test languages
	sourceLang := &lang.Language{
		ID:      lang.LanguageID{Family: "test", Branch: "source", Language: "source_lang"},
		Name:    "Source Language",
		Culture: "source_culture",
		Seed:    123,
	}

	targetLang := &lang.Language{
		ID:      lang.LanguageID{Family: "test", Branch: "target", Language: "target_lang"},
		Name:    "Target Language",
		Culture: "target_culture",
		Seed:    456,
	}

	// Create test cultural identities
	sourceIdentity := CulturalIdentity{
		CultureName:           "Source Culture",
		WritingSystemPrestige: 0.8,
		InnovationTendency:    0.7,
		PreservationInstinct:  0.3,
		MagicalTradition:      0.6,
		ReligiousInfluence:    0.5,
		EconomicPower:         0.8,
		MilitaryPower:         0.7,
		CulturalConfidence:    0.6,
	}

	targetIdentity := CulturalIdentity{
		CultureName:           "Target Culture",
		WritingSystemPrestige: 0.4,
		InnovationTendency:    0.8,
		PreservationInstinct:  0.2,
		MagicalTradition:      0.4,
		ReligiousInfluence:    0.6,
		EconomicPower:         0.5,
		MilitaryPower:         0.4,
		CulturalConfidence:    0.5,
	}

	// Test trade contact
	changes, event := engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeTrade,
		time.Hour*24*365*50, // 50 years
		sourceIdentity,
		targetIdentity,
	)

	// Verify the event structure
	if event.ID == "" {
		t.Error("Event should have an ID")
	}

	if event.SourceCulture != sourceLang.Culture {
		t.Errorf("Expected source culture %s, got %s", sourceLang.Culture, event.SourceCulture)
	}

	if event.TargetCulture != targetLang.Culture {
		t.Errorf("Expected target culture %s, got %s", targetLang.Culture, event.TargetCulture)
	}

	if event.ContactType != CulturalContactTypeTrade {
		t.Errorf("Expected contact type %s, got %s", CulturalContactTypeTrade.String(), event.ContactType.String())
	}

	if event.Intensity <= 0.0 || event.Intensity > 1.0 {
		t.Errorf("Intensity should be between 0.0 and 1.0, got %f", event.Intensity)
	}

	// Verify compatibility was calculated
	if event.Compatibility.LinguisticSimilarity == 0.0 {
		t.Error("Linguistic similarity should be calculated")
	}

	// Verify cultural identities were set
	if event.SourceIdentity.CultureName != sourceIdentity.CultureName {
		t.Error("Source identity should be set correctly")
	}

	if event.TargetIdentity.CultureName != targetIdentity.CultureName {
		t.Error("Target identity should be set correctly")
	}

	// Verify changes and features
	if len(changes) > 0 {
		event.BorrowedFeatures = engine.extractBorrowedFeatures(changes)
		if len(event.BorrowedFeatures) != len(changes) {
			t.Error("Borrowed features should match the number of changes")
		}
	}

	// Verify description was generated
	if event.Description == "" {
		t.Error("Description should be generated")
	}
}

func TestCulturalInfluenceEngine_CalculateCulturalCompatibility(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	sourceLang := &lang.Language{Culture: "source"}
	targetLang := &lang.Language{Culture: "target"}

	sourceIdentity := CulturalIdentity{
		InnovationTendency: 0.8,
		MagicalTradition:   0.7,
		ReligiousInfluence: 0.6,
		EconomicPower:      0.9,
		MilitaryPower:      0.8,
		CulturalConfidence: 0.7,
	}

	targetIdentity := CulturalIdentity{
		InnovationTendency: 0.3,
		MagicalTradition:   0.2,
		ReligiousInfluence: 0.4,
		EconomicPower:      0.4,
		MilitaryPower:      0.3,
		CulturalConfidence: 0.5,
	}

	compatibility := engine.calculateCulturalCompatibility(sourceLang, targetLang, sourceIdentity, targetIdentity)

	// Test power balance calculation
	expectedPowerBalance := float32((0.9+0.8+0.7)/3.0 - (0.4+0.3+0.5)/3.0)
	if abs(compatibility.PowerBalance-expectedPowerBalance) > 0.01 {
		t.Errorf("Expected power balance %.3f, got %.3f", expectedPowerBalance, compatibility.PowerBalance)
	}

	// Test cultural values compatibility
	expectedCulturalValues := 1.0 - abs(0.8-0.3)
	if abs(compatibility.CulturalValues-expectedCulturalValues) > 0.01 {
		t.Errorf("Expected cultural values %.3f, got %.3f", expectedCulturalValues, compatibility.CulturalValues)
	}

	// Test magical affinity
	expectedMagicalAffinity := 1.0 - abs(0.7-0.2)
	if abs(compatibility.MagicalAffinity-expectedMagicalAffinity) > 0.01 {
		t.Errorf("Expected magical affinity %.3f, got %.3f", expectedMagicalAffinity, compatibility.MagicalAffinity)
	}

	// Test religious compatibility
	expectedReligiousCompatibility := 1.0 - abs(0.6-0.4)
	if abs(compatibility.ReligiousCompatibility-expectedReligiousCompatibility) > 0.01 {
		t.Errorf("Expected religious compatibility %.3f, got %.3f", expectedReligiousCompatibility, compatibility.ReligiousCompatibility)
	}
}

func TestCulturalInfluenceEngine_CalculateBaseInfluenceIntensity(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	compatibility := CulturalCompatibility{
		LinguisticSimilarity: 0.8,
		CulturalValues:       0.7,
		GeographicProximity:  0.6,
	}

	// Test conquest (highest intensity)
	intensity := engine.calculateBaseInfluenceIntensity(
		CulturalContactTypeConquest,
		time.Hour*24*365*100, // 100 years
		compatibility,
	)

	if intensity < 0.7 {
		t.Errorf("Conquest should have high intensity, got %f", intensity)
	}

	// Test trade (moderate intensity)
	intensity = engine.calculateBaseInfluenceIntensity(
		CulturalContactTypeTrade,
		time.Hour*24*365*50, // 50 years
		compatibility,
	)

	if intensity < 0.4 || intensity > 0.8 {
		t.Errorf("Trade should have moderate intensity, got %f", intensity)
	}

	// Test diplomatic (lower intensity)
	intensity = engine.calculateBaseInfluenceIntensity(
		CulturalContactTypeDiplomatic,
		time.Hour*24*365*10, // 10 years
		compatibility,
	)

	if intensity < 0.2 || intensity > 0.6 {
		t.Errorf("Diplomatic should have lower intensity, got %f", intensity)
	}
}

func TestCulturalInfluenceEngine_ApplyCulturalModifiers(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	baseIntensity := float32(0.6)

	sourceIdentity := CulturalIdentity{
		WritingSystemPrestige: 0.8, // High prestige
		MagicalTradition:      0.9, // High magical tradition
		ReligiousInfluence:    0.8, // High religious influence
	}

	targetIdentity := CulturalIdentity{
		InnovationTendency:   0.8, // High innovation
		PreservationInstinct: 0.2, // Low preservation instinct
	}

	compatibility := CulturalCompatibility{
		PowerBalance:           -0.6, // Source dominant
		MagicalAffinity:        0.8,  // High magical affinity
		ReligiousCompatibility: 0.8,  // High religious compatibility
	}

	modifiedIntensity := engine.applyCulturalModifiers(baseIntensity, sourceIdentity, targetIdentity, compatibility)

	// Should be higher than base due to positive modifiers
	if modifiedIntensity <= baseIntensity {
		t.Errorf("Modified intensity should be higher than base, got %f vs %f", modifiedIntensity, baseIntensity)
	}

	// Should not exceed 1.0
	if modifiedIntensity > 1.0 {
		t.Errorf("Modified intensity should not exceed 1.0, got %f", modifiedIntensity)
	}
}

func TestCulturalInfluenceEngine_ShouldInfluenceMethods(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	compatibility := CulturalCompatibility{
		LinguisticSimilarity: 0.85,
	}

	// Test orthographic influence
	shouldInfluence := engine.shouldInfluenceOrthography(CulturalContactTypeTrade, 0.8, compatibility)
	if !shouldInfluence {
		t.Error("Should influence orthography with high intensity and good compatibility")
	}

	shouldInfluence = engine.shouldInfluenceOrthography(CulturalContactTypeTrade, 0.6, compatibility)
	if shouldInfluence {
		t.Error("Should not influence orthography with low intensity")
	}

	// Test lexical influence
	shouldInfluence = engine.shouldInfluenceLexical(CulturalContactTypeTrade, 0.6, compatibility)
	if !shouldInfluence {
		t.Error("Should influence lexical with moderate intensity")
	}

	// Test phonological influence
	shouldInfluence = engine.shouldInfluencePhonological(CulturalContactTypeTrade, 0.8, compatibility)
	if !shouldInfluence {
		t.Error("Should influence phonological with high intensity and good compatibility")
	}

	// Test grammatical influence
	shouldInfluence = engine.shouldInfluenceGrammatical(CulturalContactTypeTrade, 0.9, compatibility)
	if !shouldInfluence {
		t.Error("Should influence grammatical with very high intensity and excellent compatibility")
	}
}

func TestCulturalInfluenceEngine_HelperMethods(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	// Test extractBorrowedFeatures
	changes := []LinguisticChange{
		{Type: ChangeTypeOrthographic},
		{Type: ChangeTypeLexical},
		{Type: ChangeTypeSoundShift},
	}

	features := engine.extractBorrowedFeatures(changes)
	if len(features) != 3 {
		t.Errorf("Expected 3 features, got %d", len(features))
	}

	// Test generateAdaptationNotes
	compatibility := CulturalCompatibility{CulturalValues: 0.8}
	notes := engine.generateAdaptationNotes(changes, compatibility)
	if notes == "" {
		t.Error("Adaptation notes should be generated")
	}

	// Test identifyResistanceFactors
	targetIdentity := CulturalIdentity{
		PreservationInstinct: 0.8,
		CulturalConfidence:   0.9,
	}
	compatibility.PowerBalance = 0.6

	resistanceFactors := engine.identifyResistanceFactors(targetIdentity, compatibility)
	if len(resistanceFactors) == 0 {
		t.Error("Should identify resistance factors")
	}

	// Test identifySuccessFactors
	sourceIdentity := CulturalIdentity{
		WritingSystemPrestige: 0.8,
	}
	compatibility.LinguisticSimilarity = 0.8
	compatibility.GeographicProximity = 0.8

	successFactors := engine.identifySuccessFactors(sourceIdentity, compatibility)
	if len(successFactors) == 0 {
		t.Error("Should identify success factors")
	}

	// Test generateInfluenceDescription
	event := CulturalInfluenceEvent{
		ContactType:   CulturalContactTypeTrade,
		SourceCulture: "source",
		TargetCulture: "target",
		Intensity:     0.7,
	}

	description := engine.generateInfluenceDescription(event)
	if description == "" {
		t.Error("Description should be generated")
	}
}

func TestCulturalInfluenceEngine_FantasyElements(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	sourceLang := &lang.Language{Culture: "magical_source"}
	targetLang := &lang.Language{Culture: "magical_target"}

	// Test magical influence
	sourceIdentity := CulturalIdentity{
		MagicalTradition: 0.8, // High magical tradition
	}
	targetIdentity := CulturalIdentity{
		MagicalTradition: 0.5, // Moderate magical tradition
	}

	_, event := engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeTrade,
		time.Hour*24*365*50,
		sourceIdentity,
		targetIdentity,
	)

	if !event.MagicalInfluence {
		t.Error("Should detect magical influence with compatible magical traditions")
	}

	// Test religious influence
	sourceIdentity.ReligiousInfluence = 0.8
	targetIdentity.ReligiousInfluence = 0.5

	_, event = engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeReligious,
		time.Hour*24*365*50,
		sourceIdentity,
		targetIdentity,
	)

	if !event.ReligiousInfluence {
		t.Error("Should detect religious influence with compatible religious systems")
	}

	// Test ancient influence
	_, event = engine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		CulturalContactTypeAncient,
		time.Hour*24*365*50,
		sourceIdentity,
		targetIdentity,
	)

	if !event.AncientInfluence {
		t.Error("Should detect ancient influence with ancient contact type")
	}
}

func TestCulturalInfluenceEngine_ContactTypeIntensities(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewCulturalInfluenceEngine(config)

	compatibility := CulturalCompatibility{
		LinguisticSimilarity: 0.7,
		CulturalValues:       0.6,
		GeographicProximity:  0.5,
	}

	duration := time.Hour * 24 * 365 * 100 // 100 years

	// Test all contact types
	contactTypes := []CulturalContactType{
		CulturalContactTypeConquest,
		CulturalContactTypeColonization,
		CulturalContactTypeReligious,
		CulturalContactTypeTrade,
		CulturalContactTypeMigration,
		CulturalContactTypeAlliance,
		CulturalContactTypeDiplomatic,
		CulturalContactTypeCulturalExchange,
		CulturalContactTypeMagical,
		CulturalContactTypeAncient,
	}

	for _, contactType := range contactTypes {
		intensity := engine.calculateBaseInfluenceIntensity(contactType, duration, compatibility)

		if intensity <= 0.0 || intensity > 1.0 {
			t.Errorf("Contact type %s should have intensity between 0.0 and 1.0, got %f",
				contactType.String(), intensity)
		}

		// Conquest and colonization should have highest intensities
		if (contactType == CulturalContactTypeConquest || contactType == CulturalContactTypeColonization) && intensity < 0.6 {
			t.Errorf("Contact type %s should have high intensity, got %f", contactType.String(), intensity)
		}

		// Diplomatic and cultural exchange should have lower intensities
		if (contactType == CulturalContactTypeDiplomatic || contactType == CulturalContactTypeCulturalExchange) && intensity > 0.6 {
			t.Errorf("Contact type %s should have lower intensity, got %f", contactType.String(), intensity)
		}
	}
}
