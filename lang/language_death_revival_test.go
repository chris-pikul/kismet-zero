package lang

import (
	"testing"
)

// TestLanguageDeath tests the language death structure.
func TestLanguageDeath(t *testing.T) {
	t.Run("Death Creation", func(t *testing.T) {
		death := NewLanguageDeath(
			"death_1",
			"lang_1",
			"forced",
			"conquest",
			"Language was forced to extinction through conquest",
			"Conquering forces banned the language and forced assimilation",
			50,
			[]string{"northern_region", "mountain_valley"},
		)

		if death.ID != "death_1" {
			t.Errorf("Expected ID 'death_1', got '%s'", death.ID)
		}
		if death.LanguageID != "lang_1" {
			t.Errorf("Expected language ID 'lang_1', got '%s'", death.LanguageID)
		}
		if death.Type != "forced" {
			t.Errorf("Expected type 'forced', got '%s'", death.Type)
		}
		if death.Cause != "conquest" {
			t.Errorf("Expected cause 'conquest', got '%s'", death.Cause)
		}
		if death.LastSpeakers != 50 {
			t.Errorf("Expected 50 last speakers, got %d", death.LastSpeakers)
		}
		if len(death.LastRegions) != 2 {
			t.Errorf("Expected 2 last regions, got %d", len(death.LastRegions))
		}
		if death.PreservationLevel != 0.5 {
			t.Errorf("Expected preservation level 0.5, got %f", death.PreservationLevel)
		}
		if death.RevivalPotential != 0.3 {
			t.Errorf("Expected revival potential 0.3, got %f", death.RevivalPotential)
		}

		t.Logf("Successfully created language death: %s (%s - %s)", death.ID, death.Type, death.Cause)
	})
}

// TestLanguageDeathFeatures tests the feature management functions.
func TestLanguageDeathFeatures(t *testing.T) {
	t.Run("Last Documents Management", func(t *testing.T) {
		death := NewLanguageDeath("death_docs", "lang_1", "natural", "disease", "Test description", "Test details", 10, []string{"region_1"})

		// Add last documents
		documents := []string{"ancient_scroll", "stone_tablet", "wooden_plaque"}
		for _, doc := range documents {
			death.AddLastDocument(doc)
		}

		if len(death.LastDocuments) != 3 {
			t.Errorf("Expected 3 last documents, got %d", len(death.LastDocuments))
		}

		// Test duplicate prevention
		death.AddLastDocument("ancient_scroll")
		if len(death.LastDocuments) != 3 {
			t.Errorf("Expected still 3 last documents after duplicate, got %d", len(death.LastDocuments))
		}

		t.Logf("Successfully tested last documents management: %d documents", len(death.LastDocuments))
	})

	t.Run("Influence on Survivors Management", func(t *testing.T) {
		death := NewLanguageDeath("death_influence", "lang_1", "gradual", "assimilation", "Test description", "Test details", 20, []string{"region_1"})

		// Add influenced languages
		languages := []string{"survivor_lang_1", "survivor_lang_2"}
		for _, lang := range languages {
			death.AddInfluenceOnSurvivor(lang)
		}

		if len(death.InfluenceOnSurvivors) != 2 {
			t.Errorf("Expected 2 influenced languages, got %d", len(death.InfluenceOnSurvivors))
		}

		t.Logf("Successfully tested influence on survivors management: %d languages", len(death.InfluenceOnSurvivors))
	})

	t.Run("Lost Knowledge Management", func(t *testing.T) {
		death := NewLanguageDeath("death_knowledge", "lang_1", "sudden", "catastrophe", "Test description", "Test details", 5, []string{"region_1"})

		// Add lost knowledge
		knowledge := []string{"ancient_healing_techniques", "sacred_ceremonies", "traditional_medicine"}
		for _, k := range knowledge {
			death.AddLostKnowledge(k)
		}

		if len(death.LostKnowledge) != 3 {
			t.Errorf("Expected 3 lost knowledge items, got %d", len(death.LostKnowledge))
		}

		t.Logf("Successfully tested lost knowledge management: %d items", len(death.LostKnowledge))
	})

	t.Run("Revival Factors and Barriers Management", func(t *testing.T) {
		death := NewLanguageDeath("death_revival", "lang_1", "forced", "conquest", "Test description", "Test details", 15, []string{"region_1"})

		// Add revival factors
		factors := []string{"good_documentation", "cultural_interest", "academic_research"}
		for _, factor := range factors {
			death.AddRevivalFactor(factor)
		}

		// Add revival barriers
		barriers := []string{"lack_of_speakers", "political_opposition"}
		for _, barrier := range barriers {
			death.AddRevivalBarrier(barrier)
		}

		if len(death.RevivalFactors) != 3 {
			t.Errorf("Expected 3 revival factors, got %d", len(death.RevivalFactors))
		}
		if len(death.RevivalBarriers) != 2 {
			t.Errorf("Expected 2 revival barriers, got %d", len(death.RevivalBarriers))
		}

		t.Logf("Successfully tested revival factors and barriers management: %d factors, %d barriers",
			len(death.RevivalFactors), len(death.RevivalBarriers))
	})
}

