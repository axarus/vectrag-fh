package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type Model struct {
	ID          string
	Name        string
	Slug        string
	Description string
	Fields      []Field
	// Relations   []Relation
	Status        Status
	SchemaVersion int // TODO: This could be a string, what could be the best way to handle this?
}

func ValidateModel(m Model) error {
	var errors []string

	if err := validateID(m.ID); err != nil {
		errors = append(errors, fmt.Sprintf("ID: %v", err))
	}

	if strings.TrimSpace(m.Name) == "" {
		errors = append(errors, "Name: cannot be empty")
	}

	if err := validateSlug(m.Slug); err != nil {
		errors = append(errors, fmt.Sprintf("Slug: %v", err))
	}

	if err := ValidateStatus(m.Status); err != nil {
		errors = append(errors, fmt.Sprintf("Status: %v", err))
	}

	if m.SchemaVersion < 0 {
		errors = append(errors, "SchemaVersion: cannot be negative")
	}

	if len(m.Fields) == 0 {
		errors = append(errors, "Fields: model must have at least one field")
	}

	fieldIDs := make(map[string]bool)
	fieldNames := make(map[string]bool)
	for i, field := range m.Fields {
		if err := ValidateField(field); err != nil {
			errors = append(errors, fmt.Sprintf("Fields[%d]: %v", i, err))
		}

		if fieldIDs[field.ID] {
			errors = append(errors, fmt.Sprintf("Fields[%d]: duplicate field ID '%s'", i, field.ID))
		}
		fieldIDs[field.ID] = true

		if fieldNames[field.Name] {
			errors = append(errors, fmt.Sprintf("Fields[%d]: duplicate field name '%s'", i, field.Name))
		}
		fieldNames[field.Name] = true
	}

	if len(errors) > 0 {
		return &ValidationError{
			Field:   "Model",
			Message: strings.Join(errors, "; "),
		}
	}

	return nil
}

func validateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("cannot be empty")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, id)
	if !matched {
		return fmt.Errorf("must contain only alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

func validateSlug(slug string) error {
	if strings.TrimSpace(slug) == "" {
		return fmt.Errorf("cannot be empty")
	}

	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	if !matched {
		return fmt.Errorf("must contain only lowercase letters, numbers, and hyphens")
	}

	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return fmt.Errorf("cannot start or end with a hyphen")
	}

	return nil
}
