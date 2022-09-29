package structparser_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	fieldparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/field_parser"
	structparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/struct_parser"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type TestStruct struct {
	TaggedField   string `json:"taggedField" graphql:"description=This is a tagged field, decorators=[+doc(description: \"This field is tagged.\"), +requireAuthRole(role: \"admin\"))]"`
	UnTaggedField string `json:"unTaggedField"`
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
				Fields: &[]*fieldparser.Field{
					{
						Name:            "taggedField",
						Type:            "string",
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
						IncludeInOutput: true,
						ParsedTag:       nil,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := structparser.ParseStruct(reflect.TypeOf(tt.s))
			assert.Equal(t, tt.want, got)
		})
	}
}
