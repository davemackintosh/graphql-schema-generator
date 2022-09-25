package builder

import (
	"reflect"

	structparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/struct_parser"
)

type GraphQLSchemaBuilderType struct{}

type GraphQLSchemaBuilderOptions struct {
	// A callback that takes the type name and the generated schema and writes it to an ioWriter.
	Writer *func(typeName, s string) error
}

type GraphQLSchemaBuilder struct {
	Types   *[]*structparser.Struct
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

func (b *GraphQLSchemaBuilder) AddType(t interface{}) *GraphQLSchemaBuilder {
	parsed, err := structparser.ParseStruct(reflect.TypeOf(t))
	if err != nil {
		panic(err)
	}

	*b.Types = append(*b.Types, parsed)
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
