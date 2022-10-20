package typeparser_test

import (
	"reflect"
	"testing"

	typeparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/type-parser"
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
				Maps: []typeparser.Map{
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
				Maps: []typeparser.Map{
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
				Maps: []typeparser.Map{
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
				Maps: []typeparser.Map{
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
				Maps: []typeparser.Map{
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
			if !reflect.DeepEqual(test.actual, test.expected) {
				t.Errorf("got %v, want %v", test.actual, test.expected)
			}
		})
	}
}
