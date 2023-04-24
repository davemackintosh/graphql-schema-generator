package typeparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tagparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/graphql-tag-parser"
	"github.com/warpspeed-cloud/graphql-schema-generator/internal/ptr"
	typeparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/type-parser"
)

type Test struct {
	name     string
	actual   *typeparser.TypeParser
	expected *typeparser.TypeParser
}

func TestMaps(t *testing.T) {
	tests := []Test{
		{
			name:   "Basic map",
			actual: typeparser.NewTypeParser(nil).AddMap("StringString", map[string]string{}),
			expected: &typeparser.TypeParser{
				Maps: &[]typeparser.Map{
					{
						Name: "StringString",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
							Type: "string",
						},
					},
				},
			},
		},
		{
			name:   "Basic map with pointer value",
			actual: typeparser.NewTypeParser(nil).AddMap("StringPointerString", map[string]*string{}),
			expected: &typeparser.TypeParser{
				Maps: &[]typeparser.Map{
					{
						Name: "StringPointerString",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
							Type:      "string",
							IsPointer: true,
						},
					},
				},
			},
		},
		{
			name:   "Map with pointer key",
			actual: typeparser.NewTypeParser(nil).AddMap("PointerStringString", map[*string]string{}),
			expected: &typeparser.TypeParser{
				Maps: &[]typeparser.Map{
					{
						Name: "PointerStringString",
						Key: typeparser.TypeDescriptor{
							Type:      "string",
							IsPointer: true,
						},
						Val: typeparser.TypeDescriptor{
							Type: "string",
						},
					},
				},
			},
		},
		{
			name:   "Map with pointer key and value",
			actual: typeparser.NewTypeParser(nil).AddMap("PointerStringPointerString", map[*string]*string{}),
			expected: &typeparser.TypeParser{
				Maps: &[]typeparser.Map{
					{
						Name: "PointerStringPointerString",
						Key: typeparser.TypeDescriptor{
							Type:      "string",
							IsPointer: true,
						},
						Val: typeparser.TypeDescriptor{
							Type:      "string",
							IsPointer: true,
						},
					},
				},
			},
		},
		{
			name:   "Map with map value",
			actual: typeparser.NewTypeParser(nil).AddMap("StringMapStringString", map[string]map[string]string{}),
			expected: &typeparser.TypeParser{
				Maps: &[]typeparser.Map{
					{
						Name: "StringMapStringStringMap1",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
							Type: "string",
						},
					},
					{
						Name: "StringMapStringString",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
							Type: "StringMapStringStringMap1",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, test.actual)
		})
	}
}

type Roles string

const (
	RoleAdmin Roles = "admin"
	RoleUser  Roles = "user"
)

