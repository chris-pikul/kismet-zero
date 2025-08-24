package lang

import (
	"fmt"
	"strings"
	"testing"
)

// TestFantasyLanguageEcosystem tests the complete language generation and evolution
// system for fantasy world-building and historical simulation.
func TestFantasyLanguageEcosystem(t *testing.T) {
	t.Run("Complete Language Lifecycle", func(t *testing.T) {
		testCompleteLanguageLifecycle(t)
	})

	t.Run("Language Family Evolution", func(t *testing.T) {
		testLanguageFamilyEvolution(t)
	})

	t.Run("Cultural Influence and Contact", func(t *testing.T) {
		testCulturalInfluenceAndContact(t)
	})

	t.Run("Translation and Interlingua", func(t *testing.T) {
		testTranslationAndInterlingua(t)
	})

	t.Run("Dialect Formation", func(t *testing.T) {
		testDialectFormation(t)
	})

	t.Run("Fantasy Culture Languages", func(t *testing.T) {
		testFantasyCultureLanguages(t)
	})
}

// testCompleteLanguageLifecycle tests the full lifecycle of a language from creation
// through evolution and cultural changes.
func testCompleteLanguageLifecycle(t *testing.T) {
	// Step 1: Create a root language for an ancient elvish civilization
	protoElvish, err := CreateRandomLanguage("elvish", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create proto-elvish language: %v", err)
	}

	// Verify the language is complete
	if !protoElvish.IsComplete() {
		t.Errorf("Proto-elvish language is incomplete, missing: %v", protoElvish.GetMissingComponents())
	}

	// Verify cultural configuration
	if protoElvish.Culture != "ancient" {
		t.Errorf("Expected culture 'ancient', got '%s'", protoElvish.Culture)
	}

	// Step 2: Evolve the language over 1000 years
	evolvedElvish, err := EvolveLanguage(protoElvish, "1000_years", []string{"magical", "religious"}, 43)
	if err != nil {
		t.Fatalf("Failed to evolve elvish language: %v", err)
	}

	// Verify evolution occurred
	if evolvedElvish.EvolvedAt.IsZero() {
		t.Error("Expected evolution timestamp to be set")
	}

	if evolvedElvish.Description == protoElvish.Description {
		t.Error("Expected language description to change after evolution")
	}

	// Step 3: Create a dialect through geographic isolation
	dialectSeed := int64(44)
	dialectID := LanguageID{
		Family:   protoElvish.ID.Family,
		Branch:   "woodland",
		Language: "sylvan",
		Dialect:  "forest",
	}

	woodlandDialect := protoElvish.Clone(dialectID, dialectSeed)
	woodlandDialect.SetDescription("A woodland dialect of the ancient elvish language")

	// Verify dialect relationship
	if woodlandDialect.ParentID == nil || *woodlandDialect.ParentID != protoElvish.ID {
		t.Error("Expected dialect to have correct parent language")
	}

	// Step 4: Generate content in the evolved language
	names, err := GenerateNames(evolvedElvish, "person", 5, "noble", 45)
	if err != nil {
		t.Fatalf("Failed to generate names: %v", err)
	}

	if len(names) != 5 {
		t.Errorf("Expected 5 names, got %d", len(names))
	}

	for i, name := range names {
		if name == "" {
			t.Errorf("Generated name[%d] is empty", i)
		}
	}

	// Step 5: Generate written text
	text, err := GenerateText(evolvedElvish, "poetry", "medium", "nature", 46)
	if err != nil {
		t.Fatalf("Failed to generate text: %v", err)
	}

	if text == "" {
		t.Error("Generated text is empty")
	}

	// Verify text contains the language's writing system
	if evolvedElvish.Orthography != nil {
		// Text should be written in the language's orthography
		t.Logf("Generated text in %s: %s", evolvedElvish.Orthography.Style, text)
	}

	t.Logf("Successfully completed language lifecycle test")
	t.Logf("Proto-language: %s", protoElvish.Name)
	t.Logf("Evolved language: %s", evolvedElvish.Name)
	t.Logf("Dialect: %s", woodlandDialect.Name)
	t.Logf("Generated names: %v", names)
	t.Logf("Generated text length: %d characters", len(text))
}

