package lang

import (
	"testing"
)

// TestGeographicInfluenceModel tests the geographic influence model system.
func TestGeographicInfluenceModel(t *testing.T) {
	model := NewGeographicInfluenceModel()

	t.Run("Climate Effects", func(t *testing.T) {
		// Test tropical climate effects
		tropicalEffect := model.ClimateEffects["tropical"]
		if tropicalEffect.Description == "" {
			t.Error("Expected tropical climate effect to have description")
		}
		if tropicalEffect.ComplexityImpact != 0.2 {
			t.Errorf("Expected tropical complexity impact 0.2, got %f", tropicalEffect.ComplexityImpact)
		}
		if len(tropicalEffect.PhonologicalChanges) == 0 {
			t.Error("Expected tropical climate to have phonological changes")
		}
		if len(tropicalEffect.LexicalChanges) == 0 {
			t.Error("Expected tropical climate to have lexical changes")
		}
		if len(tropicalEffect.GrammaticalChanges) == 0 {
			t.Error("Expected tropical climate to have grammatical changes")
		}

		// Test arctic climate effects
		arcticEffect := model.ClimateEffects["arctic"]
		if arcticEffect.ComplexityImpact != 0.3 {
			t.Errorf("Expected arctic complexity impact 0.3, got %f", arcticEffect.ComplexityImpact)
		}
		if len(arcticEffect.LexicalChanges) == 0 {
			t.Error("Expected arctic climate to have lexical changes")
		}

		t.Logf("Climate effects: tropical=%f, arctic=%f", tropicalEffect.ComplexityImpact, arcticEffect.ComplexityImpact)
	})

	t.Run("Terrain Effects", func(t *testing.T) {
		// Test mountain terrain effects
		mountainEffect := model.TerrainEffects["mountain"]
		if mountainEffect.Description == "" {
			t.Error("Expected mountain terrain effect to have description")
		}
		if mountainEffect.ComplexityImpact != 0.1 {
			t.Errorf("Expected mountain complexity impact 0.1, got %f", mountainEffect.ComplexityImpact)
		}

		// Test coastal terrain effects
		coastalEffect := model.TerrainEffects["coastal"]
		if coastalEffect.ComplexityImpact != -0.1 {
			t.Errorf("Expected coastal complexity impact -0.1, got %f", coastalEffect.ComplexityImpact)
		}

		t.Logf("Terrain effects: mountain=%f, coastal=%f", mountainEffect.ComplexityImpact, coastalEffect.ComplexityImpact)
	})

	t.Run("Population Effects", func(t *testing.T) {
		// Test high population effects
		highEffect := model.PopulationEffects["high"]
		if highEffect.ComplexityImpact != -0.2 {
			t.Errorf("Expected high population complexity impact -0.2, got %f", highEffect.ComplexityImpact)
		}

		// Test low population effects
		lowEffect := model.PopulationEffects["low"]
		if lowEffect.ComplexityImpact != 0.3 {
			t.Errorf("Expected low population complexity impact 0.3, got %f", lowEffect.ComplexityImpact)
		}

		t.Logf("Population effects: high=%f, low=%f", highEffect.ComplexityImpact, lowEffect.ComplexityImpact)
	})

	t.Run("Urbanization Effects", func(t *testing.T) {
		// Test high urbanization effects
		highEffect := model.UrbanizationEffects["high"]
		if highEffect.ComplexityImpact != -0.3 {
			t.Errorf("Expected high urbanization complexity impact -0.3, got %f", highEffect.ComplexityImpact)
		}

		// Test low urbanization effects
		lowEffect := model.UrbanizationEffects["low"]
		if lowEffect.ComplexityImpact != 0.2 {
			t.Errorf("Expected low urbanization complexity impact 0.2, got %f", lowEffect.ComplexityImpact)
		}

		t.Logf("Urbanization effects: high=%f, low=%f", highEffect.ComplexityImpact, lowEffect.ComplexityImpact)
	})

	t.Run("Isolation Effects", func(t *testing.T) {
		// Test high isolation effects
		highEffect := model.IsolationEffects["high"]
		if highEffect.ComplexityImpact != 0.4 {
			t.Errorf("Expected high isolation complexity impact 0.4, got %f", highEffect.ComplexityImpact)
		}

		// Test low isolation effects
		lowEffect := model.IsolationEffects["low"]
		if lowEffect.ComplexityImpact != -0.2 {
			t.Errorf("Expected low isolation complexity impact -0.2, got %f", lowEffect.ComplexityImpact)
		}

		t.Logf("Isolation effects: high=%f, low=%f", highEffect.ComplexityImpact, lowEffect.ComplexityImpact)
	})

	t.Run("Trade Effects", func(t *testing.T) {
		// Test trade route effects
		tradeEffect := model.TradeEffects["trade_route"]
		if tradeEffect.ComplexityImpact != 0.1 {
			t.Errorf("Expected trade route complexity impact 0.1, got %f", tradeEffect.ComplexityImpact)
		}

		// Test no trade effects
		noTradeEffect := model.TradeEffects["no_trade"]
		if noTradeEffect.ComplexityImpact != 0.1 {
			t.Errorf("Expected no trade complexity impact 0.1, got %f", noTradeEffect.ComplexityImpact)
		}

		t.Logf("Trade effects: trade_route=%f, no_trade=%f", tradeEffect.ComplexityImpact, noTradeEffect.ComplexityImpact)
	})
}