type User struct {
	ID        string         `json:"id" graphql:"description=The ID of the user"`
	Username  string         `json:"username" graphql:"description=The username of the user,decorators=[+unique()]"`
	Password  string         `json:"-"` // This field should not appear in the graphql schema but should appear in the AST.
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

func TestBuilderStructSuite(t *testing.T) {
	expectedMeta := typeparser.Struct{
		Name: "Struct1",
		Fields: &[]typeparser.TypeDescriptor{
			{
				Name:            ptr.Of("author"),
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The author of the document",
					},
				},
			},
			{
				Name:            ptr.Of("lastModified"),
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
	expectedUserDocument := typeparser.Struct{
		Name: "UserDocument",
		Fields: &[]typeparser.TypeDescriptor{
			{
				Name:            ptr.Of("id"),
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The ID of the document",
					},
				},
			},
			{
				Name:            ptr.Of("name"),
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The name of the document",
					},
				},
			},
			{
				Name:            ptr.Of("editors"),
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
				Name:            ptr.Of("meta"),
				Type:            "UserDocumentStruct1",
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
	expectedUser := typeparser.Struct{
		Name: "User",
		Fields: &[]typeparser.TypeDescriptor{
			{
				Name:            ptr.Of("id"),
				Type:            "string",
				IncludeInOutput: true,
				ParsedTag: &tagparser.Tag{
					Options: map[string]string{
						"description": "The ID of the user",
					},
				},
			},
			{
				Name:            ptr.Of("username"),
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
				Name:            ptr.Of("Password"),
				Type:            "string",
				IncludeInOutput: false,
				ParsedTag:       nil,
			},
			{
				Name:            ptr.Of("email"),
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
				Name:            ptr.Of("phone"),
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
				Name: ptr.Of("roles"),
				// This is a string because we cannot get the named type
				// of a "enum" like type in Go (it's just a string.)
				Type:            "string",
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
				Name:            ptr.Of("documents"),
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
			actual: typeparser.NewTypeParser(nil).AddStruct(UserDocument{}, nil),
			expected: &typeparser.TypeParser{
				Structs: &[]typeparser.Struct{
					expectedUser,
					expectedMeta,
					expectedUserDocument,
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
			name:   "TestBuilder_StructWithMap",
			actual: typeparser.NewTypeParser(nil).AddStruct(project, nil),
			expected: &typeparser.TypeParser{
				Structs: &[]typeparser.Struct{
					{
						Name: "Project",
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("id"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the project",
									},
								},
							},
							{
								Name:            ptr.Of("name"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the project",
									},
								},
							},
							{
								Name:            ptr.Of("meta"),
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
				Maps: &[]typeparser.Map{
					{
						Name: "ProjectMeta",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
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

// This exists to shut go up about the param not being used.
func (t Title) HeadOfficeReference() string {
	return t.headOfficeReference
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
			actual: typeparser.NewTypeParser(nil).AddStruct(DvdStore{}, nil),
			expected: &typeparser.TypeParser{
				Structs: &[]typeparser.Struct{
					{
						Name: "Title",
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("id"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the title",
									},
								},
							},
							{
								Name:            ptr.Of("name"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the title",
									},
								},
							},
							{
								Name:            ptr.Of("rentalPrice"),
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
								Name:            ptr.Of("buyPrice"),
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
								Name:            ptr.Of("credits"),
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
								Name:            ptr.Of("headOfficeReference"),
								Type:            "string",
								IncludeInOutput: false,
								ParsedTag:       nil,
							},
						},
					},
					{
						Name: "DvdStore",
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("id"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the DVD store",
									},
								},
							},
							{
								Name:            ptr.Of("name"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the DVD store",
									},
								},
							},
							{
								Name:            ptr.Of("address"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The address of the DVD store",
									},
								},
							},
							{
								Name:            ptr.Of("phoneNumber"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The phone number of the DVD store",
									},
								},
							},
							{
								Name:            ptr.Of("availableTitles"),
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
				//Enums: &[]typeparser.Enum{
				//	{
				//		Name:   "DvdStoreAvailableTitlesShelf",
				//		Values: []typeparser.EnumKeyPairOptions{},
				//	},
				//},
				Maps: &[]typeparser.Map{
					{
						Name: "TitleCredits",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
							Type:      "string",
							IsPointer: true,
						},
					},
					{
						Name: "DvdStoreAvailableTitles",
						Key: typeparser.TypeDescriptor{
							Type: "uint",
						},
						Val: typeparser.TypeDescriptor{
							Type:    "Title",
							IsSlice: true,
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
			actual: typeparser.NewTypeParser(nil).AddStruct(&EcommerceStore{}, nil),
			expected: &typeparser.TypeParser{
				Structs: &[]typeparser.Struct{
					{
						Name: "EcommerceStore",
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("id"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the ecommerce store",
									},
								},
							},
							{
								Name:            ptr.Of("name"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the ecommerce store",
									},
								},
							},
							{
								Name:            ptr.Of("address"),
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
								Name:            ptr.Of("phoneNumber"),
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
								Name:            ptr.Of("products"),
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
					{
						Name: "ProductVariant",
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("id"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the product variant",
									},
								},
							},
							{
								Name:            ptr.Of("priceExTax"),
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
								Name:            ptr.Of("availableRegions"),
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
								Name:            ptr.Of("images"),
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
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("ThumbURL"),
								Type:            "string",
								IncludeInOutput: true,
							},
							{
								Name:            ptr.Of("Featured"),
								Type:            "bool",
								IncludeInOutput: true,
							},
							{
								Name:            ptr.Of("ThumbWidth"),
								Type:            "int",
								IncludeInOutput: true,
							},
							{
								Name:            ptr.Of("ThumbHeight"),
								Type:            "int",
								IncludeInOutput: true,
							},
							{
								Name:            ptr.Of("AltText"),
								Type:            "string",
								IncludeInOutput: true,
							},
							{
								Name:            ptr.Of("origURL"),
								Type:            "string",
								IncludeInOutput: false,
							},
						},
					},
					{
						Name: "Product",
						Fields: &[]typeparser.TypeDescriptor{
							{
								Name:            ptr.Of("id"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The ID of the product",
									},
								},
							},
							{
								Name:            ptr.Of("name"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The name of the product",
									},
								},
							},
							{
								Name:            ptr.Of("description"),
								Type:            "string",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The description of the product",
									},
								},
							},
							{
								Name:            ptr.Of("priceExTax"),
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
								Name:            ptr.Of("availableRegions"),
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
								Name:            ptr.Of("active"),
								Type:            "bool",
								IncludeInOutput: true,
								ParsedTag: &tagparser.Tag{
									Options: map[string]string{
										"description": "The active status of the product",
									},
								},
							},
							{
								Name:            ptr.Of("variants"),
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
								Name:            ptr.Of("images"),
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
				},
				Maps: &[]typeparser.Map{
					{
						Name: "AvailablePlatformRegions",
						Key: typeparser.TypeDescriptor{
							Type: "string",
						},
						Val: typeparser.TypeDescriptor{
							IncludeInOutput: true,
							IsPointer:       true,
							Type:            "string",
						},
					},
					{
						Name: "EcommerceStoreAvailablePlatformRegions",
						Key: typeparser.TypeDescriptor{
							Type: "EcommerceStorePlatformRegions",
						},
						Val: typeparser.TypeDescriptor{
							IsMap:           true,
							IncludeInOutput: true,
							IsPointer:       true,
							Type:            "EcommerceStorePlatformRegionsMap1",
						},
					},
				},
				Enums: &[]typeparser.Enum{
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
