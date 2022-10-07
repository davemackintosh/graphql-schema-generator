package builder

import (
	"reflect"
	"strings"

	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type Field struct {
	Name            string
	Type            string
	IsSlice         bool
	IsPointer       bool
	IsStruct        bool
	IncludeInOutput bool
	ParsedTag       *tagparser.Tag
}

type Struct struct {
	Name   string
	Fields *[]*Field
}

type EnumKeyPairOptions struct {
	Key         string
	Value       interface{}
	Description *string
}

type Enum struct {
	Name   string
	Values []*EnumKeyPairOptions
}

type GraphQLSchemaBuilderWriter interface {
	WriteSchema(schema string)
}

type GraphQLSchemaBuilderOptions struct {
	// A callback that takes the type name and the generated schema and writes it to an ioWriter.
	Writer GraphQLSchemaBuilderWriter
}

type GraphQLSchemaBuilder struct {
	Structs []*Struct
	Enums   []*Enum
	Options *GraphQLSchemaBuilderOptions
}

func NewGraphQLSchemaBuilder(options *GraphQLSchemaBuilderOptions) *GraphQLSchemaBuilder {
	return &GraphQLSchemaBuilder{
		Options: options,
	}
}

func (b *GraphQLSchemaBuilder) AddMutation(name, typeName string, description *string) *GraphQLSchemaBuilder {
	return b
}

func (b *GraphQLSchemaBuilder) AddQuery(name, typeName string, description *string) *GraphQLSchemaBuilder {
	return b
}

func (b *GraphQLSchemaBuilder) AddEnum(enum Enum) *GraphQLSchemaBuilder {
	b.Enums = append(b.Enums, &enum)

	return b
}

func (b *GraphQLSchemaBuilder) AddStruct(t interface{}) *GraphQLSchemaBuilder {
	structType := reflect.ValueOf(t)
	structName := structType.Type().Name()

	if structName == "" {
		panic("AddStruct struct name cannot be empty")
	}

	// Loop over the struct's fields and add them to the list of fields.
	var fields []*Field

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Type().Field(i)
		fieldType := field.Type
		fieldTypeName := fieldType.String()

		// If the field is a pointer, get the type it points to.
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			fieldTypeName = fieldType.String()
		}

		// If the field is a slice of a struct, we need to get the struct and add it.
		if structType.Field(i).Kind() == reflect.Slice && structType.Field(i).Type().Elem().Kind() == reflect.Struct {
			b.AddStruct(reflect.New(structType.Field(i).Type().Elem()).Elem().Interface())
		} else if structType.Field(i).Kind() == reflect.Struct {
			// Otherwise, if it's a struct, we need to add it.
			b.AddStruct(structType.Field(i).Interface())
		}

		fieldName := field.Name

		// If the fieldName has a period in it, it's a package name.Type and we only want the type name.
		if strings.Contains(fieldTypeName, ".") {
			fieldTypeNameParts := strings.Split(fieldTypeName, ".")
			fieldTypeName = fieldTypeNameParts[len(fieldTypeNameParts)-1]
		}

		// Get the field name from the json tag and fallback to the field name.
		jsonTag := field.Tag.Get("json")
		jsonTagParts := strings.Split(jsonTag, ",")
		jsonTagName := jsonTagParts[0]

		if jsonTagName != "" && jsonTagName != "-" {
			fieldName = jsonTagName
		}

		fields = append(fields, &Field{
			Name:            fieldName,
			Type:            fieldTypeName,
			IsPointer:       field.Type.Kind() == reflect.Ptr,
			IsSlice:         field.Type.Kind() == reflect.Slice,
			IsStruct:        field.Type.Kind() == reflect.Struct,
			ParsedTag:       tagparser.ParseTag(field.Tag.Get("graphql"), field.Name),
			IncludeInOutput: field.Tag.Get("graphql") != "-" && field.Tag.Get("json") != "-",
		})
	}

	b.Structs = append(b.Structs, &Struct{
		Name:   structName,
		Fields: &fields,
	})

	return b
}

func (b *GraphQLSchemaBuilder) Returns(typeName string, description *string) *GraphQLSchemaBuilder {
	return b
}

func (b *GraphQLSchemaBuilder) WithInputType(t interface{}) *GraphQLSchemaBuilder {
	return b
}

func (b *GraphQLSchemaBuilder) Build() string {
	return ""
}