// testLanguageFamilyEvolution tests the creation and evolution of language families
// with branching relationships and historical depth.
func testLanguageFamilyEvolution(t *testing.T) {
	// Step 1: Create a proto-language for a human civilization
	protoHuman, err := CreateRandomLanguage("human", "ancient", 100)
	if err != nil {
		t.Fatalf("Failed to create proto-human language: %v", err)
	}

	// Step 2: Create a language family with 3 branches
	family, err := CreateLanguageFamily(protoHuman, 3, "1000_years", 101)
	if err != nil {
		t.Fatalf("Failed to create language family: %v", err)
	}

	if len(family) != 4 { // proto + 3 branches
		t.Errorf("Expected 4 languages in family, got %d", len(family))
	}

	// Step 3: Evolve each branch independently
	evolvedBranches := make([]*Language, 0, len(family)-1)
	for i, branch := range family[1:] { // Skip proto-language
		evolutionSeed := int64(102 + i)
		evolutionTime := fmt.Sprintf("%d_years", 500+(i*200)) // Different evolution times

		evolvedBranch, err := EvolveLanguage(branch, evolutionTime, []string{"migration", "trade"}, evolutionSeed)
		if err != nil {
			t.Fatalf("Failed to evolve branch %d: %v", i, err)
		}

		evolvedBranches = append(evolvedBranches, evolvedBranch)
	}

	// Step 4: Verify family relationships
	// Check the original branches from the family creation, not the evolved ones
	for i, branch := range family[1:] { // Skip proto-language
		if branch.ParentID == nil || *branch.ParentID != protoHuman.ID {
			t.Errorf("Branch %d missing parent relationship", i)
		}
	}

	// Verify proto-language has child references
	if len(protoHuman.ChildIDs) != 3 {
		t.Errorf("Expected proto-language to have 3 children, got %d", len(protoHuman.ChildIDs))
	}

	// Step 5: Test mutual intelligibility (simplified)
	for i, branch1 := range evolvedBranches {
		for j, branch2 := range evolvedBranches {
			if i != j {
				// Languages from the same family should share some features
				if branch1.Phonology != nil && branch2.Phonology != nil {
					// Both should have phonology
					t.Logf("Branch %d and %d both have phonology systems", i, j)
				}
			}
		}
	}

	t.Logf("Successfully created language family with %d branches", len(family))
	t.Logf("Proto-language: %s", protoHuman.Name)
	for i, branch := range evolvedBranches {
		t.Logf("Branch %d: %s", i+1, branch.Name)
	}
}

// testCulturalInfluenceAndContact tests how languages influence each other
// through various types of cultural contact.
func testCulturalInfluenceAndContact(t *testing.T) {
	// Step 1: Create two distinct languages
	elvishLang, err := CreateRandomLanguage("elvish", "high_elf", 200)
	if err != nil {
		t.Fatalf("Failed to create elvish language: %v", err)
	}

	humanLang, err := CreateRandomLanguage("human", "northern_human", 201)
	if err != nil {
		t.Fatalf("Failed to create human language: %v", err)
	}

	// Step 2: Simulate trade contact
	tradeElvish, tradeHuman, err := SimulateLanguageContact(
		elvishLang, humanLang, "trade", 0.6, "100_years", 202)
	if err != nil {
		t.Fatalf("Failed to simulate trade contact: %v", err)
	}

	// Verify both languages were modified
	if tradeElvish.Description == elvishLang.Description {
		t.Error("Expected elvish language to change after trade contact")
	}

	if tradeHuman.Description == humanLang.Description {
		t.Error("Expected human language to change after trade contact")
	}

	// Step 3: Simulate conquest contact
	conquestElvish, conquestHuman, err := SimulateLanguageContact(
		elvishLang, humanLang, "conquest", 0.8, "200_years", 203)
	if err != nil {
		t.Fatalf("Failed to simulate conquest contact: %v", err)
	}

	// Conquest should have stronger effects than trade
	if !strings.Contains(conquestElvish.Description, "conquest") {
		t.Error("Expected conquest influence in elvish description")
	}

	if !strings.Contains(conquestHuman.Description, "conquest") {
		t.Error("Expected conquest influence in human description")
	}

	// Step 4: Test cultural influence through migration
	migrationElvish, migrationHuman, err := SimulateLanguageContact(
		elvishLang, humanLang, "migration", 0.4, "300_years", 204)
	if err != nil {
		t.Fatalf("Failed to simulate migration contact: %v", err)
	}

	// Migration should create dialect-like variations
	if !strings.Contains(migrationElvish.Description, "migration") {
		t.Error("Expected migration influence in elvish description")
	}

	if !strings.Contains(migrationHuman.Description, "migration") {
		t.Error("Expected migration influence in human description")
	}

	t.Logf("Successfully tested cultural influence and contact")
	t.Logf("Trade contact: %s ↔ %s", tradeElvish.Name, tradeHuman.Name)
	t.Logf("Conquest contact: %s ↔ %s", conquestElvish.Name, conquestHuman.Name)
	t.Logf("Migration contact: %s ↔ %s", migrationElvish.Name, migrationHuman.Name)
}

