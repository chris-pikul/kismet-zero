package lang

import (
	"testing"
)

// TestPidginLanguage tests the pidgin language structure.
func TestPidginLanguage(t *testing.T) {
	t.Run("Pidgin Creation", func(t *testing.T) {
		pidgin := NewPidginLanguage(
			"pidgin_1",
			"Trade Pidgin",
			"A simplified language for trade between different groups",
			"colonial",
			"trade",
			"marketplace",
			[]string{"lang_1", "lang_2"},
			"lang_1",
			"lang_2",
		)

		if pidgin.ID != "pidgin_1" {
			t.Errorf("Expected ID 'pidgin_1', got '%s'", pidgin.ID)
		}
		if pidgin.Name != "Trade Pidgin" {
			t.Errorf("Expected name 'Trade Pidgin', got '%s'", pidgin.Name)
		}
		if pidgin.Type != "trade" {
			t.Errorf("Expected type 'trade', got '%s'", pidgin.Type)
		}
		if pidgin.Context != "marketplace" {
			t.Errorf("Expected context 'marketplace', got '%s'", pidgin.Context)
		}
		if len(pidgin.SourceLanguages) != 2 {
			t.Errorf("Expected 2 source languages, got %d", len(pidgin.SourceLanguages))
		}
		if pidgin.PrimarySource != "lang_1" {
			t.Errorf("Expected primary source 'lang_1', got '%s'", pidgin.PrimarySource)
		}
		if pidgin.SecondarySource != "lang_2" {
			t.Errorf("Expected secondary source 'lang_2', got '%s'", pidgin.SecondarySource)
		}
		if pidgin.ComplexityLevel != 0.3 {
			t.Errorf("Expected complexity level 0.3, got %f", pidgin.ComplexityLevel)
		}
		if pidgin.StabilityLevel != 0.4 {
			t.Errorf("Expected stability level 0.4, got %f", pidgin.StabilityLevel)
		}
		if pidgin.FunctionalityLevel != 0.5 {
			t.Errorf("Expected functionality level 0.5, got %f", pidgin.FunctionalityLevel)
		}

		t.Logf("Successfully created pidgin: %s (%s - %s)", pidgin.ID, pidgin.Type, pidgin.Context)
	})
}

// TestPidginLanguageFeatures tests the feature management functions.
func TestPidginLanguageFeatures(t *testing.T) {
	t.Run("Phonological Features Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_phon", "Test Pidgin", "Test description", "test", "work", "factory", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add phonological features
		features := []string{"simplified_phonology", "basic_sounds", "reduced_phonemes"}
		for _, feature := range features {
			pidgin.AddPhonologicalFeature(feature)
		}

		if len(pidgin.PhonologicalFeatures) != 3 {
			t.Errorf("Expected 3 phonological features, got %d", len(pidgin.PhonologicalFeatures))
		}

		// Test duplicate prevention
		pidgin.AddPhonologicalFeature("simplified_phonology")
		if len(pidgin.PhonologicalFeatures) != 3 {
			t.Errorf("Expected still 3 phonological features after duplicate, got %d", len(pidgin.PhonologicalFeatures))
		}

		t.Logf("Successfully tested phonological features management: %d features", len(pidgin.PhonologicalFeatures))
	})

	t.Run("Morphological Features Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_morph", "Test Pidgin", "Test description", "test", "social", "community", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add morphological features
		features := []string{"reduced_morphology", "basic_word_formation", "simple_affixes"}
		for _, feature := range features {
			pidgin.AddMorphologicalFeature(feature)
		}

		if len(pidgin.MorphologicalFeatures) != 3 {
			t.Errorf("Expected 3 morphological features, got %d", len(pidgin.MorphologicalFeatures))
		}

		t.Logf("Successfully tested morphological features management: %d features", len(pidgin.MorphologicalFeatures))
	})

	t.Run("Syntactic Features Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_syntax", "Test Pidgin", "Test description", "test", "military", "barracks", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add syntactic features
		features := []string{"basic_syntax", "simple_sentence_structure", "reduced_clauses"}
		for _, feature := range features {
			pidgin.AddSyntacticFeature(feature)
		}

		if len(pidgin.SyntacticFeatures) != 3 {
			t.Errorf("Expected 3 syntactic features, got %d", len(pidgin.SyntacticFeatures))
		}

		t.Logf("Successfully tested syntactic features management: %d features", len(pidgin.SyntacticFeatures))
	})

	t.Run("Lexical Features Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_lex", "Test Pidgin", "Test description", "test", "trade", "port", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add lexical features
		features := []string{"core_vocabulary", "basic_terms", "essential_words"}
		for _, feature := range features {
			pidgin.AddLexicalFeature(feature)
		}

		if len(pidgin.LexicalFeatures) != 3 {
			t.Errorf("Expected 3 lexical features, got %d", len(pidgin.LexicalFeatures))
		}

		t.Logf("Successfully tested lexical features management: %d features", len(pidgin.LexicalFeatures))
	})
}

