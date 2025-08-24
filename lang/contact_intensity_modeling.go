package lang

import (
	"fmt"
	"math"
	"time"
)

// ContactIntensityModel represents the intensity modeling for language contact.
type ContactIntensityModel struct {
	ID          string      `json:"id"`
	ContactType ContactType `json:"contactType"`

	// Intensity factors
	BaseIntensity     float32 `json:"baseIntensity"`     // Base intensity for this contact type
	DurationScaling   float32 `json:"durationScaling"`   // How duration affects intensity
	FrequencyScaling  float32 `json:"frequencyScaling"`  // How contact frequency affects intensity
	GeographicScaling float32 `json:"geographicScaling"` // How geographic proximity affects intensity

	// Thresholds
	MinIntensityThreshold float32 `json:"minIntensityThreshold"` // Minimum intensity for any change
	MaxIntensityThreshold float32 `json:"maxIntensityThreshold"` // Maximum intensity cap
	AdaptationThreshold   float32 `json:"adaptationThreshold"`   // Intensity needed for structural adaptation

	// Cumulative effects
	DecayRate        float32 `json:"decayRate"`        // How quickly intensity decays over time
	AccumulationRate float32 `json:"accumulationRate"` // How quickly intensity accumulates with repeated contact
}

// ContactHistory tracks the history of contact between languages.
type ContactHistory struct {
	LanguageID    string        `json:"languageId"`
	SourceID      string        `json:"sourceId"`
	ContactType   ContactType   `json:"contactType"`
	FirstContact  time.Time     `json:"firstContact"`
	LastContact   time.Time     `json:"lastContact"`
	TotalDuration time.Duration `json:"totalDuration"`
	ContactCount  int           `json:"contactCount"`

	// Intensity tracking
	PeakIntensity    float32 `json:"peakIntensity"`
	CurrentIntensity float32 `json:"currentIntensity"`
	AverageIntensity float32 `json:"averageIntensity"`

	// Cumulative effects
	TotalComplexityChange float32 `json:"totalComplexityChange"`
	AdaptationCount       int     `json:"adaptationCount"`

	// Contact sessions
	Sessions []ContactSession `json:"sessions,omitempty"`
}

// ContactSession represents a single contact session.
type ContactSession struct {
	ID               string        `json:"id"`
	Timestamp        time.Time     `json:"timestamp"`
	Duration         time.Duration `json:"duration"`
	Intensity        float32       `json:"intensity"`
	Changes          []string      `json:"changes,omitempty"`
	ComplexityChange float32       `json:"complexityChange"`
}

// ContactIntensityCalculator calculates the effective intensity for language contact.
type ContactIntensityCalculator struct {
	models map[ContactType]*ContactIntensityModel
}

// NewContactIntensityCalculator creates a new contact intensity calculator.
func NewContactIntensityCalculator() *ContactIntensityCalculator {
	calculator := &ContactIntensityCalculator{
		models: make(map[ContactType]*ContactIntensityModel),
	}

	// Initialize default models for each contact type
	calculator.initializeDefaultModels()

	return calculator
}