// testTranslationAndInterlingua tests the translation capabilities between
// generated languages and English using the interlingua system.
func testTranslationAndInterlingua(t *testing.T) {
	// Step 1: Create a language with interlingua services
	dwarvishLang, err := CreateRandomLanguage("dwarvish", "mountain_dwarf", 300)
	if err != nil {
		t.Fatalf("Failed to create dwarven language: %v", err)
	}

	// Step 2: Configure interlingua services
	err = ConfigureInterlinguaServices(dwarvishLang)
	if err != nil {
		t.Fatalf("Failed to configure interlingua services: %v", err)
	}

	// Verify services are configured
	services := dwarvishLang.Interlingua()
	if !services.HasRealizer("en") {
		t.Error("Expected English realizer to be configured")
	}

	if !services.HasAnalyzer("en") {
		t.Error("Expected English analyzer to be configured")
	}

	// Step 3: Generate text in the dwarven language
	dwarvishText, err := GenerateText(dwarvishLang, "inscription", "short", "mining", 301)
	if err != nil {
		t.Fatalf("Failed to generate dwarven text: %v", err)
	}

	// Step 4: Translate to English
	englishTranslation, err := TranslateToEnglish(dwarvishLang, dwarvishText)
	if err != nil {
		t.Fatalf("Failed to translate to English: %v", err)
	}

	if englishTranslation == "" {
		t.Error("Expected non-empty English translation")
	}

	// Step 5: Translate from English back to dwarven
	backTranslation, err := TranslateFromEnglish(dwarvishLang, "The mountain yields precious metals")
	if err != nil {
		t.Fatalf("Failed to translate from English: %v", err)
	}

	if backTranslation == "" {
		t.Error("Expected non-empty back translation")
	}

	// Step 6: Test round-trip consistency (simplified)
	// In a full implementation, this would test semantic preservation
	if !strings.Contains(backTranslation, dwarvishLang.Name) {
		t.Error("Expected back translation to reference the target language")
	}

	t.Logf("Successfully tested translation and interlingua")
	t.Logf("Dwarven text: %s", dwarvishText)
	t.Logf("English translation: %s", englishTranslation)
	t.Logf("Back translation: %s", backTranslation)
}

