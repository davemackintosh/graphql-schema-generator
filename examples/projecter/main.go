package main

import (
	"log"

	builder "github.com/warpspeedboilerplate/graphql-schema-generator/internal/builder"
)

type Roles uint

const (
	Roles_USER  Roles = iota //nolint: revive,var-naming,stylecheck
	Roles_ADMIN Roles = 1    //nolint: revive,var-naming,stylecheck
)

type User struct {
	ID       string    `json:"id" graphql:"description=The ID of the user"`
	Username string    `json:"username" graphql:"description=The username of the user,decorators=[+unique()]"`
	Password string    `json:"-" graphql:"-"`
	Email    *string   `json:"email" graphql:"description=The email of the user"`
	Phone    *string   `json:"phone,omitempty" graphql:"description=The phone number of the user"`
	Roles    []Roles   `json:"roles" graphql:"description=The roles of the user"`
	Projects []Project `json:"projects" graphql:"description=The projects of the user"`
}

type Project struct {
	Name     string            `json:"name" graphql:"description=The name of the project"`
	Meta     map[string]string `json:"meta" graphql:"description=The meta data of the project"`
	Editors  []User            `json:"editors" graphql:"description=The editors of the project"`
	Archived bool              `json:"archived" graphql:"description=Whether the project is archived"`
}

// You need to configure the builder to use the writer you want to use.
// This example uses a simple writer that just prints the schema to the console.
type GraphQLSchemaFileWriter struct{}

func (w *GraphQLSchemaFileWriter) WriteSchema(schema string) {
	log.Println(schema)
}

func main() {
	// First we must configure the builder instance.
	// This simple project doesn't any plugins.
	schema := builder.NewGraphQLSchemaBuilder(&builder.GraphQLSchemaBuilderOptions{
		Writer: &GraphQLSchemaFileWriter{},
	})

	// Now that the schema builder is configured, we can add types to it.
	// This will take the User struct and recursively parse it and discover all the types and other
	// nested structs that it contains and add them to our builder schema.
	// GOTCHA: Due to the way Go handles enums, we need to manually add our Roles enum to the schema a bit later.
	schema.AddStruct(User{}, nil)

	// The Schema should now have both the user and project types registered to it, due to the way that
	// Go handles enums we now must also add the Roles enum to the schema.
	schema.AddEnum(builder.Enum{
		Name: "Roles",
		Values: []*builder.EnumKeyPairOptions{
			{
				Key:   "USER",
				Value: "0",
			},
			{
				Key:   "ADMIN",
				Value: "1",
			},
		},
	})

	// Now that we have added all the types to the schema, we can generate the schema.
	schema.Build()
}
