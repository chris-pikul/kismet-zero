package variation

import (
	"testing"
	"time"
)

// TestDialectFormationAPIs tests the new high-level dialect formation APIs.
func TestDialectFormationAPIs(t *testing.T) {
	// Create a base language for testing
	baseLang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create base language: %v", err)
	}

	t.Run("Geographic Dialect Creation", func(t *testing.T) {
		// Create a geographic dialect
		dialectID := LanguageID{
			Family:   "test",
			Branch:   "coastal",
			Language: "dialect",
			Dialect:  "harbor",
		}

		region := GeographicRegion{
			ID:           "coastal_harbor",
			Name:         "Coastal Harbor",
			Latitude:     45.0,
			Longitude:    -75.0,
			Climate:      "temperate",
			Terrain:      "coastal",
			Population:   5000,
			Urbanization: 0.8,
		}

		dialect, err := CreateGeographicDialect(baseLang, dialectID, "Harbor Dialect", region, "middle_evolution", 43)
		if err != nil {
			t.Fatalf("Failed to create geographic dialect: %v", err)
		}

		// Verify dialect properties
		if dialect.Name != "Harbor Dialect" {
			t.Errorf("Expected dialect name 'Harbor Dialect', got '%s'", dialect.Name)
		}

		if dialect.ParentID == nil || *dialect.ParentID != baseLang.ID {
			t.Errorf("Expected dialect to have correct parent language")
		}

		// Verify dialect formation change was recorded
		dialectChanges := dialect.GetLinguisticChangesByType(LinguisticChangeTypeDialectal)
		if len(dialectChanges) < 7 { // 6 geographic + 1 formation
			t.Errorf("Expected at least 7 dialect changes (6 geographic + 1 formation), got %d", len(dialectChanges))
		}

		// Find the formation change specifically
		var formationChange LinguisticChange
		for _, change := range dialectChanges {
			if change.Trigger == "dialect_formation" {
				formationChange = change
				break
			}
		}
		if formationChange.Trigger != "dialect_formation" {
			t.Error("Expected dialect formation change to be recorded")
		}

		// Verify parent language has dialect as child
		if len(baseLang.ChildIDs) != 1 {
			t.Errorf("Expected base language to have 1 child, got %d", len(baseLang.ChildIDs))
		}

		if baseLang.ChildIDs[0] != dialect.ID {
			t.Errorf("Expected base language to have dialect as child")
		}

		t.Logf("Successfully created geographic dialect: %s", dialect.Name)
	})

	t.Run("Social Dialect Creation", func(t *testing.T) {
		// Create a social dialect
		dialectID := LanguageID{
			Family:   "test",
			Branch:   "social",
			Language: "dialect",
			Dialect:  "noble",
		}

		dialect, err := CreateSocialDialect(baseLang, dialectID, "Noble Dialect", "upper", 0.9, "middle_evolution", 44)
		if err != nil {
			t.Fatalf("Failed to create social dialect: %v", err)
		}

		// Verify dialect properties
		if dialect.Name != "Noble Dialect" {
			t.Errorf("Expected dialect name 'Noble Dialect', got '%s'", dialect.Name)
		}

		// Verify dialect formation change was recorded
		dialectChanges := dialect.GetLinguisticChangesByType(LinguisticChangeTypeDialectal)
		if len(dialectChanges) < 1 { // At least 1 formation change
			t.Errorf("Expected at least 1 dialect formation change, got %d", len(dialectChanges))
		}

		// Find the formation change specifically
		var formationChange LinguisticChange
		for _, change := range dialectChanges {
			if change.Trigger == "dialect_formation" {
				formationChange = change
				break
			}
		}
		if formationChange.Trigger != "dialect_formation" {
			t.Error("Expected dialect formation change to be recorded")
		}

		t.Logf("Successfully created social dialect: %s", dialect.Name)
	})

	t.Run("Dialect Evolution", func(t *testing.T) {
		// Create a dialect first
		dialectID := LanguageID{
			Family:   "test",
			Branch:   "evolution",
			Language: "dialect",
			Dialect:  "evolving",
		}

		region := GeographicRegion{
			ID:           "test_region",
			Name:         "Test Region",
			Latitude:     0.0,
			Longitude:    0.0,
			Climate:      "temperate",
			Terrain:      "plains",
			Population:   1000,
			Urbanization: 0.5,
		}

		dialect, err := CreateGeographicDialect(baseLang, dialectID, "Evolving Dialect", region, "early_evolution", 45)
		if err != nil {
			t.Fatalf("Failed to create dialect for evolution: %v", err)
		}

		// Evolve the dialect
		evolvedDialect, err := EvolveDialect(dialect, "middle_evolution", 100*time.Hour, 46)
		if err != nil {
			t.Fatalf("Failed to evolve dialect: %v", err)
		}

		// Verify evolution change was recorded
		evolutionChanges := evolvedDialect.GetLinguisticChangesByType(LinguisticChangeTypeDialectal)
		if len(evolutionChanges) < 8 { // 6 geographic + 1 formation + 1 evolution
			t.Errorf("Expected at least 8 dialect changes (6 geographic + 1 formation + 1 evolution), got %d", len(evolutionChanges))
		}

		// Find the evolution change specifically
		var evolutionChange LinguisticChange
		for _, change := range evolutionChanges {
			if change.Trigger == "dialect_evolution" {
				evolutionChange = change
				break
			}
		}

		if evolutionChange.Trigger != "dialect_evolution" {
			t.Errorf("Expected evolution change to be recorded")
		}

		if evolutionChange.Era != "middle_evolution" {
			t.Errorf("Expected evolution era 'middle_evolution', got '%s'", evolutionChange.Era)
		}

		t.Logf("Successfully evolved dialect: %s", evolvedDialect.Name)
	})
}

// TestGeographicRegion tests the GeographicRegion struct.
func TestGeographicRegion(t *testing.T) {
	region := GeographicRegion{
		ID:           "test_region",
		Name:         "Test Region",
		Latitude:     45.0,
		Longitude:    -75.0,
		Climate:      "temperate",
		Terrain:      "mountain",
		Population:   10000,
		Urbanization: 0.7,
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

	t.Logf("Successfully tested GeographicRegion struct")
}
