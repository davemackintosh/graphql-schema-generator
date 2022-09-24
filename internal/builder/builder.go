package builder

type GraphQLSchemaBuilderType struct{}

type GraphQLSchemaBuilderOptions struct {
	// A callback that takes the type name and the generated schema and writes it to an ioWriter.
	Writer *func(typeName, s string) error
}

type GraphQLSchemaBuilder struct {
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
