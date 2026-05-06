package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Field struct {
	ID          string
	Name        string
	Type        FieldType
	Description string
	Unique      bool
	Required    bool
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ValidateField(f Field) error {
	var errors []string

	if strings.TrimSpace(f.ID) == "" {
		errors = append(errors, "ID: cannot be empty")
	} else {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, f.ID)
		if !matched {
			errors = append(errors, "ID: must contain only alphanumeric characters, hyphens, and underscores")
		}
	}

	if strings.TrimSpace(f.Name) == "" {
		errors = append(errors, "Name: cannot be empty")
	}

	if !IsValidType(string(f.Type)) {
		errors = append(errors, fmt.Sprintf("Type: '%s' is not a valid field type", f.Type))
	}

	if err := ValidateStatus(f.Status); err != nil {
		errors = append(errors, fmt.Sprintf("Status: %v", err))
	}

	if len(errors) > 0 {
		return &ValidationError{
			Field:   "Field",
			Message: strings.Join(errors, "; "),
		}
	}

	return nil
}
