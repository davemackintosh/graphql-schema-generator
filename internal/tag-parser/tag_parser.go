package tagparser

import (
	"reflect"
)

// Tag represents a GraphQL tag.
type Tag struct {
	Options map[string]string
}

// Get the tag from a Field.
func GetTagFromField(field reflect.StructField) *Tag {
	tag := field.Tag.Get("graphql")

	if tag == "" {
		return nil
	}

	return ParseTag(tag, field.Name)
}

// Parse a tag.
func ParseTag(tag string, fieldName string) *Tag { //nolint: cyclop
	// Loop over each character and assemble a map of tags,
	// we do this because decorators and comments can contain commas
	// and we don't want to split on those
	tagOptions := make(map[string]string)

	var (
		currentKey   string
		currentValue string
		inQuotes     bool
		escaped      bool
		inArray      bool
		prevChar     rune
	)

	for _, char := range tag {
		if string(char) == " " && string(prevChar) == "," {
			prevChar = char
			// Skip if we have a space after a comma.
			continue
		}

		if char == '[' {
			inArray = true
		}

		if char == ']' {
			inArray = false
		}

		if char == '\\' {
			escaped = true

			prevChar = char
			continue
		}

		if char == '"' && !escaped {
			inQuotes = !inQuotes
		}

		if char == ',' && !inQuotes && !inArray {
			if currentKey == "" {
				currentKey = currentValue
				currentValue = "true"
			}

			tagOptions[currentKey] = currentValue

			currentKey = ""
			currentValue = ""

			prevChar = char
			continue
		}

		if char == '=' {
			currentKey = currentValue
			currentValue = ""

			prevChar = char
			continue
		}

		currentValue += string(char)
		prevChar = char
	}

	if currentKey != "" && currentValue != "" {
		tagOptions[currentKey] = currentValue
	}

	if len(tagOptions) == 0 {
		return nil
	}

	return &Tag{
		Options: tagOptions,
	}
}
