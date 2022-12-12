package enumparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	enumparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/enum-parser"
	"github.com/warpspeed-cloud/graphql-schema-generator/internal/ptr"
)

type (
	MyStringEnum string
	MyIntEnum    int
	MyFloatEnum  float64
)

type MyPointerStruct struct {
	StringEnum *MyStringEnum
	IntEnum    *MyIntEnum
	FloatEnum  *MyFloatEnum
}

type MyStruct struct {
	StringEnum MyStringEnum
	IntEnum    MyIntEnum
	FloatEnum  MyFloatEnum
}

func TestEnumParser(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected *string
	}{
		{
			name:     "string",
			input:    MyStringEnum(""),
			expected: ptr.Of("MyStringEnum"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run("Test tag parser.", func(t *testing.T) {
			assert.Equal(t, *enumparser.Parse(test.input), *test.expected)
		})
	}
}
