package domain

import (
	"testing"
	"time"
)

func TestValidateField_Success(t *testing.T) {
	field := Field{
		ID:          "test-field-1",
		Name:        "Test Field",
		Type:        FieldString,
		Description: "A test field",
		Unique:      false,
		Required:    true,
		Status:      StatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := ValidateField(field)
	if err != nil {
		t.Errorf("ValidateField() error = %v", err)
	}
}

func TestValidateField_EmptyID(t *testing.T) {
	field := Field{
		ID:        "",
		Name:      "Test Field",
		Type:      FieldString,
		Status:    StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
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

func TestValidateField_WhitespaceID(t *testing.T) {
	field := Field{
		ID:        "   ",
		Name:      "Test Field",
		Type:      FieldString,
		Status:    StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
	if err == nil {
		t.Error("Expected validation error for whitespace-only ID")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "ID: cannot be empty") {
		t.Errorf("Expected ID validation error, got: %s", validationErr.Message)
	}
}

func TestValidateField_InvalidIDFormat(t *testing.T) {
	testCases := []struct {
		name string
		id   string
	}{
		{"ID with space", "test field"},
		{"ID with special chars", "test@field"},
		{"ID with dot", "test.field"},
		{"ID with slash", "test/field"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := Field{
				ID:        tc.id,
				Name:      "Test Field",
				Type:      FieldString,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := ValidateField(field)
			if err == nil {
				t.Errorf("Expected validation error for ID: %s", tc.id)
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Error("Expected ValidationError")
			}

			if !contains(validationErr.Message, "ID: must contain only alphanumeric characters, hyphens, and underscores") {
				t.Errorf("Expected ID format validation error, got: %s", validationErr.Message)
			}
		})
	}
}

func TestValidateField_ValidIDFormats(t *testing.T) {
	testCases := []struct {
		name string
		id   string
	}{
		{"valid alphanumeric", "testfield"},
		{"valid with hyphen", "test-field"},
		{"valid with underscore", "test_field"},
		{"valid with numbers", "test123"},
		{"valid mixed", "test_field-123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := Field{
				ID:        tc.id,
				Name:      "Test Field",
				Type:      FieldString,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := ValidateField(field)
			if err != nil {
				t.Errorf("Expected no error for valid ID '%s', got: %v", tc.id, err)
			}
		})
	}
}

func TestValidateField_EmptyName(t *testing.T) {
	field := Field{
		ID:        "test-field",
		Name:      "",
		Type:      FieldString,
		Status:    StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
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

func TestValidateField_WhitespaceName(t *testing.T) {
	field := Field{
		ID:        "test-field",
		Name:      "   ",
		Type:      FieldString,
		Status:    StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
	if err == nil {
		t.Error("Expected validation error for whitespace-only name")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "Name: cannot be empty") {
		t.Errorf("Expected name validation error, got: %s", validationErr.Message)
	}
}

func TestValidateField_InvalidType(t *testing.T) {
	field := Field{
		ID:        "test-field",
		Name:      "Test Field",
		Type:      FieldType("invalid-type"),
		Status:    StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
	if err == nil {
		t.Error("Expected validation error for invalid type")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	if !contains(validationErr.Message, "Type: 'invalid-type' is not a valid field type") {
		t.Errorf("Expected type validation error, got: %s", validationErr.Message)
	}
}

func TestValidateField_InvalidStatus(t *testing.T) {
	field := Field{
		ID:        "test-field",
		Name:      "Test Field",
		Type:      FieldString,
		Status:    Status("invalid"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
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

func TestValidateField_AllValidTypes(t *testing.T) {
	validTypes := []FieldType{
		FieldString,
		FieldText,
		FieldNumber,
		FieldBoolean,
		FieldDate,
		FieldDateTime,
		FieldRelation,
	}

	for i, fieldType := range validTypes {
		t.Run(string(fieldType), func(t *testing.T) {
			field := Field{
				ID:        "test-field",
				Name:      "Test Field",
				Type:      fieldType,
				Status:    StatusDraft,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := ValidateField(field)
			if err != nil {
				t.Errorf("Expected no error for valid type '%s', got: %v", fieldType, err)
			}
		})
		_ = i // Avoid unused variable warning
	}
}

func TestValidateField_AllValidStatuses(t *testing.T) {
	validStatuses := []Status{
		StatusDraft,
		StatusPublish,
		StatusDelete,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			field := Field{
				ID:        "test-field",
				Name:      "Test Field",
				Type:      FieldString,
				Status:    status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := ValidateField(field)
			if err != nil {
				t.Errorf("Expected no error for valid status '%s', got: %v", status, err)
			}
		})
	}
}

func TestValidateField_MultipleErrors(t *testing.T) {
	field := Field{
		ID:        "test field@",        // Invalid ID format
		Name:      "",                   // Empty name
		Type:      FieldType("invalid"), // Invalid type
		Status:    Status("invalid"),    // Invalid status
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ValidateField(field)
	if err == nil {
		t.Error("Expected validation error for multiple invalid fields")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Error("Expected ValidationError")
	}

	// Check that all expected errors are present
	message := validationErr.Message
	expectedErrors := []string{
		"ID: must contain only alphanumeric characters, hyphens, and underscores",
		"Name: cannot be empty",
		"Type: 'invalid' is not a valid field type",
		"Status:",
	}

	for _, expectedError := range expectedErrors {
		if !contains(message, expectedError) {
			t.Errorf("Expected error message to contain '%s', got: %s", expectedError, message)
		}
	}
}
