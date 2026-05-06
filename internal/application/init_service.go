package application

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed resources/init/*
var initTemplates embed.FS

// InitConfig holds the configuration for project initialization
type InitConfig struct {
	ProjectName string
	Port        string
	Database    string
}

// InitService handles project initialization
type InitService struct{}

func NewInitService() *InitService {
	return &InitService{}
}

// InitializeProject creates a new VectraG project structure
func (s *InitService) InitializeProject(path string, config InitConfig) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	workDir := filepath.Join(absPath, config.ProjectName)

	// Check if directory exists
	if info, err := os.Stat(workDir); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", workDir)
		}
		// Check if directory is empty
		entries, err := os.ReadDir(workDir)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("directory is not empty: %s", workDir)
		}
	} else if os.IsNotExist(err) {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(workDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check directory: %w", err)
	}

	// Create directory structure
	if err := s.createDirectoryStructure(workDir); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate configuration files
	if err := s.generateConfigFiles(workDir, config); err != nil {
		return fmt.Errorf("failed to generate config files: %w", err)
	}

	return nil
}

// createDirectoryStructure creates the directory structure
func (s *InitService) createDirectoryStructure(basePath string) error {
	var dirs []string

	dirs = []string{
		"models",
		"config",
		".vectrag",
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateConfigFiles generates all configuration files from templates
func (s *InitService) generateConfigFiles(basePath string, config InitConfig) error {
	// Generate vectrag.config.yaml (always in root)
	if err := s.generateFile(basePath, "vectrag.config.yaml", "vectrag.config.yaml.tmpl", config); err != nil {
		return err
	}

	if err := s.generateFile(
		filepath.Join(basePath, "config"),
		"database.config.yaml",
		"database.config.yaml.tmpl",
		config,
	); err != nil {
		return err
	}

	return nil
}

// generateFile generates a single file from a template
func (s *InitService) generateFile(targetDir, fileName, templateName string, config InitConfig) error {
	// Read template
	templateContent, err := initTemplates.ReadFile(fmt.Sprintf("resources/init/%s", templateName))
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	// Parse template
	tmpl, err := template.New(fileName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	// Create target file
	filePath := filepath.Join(targetDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}