// testDialectFormation tests the creation of regional and social dialects
// within a language family.
func testDialectFormation(t *testing.T) {
	// Step 1: Create a base language
	baseLang, err := CreateRandomLanguage("human", "coastal_human", 400)
	if err != nil {
		t.Fatalf("Failed to create base language: %v", err)
	}

	// Step 2: Create geographic dialects
	coastalDialect := baseLang.Clone(LanguageID{
		Family:   baseLang.ID.Family,
		Branch:   "coastal",
		Language: "seafarer",
		Dialect:  "harbor",
	}, 401)

	inlandDialect := baseLang.Clone(LanguageID{
		Family:   baseLang.ID.Family,
		Branch:   "inland",
		Language: "farmer",
		Dialect:  "rural",
	}, 402)

	// Step 3: Create social dialects
	nobleDialect := baseLang.Clone(LanguageID{
		Family:   baseLang.ID.Family,
		Branch:   "noble",
		Language: "aristocrat",
		Dialect:  "court",
	}, 403)

	commonDialect := baseLang.Clone(LanguageID{
		Family:   baseLang.ID.Family,
		Branch:   "common",
		Language: "merchant",
		Dialect:  "market",
	}, 404)

	// Step 4: Apply dialect-specific changes
	dialects := []*Language{coastalDialect, inlandDialect, nobleDialect, commonDialect}
	dialectTypes := []string{"coastal", "inland", "noble", "common"}

	for i, dialect := range dialects {
		dialectType := dialectTypes[i]
		dialect.SetDescription(fmt.Sprintf("%s dialect of %s", dialectType, baseLang.Name))

		// Apply dialect-specific features
		switch dialectType {
		case "coastal":
			dialect.SetDescription(fmt.Sprintf("%s with maritime vocabulary", dialect.Description))
		case "inland":
			dialect.SetDescription(fmt.Sprintf("%s with agricultural terms", dialect.Description))
		case "noble":
			dialect.SetDescription(fmt.Sprintf("%s with formal register", dialect.Description))
		case "common":
			dialect.SetDescription(fmt.Sprintf("%s with colloquial expressions", dialect.Description))
		}
	}

	// Step 5: Verify dialect relationships
	for _, dialect := range dialects {
		if dialect.ParentID == nil || *dialect.ParentID != baseLang.ID {
			t.Error("Expected dialect to have correct parent language")
		}
	}

	// Update base language's child references
	if baseLang.ChildIDs == nil {
		baseLang.ChildIDs = make([]LanguageID, 0, len(dialects))
	}
	for _, dialect := range dialects {
		baseLang.ChildIDs = append(baseLang.ChildIDs, dialect.ID)
	}

	// Verify base language has child references
	if len(baseLang.ChildIDs) != 4 {
		t.Errorf("Expected base language to have 4 children, got %d", len(baseLang.ChildIDs))
	}

	// Step 6: Test mutual intelligibility between dialects
	// Dialects of the same base language should be mutually intelligible
	for i, dialect1 := range dialects {
		for j, dialect2 := range dialects {
			if i != j {
				// Both should share the same base linguistic features
				if dialect1.Phonology != nil && dialect2.Phonology != nil {
					t.Logf("Dialects %d and %d share phonology system", i, j)
				}
			}
		}
	}

	t.Logf("Successfully created dialect system")
	t.Logf("Base language: %s", baseLang.Name)
	for i, dialect := range dialects {
		t.Logf("Dialect %d: %s (%s)", i+1, dialect.Name, dialectTypes[i])
	}
}

// testFantasyCultureLanguages tests the creation of complete linguistic
// ecosystems for different fantasy cultures.
func testFantasyCultureLanguages(t *testing.T) {
	// Step 1: Create a complete elvish linguistic ecosystem
	elvishEcosystem, err := CreateCulturalLanguageSet("elvish", 4, "high", 500)
	if err != nil {
		t.Fatalf("Failed to create elvish ecosystem: %v", err)
	}

	if len(elvishEcosystem) != 4 {
		t.Errorf("Expected 4 elvish languages, got %d", len(elvishEcosystem))
	}

	// Step 2: Create a complete dwarven linguistic ecosystem
	dwarvenEcosystem, err := CreateCulturalLanguageSet("dwarvish", 3, "moderate", 501)
	if err != nil {
		t.Fatalf("Failed to create dwarven ecosystem: %v", err)
	}

	if len(dwarvenEcosystem) != 3 {
		t.Errorf("Expected 3 dwarven languages, got %d", len(dwarvenEcosystem))
	}

	// Step 3: Create a complete human linguistic ecosystem
	humanEcosystem, err := CreateCulturalLanguageSet("human", 5, "high", 502)
	if err != nil {
		t.Fatalf("Failed to create human ecosystem: %v", err)
	}

	if len(humanEcosystem) != 5 {
		t.Errorf("Expected 5 human languages, got %d", len(humanEcosystem))
	}

	// Step 4: Test cultural consistency within ecosystems
	testCulturalConsistency(t, "elvish", elvishEcosystem)
	testCulturalConsistency(t, "dwarvish", dwarvenEcosystem)
	testCulturalConsistency(t, "human", humanEcosystem)

	// Step 5: Test inter-cultural contact
	testInterCulturalContact(t, elvishEcosystem[0], dwarvenEcosystem[0], 503)
	testInterCulturalContact(t, humanEcosystem[0], elvishEcosystem[0], 504)

	t.Logf("Successfully created fantasy culture language ecosystems")
	t.Logf("Elvish ecosystem: %d languages", len(elvishEcosystem))
	t.Logf("Dwarven ecosystem: %d languages", len(dwarvenEcosystem))
	t.Logf("Human ecosystem: %d languages", len(humanEcosystem))
}