// TestGeographicInfluenceApplication tests the application of geographic influence to languages.
func TestGeographicInfluenceApplication(t *testing.T) {
	model := NewGeographicInfluenceModel()

	// Create a test language
	baseLang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	t.Run("Tropical Coastal Region", func(t *testing.T) {
		region := GeographicRegion{
			ID:           "tropical_coastal",
			Name:         "Tropical Coastal",
			Latitude:     10.0,
			Longitude:    80.0,
			Climate:      "tropical",
			Terrain:      "coastal",
			Population:   5000,
			Urbanization: 0.6,
			Isolation:    0.3,
			TradeRoutes:  true,
			PortAccess:   true,
		}

		changes := model.ApplyGeographicInfluence(baseLang, region)

		// Should have changes for climate, terrain, population, urbanization, isolation, and trade
		expectedChanges := 6
		if len(changes) != expectedChanges {
			t.Errorf("Expected %d geographic changes, got %d", expectedChanges, len(changes))
		}

		// Verify climate change
		var climateChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "climate_influence" {
				climateChange = change
				break
			}
		}
		if climateChange.Trigger != "climate_influence" {
			t.Error("Expected climate influence change to be applied")
		}
		if climateChange.Description == "" {
			t.Error("Expected climate change to have description")
		}

		// Verify terrain change
		var terrainChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "terrain_influence" {
				terrainChange = change
				break
			}
		}
		if terrainChange.Trigger != "terrain_influence" {
			t.Error("Expected terrain influence change to be applied")
		}

		t.Logf("Applied %d geographic changes for tropical coastal region", len(changes))
	})

	t.Run("Arctic Mountain Region", func(t *testing.T) {
		region := GeographicRegion{
			ID:           "arctic_mountain",
			Name:         "Arctic Mountain",
			Latitude:     70.0,
			Longitude:    -150.0,
			Climate:      "arctic",
			Terrain:      "mountain",
			Population:   500,
			Urbanization: 0.1,
			Isolation:    0.9,
			TradeRoutes:  false,
			PortAccess:   false,
		}

		changes := model.ApplyGeographicInfluence(baseLang, region)

		// Should have changes for climate, terrain, population, urbanization, isolation, and trade
		expectedChanges := 6
		if len(changes) != expectedChanges {
			t.Errorf("Expected %d geographic changes, got %d", expectedChanges, len(changes))
		}

		// Verify high isolation effect
		var isolationChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "isolation_influence" {
				isolationChange = change
				break
			}
		}
		if isolationChange.Trigger != "isolation_influence" {
			t.Error("Expected isolation influence change to be applied")
		}

		// Verify arctic climate effect
		var climateChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "climate_influence" {
				climateChange = change
				break
			}
		}
		if climateChange.Trigger != "climate_influence" {
			t.Error("Expected climate influence change to be applied")
		}

		t.Logf("Applied %d geographic changes for arctic mountain region", len(changes))
	})

	t.Run("Desert Plains Region", func(t *testing.T) {
		region := GeographicRegion{
			ID:           "desert_plains",
			Name:         "Desert Plains",
			Latitude:     25.0,
			Longitude:    45.0,
			Climate:      "desert",
			Terrain:      "plains",
			Population:   2000,
			Urbanization: 0.4,
			Isolation:    0.5,
			TradeRoutes:  true,
			PortAccess:   false,
		}

		changes := model.ApplyGeographicInfluence(baseLang, region)

		// Should have changes for climate, terrain, population, urbanization, isolation, and trade
		expectedChanges := 6
		if len(changes) != expectedChanges {
			t.Errorf("Expected %d geographic changes, got %d", expectedChanges, len(changes))
		}

		// Verify desert climate effect
		var climateChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "climate_influence" {
				climateChange = change
				break
			}
		}
		if climateChange.Trigger != "climate_influence" {
			t.Error("Expected climate influence change to be applied")
		}

		t.Logf("Applied %d geographic changes for desert plains region", len(changes))
	})
}