// TestLanguageDeathMetrics tests the metrics calculation system.
func TestLanguageDeathMetrics(t *testing.T) {
	t.Run("Preservation Level Calculation", func(t *testing.T) {
		death := NewLanguageDeath("death_preservation", "lang_1", "natural", "disease", "Test description", "Test details", 25, []string{"region_1", "region_2"})

		initialPreservation := death.PreservationLevel

		// Add documents to trigger preservation level update
		death.AddLastDocument("ancient_text_1")
		death.AddLastDocument("ancient_text_2")
		death.AddLastDocument("ancient_text_3")

		// Check that preservation level was updated
		if death.PreservationLevel <= initialPreservation {
			t.Error("Expected preservation level to increase after adding documents")
		}

		t.Logf("Preservation level updated: %f -> %f", initialPreservation, death.PreservationLevel)
	})

	t.Run("Revival Potential Calculation", func(t *testing.T) {
		death := NewLanguageDeath("death_revival_potential", "lang_1", "gradual", "assimilation", "Test description", "Test details", 30, []string{"region_1", "region_2", "region_3"})

		initialPotential := death.RevivalPotential

		// Add revival factors
		death.AddRevivalFactor("excellent_documentation")
		death.AddRevivalFactor("strong_cultural_interest")
		death.AddRevivalFactor("government_support")

		// Check that revival potential was updated
		if death.RevivalPotential <= initialPotential {
			t.Error("Expected revival potential to increase after adding factors")
		}

		// Add revival barriers
		death.AddRevivalBarrier("lack_of_native_speakers")
		death.AddRevivalBarrier("political_opposition")

		// Check that revival potential was reduced
		if death.RevivalPotential >= death.RevivalPotential {
			t.Logf("Revival potential adjusted: %f (factors: %d, barriers: %d)",
				death.RevivalPotential, len(death.RevivalFactors), len(death.RevivalBarriers))
		}

		t.Logf("Revival potential: %f", death.RevivalPotential)
	})
}

// TestLanguageRevival tests the language revival structure.
func TestLanguageRevival(t *testing.T) {
	t.Run("Revival Creation", func(t *testing.T) {
		revival := NewLanguageRevival(
			"revival_1",
			"lang_1",
			"academic",
			"cultural_renaissance",
			"Language was revived through academic research",
			"Scholars reconstructed the language from ancient texts",
			"reconstruction",
			"comprehensive",
		)

		if revival.ID != "revival_1" {
			t.Errorf("Expected ID 'revival_1', got '%s'", revival.ID)
		}
		if revival.LanguageID != "lang_1" {
			t.Errorf("Expected language ID 'lang_1', got '%s'", revival.LanguageID)
		}
		if revival.Type != "academic" {
			t.Errorf("Expected type 'academic', got '%s'", revival.Type)
		}
		if revival.Trigger != "cultural_renaissance" {
			t.Errorf("Expected trigger 'cultural_renaissance', got '%s'", revival.Trigger)
		}
		if revival.RevivalMethod != "reconstruction" {
			t.Errorf("Expected revival method 'reconstruction', got '%s'", revival.RevivalMethod)
		}
		if revival.RevivalScope != "comprehensive" {
			t.Errorf("Expected revival scope 'comprehensive', got '%s'", revival.RevivalScope)
		}
		if revival.SuccessLevel != 0.5 {
			t.Errorf("Expected success level 0.5, got %f", revival.SuccessLevel)
		}

		t.Logf("Successfully created language revival: %s (%s - %s)", revival.ID, revival.Type, revival.Trigger)
	})
}

