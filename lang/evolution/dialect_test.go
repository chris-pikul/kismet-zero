package evolution

import (
	"testing"
	"time"
)

func TestNewDialectFormationEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	if engine == nil {
		t.Fatal("Dialect formation engine should not be nil")
	}

	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}
}

func TestCreateGeographicDialect(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	// Create a test parent language
	parentLang := createTestLanguage("test_lang", "test_culture")

	// Create a geographic region
	region := GeographicRegion{
		ID:           "test_region",
		Name:         "Mountain Valley",
		Latitude:     45.0,
		Longitude:    -120.0,
		Climate:      "temperate",
		Terrain:      "mountain",
		Population:   5000,
		Urbanization: 0.3,
	}

	dialect, formationEvent, err := engine.CreateGeographicDialect(
		parentLang,
		"mountain_dialect",
		"Mountain Valley Dialect",
		region,
		"modern",
	)

	if err != nil {
		t.Fatalf("Failed to create geographic dialect: %v", err)
	}

	if dialect == nil {
		t.Fatal("Dialect should not be nil")
	}

	if dialect.ID != "mountain_dialect" {
		t.Errorf("Expected dialect ID 'mountain_dialect', got '%s'", dialect.ID)
	}

	if dialect.Name != "Mountain Valley Dialect" {
		t.Errorf("Expected dialect name 'Mountain Valley Dialect', got '%s'", dialect.Name)
	}

	if dialect.Type != DialectTypeGeographic {
		t.Errorf("Expected dialect type %v, got %v", DialectTypeGeographic, dialect.Type)
	}

	if dialect.ParentLang != parentLang.ID.String() {
		t.Errorf("Expected parent language '%s', got '%s'", parentLang.ID.String(), dialect.ParentLang)
	}

	if dialect.Region.ID != "test_region" {
		t.Errorf("Expected region ID 'test_region', got '%s'", dialect.Region.ID)
	}

	// Check formation event
	if formationEvent.ParentLang != parentLang.ID.String() {
		t.Errorf("Expected formation event parent language '%s', got '%s'", parentLang.ID.String(), formationEvent.ParentLang)
	}

	if formationEvent.NewDialect != "mountain_dialect" {
		t.Errorf("Expected formation event new dialect 'mountain_dialect', got '%s'", formationEvent.NewDialect)
	}

	// Check that geographic factors were calculated
	if len(formationEvent.GeographicFactors) == 0 {
		t.Error("Formation event should have geographic factors")
	}

	// Check dialect features
	if dialect.Features.IntelligibilityScore <= 0 || dialect.Features.IntelligibilityScore > 1 {
		t.Errorf("Intelligibility score should be between 0 and 1, got %f", dialect.Features.IntelligibilityScore)
	}

	// Check that mountain-specific features were added
	foundMountainFeature := false
	for _, feature := range dialect.Features.PhonologicalFeatures {
		if feature == "consonant_cluster_simplification" {
			foundMountainFeature = true
			break
		}
	}
	if !foundMountainFeature {
		t.Error("Mountain dialect should have mountain-specific phonological features")
	}
}

func TestCreateSocialDialect(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	// Create a test parent language
	parentLang := createTestLanguage("test_lang", "test_culture")

	dialect, formationEvent, err := engine.CreateSocialDialect(
		parentLang,
		"working_class_dialect",
		"Working Class Dialect",
		"working",
		0.8, // High urbanization
		"modern",
	)

	if err != nil {
		t.Fatalf("Failed to create social dialect: %v", err)
	}

	if dialect == nil {
		t.Fatal("Dialect should not be nil")
	}

	if dialect.ID != "working_class_dialect" {
		t.Errorf("Expected dialect ID 'working_class_dialect', got '%s'", dialect.ID)
	}

	if dialect.Type != DialectTypeUrban {
		t.Errorf("Expected dialect type %v for high urbanization, got %v", DialectTypeUrban, dialect.Type)
	}

	// Check that social factors were calculated
	if len(formationEvent.SocialFactors) == 0 {
		t.Error("Formation event should have social factors")
	}

	// Check that urban features were added
	foundUrbanFeature := false
	for _, feature := range dialect.Features.PhonologicalFeatures {
		if feature == "fast_speech_patterns" {
			foundUrbanFeature = true
			break
		}
	}
	if !foundUrbanFeature {
		t.Error("Urban dialect should have urban-specific phonological features")
	}
}

