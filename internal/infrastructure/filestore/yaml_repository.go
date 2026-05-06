package filestore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/axarus/vectrag/internal/domain"
	"gopkg.in/yaml.v3"
)

type YamlRepository struct {
	basePath string
}

func NewYamlRepository(basePath string) (*YamlRepository, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	return &YamlRepository{basePath: basePath}, nil
}

func (r *YamlRepository) modelFilePath(slug string) string {
	return filepath.Join(r.basePath, fmt.Sprintf("%s.yaml", slug))
}

func (r *YamlRepository) fallbackModelFilePath(slug string) string {
	return filepath.Join(r.basePath, fmt.Sprintf("%s.yml", slug))
}

func (r *YamlRepository) CreateModel(model domain.Model) error {
	filePath := r.modelFilePath(model.Slug)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("model with slug %s already exists", model.Slug)
	}
	if _, err := os.Stat(r.fallbackModelFilePath(model.Slug)); err == nil {
		return fmt.Errorf("model with slug %s already exists", model.Slug)
	}
	return r.saveModel(model)
}


func (r *YamlRepository) UpdateModel(model domain.Model) error {
	yamlPath := r.modelFilePath(model.Slug)
	if _, err := os.Stat(yamlPath); err == nil {
		return r.saveModel(model)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	ymlPath := r.fallbackModelFilePath(model.Slug)
	if _, err := os.Stat(ymlPath); err == nil {
		return r.saveModel(model)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat fallback file: %w", err)
	}

	return fmt.Errorf("model with slug %s does not exist", model.Slug)
}

func (r *YamlRepository) DeleteModel(slug string) error {
	filePath := r.modelFilePath(slug)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ymlPath := r.fallbackModelFilePath(slug)
		if _, ymlErr := os.Stat(ymlPath); os.IsNotExist(ymlErr) {
			return fmt.Errorf("model with slug %s does not exist", slug)
		}
		filePath = ymlPath
	}
	return os.Remove(filePath)
}

func (r *YamlRepository) GetModel(slug string) (domain.Model, error) {
	filePath := r.modelFilePath(slug)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ymlPath := r.fallbackModelFilePath(slug)
			ymlData, ymlErr := os.ReadFile(ymlPath)
			if ymlErr != nil {
				if os.IsNotExist(ymlErr) {
					return domain.Model{}, fmt.Errorf("model with slug %s not found", slug)
				}
				return domain.Model{}, fmt.Errorf("failed to read file: %w", ymlErr)
			}
			data = ymlData
		}
		if data == nil {
			return domain.Model{}, fmt.Errorf("failed to read file: %w", err)
		}
	}

	var dto modelDTO
	if err := yaml.Unmarshal(data, &dto); err != nil {
		return domain.Model{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return dto.toDomain(), nil
}

func (r *YamlRepository) GetModels() ([]domain.Model, error) {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var models []domain.Model
	seen := make(map[string]struct{})
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}
		model, err := r.GetModel(slug)
		if err != nil {
			continue
		}
		models = append(models, model)
	}

	return models, nil
}

func (r *YamlRepository) saveModel(model domain.Model) error {
	dto := modelDTOFromDomain(model)
	data, err := yaml.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal model: %w", err)
	}

	filePath := r.modelFilePath(model.Slug)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
