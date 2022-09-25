package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warpspeedboilerplate/graphql-schema-generator/internal/builder"
	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/field_parser"
	structparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/struct_parser"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type Roles string

const (
	RoleAdmin Roles = "admin"
	RoleUser  Roles = "user"
)

type User struct {
	ID       string  `graphql:"id, description=The ID of the user"`
	Username string  `graphql:"username, description=The username of the user,decorators=[+unique()]"`
	Password string  // This field should not appear in the graphql schema.
	Email    *string `graphql:"email, description=The email of the user"`
	Phone    *string `graphql:"phone, description=The phone number of the user"`
	Roles    Roles   `graphql:"roles, description=The roles of the user"`
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name     string
		expected *builder.GraphQLSchemaBuilder
		actual   *builder.GraphQLSchemaBuilder
	}{
		{
			name: "TestBuilder_Struct",
			actual: builder.NewGraphQLSchemaBuilder(nil).
				AddType(User{}),
			expected: &builder.GraphQLSchemaBuilder{
				Options: nil,
				Types: &[]*structparser.Struct{
					{
						Name: "User",
						Fields: &[]*fieldparser.Field{
							{
								Name: "id",
								Type: "String",
								ParsedTag: &tagparser.Tag{
									Name: "id",
									Options: &map[string]string{
										"description": "The ID of the user",
									},
								},
							},
							{
								Name: "username",
								Type: "String",
								ParsedTag: &tagparser.Tag{
									Name: "username",
									Options: &map[string]string{
										"description": "The username of the user",
										"decorators":  "[+unique()]",
									},
								},
							},
							{
								Name: "email",
								Type: "String",
								ParsedTag: &tagparser.Tag{
									Name: "email",
									Options: &map[string]string{
										"description": "The email of the user",
									},
								},
							},
							{
								Name: "phone",
								Type: "String",
								ParsedTag: &tagparser.Tag{
									Name: "phone",
									Options: &map[string]string{
										"description": "The phone number of the user",
									},
								},
							},
							{
								Name: "roles",
								Type: "Roles",
								ParsedTag: &tagparser.Tag{
									Name: "roles",
									Options: &map[string]string{
										"description": "The roles of the user",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual)
		})
	}
}
