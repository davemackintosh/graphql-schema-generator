package fieldparser

import (
	"reflect"
	"strings"

	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type Field struct {
	Name            string
	Type            string
	IsArray         bool
	IsPointer       bool
	IncludeInOutput bool
	ParsedTag       *tagparser.Tag
}

// Parse a struct field into its name and type with parsed tags as a Field.
func ParseField(field reflect.StructField) *Field {
	var fieldName string

	jsonTag := field.Tag.Get("json")
	jsonTagParts := strings.Split(jsonTag, ",")
	jsonTagName := jsonTagParts[0]

	if jsonTagName != "" && jsonTagName != "-" {
		fieldName = jsonTagName
	} else {
		fieldName = field.Name
	}

	fieldParserType := &Field{
		Name:            fieldName,
		Type:            field.Type.Name(),
		IsArray:         field.Type.Kind() == reflect.Slice,
		IsPointer:       field.Type.Kind() == reflect.Ptr,
		ParsedTag:       tagparser.GetTagFromField(field),
		IncludeInOutput: fieldName != "-",
	}

	return fieldParserType
}

// GetFieldsFromStruct returns a map of field names to Fields for a given struct.
func GetFieldsFromStruct(structType reflect.Type) *[]*Field {
	var fields []*Field

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldParserType := ParseField(field)

		fields = append(fields, fieldParserType)
	}

	return &fields
}
