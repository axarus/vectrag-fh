package database

import (
	"fmt"

	"github.com/axarus/vectrag/internal/domain"
)

// FieldTypeMapper converts domain FieldType to SQL types
type FieldTypeMapper struct{}

// NewFieldTypeMapper creates a new FieldTypeMapper instanhce
func NewFieldTypeMapper() *FieldTypeMapper {
	return &FieldTypeMapper{}
}

// ToSQLType converts a domain FieldType to its corresponding SQL type
func (m *FieldTypeMapper) ToSQLType(ft domain.FieldType) (string, error) {
	switch ft {
	case domain.FieldString:
		return "VARCHAR(255)", nil
	case domain.FieldText:
		return "TEXT", nil
	case domain.FieldNumber:
		return "DOUBLE PRECISION", nil
	case domain.FieldBoolean:
		return "BOOLEAN", nil
	case domain.FieldDate:
		return "DATE", nil
	case domain.FieldDateTime:
		return "TIMESTAMP", nil
	case domain.FieldRelation:
		return "VARCHAR(255)", nil
	default:
		return "", fmt.Errorf("unsupported field type: %s", ft)
	}
}

// ToSQLTypeString is a convenience function that works with FieldType as string
func ToSQLTypeString(fieldType string) (string, error) {
	mapper := NewFieldTypeMapper()
	return mapper.ToSQLType(domain.FieldType(fieldType))
}