// initializeDefaultModels sets up default intensity models for each contact type.
func (c *ContactIntensityCalculator) initializeDefaultModels() {
	// Trade contact - moderate intensity, builds up over time
	c.models[ContactTypeTrade] = &ContactIntensityModel{
		ID:                    "trade_intensity",
		ContactType:           ContactTypeTrade,
		BaseIntensity:         0.4,
		DurationScaling:       0.3,
		FrequencyScaling:      0.4,
		GeographicScaling:     0.2,
		MinIntensityThreshold: 0.1,
		MaxIntensityThreshold: 0.8,
		AdaptationThreshold:   0.6,
		DecayRate:             0.1,
		AccumulationRate:      0.2,
	}

	// Conquest contact - high intensity, immediate impact
	c.models[ContactTypeConquest] = &ContactIntensityModel{
		ID:                    "conquest_intensity",
		ContactType:           ContactTypeConquest,
		BaseIntensity:         0.8,
		DurationScaling:       0.1,
		FrequencyScaling:      0.1,
		GeographicScaling:     0.0,
		MinIntensityThreshold: 0.5,
		MaxIntensityThreshold: 1.0,
		AdaptationThreshold:   0.7,
		DecayRate:             0.05,
		AccumulationRate:      0.1,
	}

	// Migration contact - moderate-high intensity, gradual build-up
	c.models[ContactTypeMigration] = &ContactIntensityModel{
		ID:                    "migration_intensity",
		ContactType:           ContactTypeMigration,
		BaseIntensity:         0.6,
		DurationScaling:       0.5,
		FrequencyScaling:      0.3,
		GeographicScaling:     0.4,
		MinIntensityThreshold: 0.2,
		MaxIntensityThreshold: 0.9,
		AdaptationThreshold:   0.6,
		DecayRate:             0.15,
		AccumulationRate:      0.25,
	}

	// Cultural contact - low-moderate intensity, steady accumulation
	c.models[ContactTypeCultural] = &ContactIntensityModel{
		ID:                    "cultural_intensity",
		ContactType:           ContactTypeCultural,
		BaseIntensity:         0.3,
		DurationScaling:       0.6,
		FrequencyScaling:      0.5,
		GeographicScaling:     0.3,
		MinIntensityThreshold: 0.05,
		MaxIntensityThreshold: 0.7,
		AdaptationThreshold:   0.5,
		DecayRate:             0.2,
		AccumulationRate:      0.3,
	}

	// Religious contact - moderate intensity, strong persistence
	c.models[ContactTypeReligious] = &ContactIntensityModel{
		ID:                    "religious_intensity",
		ContactType:           ContactTypeReligious,
		BaseIntensity:         0.5,
		DurationScaling:       0.4,
		FrequencyScaling:      0.6,
		GeographicScaling:     0.2,
		MinIntensityThreshold: 0.1,
		MaxIntensityThreshold: 0.8,
		AdaptationThreshold:   0.6,
		DecayRate:             0.08,
		AccumulationRate:      0.35,
	}

	// Educational contact - low intensity, steady accumulation
	c.models[ContactTypeEducational] = &ContactIntensityModel{
		ID:                    "educational_intensity",
		ContactType:           ContactTypeEducational,
		BaseIntensity:         0.2,
		DurationScaling:       0.7,
		FrequencyScaling:      0.6,
		GeographicScaling:     0.1,
		MinIntensityThreshold: 0.05,
		MaxIntensityThreshold: 0.6,
		AdaptationThreshold:   0.5,
		DecayRate:             0.25,
		AccumulationRate:      0.4,
	}
}

// CalculateEffectiveIntensity calculates the effective intensity for a contact event.
func (c *ContactIntensityCalculator) CalculateEffectiveIntensity(
	contactType ContactType,
	baseIntensity float32,
	duration time.Duration,
	history *ContactHistory,
	geographicProximity float32,
) float32 {
	model, exists := c.models[contactType]
	if !exists {
		return baseIntensity
	}

	// Start with base intensity
	effectiveIntensity := model.BaseIntensity

	// Apply duration scaling
	durationHours := float32(duration.Hours())
	durationFactor := float32(math.Min(float64(durationHours/24.0), 30.0)) // Cap at 30 days
	effectiveIntensity += model.DurationScaling * durationFactor * 0.1

	// Apply frequency scaling from history
	if history != nil {
		frequencyFactor := float32(math.Min(float64(history.ContactCount), 10.0)) // Cap at 10 contacts
		effectiveIntensity += model.FrequencyScaling * frequencyFactor * 0.05

		// Apply accumulation from previous contacts
		timeSinceLastContact := time.Since(history.LastContact)
		decayFactor := math.Exp(-float64(model.DecayRate) * timeSinceLastContact.Hours() / 24.0)
		accumulatedIntensity := history.PeakIntensity * float32(decayFactor)
		effectiveIntensity += model.AccumulationRate * accumulatedIntensity * 0.3
	}

	// Apply geographic proximity scaling
	effectiveIntensity += model.GeographicScaling * geographicProximity * 0.2

	// Apply base intensity modifier
	effectiveIntensity += baseIntensity * 0.4

	// Clamp to thresholds
	effectiveIntensity = float32(math.Max(float64(effectiveIntensity), float64(model.MinIntensityThreshold)))
	effectiveIntensity = float32(math.Min(float64(effectiveIntensity), float64(model.MaxIntensityThreshold)))

	return effectiveIntensity
}

// ShouldApplyAdaptation determines if structural adaptation should occur based on intensity.
func (c *ContactIntensityCalculator) ShouldApplyAdaptation(
	contactType ContactType,
	effectiveIntensity float32,
) bool {
	model, exists := c.models[contactType]
	if !exists {
		return false
	}

	return effectiveIntensity >= model.AdaptationThreshold
}