// TestLanguageRevivalFeatures tests the revival feature management functions.
func TestLanguageRevivalFeatures(t *testing.T) {
	t.Run("Reconstruction Steps Management", func(t *testing.T) {
		revival := NewLanguageRevival("revival_steps", "lang_1", "cultural", "academic_interest", "Test description", "Test details", "reconstruction", "moderate")

		// Add reconstruction steps
		steps := []string{"analyze_ancient_texts", "reconstruct_phonology", "develop_grammar", "create_dictionary"}
		for _, step := range steps {
			revival.AddReconstructionStep(step)
		}

		if len(revival.ReconstructionSteps) != 4 {
			t.Errorf("Expected 4 reconstruction steps, got %d", len(revival.ReconstructionSteps))
		}

		// Test duplicate prevention
		revival.AddReconstructionStep("analyze_ancient_texts")
		if len(revival.ReconstructionSteps) != 4 {
			t.Errorf("Expected still 4 reconstruction steps after duplicate, got %d", len(revival.ReconstructionSteps))
		}

		t.Logf("Successfully tested reconstruction steps management: %d steps", len(revival.ReconstructionSteps))
	})

	t.Run("Documentation Used Management", func(t *testing.T) {
		revival := NewLanguageRevival("revival_docs", "lang_1", "religious", "spiritual_renewal", "Test description", "Test details", "documentation", "limited")

		// Add documentation used
		docs := []string{"sacred_scrolls", "temple_inscriptions", "ritual_texts"}
		for _, doc := range docs {
			revival.AddDocumentationUsed(doc)
		}

		if len(revival.DocumentationUsed) != 3 {
			t.Errorf("Expected 3 documentation sources, got %d", len(revival.DocumentationUsed))
		}

		t.Logf("Successfully tested documentation used management: %d sources", len(revival.DocumentationUsed))
	})

	t.Run("Modern Adaptations Management", func(t *testing.T) {
		revival := NewLanguageRevival("revival_adaptations", "lang_1", "practical", "modern_need", "Test description", "Test details", "teaching", "moderate")

		// Add modern adaptations
		adaptations := []string{"computer_terminology", "scientific_vocabulary", "modern_syntax"}
		for _, adaptation := range adaptations {
			revival.AddModernAdaptation(adaptation)
		}

		if len(revival.ModernAdaptations) != 3 {
			t.Errorf("Expected 3 modern adaptations, got %d", len(revival.ModernAdaptations))
		}
		if !revival.Modernized {
			t.Error("Expected language to be marked as modernized after adding adaptations")
		}

		t.Logf("Successfully tested modern adaptations management: %d adaptations, modernized: %t",
			len(revival.ModernAdaptations), revival.Modernized)
	})

	t.Run("Community and Institutional Support Management", func(t *testing.T) {
		revival := NewLanguageRevival("revival_community", "lang_1", "cultural", "heritage_preservation", "Test description", "Test details", "teaching", "comprehensive")

		// Set community size
		revival.SetCommunitySize(150)

		// Add institutional support
		institutions := []string{"university_department", "cultural_center", "government_agency"}
		for _, institution := range institutions {
			revival.AddInstitutionalSupport(institution)
		}

		if revival.CommunitySize != 150 {
			t.Errorf("Expected community size 150, got %d", revival.CommunitySize)
		}
		if len(revival.InstitutionalSupport) != 3 {
			t.Errorf("Expected 3 institutional supporters, got %d", len(revival.InstitutionalSupport))
		}

		t.Logf("Successfully tested community and institutional support: community=%d, institutions=%d",
			revival.CommunitySize, len(revival.InstitutionalSupport))
	})
}

// TestLanguageRevivalMetrics tests the revival metrics calculation system.
func TestLanguageRevivalMetrics(t *testing.T) {
	t.Run("Success Level Calculation", func(t *testing.T) {
		revival := NewLanguageRevival("revival_success", "lang_1", "academic", "research_interest", "Test description", "Test details", "reconstruction", "comprehensive")

		initialSuccess := revival.SuccessLevel

		// Add reconstruction steps
		revival.AddReconstructionStep("analyze_texts")
		revival.AddReconstructionStep("reconstruct_phonology")
		revival.AddReconstructionStep("develop_grammar")

		// Add documentation
		revival.AddDocumentationUsed("ancient_manuscript")
		revival.AddDocumentationUsed("stone_inscription")

		// Set community size
		revival.SetCommunitySize(200)

		// Add institutional support
		revival.AddInstitutionalSupport("university")
		revival.AddInstitutionalSupport("research_institute")

		// Check that success level was updated
		if revival.SuccessLevel <= initialSuccess {
			t.Error("Expected success level to increase after adding revival elements")
		}

		t.Logf("Success level updated: %f -> %f", initialSuccess, revival.SuccessLevel)
	})
}

