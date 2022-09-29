package fieldparser_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/field_parser"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type TestStruct struct {
	TaggedField   string `json:"taggedField" graphql:"description=This is a tagged field, decorators=[+doc(description: \"This field is tagged.\"), +requireAuthRole(role: \"admin\"))]"`
	UnTaggedField string `json:"unTaggedField"`
	NormField     string
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
					Name:            "taggedField",
					Type:            "string",
					IsArray:         false,
					IsPointer:       false,
					IncludeInOutput: true,
					ParsedTag: &tagparser.Tag{
						Options: map[string]string{
							"description": "This is a tagged field",
							"decorators":  "[+doc(description: \"This field is tagged.\"),+requireAuthRole(role: \"admin\"))]",
						},
					},
				},
				{
					Name:            "unTaggedField",
					Type:            "string",
					IsArray:         false,
					IsPointer:       false,
					IncludeInOutput: true,
					ParsedTag:       nil,
				},
				{
					Name:            "NormField",
					Type:            "string",
					IsArray:         false,
					IsPointer:       false,
					IncludeInOutput: true,
					ParsedTag:       nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := fieldparser.GetFieldsFromStruct(reflect.TypeOf(tt.s))
			assert.Equal(t, tt.want, got)
		})
	}
}