// UpdateContactHistory updates the contact history with a new contact event.
func (c *ContactIntensityCalculator) UpdateContactHistory(
	history *ContactHistory,
	contactType ContactType,
	duration time.Duration,
	intensity float32,
	changes []string,
	complexityChange float32,
) *ContactHistory {
	if history == nil {
		history = &ContactHistory{
			LanguageID:            "", // Will be set by caller
			SourceID:              "", // Will be set by caller
			ContactType:           contactType,
			FirstContact:          time.Now(),
			LastContact:           time.Now(),
			TotalDuration:         duration,
			ContactCount:          1,
			PeakIntensity:         intensity,
			CurrentIntensity:      intensity,
			AverageIntensity:      intensity,
			TotalComplexityChange: complexityChange,
			AdaptationCount:       0,
			Sessions:              make([]ContactSession, 0),
		}
	} else {
		// Update existing history
		history.LastContact = time.Now()
		history.TotalDuration += duration
		history.ContactCount++

		// Update intensity tracking
		if intensity > history.PeakIntensity {
			history.PeakIntensity = intensity
		}
		history.CurrentIntensity = intensity

		// Update average intensity
		totalIntensity := history.AverageIntensity*float32(history.ContactCount-1) + intensity
		history.AverageIntensity = totalIntensity / float32(history.ContactCount)

		// Update cumulative effects
		history.TotalComplexityChange += complexityChange

		// Count adaptations (changes that affect structure)
		if len(changes) > 0 {
			// Simple heuristic: if we have structural changes, count as adaptation
			history.AdaptationCount++
		}
	}

	// Add new session
	session := ContactSession{
		ID:               fmt.Sprintf("session_%d", len(history.Sessions)+1),
		Timestamp:        time.Now(),
		Duration:         duration,
		Intensity:        intensity,
		Changes:          changes,
		ComplexityChange: complexityChange,
	}
	history.Sessions = append(history.Sessions, session)

	return history
}

// GetContactHistory retrieves or creates contact history for a language pair.
func (c *ContactIntensityCalculator) GetContactHistory(
	targetLang *Language,
	sourceLang *Language,
	contactType ContactType,
) *ContactHistory {
	// For now, we'll create a simple history tracking
	// In a full implementation, this would be stored persistently

	// In reality, this would query a database or persistent storage using:
	// historyKey := fmt.Sprintf("%s_%s_%s",
	// 	targetLang.ID.String(),
	// 	sourceLang.ID.String(),
	// 	contactType.String())

	return &ContactHistory{
		LanguageID:            targetLang.ID.String(),
		SourceID:              sourceLang.ID.String(),
		ContactType:           contactType,
		FirstContact:          time.Now(),
		LastContact:           time.Now(),
		TotalDuration:         0,
		ContactCount:          0,
		PeakIntensity:         0.0,
		CurrentIntensity:      0.0,
		AverageIntensity:      0.0,
		TotalComplexityChange: 0.0,
		AdaptationCount:       0,
		Sessions:              make([]ContactSession, 0),
	}
}

// CalculateGeographicProximity calculates geographic proximity between languages.
// This is a simplified implementation - in reality, this would use actual geographic data.
func (c *ContactIntensityCalculator) CalculateGeographicProximity(
	lang1, lang2 *Language,
) float32 {
	// Simplified geographic proximity calculation
	// In reality, this would use actual geographic coordinates and calculate distances

	// For now, we'll use a random but consistent proximity based on language IDs
	// This simulates different geographic relationships

	hash1 := hashString(lang1.ID.String())
	hash2 := hashString(lang2.ID.String())

	// Create a pseudo-random but consistent proximity value
	proximity := float32((hash1+hash2)%100) / 100.0

	// Adjust proximity based on language family relationships (simplified)
	// In reality, this would check actual family relationships
	lang1Str := lang1.ID.String()
	lang2Str := lang2.ID.String()

	if len(lang1Str) >= 3 && len(lang2Str) >= 3 && lang1Str[:3] == lang2Str[:3] {
		proximity += 0.3 // Similar IDs = closer proximity (simulating family relationship)
	}

	// Clamp to valid range
	if proximity > 1.0 {
		proximity = 1.0
	}

	return proximity
}

// hashString creates a simple hash for consistent pseudo-random values.
func hashString(s string) int {
	hash := 0
	for _, char := range s {
		hash = (hash*31 + int(char)) % 1000000
	}
	return hash
}
