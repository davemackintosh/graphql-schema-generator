package structparser_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/cmd/v1/graphql-generator/internal/field_parser"
	structparser "github.com/warpspeedboilerplate/graphql-schema-generator/cmd/v1/graphql-generator/internal/struct_parser"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/cmd/v1/graphql-generator/internal/tag-parser"
)

type TestStruct struct {
	TaggedField   string `graphql:"taggedField, omitempty, description=This is a tagged field, decorators=[+doc(description: \"This field is tagged.\"), +requireAuthRole(role: \"admin\"))]"`
	UnTaggedField string
}

func TestParseStruct(t *testing.T) {
	tests := []struct {
		name    string
		s       interface{}
		want    *structparser.Struct
		wantErr bool
	}{
		{
			name: "TestParseStruct",
			s:    TestStruct{},
			want: &structparser.Struct{
				Name: "TestStruct",
				Fields: &map[string]*fieldparser.Field{
					"TaggedField": {
						Name: "TaggedField",
						Type: "string",
						ParsedTag: &tagparser.Tag{
							Name: "taggedField",
							Options: map[string]string{
								"omitempty":   "true",
								"description": "This is a tagged field",
								"decorators":  "[+doc(description: \"This field is tagged.\"),+requireAuthRole(role: \"admin\"))]",
								"name":        "taggedField",
							},
						},
					},
					"UnTaggedField": {
						Name:      "UnTaggedField",
						Type:      "string",
						ParsedTag: nil,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := structparser.ParseStruct(reflect.TypeOf(tt.s))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStruct() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
