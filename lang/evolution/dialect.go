package evolution

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// DialectFormationEngine manages the creation and evolution of dialects.
type DialectFormationEngine struct {
	config EvolutionConfig
	rng    *rand.Rand
}

// NewDialectFormationEngine creates a new dialect formation engine.
func NewDialectFormationEngine(config EvolutionConfig) *DialectFormationEngine {
	rng := rand.New(rand.NewPCG(uint64(config.Seed), 0))

	return &DialectFormationEngine{
		config: config,
		rng:    rng,
	}
}

// CreateGeographicDialect creates a new geographic dialect based on regional factors.
func (dfe *DialectFormationEngine) CreateGeographicDialect(
	parentLang *lang.Language,
	dialectID string,
	dialectName string,
	region GeographicRegion,
	era string,
) (*Dialect, DialectFormationEvent, error) {

	if parentLang == nil {
		return nil, DialectFormationEvent{}, fmt.Errorf("parent language cannot be nil")
	}

	// Calculate formation factors based on geographic isolation
	geographicFactors := dfe.calculateGeographicFactors(region)

	// Create dialect features based on geographic factors
	features := dfe.generateDialectFeatures(parentLang, geographicFactors, region)

	// Create the dialect
	dialect := &Dialect{
		ID:            dialectID,
		Name:          dialectName,
		Type:          DialectTypeGeographic,
		ParentLang:    parentLang.ID.String(),
		Region:        region,
		Features:      features,
		FormationDate: time.Now(),
		Era:           era,
		Status:        "active",
		Seed:          dfe.rng.Int64(),
	}

	// Create formation event
	formationEvent := DialectFormationEvent{
		ID:                fmt.Sprintf("dialect_formation_%s_%d", dialectID, time.Now().Unix()),
		Timestamp:         time.Now(),
		Era:               era,
		ParentLang:        parentLang.ID.String(),
		NewDialect:        dialectID,
		GeographicFactors: geographicFactors,
		Intensity:         dfe.calculateFormationIntensity(geographicFactors),
		Description:       fmt.Sprintf("Geographic dialect formed in %s region", region.Name),
		Seed:              dfe.rng.Int64(),
	}

	return dialect, formationEvent, nil
}

// CreateSocialDialect creates a new social dialect based on social stratification.
func (dfe *DialectFormationEngine) CreateSocialDialect(
	parentLang *lang.Language,
	dialectID string,
	dialectName string,
	socialClass string,
	urbanization float32,
	era string,
) (*Dialect, DialectFormationEvent, error) {

	if parentLang == nil {
		return nil, DialectFormationEvent{}, fmt.Errorf("parent language cannot be nil")
	}

	// Create a geographic region for the social dialect
	region := GeographicRegion{
		ID:           fmt.Sprintf("social_%s", socialClass),
		Name:         fmt.Sprintf("%s %s", socialClass, "region"),
		Latitude:     0.0, // Placeholder
		Longitude:    0.0, // Placeholder
		Urbanization: urbanization,
		Population:   1000, // Placeholder
	}

	// Calculate social factors
	socialFactors := dfe.calculateSocialFactors(socialClass, urbanization)

	// Create dialect features based on social factors
	features := dfe.generateDialectFeatures(parentLang, socialFactors, region)

	// Determine dialect type based on urbanization
	dialectType := DialectTypeSocial
	if urbanization > 0.7 {
		dialectType = DialectTypeUrban
	}

	// Create the dialect
	dialect := &Dialect{
		ID:            dialectID,
		Name:          dialectName,
		Type:          dialectType,
		ParentLang:    parentLang.ID.String(),
		Region:        region,
		Features:      features,
		FormationDate: time.Now(),
		Era:           era,
		Status:        "active",
		Seed:          dfe.rng.Int64(),
	}

	// Create formation event
	formationEvent := DialectFormationEvent{
		ID:            fmt.Sprintf("dialect_formation_%s_%d", dialectID, time.Now().Unix()),
		Timestamp:     time.Now(),
		Era:           era,
		ParentLang:    parentLang.ID.String(),
		NewDialect:    dialectID,
		SocialFactors: socialFactors,
		Intensity:     dfe.calculateFormationIntensity(socialFactors),
		Description:   fmt.Sprintf("%s dialect formed for %s class", dialectType.String(), socialClass),
		Seed:          dfe.rng.Int64(),
	}

	return dialect, formationEvent, nil
}

