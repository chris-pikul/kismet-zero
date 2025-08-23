package interlingua

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error with a code and message.
type ValidationError struct {
	Code    string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error [%s]: %s", e.Code, e.Message)
}

// Validate checks if a Document is well-formed and consistent.
func Validate(doc Document) error {
	// Check required fields
	if doc.SchemaVersion == "" {
		return ValidationError{
			Code:    "missing.schema_version",
			Message: "document must have a schema version",
		}
	}

	if doc.SchemaVersion != SchemaVersion {
		return ValidationError{
			Code:    "unsupported.schema_version",
			Message: fmt.Sprintf("unsupported schema version: %s (expected %s)", doc.SchemaVersion, SchemaVersion),
		}
	}

	// Build entity ID set for validation
	entityIDs := make(map[string]bool)
	for _, entity := range doc.Entities {
		if entity.ID == "" {
			return ValidationError{
				Code:    "missing.entity_id",
				Message: "all entities must have an ID",
			}
		}
		if entity.Concept == "" {
			return ValidationError{
				Code:    "missing.entity_concept",
				Message: fmt.Sprintf("entity %s must have a concept", entity.ID),
			}
		}
		entityIDs[entity.ID] = true
	}

	// Validate events
	for _, event := range doc.Events {
		if err := validateEvent(event, entityIDs); err != nil {
			return err
		}
	}

	return nil
}

// validateEvent validates a single event for consistency.
func validateEvent(event Event, entityIDs map[string]bool) error {
	if event.ID == "" {
		return ValidationError{
			Code:    "missing.event_id",
			Message: "all events must have an ID",
		}
	}

	if event.Predicate == "" {
		return ValidationError{
			Code:    "missing.event_predicate",
			Message: fmt.Sprintf("event %s must have a predicate", event.ID),
		}
	}

	// Validate role references
	if err := validateRoleReferences(event, entityIDs, "roles"); err != nil {
		return err
	}

	// Validate adjunct references
	if err := validateRoleReferences(event, entityIDs, "adjuncts"); err != nil {
		return err
	}

	// Validate valency consistency
	if err := validateValency(event); err != nil {
		return err
	}

	return nil
}

// validateRoleReferences ensures all role/adjunct references point to valid entities.
func validateRoleReferences(event Event, entityIDs map[string]bool, fieldType string) error {
	var refs map[Role]string
	switch fieldType {
	case "roles":
		refs = event.Roles
	case "adjuncts":
		refs = event.Adjuncts
	default:
		return ValidationError{
			Code:    "internal.error",
			Message: fmt.Sprintf("unknown field type: %s", fieldType),
		}
	}

	for role, entityID := range refs {
		if entityID == "" {
			return ValidationError{
				Code:    "missing.role_entity",
				Message: fmt.Sprintf("event %s has empty entity ID for role %s", event.ID, role),
			}
		}

		if !entityIDs[entityID] {
			return ValidationError{
				Code:    "invalid.role_reference",
				Message: fmt.Sprintf("event %s references non-existent entity %s for role %s", event.ID, entityID, role),
			}
		}
	}

	return nil
}

// validateValency checks if the event's valency is consistent with its roles.
func validateValency(event Event) error {
	if event.Valency <= 0 {
		return nil // Valency is optional, skip validation if not set
	}

	roleCount := len(event.Roles)
	if roleCount != event.Valency {
		return ValidationError{
			Code:    "inconsistent.valency",
			Message: fmt.Sprintf("event %s has valency %d but %d roles", event.ID, event.Valency, roleCount),
		}
	}

	return nil
}

// ValidateRoleCompatibility checks if a role is appropriate for a given predicate type.
func ValidateRoleCompatibility(predicate ConceptID, role Role) error {
	// Core semantic roles that are generally valid
	coreRoles := map[Role]bool{
		RoleAgent:       true,
		RolePatient:     true,
		RoleExperiencer: true,
		RoleStimulus:    true,
		RoleInstrument:  true,
		RoleLocation:    true,
		RoleTime:        true,
		RoleManner:      true,
		RoleCause:       true,
	}

	if coreRoles[role] {
		return nil
	}

	// Specialized roles that require specific predicates
	switch role {
	case RoleRecipient:
		// Recipient typically requires ditransitive predicates like "give", "send", "tell"
		if !isDitransitivePredicate(predicate) {
			return ValidationError{
				Code:    "incompatible.role",
				Message: fmt.Sprintf("role 'recipient' is not appropriate for predicate %s", predicate),
			}
		}
	case RoleSource, RoleGoal:
		// Source/Goal typically require motion or transfer predicates
		if !isMotionOrTransferPredicate(predicate) {
			return ValidationError{
				Code:    "incompatible.role",
				Message: fmt.Sprintf("role '%s' is not appropriate for predicate %s", role, predicate),
			}
		}
	}

	return nil
}

// isDitransitivePredicate checks if a predicate typically takes three arguments.
func isDitransitivePredicate(predicate ConceptID) bool {
	ditransitivePredicates := []string{
		"give", "send", "tell", "show", "teach", "bring", "take", "hand", "pass",
		"lend", "rent", "sell", "buy", "owe", "promise", "offer", "award",
	}

	predStr := string(predicate)
	for _, dp := range ditransitivePredicates {
		if strings.HasPrefix(predStr, dp) {
			return true
		}
	}
	return false
}

// isMotionOrTransferPredicate checks if a predicate involves motion or transfer.
func isMotionOrTransferPredicate(predicate ConceptID) bool {
	motionPredicates := []string{
		"move", "go", "come", "walk", "run", "fly", "swim", "jump", "climb",
		"carry", "bring", "take", "send", "receive", "give", "get", "put",
	}

	predStr := string(predicate)
	for _, mp := range motionPredicates {
		if strings.HasPrefix(predStr, mp) {
			return true
		}
	}
	return false
}
