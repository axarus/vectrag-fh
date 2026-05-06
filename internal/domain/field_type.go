package domain

type FieldType string

const (
	FieldString   FieldType = "string"
	FieldText     FieldType = "text"
	FieldNumber   FieldType = "number"
	FieldBoolean  FieldType = "boolean"
	FieldDate     FieldType = "date"
	FieldDateTime FieldType = "datetime"
	FieldRelation FieldType = "relation"
)

var fieldTypeRegistry = map[FieldType]struct{}{
	FieldString:   {},
	FieldText:     {},
	FieldNumber:   {},
	FieldBoolean:  {},
	FieldDate:     {},
	FieldDateTime: {},
	FieldRelation: {},
}

func IsValidType(t string) bool {
	_, ok := fieldTypeRegistry[FieldType(t)]
	return ok
}

func ListFieldTypes() []string {
	types := make([]string, 0, len(fieldTypeRegistry))
	for ft := range fieldTypeRegistry {
		types = append(types, string(ft))
	}

	return types
}