// EvolveDialect evolves a dialect over time, potentially creating new sub-dialects.
func (dfe *DialectFormationEngine) EvolveDialect(
	dialect *Dialect,
	era string,
	duration time.Duration,
) (*Dialect, EvolutionEvent, error) {

	if dialect == nil {
		return nil, EvolutionEvent{}, fmt.Errorf("dialect cannot be nil")
	}

	// Create a copy for evolution
	evolvedDialect := *dialect

	// Apply dialect-specific changes
	changes := dfe.applyDialectalChanges(&evolvedDialect, era)

	// Update dialect features
	evolvedDialect.Features = dfe.evolveDialectFeatures(dialect.Features, changes)

	// Create evolution event
	evolutionEvent := EvolutionEvent{
		ID:          fmt.Sprintf("dialect_evolution_%s_%d", dialect.ID, time.Now().Unix()),
		Name:        fmt.Sprintf("Evolution of %s", dialect.Name),
		Description: fmt.Sprintf("Dialect evolved over %v during %s era", duration, era),
		Timestamp:   time.Now(),
		Era:         era,
		Changes:     changes,
		TriggerType: "dialectal_evolution",
		Intensity:   dfe.calculateOverallIntensity(changes),
		Seed:        dfe.rng.Int64(),
	}

	// Add to evolution history
	evolvedDialect.EvolutionHistory = append(evolvedDialect.EvolutionHistory, evolutionEvent)

	return &evolvedDialect, evolutionEvent, nil
}

// calculateGeographicFactors determines what geographic factors influence dialect formation.
func (dfe *DialectFormationEngine) calculateGeographicFactors(region GeographicRegion) []string {
	var factors []string

	// Climate influence
	if region.Climate != "" {
		factors = append(factors, fmt.Sprintf("climate_%s", region.Climate))
	}

	// Terrain influence
	if region.Terrain != "" {
		factors = append(factors, fmt.Sprintf("terrain_%s", region.Terrain))
	}

	// Population density
	if region.Population > 10000 {
		factors = append(factors, "high_population_density")
	} else if region.Population < 1000 {
		factors = append(factors, "low_population_density")
	}

	// Urbanization level
	if region.Urbanization > 0.8 {
		factors = append(factors, "high_urbanization")
	} else if region.Urbanization < 0.2 {
		factors = append(factors, "low_urbanization")
	}

	// Geographic isolation (distance from parent language center)
	if region.Latitude != 0 || region.Longitude != 0 {
		factors = append(factors, "geographic_isolation")
	}

	return factors
}

// calculateSocialFactors determines what social factors influence dialect formation.
func (dfe *DialectFormationEngine) calculateSocialFactors(socialClass string, urbanization float32) []string {
	var factors []string

	// Social class
	factors = append(factors, fmt.Sprintf("social_class_%s", socialClass))

	// Urbanization level
	if urbanization > 0.7 {
		factors = append(factors, "high_urbanization")
	} else if urbanization < 0.3 {
		factors = append(factors, "low_urbanization")
	}

	// Education level (inferred from social class)
	if socialClass == "upper" || socialClass == "educated" {
		factors = append(factors, "high_education")
	} else if socialClass == "working" || socialClass == "rural" {
		factors = append(factors, "practical_education")
	}

	return factors
}

// generateDialectFeatures creates distinctive features for a new dialect.
func (dfe *DialectFormationEngine) generateDialectFeatures(
	parentLang *lang.Language,
	factors []string,
	region GeographicRegion,
) DialectFeatures {

	features := DialectFeatures{
		ID:                   fmt.Sprintf("features_%d", time.Now().UnixNano()),
		PhonologicalFeatures: make([]string, 0),
		LexicalFeatures:      make([]string, 0),
		GrammaticalFeatures:  make([]string, 0),
		PragmaticFeatures:    make([]string, 0),
		IntelligibilityScore: 0.85, // Start with high intelligibility
	}

	// Generate phonological features based on factors
	for _, factor := range factors {
		switch {
		case factor == "climate_tropical":
			features.PhonologicalFeatures = append(features.PhonologicalFeatures, "vowel_harmony_shift")
		case factor == "terrain_mountain":
			features.PhonologicalFeatures = append(features.PhonologicalFeatures, "consonant_cluster_simplification")
		case factor == "high_urbanization":
			features.PhonologicalFeatures = append(features.PhonologicalFeatures, "fast_speech_patterns")
		case factor == "low_urbanization":
			features.PhonologicalFeatures = append(features.PhonologicalFeatures, "conservative_pronunciation")
		}
	}

	// Generate lexical features
	for _, factor := range factors {
		switch {
		case factor == "terrain_coastal":
			features.LexicalFeatures = append(features.LexicalFeatures, "maritime_vocabulary")
		case factor == "terrain_mountain":
			features.LexicalFeatures = append(features.LexicalFeatures, "mountain_terminology")
		case factor == "high_urbanization":
			features.LexicalFeatures = append(features.LexicalFeatures, "urban_slang")
		case factor == "social_class_working":
			features.LexicalFeatures = append(features.LexicalFeatures, "occupational_terms")
		}
	}

	// Generate grammatical features
	for _, factor := range factors {
		switch {
		case factor == "high_population_density":
			features.GrammaticalFeatures = append(features.GrammaticalFeatures, "simplified_grammar")
		case factor == "low_population_density":
			features.GrammaticalFeatures = append(features.GrammaticalFeatures, "conservative_grammar")
		case factor == "social_class_educated":
			features.GrammaticalFeatures = append(features.GrammaticalFeatures, "formal_register")
		}
	}

	// Adjust intelligibility based on factors
	features.IntelligibilityScore = dfe.calculateIntelligibility(factors)

	return features
}

