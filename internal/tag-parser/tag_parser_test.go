package tagparser_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type TestStruct struct {
	TaggedField   string `graphql:"taggedField, omitempty, description=This is a tagged field, decorators=[+doc(description: \"This field is tagged.\"), +requireAuthRole(role: \"admin\"))]"`
	UnTaggedField string
}

func TestGetTagsFromStruct(t *testing.T) {
	tests := []struct {
		name    string
		s       interface{}
		want    map[string]tagparser.Tag
		wantErr bool
	}{
		{
			name: "TestGetTagsFromStruct",
			s:    TestStruct{},
			want: map[string]tagparser.Tag{
				"taggedField": {
					Name: "taggedField",
					Options: &map[string]string{
						"omitempty":   "true",
						"description": "This is a tagged field",
						"decorators":  "[+doc(description: \"This field is tagged.\"),+requireAuthRole(role: \"admin\"))]",
						"name":        "taggedField",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := tagparser.GetTagsFromStruct(reflect.TypeOf(tt.s))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsFromStruct() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
