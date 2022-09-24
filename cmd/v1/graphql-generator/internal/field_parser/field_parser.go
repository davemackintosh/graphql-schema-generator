package fieldparser

import (
	"fmt"
	"reflect"

	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/cmd/v1/graphql-generator/internal/tag-parser"
)

type Field struct {
	Name      string
	Type      string
	IsArray   bool
	IsPointer bool
	ParsedTag *tagparser.Tag
}

// Parse a struct field into its name and type with parsed tags as a Field.
func ParseField(field reflect.StructField) (*Field, error) {
	fieldParserType := &Field{
		Name:      field.Name,
		Type:      field.Type.String(),
		IsArray:   field.Type.Kind() == reflect.Slice,
		IsPointer: field.Type.Kind() == reflect.Ptr,
	}

	parsedTag, err := tagparser.GetTagFromField(field)
	if err != nil {
		return nil, fmt.Errorf("error parsing tag for field %s: %w", field.Name, err)
	}

	fieldParserType.ParsedTag = parsedTag

	return fieldParserType, nil
}

// GetFieldsFromStruct returns a map of field names to Fields for a given struct.
func GetFieldsFromStruct(structType reflect.Type) (*map[string]*Field, error) {
	fields := make(map[string]*Field)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		fieldParserType, err := ParseField(field)
		if err != nil {
			return nil, err
		}

		fields[field.Name] = fieldParserType
	}

	return &fields, nil
}
