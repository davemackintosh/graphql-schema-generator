package builder_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warpspeedboilerplate/graphql-schema-generator/cmd/v1/graphql-generator/internal/builder"
	"github.com/warpspeedboilerplate/graphql-schema-generator/cmd/v1/graphql-generator/internal/ptr"
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

const expectedUserSchema = `enum Roles {
	ADMIN
	USER
}

type User {
	// The ID of the User
	id: String! @doc(description: "The ID of the user")
	// The username of the User
	username: String! @doc(description: "The username of the user") @unique()
	// The email of the User
	email: String @doc(description: "The email of the user")
	// The phone number of the User
	phone: String @doc(description: "The phone number of the user")
	// The roles of the User
	roles: Roles! @doc(description: "The roles of the user")
}`

func TestBuilder(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		actual   string
	}{
		{
			name: "TestBuilder_Struct",
			actual: builder.NewGraphQLSchemaBuilder(nil).
				AddType(User{}).
				Build(),
			expected: expectedUserSchema,
		},
		{
			name: "TestBuilder_Struct_Add_Query",
			actual: builder.NewGraphQLSchemaBuilder(nil).
				AddType(User{}).
				AddQuery("currentUser", "User", ptr.Of("Get a user by ID")).
				Returns("User", ptr.Of("The user")).
				Build(),
			expected: fmt.Sprintf(`%s

type Query {
	// Get a user by ID
	currentUser: User
}`, expectedUserSchema),
		},
		{
			name: "TestBuilder_Struct_Add_Mutation",
			actual: builder.NewGraphQLSchemaBuilder(nil).
				AddType(User{}).
				AddMutation("createUser", "User", ptr.Of("Create a new user")).
				Returns("User", ptr.Of("The user")).
				Build(),
			expected: fmt.Sprintf(`%s

type Mutation {
	// Create a new user
	createUser($input: UserInput!): User
}`, expectedUserSchema),
		},
		{
			name: "TestBuilder_Struct_Add_Options",
			actual: builder.NewGraphQLSchemaBuilder(&builder.GraphQLSchemaBuilderOptions{
				Writer: ptr.Of(func(typeName, s string) error {
					return nil
				}),
			}).
				AddType(User{}).
				Build(),
			expected: fmt.Sprintf(`%s

type Mutation {
	// Create a new user
	createUser($input: UserInput!): User
}`, expectedUserSchema),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.actual, tt.expected)
		})
	}
}
