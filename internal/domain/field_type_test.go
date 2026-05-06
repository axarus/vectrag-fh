package domain

import (
	"reflect"
	"testing"
)

func TestIsValidType(t *testing.T) {
	testCases := []struct {
		name     string
		typeStr  string
		expected bool
	}{
		{"valid string", "string", true},
		{"valid text", "text", true},
		{"valid number", "number", true},
		{"valid boolean", "boolean", true},
		{"valid date", "date", true},
		{"valid datetime", "datetime", true},
		{"valid relation", "relation", true},
		{"invalid type", "invalid", false},
		{"empty type", "", false},
		{"case sensitive", "String", false},
		{"partial match", "strin", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidType(tc.typeStr)
			if result != tc.expected {
				t.Errorf("IsValidType(%s) = %v, expected %v", tc.typeStr, result, tc.expected)
			}
		})
	}
}

func TestListFieldTypes(t *testing.T) {
	types := ListFieldTypes()

	if len(types) == 0 {
		t.Error("ListFieldTypes() returned empty slice")
	}

	expectedTypes := []string{
		"string",
		"text",
		"number",
		"boolean",
		"date",
		"datetime",
		"relation",
	}

	if len(types) != len(expectedTypes) {
		t.Errorf("Expected %d types, got %d", len(expectedTypes), len(types))
	}

	// Create a map for easier comparison
	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	for _, expectedType := range expectedTypes {
		if !typeMap[expectedType] {
			t.Errorf("Expected type '%s' not found in result", expectedType)
		}
	}

	// Verify all returned types are valid
	for _, returnedType := range types {
		if !IsValidType(returnedType) {
			t.Errorf("Returned type '%s' is not valid", returnedType)
		}
	}
}

func TestFieldTypeConstants(t *testing.T) {
	// Test that all field type constants have valid string values
	testCases := []struct {
		name      string
		fieldType FieldType
		expected  string
	}{
		{"FieldString", FieldString, "string"},
		{"FieldText", FieldText, "text"},
		{"FieldNumber", FieldNumber, "number"},
		{"FieldBoolean", FieldBoolean, "boolean"},
		{"FieldDate", FieldDate, "date"},
		{"FieldDateTime", FieldDateTime, "datetime"},
		{"FieldRelation", FieldRelation, "relation"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(tc.fieldType)
			if actual != tc.expected {
				t.Errorf("Expected %s = '%s', got '%s'", tc.name, tc.expected, actual)
			}

			// Verify the constant is registered as valid
			if !IsValidType(actual) {
				t.Errorf("Field type constant %s ('%s') is not registered as valid", tc.name, actual)
			}
		})
	}
}

func TestFieldTypeRegistryCompleteness(t *testing.T) {
	// Use reflection to ensure all constants are in the registry
	fieldTypeType := reflect.TypeOf(FieldType(""))

	var constants []FieldType
	for i := 0; i < fieldTypeType.NumMethod(); i++ {
		// This is a bit of a hack since we can't easily get constants via reflection
		// in Go. We'll rely on the manual test above.
	}
	_ = constants // Avoid unused variable warning

	// Instead, let's verify that the registry contains exactly the expected types
	expectedCount := 7 // string, text, number, boolean, date, datetime, relation

	actualTypes := ListFieldTypes()
	if len(actualTypes) != expectedCount {
		t.Errorf("Expected %d field types in registry, got %d", expectedCount, len(actualTypes))
	}
}

func TestFieldTypeRegistryImmutability(t *testing.T) {
	// Test that calling ListFieldTypes multiple times returns consistent results
	types1 := ListFieldTypes()
	types2 := ListFieldTypes()

	if len(types1) != len(types2) {
		t.Error("ListFieldTypes() returned different lengths on subsequent calls")
	}

	// Create maps for comparison
	map1 := make(map[string]bool)
	map2 := make(map[string]bool)

	for _, t := range types1 {
		map1[t] = true
	}

	for _, t := range types2 {
		map2[t] = true
	}

	for typeName := range map1 {
		if !map2[typeName] {
			t.Errorf("Type '%s' missing from second call", typeName)
		}
	}

	for typeName := range map2 {
		if !map1[typeName] {
			t.Errorf("Type '%s' missing from first call", typeName)
		}
	}
}