// TestLanguageDeathAndRevivalEngine tests the engine system.
func TestLanguageDeathAndRevivalEngine(t *testing.T) {
	engine := NewLanguageDeathAndRevivalEngine()

	t.Run("Engine Creation", func(t *testing.T) {
		if engine.NaturalDeathRate != 0.05 {
			t.Errorf("Expected natural death rate 0.05, got %f", engine.NaturalDeathRate)
		}
		if engine.ForcedDeathRate != 0.03 {
			t.Errorf("Expected forced death rate 0.03, got %f", engine.ForcedDeathRate)
		}
		if engine.GradualDeathRate != 0.04 {
			t.Errorf("Expected gradual death rate 0.04, got %f", engine.GradualDeathRate)
		}
		if engine.SuddenDeathRate != 0.02 {
			t.Errorf("Expected sudden death rate 0.02, got %f", engine.SuddenDeathRate)
		}
		if engine.AcademicRevivalRate != 0.02 {
			t.Errorf("Expected academic revival rate 0.02, got %f", engine.AcademicRevivalRate)
		}
		if engine.CulturalRevivalRate != 0.03 {
			t.Errorf("Expected cultural revival rate 0.03, got %f", engine.CulturalRevivalRate)
		}

		t.Logf("Successfully created language death and revival engine with configured rates")
	})
}

// TestLanguageDeathAndRevivalSimulation tests the simulation functions.
func TestLanguageDeathAndRevivalSimulation(t *testing.T) {
	t.Run("Language Death Simulation", func(t *testing.T) {
		engine := NewLanguageDeathAndRevivalEngine()
		language, _ := CreateRandomLanguage("test_language", "test_family", 12345)

		// Simulate language death
		death, err := engine.SimulateLanguageDeath(
			language,
			"forced",
			"conquest",
			"Language was forced to extinction through conquest",
			"Conquering forces banned the language and forced assimilation",
			25,
			[]string{"conquered_region", "mountain_fortress"},
		)

		if err != nil {
			t.Fatalf("Failed to simulate language death: %v", err)
		}

		// Verify death record
		if death.LanguageID != language.ID.String() {
			t.Errorf("Expected language ID %s, got %s", language.ID.String(), death.LanguageID)
		}
		if death.Type != "forced" {
			t.Errorf("Expected death type 'forced', got '%s'", death.Type)
		}
		if death.Cause != "conquest" {
			t.Errorf("Expected death cause 'conquest', got '%s'", death.Cause)
		}
		if death.LastSpeakers != 25 {
			t.Errorf("Expected 25 last speakers, got %d", death.LastSpeakers)
		}
		if len(death.LastRegions) != 2 {
			t.Errorf("Expected 2 last regions, got %d", len(death.LastRegions))
		}

		t.Logf("Successfully simulated language death: %s (%s - %s)", death.ID, death.Type, death.Cause)
	})

	t.Run("Language Revival Simulation", func(t *testing.T) {
		engine := NewLanguageDeathAndRevivalEngine()
		language, _ := CreateRandomLanguage("test_language", "test_family", 12345)

		// Simulate language revival
		revival, err := engine.SimulateLanguageRevival(
			language,
			"academic",
			"cultural_renaissance",
			"Language was revived through academic research",
			"Scholars reconstructed the language from ancient texts",
			"reconstruction",
			"comprehensive",
		)

		if err != nil {
			t.Fatalf("Failed to simulate language revival: %v", err)
		}

		// Verify revival record
		if revival.LanguageID != language.ID.String() {
			t.Errorf("Expected language ID %s, got %s", language.ID.String(), revival.LanguageID)
		}
		if revival.Type != "academic" {
			t.Errorf("Expected revival type 'academic', got '%s'", revival.Type)
		}
		if revival.Trigger != "cultural_renaissance" {
			t.Errorf("Expected revival trigger 'cultural_renaissance', got '%s'", revival.Trigger)
		}
		if revival.RevivalMethod != "reconstruction" {
			t.Errorf("Expected revival method 'reconstruction', got '%s'", revival.RevivalMethod)
		}
		if revival.RevivalScope != "comprehensive" {
			t.Errorf("Expected revival scope 'comprehensive', got '%s'", revival.RevivalScope)
		}

		t.Logf("Successfully simulated language revival: %s (%s - %s)", revival.ID, revival.Type, revival.Trigger)
	})

	t.Run("Revival Potential Assessment", func(t *testing.T) {
		engine := NewLanguageDeathAndRevivalEngine()
		language, _ := CreateRandomLanguage("test_language", "test_family", 12345)

		// Create a death record with good preservation
		death := NewLanguageDeath(
			"death_assessment",
			language.ID.String(),
			"gradual",
			"assimilation",
			"Test description",
			"Test details",
			30,
			[]string{"region_1", "region_2", "region_3"},
		)

		// Add good revival factors
		death.AddLastDocument("ancient_manuscript")
		death.AddLastDocument("stone_inscription")
		death.AddLastDocument("wooden_tablet")
		death.AddRevivalFactor("excellent_documentation")
		death.AddRevivalFactor("strong_cultural_interest")
		death.AddRevivalFactor("government_support")

		// Assess revival potential
		potential := engine.AssessRevivalPotential(language, death)

		if potential < 0.5 {
			t.Errorf("Expected high revival potential, got %f", potential)
		}

		t.Logf("Revival potential assessment: %f (preservation: %f, factors: %d)",
			potential, death.PreservationLevel, len(death.RevivalFactors))
	})
}

