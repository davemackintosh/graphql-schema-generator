package structparser

import (
	"fmt"
	"reflect"

	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/field_parser"
)

type Struct struct {
	Name   string
	Fields *[]*fieldparser.Field
}

// Parse a struct by reflection into a Struct.
func ParseStruct(structType reflect.Type) (*Struct, error) {
	fields, err := fieldparser.GetFieldsFromStruct(structType)
	if err != nil {
		return nil, fmt.Errorf("error getting fields from struct: %w", err)
	}

	return &Struct{
		Name:   structType.Name(),
		Fields: fields,
	}, nil
}