// calculateFormationIntensity calculates the intensity of dialect formation.
func (dfe *DialectFormationEngine) calculateFormationIntensity(factors []string) float32 {
	baseIntensity := 0.3

	// Geographic isolation increases intensity
	for _, factor := range factors {
		switch {
		case factor == "geographic_isolation":
			baseIntensity += 0.2
		case factor == "high_population_density":
			baseIntensity += 0.1
		case factor == "low_population_density":
			baseIntensity += 0.15
		case factor == "high_urbanization":
			baseIntensity += 0.1
		case factor == "low_urbanization":
			baseIntensity += 0.1
		}
	}

	if baseIntensity > 1.0 {
		baseIntensity = 1.0
	}

	return float32(baseIntensity)
}

// calculateIntelligibility calculates mutual intelligibility with parent language.
func (dfe *DialectFormationEngine) calculateIntelligibility(factors []string) float32 {
	baseIntelligibility := 0.9

	// Factors that reduce intelligibility
	for _, factor := range factors {
		switch {
		case factor == "geographic_isolation":
			baseIntelligibility -= 0.1
		case factor == "high_urbanization":
			baseIntelligibility -= 0.05
		case factor == "social_class_working":
			baseIntelligibility -= 0.03
		case factor == "terrain_mountain":
			baseIntelligibility -= 0.05
		}
	}

	if baseIntelligibility < 0.5 {
		baseIntelligibility = 0.5
	}

	return float32(baseIntelligibility)
}

// applyDialectalChanges applies dialect-specific changes over time.
func (dfe *DialectFormationEngine) applyDialectalChanges(dialect *Dialect, era string) []LinguisticChange {
	var changes []LinguisticChange

	// Apply changes based on dialect type and region
	if dfe.rng.Float32() < dfe.config.DialectFormationRate {
		// Phonological drift
		if len(dialect.Features.PhonologicalFeatures) > 0 {
			change := LinguisticChange{
				ID:          fmt.Sprintf("dialectal_phonological_%d", time.Now().UnixNano()),
				Type:        ChangeTypeDialectal,
				Direction:   ChangeDirectionModifying,
				Description: "Phonological drift in dialect",
				Details:     "Gradual sound changes specific to this dialect",
				Timestamp:   time.Now(),
				Era:         era,
				Trigger:     "dialectal_evolution",
				Intensity:   0.2,
			}
			changes = append(changes, change)
		}

		// Lexical innovation
		if dfe.rng.Float32() < 0.3 {
			change := LinguisticChange{
				ID:          fmt.Sprintf("dialectal_lexical_%d", time.Now().UnixNano()),
				Type:        ChangeTypeDialectal,
				Direction:   ChangeDirectionAdditive,
				Description: "New dialect-specific vocabulary",
				Details:     "Local terms and expressions developed",
				Timestamp:   time.Now(),
				Era:         era,
				Trigger:     "dialectal_innovation",
				Intensity:   0.15,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

// evolveDialectFeatures evolves the features of a dialect based on changes.
func (dfe *DialectFormationEngine) evolveDialectFeatures(
	features DialectFeatures,
	changes []LinguisticChange,
) DialectFeatures {

	evolvedFeatures := features

	// Apply changes to features
	for _, change := range changes {
		switch change.Type {
		case ChangeTypeDialectal:
			switch change.Direction {
			case ChangeDirectionAdditive:
				if change.Description == "New dialect-specific vocabulary" {
					evolvedFeatures.LexicalFeatures = append(evolvedFeatures.LexicalFeatures, "new_local_term")
				}
			case ChangeDirectionModifying:
				if change.Description == "Phonological drift in dialect" {
					evolvedFeatures.PhonologicalFeatures = append(evolvedFeatures.PhonologicalFeatures, "phonological_drift")
				}
			}
		}
	}

	// Gradually reduce intelligibility over time
	evolvedFeatures.IntelligibilityScore *= 0.995

	return evolvedFeatures
}

// calculateOverallIntensity calculates the overall intensity of dialect changes.
func (dfe *DialectFormationEngine) calculateOverallIntensity(changes []LinguisticChange) float32 {
	if len(changes) == 0 {
		return 0.0
	}

	totalIntensity := float32(0.0)
	for _, change := range changes {
		totalIntensity += change.Intensity
	}

	return totalIntensity / float32(len(changes))
}
