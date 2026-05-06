package domain

import (
	"testing"
)

func TestValidateStatus(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"valid draft", StatusDraft, true},
		{"valid publish", StatusPublish, true},
		{"valid delete", StatusDelete, true},
		{"invalid status", Status("invalid"), false},
		{"empty status", Status(""), false},
		{"case sensitive draft", Status("Draft"), false},
		{"case sensitive publish", Status("Publish"), false},
		{"case sensitive delete", Status("Delete"), false},
		{"partial match", Status("dra"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateStatus(tc.status)
			if tc.expected && err != nil {
				t.Errorf("ValidateStatus(%s) error = %v, expected no error", tc.status, err)
			}
			if !tc.expected && err == nil {
				t.Errorf("ValidateStatus(%s) expected error, got nil", tc.status)
			}
		})
	}
}

func TestStatusConstants(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected string
	}{
		{"StatusDraft", StatusDraft, "draft"},
		{"StatusPublish", StatusPublish, "publish"},
		{"StatusDelete", StatusDelete, "delete"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(tc.status)
			if actual != tc.expected {
				t.Errorf("Expected %s = '%s', got '%s'", tc.name, tc.expected, actual)
			}
			
			// Verify the constant is valid
			err := ValidateStatus(tc.status)
			if err != nil {
				t.Errorf("Status constant %s ('%s') is not valid: %v", tc.name, actual, err)
			}
		})
	}
}

func TestStatusErrorMessages(t *testing.T) {
	invalidStatus := Status("invalid")
	err := ValidateStatus(invalidStatus)
	
	if err == nil {
		t.Error("Expected error for invalid status")
	}
	
	expectedMessage := "must be either 'draft', 'publish', or 'delete'"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}
