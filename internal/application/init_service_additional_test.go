package application

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitService_generateConfigFiles_Error(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()
	
	// Create a config that will cause an error
	config := InitConfig{
		ProjectName: "test-project",
		Port:        "8080",
		Database:    "sqlite",
	}

	// Test with invalid template path by temporarily modifying the service
	// We'll test the error handling in generateFile directly
	err := service.generateFile(tempDir, "test.txt", "nonexistent.tmpl", config)
	if err == nil {
		t.Error("Expected error for nonexistent template, got nil")
	}
}

func TestInitService_generateFile_TemplateParseError(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()
	config := InitConfig{
		ProjectName: "test",
		Port:        "8080",
		Database:    "sqlite",
	}

	// Create a temporary invalid template file
	invalidTemplate := `{{.InvalidField}}`
	templateDir := filepath.Join(tempDir, "resources", "init")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	
	templateFile := filepath.Join(templateDir, "invalid.tmpl")
	if err := os.WriteFile(templateFile, []byte(invalidTemplate), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// This should cause a template execution error
	err := service.generateFile(tempDir, "test.txt", "invalid.tmpl", config)
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}
}

func TestInitService_generateFile_CreateFileError(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()
	config := InitConfig{
		ProjectName: "test",
		Port:        "8080",
		Database:    "sqlite",
	}

	// Try to create a file in a non-existent directory
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	err := service.generateFile(nonExistentDir, "test.txt", "vectrag.config.yaml.tmpl", config)
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestInitService_InitializeProject_InvalidPath(t *testing.T) {
	service := NewInitService()
	config := InitConfig{
		ProjectName: "test-project",
		Port:        "8080",
		Database:    "sqlite",
	}

	// Use an invalid path that should cause filepath.Abs to fail
	invalidPath := string([]byte{0x00}) // null byte
	err := service.InitializeProject(invalidPath, config)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
