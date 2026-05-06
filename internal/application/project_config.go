package application

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Paths       ProjectPaths       `yaml:"paths"`
	Development ProjectDevelopment `yaml:"development"`
}

type ProjectPaths struct {
	Models string `yaml:"models"`
	Config string `yaml:"config"`
	Admin  string `yaml:"admin"`
}

type ProjectDevelopment struct {
	EnableCORS bool `yaml:"enableCORS"`
}

func FindProjectRoot(startDir string) (string, error) {
	if startDir == "" {
		return "", fmt.Errorf("start directory is empty")
	}

	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve start directory: %w", err)
	}

	for {
		candidate := filepath.Join(dir, "vectrag.config.yaml")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("vectrag.config.yaml not found in %s or any parent directory", startDir)
}

func LoadProjectConfig(projectRoot string) (ProjectConfig, error) {
	if projectRoot == "" {
		return ProjectConfig{}, fmt.Errorf("project root is empty")
	}

	path := filepath.Join(projectRoot, "vectrag.config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ProjectConfig{}, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	if cfg.Paths.Models == "" {
		return ProjectConfig{}, errors.New("paths.models is required")
	}

	return cfg, nil
}

func ResolveModelsDir(projectRoot string, cfg ProjectConfig) (string, error) {
	modelsPath := cfg.Paths.Models
	if filepath.IsAbs(modelsPath) {
		return modelsPath, nil
	}

	abs := filepath.Join(projectRoot, modelsPath)
	abs, err := filepath.Abs(abs)
	if err != nil {
		return "", fmt.Errorf("failed to resolve models dir: %w", err)
	}
	return abs, nil
}
