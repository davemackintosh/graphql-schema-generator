package jsontagparser_test

import (
	"testing"

	jsontagparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/json-tag-parser"

	"github.com/stretchr/testify/assert"
)

func Test_JSONTagParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *jsontagparser.JSONTag
	}{
		{
			name:     "empty",
			input:    "",
			expected: nil,
		},
		{
			name:  "one",
			input: "one",
			expected: &jsontagparser.JSONTag{
				Name:       "one",
				TargetType: jsontagparser.String,
			},
		},
		{
			name:     "none",
			input:    "-",
			expected: nil,
		},
		{
			name:  "one, omitempty",
			input: "one,omitempty",
			expected: &jsontagparser.JSONTag{
				Name:       "one",
				OmitEmpty:  true,
				TargetType: jsontagparser.String,
			},
		},
		{
			name:  "one, omitempty, string",
			input: "one,omitempty,string",
			expected: &jsontagparser.JSONTag{
				Name:       "one",
				OmitEmpty:  true,
				TargetType: jsontagparser.String,
			},
		},
		{
			name:  "one, omitempty, int",
			input: "one,omitempty,int",
			expected: &jsontagparser.JSONTag{
				Name:       "one",
				OmitEmpty:  true,
				TargetType: jsontagparser.Int,
			},
		},
		{
			name:  "one, omitempty, int64",
			input: "one,omitempty,int64",
			expected: &jsontagparser.JSONTag{
				Name:       "one",
				OmitEmpty:  true,
				TargetType: jsontagparser.String,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := jsontagparser.Parse(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
