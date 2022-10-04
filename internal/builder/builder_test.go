package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warpspeedboilerplate/graphql-schema-generator/internal/builder"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type Roles string

const (
	RoleAdmin Roles = "admin"
	RoleUser  Roles = "user"
)

type User struct {
	ID       string  `json:"id" graphql:"description=The ID of the user"`
	Username string  `json:"username" graphql:"description=The username of the user,decorators=[+unique()]"`
	Password string  `json:"-"` // This field should not appear in the graphql schema.
	Email    *string `json:"email" graphql:"description=The email of the user"`
	Phone    *string `json:"phone" graphql:"description=The phone number of the user"`
	Roles    []Roles `json:"roles" graphql:"description=The roles of the user"`
}

type UserDocument struct {
	ID      string `json:"id" graphql:"description=The ID of the document"`
	Name    string `json:"name" graphql:"description=The name of the document"`
	Editors []User `json:"editors" graphql:"description=The editors of the document"`
}

func TestBuilder(t *testing.T) {
	expectedUser := builder.Struct{
		Name: "User",
		Fields: &[]*builder.Field{
			{
				Name:            "id",
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The ID of the user",
					},
				},
			},
			{
				Name:            "username",
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The username of the user",
						"decorators":  "[+unique()]",
					},
				},
			},
			{
				Name:            "Password",
				Type:            "string",
				IncludeInOutput: false,
				ParsedTag:       nil,
			},
			{
				Name:            "email",
				Type:            "string",
				IsPointer:       true,
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The email of the user",
					},
				},
			},
			{
				Name:            "phone",
				Type:            "string",
				IsPointer:       true,
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The phone number of the user",
					},
				},
			},
			{
				Name:            "roles",
				Type:            "Roles",
				IsPointer:       false,
				IsSlice:         true,
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The roles of the user",
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		expected builder.GraphQLSchemaBuilder
		actual   builder.GraphQLSchemaBuilder
	}{
		/*{
			name: "TestBuilder_NestedFieldStructs",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddStruct(UserDocument{
				ID:   "123",
				Name: "Test",
				Editors: []User{
					{},
				},
			}),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Structs: []*builder.Struct{
					{
						Name: "UserDocument",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the document",
									},
								},
							},
							{
								Name:            "name",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the document",
									},
								},
							},
							{
								Name:            "editors",
								Type:            "User",
								IsPointer:       false,
								IsSlice:         true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The editors of the document",
									},
								},
							},
						},
					},
					&expectedUser,
				},
			},
		},*/
		{
			name: "TestBuilder_Struct",
			actual: *builder.NewGraphQLSchemaBuilder(nil).
				AddStruct(User{}),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Structs: []*builder.Struct{
					&expectedUser,
				},
			},
		},
		{
			name: "TestBuilder_Enum",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddEnum(builder.Enum{
				Name: "Roles",
				Values: []*builder.EnumKeyPairOptions{
					{
						Key:   "ADMIN",
						Value: RoleAdmin,
					},
					{
						Key:   "USER",
						Value: RoleUser,
					},
				},
			}),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Enums: []*builder.Enum{
					{
						Name: "Roles",
						Values: []*builder.EnumKeyPairOptions{
							{
								Key:   "ADMIN",
								Value: RoleAdmin,
							},
							{
								Key:   "USER",
								Value: RoleUser,
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