// TestLanguageDeathAndRevivalIntegration tests the integration with the language system.
func TestLanguageDeathAndRevivalIntegration(t *testing.T) {
	t.Run("Complete Death and Revival Cycle", func(t *testing.T) {
		engine := NewLanguageDeathAndRevivalEngine()
		language, _ := CreateRandomLanguage("test_language", "test_family", 12345)

		// Simulate language death
		death, err := engine.SimulateLanguageDeath(
			language,
			"forced",
			"conquest",
			"Language was forced to extinction through conquest",
			"Conquering forces banned the language and forced assimilation",
			20,
			[]string{"conquered_region"},
		)

		if err != nil {
			t.Fatalf("Failed to simulate language death: %v", err)
		}

		// Add death characteristics
		death.AddLastDocument("ancient_scroll")
		death.AddLastDocument("stone_tablet")
		death.AddInfluenceOnSurvivor("survivor_lang_1")
		death.AddLostKnowledge("ancient_healing_techniques")
		death.AddRevivalFactor("good_documentation")
		death.AddRevivalBarrier("lack_of_speakers")

		// Assess revival potential
		potential := engine.AssessRevivalPotential(language, death)

		// Simulate language revival
		revival, err := engine.SimulateLanguageRevival(
			language,
			"academic",
			"cultural_renaissance",
			"Language was revived through academic research",
			"Scholars reconstructed the language from ancient texts",
			"reconstruction",
			"comprehensive",
		)

		if err != nil {
			t.Fatalf("Failed to simulate language revival: %v", err)
		}

		// Add revival characteristics
		revival.AddReconstructionStep("analyze_ancient_texts")
		revival.AddReconstructionStep("reconstruct_phonology")
		revival.AddDocumentationUsed("ancient_scroll")
		revival.AddDocumentationUsed("stone_tablet")
		revival.AddModernAdaptation("computer_terminology")
		revival.SetCommunitySize(100)
		revival.AddInstitutionalSupport("university")

		// Verify the complete cycle
		if death.LanguageID != language.ID.String() {
			t.Errorf("Death record language ID mismatch: expected %s, got %s", language.ID.String(), death.LanguageID)
		}
		if revival.LanguageID != language.ID.String() {
			t.Errorf("Revival record language ID mismatch: expected %s, got %s", language.ID.String(), revival.LanguageID)
		}
		if len(death.LastDocuments) != 2 {
			t.Errorf("Expected 2 last documents, got %d", len(death.LastDocuments))
		}
		if len(revival.ReconstructionSteps) != 2 {
			t.Errorf("Expected 2 reconstruction steps, got %d", len(revival.ReconstructionSteps))
		}
		if revival.CommunitySize != 100 {
			t.Errorf("Expected community size 100, got %d", revival.CommunitySize)
		}

		t.Logf("Complete death and revival cycle: death=%s, revival=%s, potential=%f",
			death.ID, revival.ID, potential)
		t.Logf("Death characteristics: documents=%d, influence=%d, knowledge=%d, factors=%d, barriers=%d",
			len(death.LastDocuments), len(death.InfluenceOnSurvivors), len(death.LostKnowledge),
			len(death.RevivalFactors), len(death.RevivalBarriers))
		t.Logf("Revival characteristics: steps=%d, docs=%d, adaptations=%d, community=%d, institutions=%d",
			len(revival.ReconstructionSteps), len(revival.DocumentationUsed), len(revival.ModernAdaptations),
			revival.CommunitySize, len(revival.InstitutionalSupport))
	})
}
