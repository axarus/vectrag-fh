package domain

import (
	"testing"
	"time"
)

func TestValidateModel_Success(t *testing.T) {
	model := Model{
		ID:          "test-model-1",
		Name:        "Test Model",
		Slug:        "test-model",
		Description: "A test model",
		Fields: []Field{
			{
				ID:        "field1",
				Name:      "Test Field",
				Type:      FieldString,
				Required:  true,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Status:        StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err != nil {
		t.Errorf("ValidateModel() error = %v", err)
	}
}

func TestValidateModel_EmptyID(t *testing.T) {
	model := Model{
		ID:          "",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields:      []Field{},
		Status:      StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for empty ID")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "ID: cannot be empty") {
		t.Errorf("Expected ID validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_EmptyName(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "",
		Slug:        "test-model",
		Fields:      []Field{},
		Status:      StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "Name: cannot be empty") {
		t.Errorf("Expected name validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_InvalidSlug(t *testing.T) {
	testCases := []struct {
		name string
		slug string
	}{
		{"empty slug", ""},
		{"slug with uppercase", "Test-Model"},
		{"slug with underscore", "test_model"},
		{"slug starting with hyphen", "-test"},
		{"slug ending with hyphen", "test-"},
		{"slug with special chars", "test@model"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := Model{
				ID:          "test-model",
				Name:        "Test Model",
				Slug:        tc.slug,
				Fields:      []Field{},
				Status:      StatusDraft,
				SchemaVersion: 1,
			}

			err := ValidateModel(model)
			if err == nil {
				t.Errorf("Expected validation error for slug: %s", tc.slug)
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Error("Expected ValidationError")
			}

			if !contains(validationErr.Message, "Slug:") {
				t.Errorf("Expected slug validation error, got: %s", validationErr.Message)
			}
		})
	}
}

func TestValidateModel_InvalidStatus(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields:      []Field{},
		Status:      Status("invalid"),
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for invalid status")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "Status:") {
		t.Errorf("Expected status validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_NegativeSchemaVersion(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields:      []Field{},
		Status:      StatusDraft,
		SchemaVersion: -1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for negative schema version")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "SchemaVersion: cannot be negative") {
		t.Errorf("Expected schema version validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_NoFields(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields:      []Field{},
		Status:      StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for no fields")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "Fields: model must have at least one field") {
		t.Errorf("Expected fields validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_DuplicateFieldIDs(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields: []Field{
			{
				ID:        "duplicate",
				Name:      "Field 1",
				Type:      FieldString,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "duplicate",
				Name:      "Field 2",
				Type:      FieldString,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Status:        StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for duplicate field IDs")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "duplicate field ID") {
		t.Errorf("Expected duplicate field ID validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_DuplicateFieldNames(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields: []Field{
			{
				ID:        "field1",
				Name:      "Duplicate Name",
				Type:      FieldString,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "field2",
				Name:      "Duplicate Name",
				Type:      FieldString,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Status:        StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for duplicate field names")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "duplicate field name") {
		t.Errorf("Expected duplicate field name validation error, got: %s", validationErr.Message)
	}
}

func TestValidateModel_InvalidField(t *testing.T) {
	model := Model{
		ID:          "test-model",
		Name:        "Test Model",
		Slug:        "test-model",
		Fields: []Field{
			{
				ID:        "",
				Name:      "",
				Type:      FieldType("invalid"),
				Status:    Status("invalid"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Status:        StatusDraft,
		SchemaVersion: 1,
	}

	err := ValidateModel(model)
	if err == nil {
		t.Error("Expected validation error for invalid field")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "Fields[0]:") {
		t.Errorf("Expected field validation error, got: %s", validationErr.Message)
	}
}

func TestValidateID(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid ID", "test-model-1", true},
		{"empty ID", "", false},
		{"whitespace only", "   ", false},
		{"valid with underscore", "test_model", true},
		{"valid with hyphen", "test-model", true},
		{"valid alphanumeric", "test123", true},
		{"invalid with space", "test model", false},
		{"invalid with special chars", "test@model", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateID(tc.id)
			if tc.expected && err != nil {
				t.Errorf("Expected no error for ID '%s', got: %v", tc.id, err)
			}
			if !tc.expected && err == nil {
				t.Errorf("Expected error for ID '%s'", tc.id)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	testCases := []struct {
		name     string
		slug     string
		expected bool
	}{
		{"valid slug", "test-model", true},
		{"empty slug", "", false},
		{"whitespace only", "   ", false},
		{"valid lowercase", "testmodel", true},
		{"valid with numbers", "test-123", true},
		{"invalid uppercase", "Test-Model", false},
		{"invalid underscore", "test_model", false},
		{"invalid starting hyphen", "-test", false},
		{"invalid ending hyphen", "test-", false},
		{"invalid special chars", "test@model", false},
		{"invalid space", "test model", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSlug(tc.slug)
			if tc.expected && err != nil {
				t.Errorf("Expected no error for slug '%s', got: %v", tc.slug, err)
			}
			if !tc.expected && err == nil {
				t.Errorf("Expected error for slug '%s'", tc.slug)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