// TestGeographicInfluenceIntegration tests the integration of geographic influence with dialect creation.
func TestGeographicInfluenceIntegration(t *testing.T) {
	t.Run("Geographic Dialect with Full Region", func(t *testing.T) {
		// Create a base language
		baseLang, err := CreateRandomLanguage("test", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create base language: %v", err)
		}

		// Create a comprehensive geographic region
		region := GeographicRegion{
			ID:           "comprehensive_test",
			Name:         "Comprehensive Test Region",
			Latitude:     45.0,
			Longitude:    -75.0,
			Climate:      "temperate",
			Terrain:      "forest",
			Population:   8000,
			Urbanization: 0.7,
			Isolation:    0.4,
			TradeRoutes:  true,
			PortAccess:   false,
		}

		// Create the geographic dialect
		dialectID := LanguageID{
			Family:   "test",
			Branch:   "geographic",
			Language: "dialect",
			Dialect:  "comprehensive",
		}

		dialect, err := CreateGeographicDialect(baseLang, dialectID, "Comprehensive Dialect", region, "middle_evolution", 43)
		if err != nil {
			t.Fatalf("Failed to create geographic dialect: %v", err)
		}

		// Verify dialect properties
		if dialect.Name != "Comprehensive Dialect" {
			t.Errorf("Expected dialect name 'Comprehensive Dialect', got '%s'", dialect.Name)
		}

		// Verify geographic influence changes were applied
		geographicChanges := dialect.GetLinguisticChanges(LinguisticChangeTypeDialectal)
		if len(geographicChanges) < 7 { // 6 geographic + 1 formation
			t.Errorf("Expected at least 7 dialect changes, got %d", len(geographicChanges))
		}

		// Verify specific geographic changes
		var climateChange, terrainChange, populationChange LinguisticChange
		for _, change := range geographicChanges {
			switch change.Trigger {
			case "climate_influence":
				climateChange = change
			case "terrain_influence":
				terrainChange = change
			case "population_influence":
				populationChange = change
			}
		}

		if climateChange.Trigger != "climate_influence" {
			t.Error("Expected climate influence change to be applied")
		}
		if terrainChange.Trigger != "terrain_influence" {
			t.Error("Expected terrain influence change to be applied")
		}
		if populationChange.Trigger != "population_influence" {
			t.Error("Expected population influence change to be applied")
		}

		// Verify complexity changes
		totalComplexityChange := float32(0.0)
		for _, change := range geographicChanges {
			totalComplexityChange += change.ComplexityChange
		}

		t.Logf("Dialect created with %d geographic changes, total complexity impact: %f", len(geographicChanges), totalComplexityChange)
		t.Logf("Climate change: %s", climateChange.Description)
		t.Logf("Terrain change: %s", terrainChange.Description)
		t.Logf("Population change: %s", populationChange.Description)
	})
}

// TestGeographicRegionStruct tests the GeographicRegion struct.
func TestGeographicRegionStruct(t *testing.T) {
	region := GeographicRegion{
		ID:           "test_region",
		Name:         "Test Region",
		Latitude:     45.0,
		Longitude:    -75.0,
		Climate:      "temperate",
		Terrain:      "mountain",
		Population:   10000,
		Urbanization: 0.7,
		Isolation:    0.4,
		TradeRoutes:  true,
		PortAccess:   false,
	}

	// Verify region properties
	if region.ID != "test_region" {
		t.Errorf("Expected region ID 'test_region', got '%s'", region.ID)
	}
	if region.Name != "Test Region" {
		t.Errorf("Expected region name 'Test Region', got '%s'", region.Name)
	}
	if region.Latitude != 45.0 {
		t.Errorf("Expected latitude 45.0, got %f", region.Latitude)
	}
	if region.Longitude != -75.0 {
		t.Errorf("Expected longitude -75.0, got %f", region.Longitude)
	}
	if region.Climate != "temperate" {
		t.Errorf("Expected climate 'temperate', got '%s'", region.Climate)
	}
	if region.Terrain != "mountain" {
		t.Errorf("Expected terrain 'mountain', got '%s'", region.Terrain)
	}
	if region.Population != 10000 {
		t.Errorf("Expected population 10000, got %d", region.Population)
	}
	if region.Urbanization != 0.7 {
		t.Errorf("Expected urbanization 0.7, got %f", region.Urbanization)
	}
	if region.Isolation != 0.4 {
		t.Errorf("Expected isolation 0.4, got %f", region.Isolation)
	}
	if !region.TradeRoutes {
		t.Error("Expected trade routes to be true")
	}
	if region.PortAccess {
		t.Error("Expected port access to be false")
	}

	t.Logf("Successfully tested GeographicRegion struct with all fields")
}
