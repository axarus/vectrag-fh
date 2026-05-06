package filestore

import (
	"time"

	"github.com/axarus/vectrag/internal/domain"
)

type modelDTO struct {
	ID            string     `yaml:"id"`
	Name          string     `yaml:"name"`
	Slug          string     `yaml:"slug"`
	Description   string     `yaml:"description,omitempty"`
	Fields        []fieldDTO `yaml:"fields"`
	Status        string     `yaml:"status"`
	SchemaVersion int        `yaml:"schemaVersion"`
}

type fieldDTO struct {
	ID          string    `yaml:"id"`
	Name        string    `yaml:"name"`
	Type        string    `yaml:"type"`
	Description string    `yaml:"description,omitempty"`
	Unique      bool      `yaml:"unique,omitempty"`
	Required    bool      `yaml:"required,omitempty"`
	Status      string    `yaml:"status"`
	CreatedAt   time.Time `yaml:"createdAt,omitempty"`
	UpdatedAt   time.Time `yaml:"updatedAt,omitempty"`
}

func modelDTOFromDomain(m domain.Model) modelDTO {
	fields := make([]fieldDTO, len(m.Fields))
	for i, f := range m.Fields {
		fields[i] = fieldDTO{
			ID:          f.ID,
			Name:        f.Name,
			Type:        string(f.Type),
			Description: f.Description,
			Unique:      f.Unique,
			Required:    f.Required,
			Status:      string(f.Status),
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		}
	}

	return modelDTO{
		ID:            m.ID,
		Name:          m.Name,
		Slug:          m.Slug,
		Description:   m.Description,
		Fields:        fields,
		Status:        string(m.Status),
		SchemaVersion: m.SchemaVersion,
	}
}

func (dto modelDTO) toDomain() domain.Model {
	fields := make([]domain.Field, len(dto.Fields))
	for i, f := range dto.Fields {
		fields[i] = domain.Field{
			ID:          f.ID,
			Name:        f.Name,
			Type:        domain.FieldType(f.Type),
			Description: f.Description,
			Unique:      f.Unique,
			Required:    f.Required,
			Status:      domain.Status(f.Status),
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		}
	}

	return domain.Model{
		ID:            dto.ID,
		Name:          dto.Name,
		Slug:          dto.Slug,
		Description:   dto.Description,
		Fields:        fields,
		Status:        domain.Status(dto.Status),
		SchemaVersion: dto.SchemaVersion,
	}
}