// TestPidginLanguageUsage tests the usage management functions.
func TestPidginLanguageUsage(t *testing.T) {
	t.Run("Speaker Count Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_speakers", "Test Pidgin", "Test description", "test", "work", "plantation", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Set speaker count
		pidgin.SetSpeakerCount(150)

		if pidgin.SpeakerCount != 150 {
			t.Errorf("Expected speaker count 150, got %d", pidgin.SpeakerCount)
		}

		t.Logf("Successfully tested speaker count management: %d speakers", pidgin.SpeakerCount)
	})

	t.Run("Usage Regions Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_regions", "Test Pidgin", "Test description", "test", "trade", "marketplace", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add usage regions
		regions := []string{"port_city", "trading_post", "border_region"}
		for _, region := range regions {
			pidgin.AddUsageRegion(region)
		}

		if len(pidgin.UsageRegions) != 3 {
			t.Errorf("Expected 3 usage regions, got %d", len(pidgin.UsageRegions))
		}

		t.Logf("Successfully tested usage regions management: %d regions", len(pidgin.UsageRegions))
	})

	t.Run("Usage Domains Management", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_domains", "Test Pidgin", "Test description", "test", "social", "community", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add usage domains
		domains := []string{"basic_communication", "trade_negotiation", "simple_instructions"}
		for _, domain := range domains {
			pidgin.AddUsageDomain(domain)
		}

		if len(pidgin.UsageDomains) != 3 {
			t.Errorf("Expected 3 usage domains, got %d", len(pidgin.UsageDomains))
		}

		t.Logf("Successfully tested usage domains management: %d domains", len(pidgin.UsageDomains))
	})
}

