package fieldparser_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/field_parser"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type TestStruct struct {
	TaggedField   string `graphql:"taggedField, omitempty, description=This is a tagged field, decorators=[+doc(description: \"This field is tagged.\"), +requireAuthRole(role: \"admin\"))]"`
	UnTaggedField string
}

func TestGetFieldsFromStruct(t *testing.T) {
	tests := []struct {
		name    string
		s       interface{}
		want    *[]*fieldparser.Field
		wantErr bool
	}{
		{
			name: "TestGetFieldsFromStruct",
			s:    TestStruct{},
			want: &[]*fieldparser.Field{
				{
					Name:      "TaggedField",
					Type:      "string",
					IsArray:   false,
					IsPointer: false,
					ParsedTag: &tagparser.Tag{
						Name: "taggedField",
						Options: &map[string]string{
							"omitempty":   "true",
							"description": "This is a tagged field",
							"decorators":  "[+doc(description: \"This field is tagged.\"),+requireAuthRole(role: \"admin\"))]",
							"name":        "taggedField",
						},
					},
				},
				{
					Name:      "UnTaggedField",
					Type:      "string",
					IsArray:   false,
					IsPointer: false,
					ParsedTag: nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := fieldparser.GetFieldsFromStruct(reflect.TypeOf(tt.s))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldsFromStruct() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
