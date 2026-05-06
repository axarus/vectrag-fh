package application

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewInitService(t *testing.T) {
	service := NewInitService()
	if service == nil {
		t.Fatal("NewInitService() returned nil")
	}
}

func TestInitService_InitializeProject_Success(t *testing.T) {
	service := NewInitService()

	// Create temporary directory for testing
	tempDir := t.TempDir()
	config := InitConfig{
		ProjectName: "test-project",
		Port:        "8080",
		Database:    "sqlite",
	}

	err := service.InitializeProject(tempDir, config)
	if err != nil {
		t.Fatalf("InitializeProject() error = %v", err)
	}

	// Verify directory structure was created
	projectPath := filepath.Join(tempDir, "test-project")

	// Check main directories
	expectedDirs := []string{"models", "config", ".vectrag"}
	for _, dir := range expectedDirs {
		fullPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", fullPath)
		}
	}

	// Check config files
	expectedFiles := []string{
		"vectrag.config.yaml",
		"config/database.config.yaml",
	}
	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", fullPath)
		}
	}
}

func TestInitService_InitializeProject_EmptyProjectName(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()
	config := InitConfig{
		ProjectName: "",
		Port:        "8080",
		Database:    "sqlite",
	}

	err := service.InitializeProject(tempDir, config)
	if err != nil {
		t.Errorf("Expected no error for empty project name, got: %v", err)
	}

	// Verify that it creates the project structure in the temp directory itself
	expectedDirs := []string{"models", "config", ".vectrag"}
	for _, dir := range expectedDirs {
		fullPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", fullPath)
		}
	}
}

func TestInitService_InitializeProject_NonEmptyDirectory(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()

	// Create a non-empty directory
	projectPath := filepath.Join(tempDir, "existing-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a file in the directory
	testFile := filepath.Join(projectPath, "existing.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := InitConfig{
		ProjectName: "existing-project",
		Port:        "8080",
		Database:    "sqlite",
	}

	err := service.InitializeProject(tempDir, config)
	if err == nil {
		t.Error("Expected error for non-empty directory, got nil")
	}
}

func TestInitService_InitializeProject_ExistingFileAtPath(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()

	// Create a file where we want to create the directory
	existingFile := filepath.Join(tempDir, "file-project")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := InitConfig{
		ProjectName: "file-project",
		Port:        "8080",
		Database:    "sqlite",
	}

	err := service.InitializeProject(tempDir, config)
	if err == nil {
		t.Error("Expected error when path exists but is not a directory, got nil")
	}
}

func TestInitService_createDirectoryStructure(t *testing.T) {
	service := NewInitService()
	tempDir := t.TempDir()

	err := service.createDirectoryStructure(tempDir)
	if err != nil {
		t.Fatalf("createDirectoryStructure() error = %v", err)
	}

	expectedDirs := []string{"models", "config", ".vectrag"}
	for _, dir := range expectedDirs {
		fullPath := filepath.Join(tempDir, dir)
		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", fullPath)
		} else if !info.IsDir() {
			t.Errorf("Expected %s to be a directory", fullPath)
		}
	}
}

func TestInitService_createDirectoryStructure_Error(t *testing.T) {
	service := NewInitService()

	// Use a non-existent parent directory that we can't create
	invalidPath := "/root/nonexistent/path"

	err := service.createDirectoryStructure(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
