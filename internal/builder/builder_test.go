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
// X test that the builder can handle a struct with a field that is a slice of
//   pointers to structs
// X test that the builder can handle a struct with a field that is a map of
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

type Project struct {
	ID   string            `json:"id" graphql:"description=The ID of the project"`
	Name string            `json:"name" graphql:"description=The name of the project"`
	Meta map[string]string `json:"meta" graphql:"description=The meta data of the project"`
}

func TestBuilder_Map(t *testing.T) {
	project := Project{}
	tests := []Test{
		{
			name:   "TestBuilder_Map",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddMap("ProjectMeta", project.Meta),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Maps: []*builder.Map{
					{
						Name:    "ProjectMeta",
						KeyType: "string",
						Field: builder.Field{
							Name:            "",
							Type:            "string",
							IncludeInOutput: true,
						},
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
								Type:            "ProjectMeta",
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
						Name:    "ProjectMeta",
						KeyType: "string",
						Field: builder.Field{
							Name:            "",
							Type:            "string",
							IncludeInOutput: true,
							ParsedTag:       nil,
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

type Title struct {
	ID                  string              `json:"id" graphql:"description=The ID of the title"`
	Name                string              `json:"name" graphql:"description=The name of the title"`
	RentalPrice         *int                `json:"rentalPrice" graphql:"description=The rental price of the title"`
	BuyPrice            *int                `json:"buyPrice" graphql:"description=The buy price of the title"`
	Credits             *map[string]*string `json:"credits" graphql:"description=The credits of the title"`
	headOfficeReference string
}

type Shelf uint

const (
	ShelfAction Shelf = iota
	ShelfAdventure
	ShelfComedy
	ShelfDrama
	ShelfHorror
)

type DvdStore struct {
	ID              string            `json:"id" graphql:"description=The ID of the DVD store"`
	Name            string            `json:"name" graphql:"description=The name of the DVD store"`
	Address         string            `json:"address" graphql:"description=The address of the DVD store"`
	PhoneNumber     string            `json:"phoneNumber" graphql:"description=The phone number of the DVD store"`
	AvailableTitles map[Shelf][]Title `json:"availableTitles" graphql:"description=The available titles of the DVD store"`
}

func TestBuilder_ComplexDVDStore(t *testing.T) {
	tests := []Test{
		{
			name:   "TestBuilder_DVDStore",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddStruct(DvdStore{}, nil),
			expected: builder.GraphQLSchemaBuilder{
				Options: nil,
				Structs: []*builder.Struct{
					{
						Name: "Title",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the title",
									},
								},
							},
							{
								Name:            "name",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the title",
									},
								},
							},
							{
								Name:            "rentalPrice",
								Type:            "int",
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The rental price of the title",
									},
								},
							},
							{
								Name:            "buyPrice",
								Type:            "int",
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The buy price of the title",
									},
								},
							},
							{
								Name:            "credits",
								Type:            "TitleCredits",
								IsPointer:       true,
								IsMap:           true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The credits of the title",
									},
								},
							},
							{
								Name:            "headOfficeReference",
								Type:            "string",
								IncludeInOutput: false,
								ParsedTag:       nil,
							},
						},
					},
					{
						Name: "DvdStore",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the DVD store",
									},
								},
							},
							{
								Name:            "name",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the DVD store",
									},
								},
							},
							{
								Name:            "address",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The address of the DVD store",
									},
								},
							},
							{
								Name:            "phoneNumber",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The phone number of the DVD store",
									},
								},
							},
							{
								Name:            "availableTitles",
								Type:            "DvdStoreAvailableTitles",
								IsMap:           true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The available titles of the DVD store",
									},
								},
							},
						},
					},
				},
				Enums: []*builder.Enum{
					{
						Name:   "DvdStoreAvailableTitlesShelf",
						Values: []*builder.EnumKeyPairOptions{},
					},
				},
				Maps: []*builder.Map{
					{
						Name:    "TitleCredits",
						KeyType: "string",
						Field: builder.Field{
							Type:            "string",
							IsPointer:       true,
							IncludeInOutput: true,
						},
					},
					{
						Name:    "DvdStoreAvailableTitles",
						KeyType: "Shelf",
						Field: builder.Field{
							Type:            "Title",
							IsSlice:         true,
							IncludeInOutput: true,
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

type PlatformRegions string

const (
	PlatformRegionUS PlatformRegions = "US"
	PlatformRegionUK PlatformRegions = "UK"
)

type ProductImage struct {
	ThumbURL    string `json:"thumbUrl"`
	Featured    bool   `json:"featured"`
	ThumbWidth  int    `json:"thumbWidth"`
	ThumbHeight int    `json:"thumbHeight"`
	AltText     string `json:"altText"`
}

type ProductVariant struct {
	ID               string                                   `json:"id" description:"The ID of the product variant"`
	PriceExTax       *int                                     `json:"priceExTax" description:"The price of the product variant excluding tax"`
	AvailableRegions *map[PlatformRegions]*map[string]*string `json:"availableRegions" description:"The available regions of the product"`
	Images           []*ProductImage                          `json:"images" description:"The images of the product variant"`
}

type Product struct {
	ID               string                                   `json:"id" description:"The ID of the product"`
	Name             string                                   `json:"name" description:"The name of the product"`
	Description      string                                   `json:"description" description:"The description of the product"`
	PriceExTax       *int                                     `json:"priceExTax" description:"The price of the product excluding tax"`
	AvailableRegions *map[PlatformRegions]*map[string]*string `json:"availableRegions" description:"The available regions of the product"`
	Active           bool                                     `json:"active" description:"The active status of the product"`
	Variants         *[]*ProductVariant                       `json:"variants" description:"The variants of the product"`
	Images           *[]*ProductImage                         `json:"images" description:"The images of the product"`
}

type EcommerceStore struct {
	ID          string      `json:"id" description:"The ID of the ecommerce store"`
	Name        string      `json:"name" description:"The name of the ecommerce store"`
	Address     *string     `json:"address" description:"The address of the ecommerce store"`
	PhoneNumber *string     `json:"phoneNumber" description:"The phone number of the ecommerce store"`
	Products    *[]*Product `json:"products" description:"The products of the ecommerce store"`
}

func TestEcommerceStore(t *testing.T) {
	tests := []Test{
		{
			name:   "EcommerceStore",
			actual: *builder.NewGraphQLSchemaBuilder(nil).AddStruct(&EcommerceStore{}, nil),
			expected: builder.GraphQLSchemaBuilder{
				Structs: []*builder.Struct{
					{
						Name: "ProductVariant",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the product variant",
									},
								},
							},
							{
								Name:            "priceExTax",
								Type:            "int",
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The price of the product variant excluding tax",
									},
								},
							},
							{
								Name:            "availableRegions",
								Type:            "PlatformRegions",
								IsMap:           true,
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The available regions of the product",
									},
								},
							},
							{
								Name:            "images",
								Type:            "ProductImage",
								IsSlice:         true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The images of the product variant",
									},
								},
							},
						},
					},
					{
						Name: "ProductImage",
						Fields: &[]*builder.Field{
							{
								Name:            "ThumbURL",
								Type:            "string",
								IncludeInOutput: true,
							},
							{
								Name:            "Featured",
								Type:            "bool",
								IncludeInOutput: true,
							},
							{
								Name:            "ThumbWidth",
								Type:            "int",
								IncludeInOutput: true,
							},
							{
								Name:            "ThumbHeight",
								Type:            "int",
								IncludeInOutput: true,
							},
							{
								Name:            "AltText",
								Type:            "string",
								IncludeInOutput: true,
							},
							{
								Name:            "origURL",
								Type:            "string",
								IncludeInOutput: false,
							},
						},
					},
					{
						Name: "Product",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the product",
									},
								},
							},
							{
								Name:            "name",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the product",
									},
								},
							},
							{
								Name:            "description",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The description of the product",
									},
								},
							},
							{
								Name:            "priceExTax",
								Type:            "int",
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The price of the product excluding tax",
									},
								},
							},
							{
								Name:            "availableRegions",
								Type:            "ProductAvailableRegions",
								IsMap:           true,
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The available regions of the product",
									},
								},
							},
							{
								Name:            "active",
								Type:            "bool",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The active status of the product",
									},
								},
							},
							{
								Name:            "variants",
								Type:            "ProductVariant",
								IsSlice:         true,
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The variants of the product",
									},
								},
							},
							{
								Name:            "images",
								Type:            "ProductImage",
								IsSlice:         true,
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The images of the product",
									},
								},
							},
						},
					},
					{
						Name: "EcommerceStore",
						Fields: &[]*builder.Field{
							{
								Name:            "id",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the ecommerce store",
									},
								},
							},
							{
								Name:            "name",
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the ecommerce store",
									},
								},
							},
							{
								Name:            "address",
								Type:            "string",
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The address of the ecommerce store",
									},
								},
							},
							{
								Name:            "phoneNumber",
								Type:            "string",
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The phone number of the ecommerce store",
									},
								},
							},
							{
								Name:            "products",
								Type:            "Product",
								IsSlice:         true,
								IsPointer:       true,
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The products of the ecommerce store",
									},
								},
							},
						},
					},
				},
				Maps: []*builder.Map{
					{
						Name:    "AvailablePlatformRegions",
						KeyType: "string",
						Field: builder.Field{
							IncludeInOutput: true,
							IsPointer:       true,
							Type:            "string",
						},
					},
					{
						Name:    "EcommerceStoreAvailablePlatformRegions",
						KeyType: "EcommerceStorePlatformRegions",
						Field: builder.Field{
							IsMap:           true,
							IncludeInOutput: true,
							IsPointer:       true,
							Type:            "EcommerceStorePlatformRegionsMap1",
						},
					},
				},
				Enums: []*builder.Enum{
					{
						Name: "EcommerceStorePlatformRegions",
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
