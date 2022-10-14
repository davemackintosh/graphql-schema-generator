package builder_test

// Checklist
// X test that the builder is actually doing the right thing (i.e. that it's
//   actually building the right thing) with a variety of different inputs
// X test nested named structs get added to the schema
// X test that anonymous structs get added to the schema and that they are
//   named TheStructTheFieldIsEmbeddedIn_FieldName
// X test recursive structs get added to the schema and don't cause infinite
//   recursion/failures
// X test that the builder can handle a struct with a field that is a pointer
//   to a struct
// X test that the builder can handle a struct with a field that is a slice of
//   structs
// - test that the builder can handle a non-concrete type (i.e. an interface{})
// - test that the builder can handle a struct with a field that is a map of
//   structs
// - test that the builder can handle a struct with a field that is a slice of
//   pointers to structs
// - test that the builder can handle a struct with a field that is a map of
//   pointers to structs

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
	ID        string         `json:"id" graphql:"description=The ID of the user"`
	Username  string         `json:"username" graphql:"description=The username of the user,decorators=[+unique()]"`
	Password  string         `json:"-"` // This field should not appear in the graphql schema.
	Email     *string        `json:"email" graphql:"description=The email of the user"`
	Phone     *string        `json:"phone" graphql:"description=The phone number of the user"`
	Roles     []Roles        `json:"roles" graphql:"description=The roles of the user"`
	Documents []UserDocument `json:"documents" graphql:"description=The documents of the user"`
}

type UserDocument struct {
	ID      string `json:"id" graphql:"description=The ID of the document"`
	Name    string `json:"name" graphql:"description=The name of the document"`
	Editors []User `json:"editors" graphql:"description=The editors of the document"`
	Meta    struct {
		Author       string `json:"author" graphql:"description=The author of the document"`
		LastModified string `json:"lastModified" graphql:"description=The last time the document was modified"`
	} `json:"meta" graphql:"description=The meta data of the document"`
}

type Project struct {
	ID   string            `json:"id" graphql:"description=The ID of the project"`
	Name string            `json:"name" graphql:"description=The name of the project"`
	Meta map[string]string `json:"meta" graphql:"description=The meta data of the project"`
}

type Test struct {
	name     string
	expected builder.GraphQLSchemaBuilder
	actual   builder.GraphQLSchemaBuilder
}

func TestBuilderStructSuite(t *testing.T) {
	expectedMeta := builder.Struct{
		Name: "UserDocument_Meta",
		Fields: &[]*builder.Field{
			{
				Name:            "author",
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The author of the document",
					},
				},
			},
			{
				Name:            "lastModified",
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The last time the document was modified",
					},
				},
			},
		},
	}
	expectedUserDocument := builder.Struct{
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
			{
				Name:            "meta",
				Type:            "UserDocument_Meta",
				IsPointer:       false,
				IncludeInOutput: true,
				IsStruct:        true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The meta data of the document",
					},
				},
			},
		},
	}
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
			{
				Name:            "documents",
				Type:            "UserDocument",
				IsPointer:       false,
				IsSlice:         true,
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The documents of the user",
					},
				},
			},
		},
	}

	tests := []Test{
		{
			name:   "TestBuilder_NestedFieldStructs",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddStruct(UserDocument{}, nil),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Structs: []*builder.Struct{
					&expectedUser,
					&expectedMeta,
					&expectedUserDocument,
				},
			},
		},
		{
			name: "TestBuilder_Enum",
			actual: builder.NewGraphQLSchemaBuilder(nil).AddEnum(builder.Enum{
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

func TestBuilder_Map(t *testing.T) {
	project := Project{}
	tests := []Test{
		{
			name:   "TestBuilder_Map",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddMap("ProjectMeta", project.Meta, nil),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Maps: []*builder.Map{
					{
						Name:   "ProjectMeta",
						Fields: &[]*builder.Field{},
					},
				},
			},
		},
		{
			name:   "TestBuilder_StructWithMap",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddStruct(project, nil),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Structs: []*builder.Struct{
					{
						Name: "Project",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the project",
									},
								},
							},
							{
								Name:            "name",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the project",
									},
								},
							},
							{
								Name:            "meta",
								Type:            "string",
								IsMap:           true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The meta data of the project",
									},
								},
							},
						},
					},
				},
				Maps: []*builder.Map{
					{
						Name:   "ProjectMeta",
						Fields: &[]*builder.Field{},
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