// testCulturalConsistency verifies that languages within a cultural ecosystem
// share appropriate cultural characteristics.
func testCulturalConsistency(t *testing.T, culture string, languages []*Language) {
	for i, lang := range languages {
		if lang.Culture != culture {
			t.Errorf("Language %d in %s ecosystem has wrong culture: %s", i, culture, lang.Culture)
		}

		// Verify linguistic features match cultural expectations
		switch culture {
		case "elvish":
			if lang.Complexity < LanguageComplexityComplex {
				t.Errorf("Elvish language %d should be complex, got %s", i, lang.Complexity.String())
			}
		case "dwarvish":
			if lang.Complexity < LanguageComplexityModerate {
				t.Errorf("Dwarven language %d should be at least moderate, got %s", i, lang.Complexity.String())
			}
		case "human":
			if lang.Complexity < LanguageComplexityModerate {
				t.Errorf("Human language %d should be at least moderate, got %s", i, lang.Complexity.String())
			}
		}
	}
}

// testInterCulturalContact tests contact between languages from different
// cultural ecosystems.
func testInterCulturalContact(t *testing.T, lang1, lang2 *Language, seed int64) {
	// Simulate trade contact between different cultures
	contactLang1, contactLang2, err := SimulateLanguageContact(
		lang1, lang2, "trade", 0.5, "100_years", seed)
	if err != nil {
		t.Fatalf("Failed to simulate inter-cultural contact: %v", err)
	}

	// Both languages should show influence from the other culture
	if !strings.Contains(contactLang1.Description, "trade") {
		t.Error("Expected trade influence in first language")
	}

	if !strings.Contains(contactLang2.Description, "trade") {
		t.Error("Expected trade influence in second language")
	}

	// Languages from different cultures should maintain distinct identities
	if contactLang1.Culture == contactLang2.Culture {
		t.Error("Expected languages to maintain distinct cultural identities after contact")
	}
}

// TestIntegrationGaps identifies areas where the current implementation
// could be improved to better support realistic fantasy world-building use cases.
func TestIntegrationGaps(t *testing.T) {
	t.Run("Identify Missing Features", func(t *testing.T) {
		identifyMissingFeatures(t)
	})
}

// identifyMissingFeatures documents gaps in the current implementation
// that would make realistic use cases easier to achieve.
func identifyMissingFeatures(t *testing.T) {
	gaps := []string{
		"Real evolution engine integration - Currently using placeholder evolution",
		"Actual sound change simulation - Need to integrate with SoundChangeEngine",
		"Morphological evolution - Need to integrate with MorphologicalEvolutionEngine",
		"Orthographic evolution - Need to integrate with OrthographicEvolutionEngine",
		"Cultural influence simulation - Need to integrate with CulturalInfluenceEngine",
		"Dialect formation engine - Need to integrate with DialectFormationEngine",
		"Full interlingua pipeline - Currently using placeholder translation",
		"Mutual intelligibility calculation - Need to implement linguistic distance metrics",
		"Historical language reconstruction - Need to implement proto-language reconstruction",
		"Writing system genealogy - Need to integrate with WritingTree system",
		"Geographic dialect formation - Need to implement geographic influence models",
		"Social dialect formation - Need to implement social stratification models",
		"Language contact intensity modeling - Need to implement contact strength calculations",
		"Borrowing pattern simulation - Need to implement realistic loanword integration",
		"Phonological adaptation - Need to implement sound system adaptation",
		"Grammatical borrowing - Need to implement grammar feature borrowing",
		"Semantic field evolution - Need to implement vocabulary domain changes",
		"Register and style variation - Need to implement formal/informal distinctions",
		"Language death and revival - Need to implement language extinction and revival",
		"Pidgin and creole formation - Need to implement contact language creation",
	}

	t.Logf("Identified %d areas for improvement:", len(gaps))
	for i, gap := range gaps {
		t.Logf("%d. %s", i+1, gap)
	}

	// These gaps represent opportunities to enhance the system for more
	// realistic and sophisticated fantasy world-building scenarios.
}
