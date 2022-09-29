package structparser

import (
	"reflect"

	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/field_parser"
)

type Struct struct {
	Name   string
	Fields *[]*fieldparser.Field
}

// Parse a struct by reflection into a Struct.
func ParseStruct(structType reflect.Type) *Struct {
	return &Struct{
		Name:   structType.Name(),
		Fields: fieldparser.GetFieldsFromStruct(structType),
	}
}