func TestEvolveDialect(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	// Create a test dialect
	dialect := &Dialect{
		ID:   "test_dialect",
		Name: "Test Dialect",
		Type: DialectTypeGeographic,
		Features: DialectFeatures{
			ID:                   "test_features",
			PhonologicalFeatures: []string{"existing_feature"},
			IntelligibilityScore: 0.9,
		},
		Status: "active",
	}

	evolvedDialect, evolutionEvent, err := engine.EvolveDialect(
		dialect,
		"modern",
		time.Hour*24*365*50, // 50 years
	)

	if err != nil {
		t.Fatalf("Failed to evolve dialect: %v", err)
	}

	if evolvedDialect == nil {
		t.Fatal("Evolved dialect should not be nil")
	}

	if len(evolvedDialect.EvolutionHistory) == 0 {
		t.Error("Evolved dialect should have evolution history")
	}

	if evolutionEvent.TriggerType != "dialectal_evolution" {
		t.Errorf("Expected trigger type 'dialectal_evolution', got '%s'", evolutionEvent.TriggerType)
	}

	// Check that intelligibility decreased slightly
	if evolvedDialect.Features.IntelligibilityScore >= dialect.Features.IntelligibilityScore {
		t.Error("Dialect intelligibility should decrease over time")
	}
}

func TestCalculateGeographicFactors(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	// Test mountain region
	mountainRegion := GeographicRegion{
		ID:           "mountain",
		Name:         "Mountain Region",
		Latitude:     45.0,
		Longitude:    -120.0,
		Climate:      "temperate",
		Terrain:      "mountain",
		Population:   800,  // Below 1000 threshold for low_population_density
		Urbanization: 0.15, // Below 0.2 threshold for low_urbanization
	}

	factors := engine.calculateGeographicFactors(mountainRegion)

	expectedFactors := []string{
		"climate_temperate",
		"terrain_mountain",
		"low_population_density",
		"low_urbanization",
		"geographic_isolation",
	}

	for _, expected := range expectedFactors {
		found := false
		for _, factor := range factors {
			if factor == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected factor '%s' not found in %v", expected, factors)
		}
	}
}

func TestCalculateSocialFactors(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	factors := engine.calculateSocialFactors("educated", 0.8)

	expectedFactors := []string{
		"social_class_educated",
		"high_urbanization",
		"high_education",
	}

	for _, expected := range expectedFactors {
		found := false
		for _, factor := range factors {
			if factor == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected factor '%s' not found in %v", expected, factors)
		}
	}
}

func TestGenerateDialectFeatures(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	parentLang := createTestLanguage("test_lang", "test_culture")
	region := GeographicRegion{
		ID:           "coastal",
		Name:         "Coastal Region",
		Climate:      "tropical",
		Terrain:      "coastal",
		Population:   15000,
		Urbanization: 0.9,
	}

	factors := []string{"climate_tropical", "terrain_coastal", "high_population_density", "high_urbanization"}

	features := engine.generateDialectFeatures(parentLang, factors, region)

	if features.IntelligibilityScore <= 0 || features.IntelligibilityScore > 1 {
		t.Errorf("Intelligibility score should be between 0 and 1, got %f", features.IntelligibilityScore)
	}

	// Check that coastal features were added
	foundCoastalFeature := false
	for _, feature := range features.LexicalFeatures {
		if feature == "maritime_vocabulary" {
			foundCoastalFeature = true
			break
		}
	}
	if !foundCoastalFeature {
		t.Error("Coastal dialect should have maritime vocabulary")
	}

	// Check that urban features were added
	foundUrbanFeature := false
	for _, feature := range features.PhonologicalFeatures {
		if feature == "fast_speech_patterns" {
			foundUrbanFeature = true
			break
		}
	}
	if !foundUrbanFeature {
		t.Error("Urban dialect should have fast speech patterns")
	}
}

func TestCalculateFormationIntensity(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	// Test with geographic isolation
	factors := []string{"geographic_isolation", "low_population_density"}
	intensity := engine.calculateFormationIntensity(factors)

	expectedIntensity := float32(0.3 + 0.2 + 0.15) // base + isolation + low population
	if intensity != expectedIntensity {
		t.Errorf("Expected intensity %f, got %f", expectedIntensity, intensity)
	}

	// Test that intensity doesn't exceed 1.0
	highFactors := []string{"geographic_isolation", "geographic_isolation", "geographic_isolation"}
	highIntensity := engine.calculateFormationIntensity(highFactors)
	if highIntensity > 1.0 {
		t.Errorf("Intensity should not exceed 1.0, got %f", highIntensity)
	}
}

func TestCalculateIntelligibility(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewDialectFormationEngine(config)

	// Test with factors that reduce intelligibility
	factors := []string{"geographic_isolation", "high_urbanization", "terrain_mountain"}
	intelligibility := engine.calculateIntelligibility(factors)

	expectedIntelligibility := float32(0.9 - 0.1 - 0.05 - 0.05) // base - isolation - urbanization - terrain
	if intelligibility != expectedIntelligibility {
		t.Errorf("Expected intelligibility %f, got %f", expectedIntelligibility, intelligibility)
	}

	// Test that intelligibility doesn't go below 0.5
	veryLowFactors := []string{"geographic_isolation", "geographic_isolation", "geographic_isolation"}
	veryLowIntelligibility := engine.calculateIntelligibility(veryLowFactors)
	if veryLowIntelligibility < 0.5 {
		t.Errorf("Intelligibility should not go below 0.5, got %f", veryLowIntelligibility)
	}
}