// TestPidginLanguageMetrics tests the metrics calculation system.
func TestPidginLanguageMetrics(t *testing.T) {
	t.Run("Metrics Calculation", func(t *testing.T) {
		pidgin := NewPidginLanguage("pidgin_metrics", "Test Pidgin", "Test description", "test", "work", "factory", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		initialComplexity := pidgin.ComplexityLevel
		initialStability := pidgin.StabilityLevel
		initialFunctionality := pidgin.FunctionalityLevel

		// Add features to trigger metrics update
		pidgin.AddPhonologicalFeature("feature_1")
		pidgin.AddPhonologicalFeature("feature_2")
		pidgin.AddMorphologicalFeature("feature_3")
		pidgin.AddSyntacticFeature("feature_4")
		pidgin.AddLexicalFeature("feature_5")

		// Set speaker count
		pidgin.SetSpeakerCount(200)

		// Add usage information
		pidgin.AddUsageDomain("domain_1")
		pidgin.AddUsageDomain("domain_2")

		// Check that metrics were updated
		if pidgin.ComplexityLevel <= initialComplexity {
			t.Error("Expected complexity level to increase after adding features")
		}
		if pidgin.StabilityLevel <= initialStability {
			t.Error("Expected stability level to increase after adding speakers")
		}
		if pidgin.FunctionalityLevel <= initialFunctionality {
			t.Error("Expected functionality level to increase after adding domains")
		}

		t.Logf("Metrics updated: complexity=%f, stability=%f, functionality=%f",
			pidgin.ComplexityLevel, pidgin.StabilityLevel, pidgin.FunctionalityLevel)
	})
}

// TestCreoleLanguage tests the creole language structure.
func TestCreoleLanguage(t *testing.T) {
	t.Run("Creole Creation", func(t *testing.T) {
		creole := NewCreoleLanguage(
			"creole_1",
			"Plantation Creole",
			"A creole language developed from a pidgin",
			"colonial",
			"plantation",
			"agricultural",
			"pidgin_1",
			[]string{"lang_1", "lang_2"},
			"lang_1",
			"lang_2",
		)

		if creole.ID != "creole_1" {
			t.Errorf("Expected ID 'creole_1', got '%s'", creole.ID)
		}
		if creole.Name != "Plantation Creole" {
			t.Errorf("Expected name 'Plantation Creole', got '%s'", creole.Name)
		}
		if creole.Type != "plantation" {
			t.Errorf("Expected type 'plantation', got '%s'", creole.Type)
		}
		if creole.Context != "agricultural" {
			t.Errorf("Expected context 'agricultural', got '%s'", creole.Context)
		}
		if creole.PidginID != "pidgin_1" {
			t.Errorf("Expected pidgin ID 'pidgin_1', got '%s'", creole.PidginID)
		}
		if len(creole.SourceLanguages) != 2 {
			t.Errorf("Expected 2 source languages, got %d", len(creole.SourceLanguages))
		}
		if creole.ComplexityLevel != 0.6 {
			t.Errorf("Expected complexity level 0.6, got %f", creole.ComplexityLevel)
		}
		if creole.StabilityLevel != 0.7 {
			t.Errorf("Expected stability level 0.7, got %f", creole.StabilityLevel)
		}
		if creole.FunctionalityLevel != 0.8 {
			t.Errorf("Expected functionality level 0.8, got %f", creole.FunctionalityLevel)
		}

		t.Logf("Successfully created creole: %s (%s - %s)", creole.ID, creole.Type, creole.Context)
	})
}

// TestCreoleLanguageFeatures tests the creole feature management functions.
func TestCreoleLanguageFeatures(t *testing.T) {
	t.Run("Feature Management", func(t *testing.T) {
		creole := NewCreoleLanguage("creole_features", "Test Creole", "Test description", "test", "urban", "city", "pidgin_1", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add features
		creole.AddPhonologicalFeature("developed_phonology")
		creole.AddMorphologicalFeature("developed_morphology")
		creole.AddSyntacticFeature("developed_syntax")
		creole.AddLexicalFeature("expanded_vocabulary")

		if len(creole.PhonologicalFeatures) != 1 {
			t.Errorf("Expected 1 phonological feature, got %d", len(creole.PhonologicalFeatures))
		}
		if len(creole.MorphologicalFeatures) != 1 {
			t.Errorf("Expected 1 morphological feature, got %d", len(creole.MorphologicalFeatures))
		}
		if len(creole.SyntacticFeatures) != 1 {
			t.Errorf("Expected 1 syntactic feature, got %d", len(creole.SyntacticFeatures))
		}
		if len(creole.LexicalFeatures) != 1 {
			t.Errorf("Expected 1 lexical feature, got %d", len(creole.LexicalFeatures))
		}

		t.Logf("Successfully tested creole feature management: %d total features",
			len(creole.PhonologicalFeatures)+len(creole.MorphologicalFeatures)+len(creole.SyntacticFeatures)+len(creole.LexicalFeatures))
	})
}

// TestCreoleLanguageUsage tests the creole usage management functions.
func TestCreoleLanguageUsage(t *testing.T) {
	t.Run("Speaker Count Management", func(t *testing.T) {
		creole := NewCreoleLanguage("creole_speakers", "Test Creole", "Test description", "test", "urban", "city", "pidgin_1", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Set speaker counts
		creole.SetNativeSpeakerCount(500)
		creole.SetTotalSpeakerCount(800)

		if creole.NativeSpeakerCount != 500 {
			t.Errorf("Expected native speaker count 500, got %d", creole.NativeSpeakerCount)
		}
		if creole.TotalSpeakerCount != 800 {
			t.Errorf("Expected total speaker count 800, got %d", creole.TotalSpeakerCount)
		}

		t.Logf("Successfully tested speaker count management: native=%d, total=%d",
			creole.NativeSpeakerCount, creole.TotalSpeakerCount)
	})

	t.Run("Usage Information Management", func(t *testing.T) {
		creole := NewCreoleLanguage("creole_usage", "Test Creole", "Test description", "test", "urban", "city", "pidgin_1", []string{"lang_1", "lang_2"}, "lang_1", "lang_2")

		// Add usage information
		creole.AddUsageRegion("urban_center")
		creole.AddUsageRegion("suburban_area")
		creole.AddUsageDomain("daily_communication")
		creole.AddUsageDomain("education")
		creole.AddUsageDomain("media")

		if len(creole.UsageRegions) != 2 {
			t.Errorf("Expected 2 usage regions, got %d", len(creole.UsageRegions))
		}
		if len(creole.UsageDomains) != 3 {
			t.Errorf("Expected 3 usage domains, got %d", len(creole.UsageDomains))
		}

		t.Logf("Successfully tested usage information management: regions=%d, domains=%d",
			len(creole.UsageRegions), len(creole.UsageDomains))
	})
}

// TestPidginAndCreoleFormationEngine tests the formation engine.
func TestPidginAndCreoleFormationEngine(t *testing.T) {
	engine := NewPidginAndCreoleFormationEngine()

	t.Run("Engine Creation", func(t *testing.T) {
		if engine.PidginFormationRate != 0.3 {
			t.Errorf("Expected pidgin formation rate 0.3, got %f", engine.PidginFormationRate)
		}
		if engine.CreoleFormationRate != 0.1 {
			t.Errorf("Expected creole formation rate 0.1, got %f", engine.CreoleFormationRate)
		}
		if engine.TradeContactWeight != 0.4 {
			t.Errorf("Expected trade contact weight 0.4, got %f", engine.TradeContactWeight)
		}
		if engine.WorkContactWeight != 0.3 {
			t.Errorf("Expected work contact weight 0.3, got %f", engine.WorkContactWeight)
		}
		if engine.SocialContactWeight != 0.2 {
			t.Errorf("Expected social contact weight 0.2, got %f", engine.SocialContactWeight)
		}
		if engine.MilitaryContactWeight != 0.1 {
			t.Errorf("Expected military contact weight 0.1, got %f", engine.MilitaryContactWeight)
		}

		t.Logf("Successfully created pidgin and creole formation engine with configured rates")
	})
}

// TestPidginAndCreoleFormationSimulation tests the formation simulation functions.
func TestPidginAndCreoleFormationSimulation(t *testing.T) {
	t.Run("Pidgin Formation Simulation", func(t *testing.T) {
		engine := NewPidginAndCreoleFormationEngine()

		// Simulate pidgin formation
		pidgin, err := engine.SimulatePidginFormation(
			[]string{"lang_1", "lang_2"},
			"lang_1",
			"lang_2",
			"trade",
			"marketplace",
			"colonial",
		)

		if err != nil {
			t.Fatalf("Failed to simulate pidgin formation: %v", err)
		}

		// Verify pidgin
		if pidgin.Type != "trade" {
			t.Errorf("Expected pidgin type 'trade', got '%s'", pidgin.Type)
		}
		if pidgin.Context != "marketplace" {
			t.Errorf("Expected pidgin context 'marketplace', got '%s'", pidgin.Context)
		}
		if len(pidgin.SourceLanguages) != 2 {
			t.Errorf("Expected 2 source languages, got %d", len(pidgin.SourceLanguages))
		}
		if pidgin.PrimarySource != "lang_1" {
			t.Errorf("Expected primary source 'lang_1', got '%s'", pidgin.PrimarySource)
		}
		if pidgin.SecondarySource != "lang_2" {
			t.Errorf("Expected secondary source 'lang_2', got '%s'", pidgin.SecondarySource)
		}

		t.Logf("Successfully simulated pidgin formation: %s (%s - %s)", pidgin.ID, pidgin.Type, pidgin.Context)
	})

	t.Run("Creole Formation Simulation", func(t *testing.T) {
		engine := NewPidginAndCreoleFormationEngine()

		// Create a pidgin first
		pidgin, err := engine.SimulatePidginFormation(
			[]string{"lang_1", "lang_2"},
			"lang_1",
			"lang_2",
			"work",
			"plantation",
			"colonial",
		)

		if err != nil {
			t.Fatalf("Failed to simulate pidgin formation: %v", err)
		}

		// Simulate creole formation from pidgin
		creole, err := engine.SimulateCreoleFormation(
			pidgin,
			"plantation",
			"agricultural",
			"colonial",
		)

		if err != nil {
			t.Fatalf("Failed to simulate creole formation: %v", err)
		}

		// Verify creole
		if creole.Type != "plantation" {
			t.Errorf("Expected creole type 'plantation', got '%s'", creole.Type)
		}
		if creole.Context != "agricultural" {
			t.Errorf("Expected creole context 'agricultural', got '%s'", creole.Context)
		}
		if creole.PidginID != pidgin.ID {
			t.Errorf("Expected pidgin ID %s, got %s", pidgin.ID, creole.PidginID)
		}
		if len(creole.SourceLanguages) != 2 {
			t.Errorf("Expected 2 source languages, got %d", len(creole.SourceLanguages))
		}

		t.Logf("Successfully simulated creole formation: %s (%s - %s)", creole.ID, creole.Type, creole.Context)
	})

	t.Run("Creolization Potential Assessment", func(t *testing.T) {
		engine := NewPidginAndCreoleFormationEngine()

		// Create a pidgin with good creolization potential
		pidgin, err := engine.SimulatePidginFormation(
			[]string{"lang_1", "lang_2"},
			"lang_1",
			"lang_2",
			"work",
			"plantation",
			"colonial",
		)

		if err != nil {
			t.Fatalf("Failed to simulate pidgin formation: %v", err)
		}

		// Set good creolization factors
		pidgin.SetSpeakerCount(200)
		pidgin.AddUsageDomain("daily_communication")
		pidgin.AddUsageDomain("work_instructions")

		// Assess creolization potential
		potential := engine.AssessCreolizationPotential(pidgin)

		if potential < 0.5 {
			t.Errorf("Expected high creolization potential, got %f", potential)
		}

		t.Logf("Creolization potential assessment: %f (speakers=%d, domains=%d)",
			potential, pidgin.SpeakerCount, len(pidgin.UsageDomains))
	})
}

// TestPidginAndCreoleIntegration tests the integration with the language system.
func TestPidginAndCreoleIntegration(t *testing.T) {
	t.Run("Complete Pidgin to Creole Cycle", func(t *testing.T) {
		engine := NewPidginAndCreoleFormationEngine()

		// Simulate pidgin formation
		pidgin, err := engine.SimulatePidginFormation(
			[]string{"lang_1", "lang_2"},
			"lang_1",
			"lang_2",
			"trade",
			"port_city",
			"colonial",
		)

		if err != nil {
			t.Fatalf("Failed to simulate pidgin formation: %v", err)
		}

		// Add pidgin characteristics
		pidgin.SetSpeakerCount(300)
		pidgin.AddUsageRegion("port_area")
		pidgin.AddUsageRegion("market_district")
		pidgin.AddUsageDomain("trade_negotiation")
		pidgin.AddUsageDomain("basic_communication")

		// Assess creolization potential
		potential := engine.AssessCreolizationPotential(pidgin)

		// Simulate creole formation
		creole, err := engine.SimulateCreoleFormation(
			pidgin,
			"urban",
			"commercial",
			"colonial",
		)

		if err != nil {
			t.Fatalf("Failed to simulate creole formation: %v", err)
		}

		// Add creole characteristics
		creole.SetNativeSpeakerCount(800)
		creole.SetTotalSpeakerCount(1200)
		creole.AddUsageRegion("urban_center")
		creole.AddUsageRegion("suburban_area")
		creole.AddUsageDomain("daily_communication")
		creole.AddUsageDomain("education")
		creole.AddUsageDomain("media")

		// Verify the complete cycle
		if pidgin.Type != "trade" {
			t.Errorf("Expected pidgin type 'trade', got '%s'", pidgin.Type)
		}
		if creole.Type != "urban" {
			t.Errorf("Expected creole type 'urban', got '%s'", creole.Type)
		}
		if creole.PidginID != pidgin.ID {
			t.Errorf("Expected creole pidgin ID %s, got %s", pidgin.ID, creole.PidginID)
		}
		if pidgin.SpeakerCount != 300 {
			t.Errorf("Expected pidgin speaker count 300, got %d", pidgin.SpeakerCount)
		}
		if creole.NativeSpeakerCount != 800 {
			t.Errorf("Expected creole native speaker count 800, got %d", creole.NativeSpeakerCount)
		}

		t.Logf("Complete pidgin to creole cycle: pidgin=%s, creole=%s, potential=%f",
			pidgin.ID, creole.ID, potential)
		t.Logf("Pidgin characteristics: speakers=%d, regions=%d, domains=%d",
			pidgin.SpeakerCount, len(pidgin.UsageRegions), len(pidgin.UsageDomains))
		t.Logf("Creole characteristics: native_speakers=%d, total_speakers=%d, regions=%d, domains=%d",
			creole.NativeSpeakerCount, creole.TotalSpeakerCount, len(creole.UsageRegions), len(creole.UsageDomains))
	})
}
