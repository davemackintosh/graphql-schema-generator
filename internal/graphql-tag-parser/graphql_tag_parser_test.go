package tagparser_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	tagparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/graphql-tag-parser"
)

type TestStruct struct {
	TaggedField   string `json:"taggedField" graphql:"description=This is a tagged field, decorators=[+doc(description: \"This field is tagged.\"), +requireAuthRole(role: \"admin\"))]"`
	UnTaggedField string `json:"unTaggedField"`
}

func TestGetTagsFromStruct(t *testing.T) {
	t.Run("Test tag parser.", func(t *testing.T) {
		target := &TestStruct{}
		fields := reflect.TypeOf(target).Elem()

		field, exists := fields.FieldByName("TaggedField")
		assert.Equal(t, true, exists)
		got := tagparser.GetTagFromField(field)

		assert.Equal(t, tagparser.Tag{
			Options: map[string]string{
				"description": "This is a tagged field",
				"decorators":  "[+doc(description: \"This field is tagged.\"),+requireAuthRole(role: \"admin\"))]",
			},
		}, *got)
	})
}
