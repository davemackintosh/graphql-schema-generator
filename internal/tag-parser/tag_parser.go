package tagparser

import (
	"reflect"
)

// Tag represents a GraphQL tag.
type Tag struct {
	Name    string
	Options map[string]string
}

// Get all the tags from a struct and deduce which are GraphQL related.
func GetTagsFromStruct(structType reflect.Type) (map[string]Tag, error) {
	tags := make(map[string]Tag)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag, err := GetTagFromField(field)
		if err != nil {
			return nil, err
		}

		if tag != nil {
			tags[tag.Name] = *tag
		}
	}

	return tags, nil
}

// Get the tag from a Field.
func GetTagFromField(field reflect.StructField) (*Tag, error) {
	tag := field.Tag.Get("graphql")

	if tag == "" {
		return nil, nil
	}

	return ParseTag(tag, field.Name)
}

// Parse a tag.
func ParseTag(tag string, fieldName string) (*Tag, error) {
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

			if _, ok := tagOptions["name"]; !ok {
				tagOptions["name"] = currentKey
			} else {
				tagOptions[currentKey] = currentValue
			}

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

	tagOptions[currentKey] = currentValue

	// If the tag is empty, return nil
	if len(tagOptions) == 0 {
		return &Tag{Name: fieldName}, nil
	}

	// If the tag doesn't have a name, use the field field
	if tagOptions["name"] == "" {
		tagOptions["name"] = fieldName
	}

	return &Tag{
		Name:    tagOptions["name"],
		Options: tagOptions,
	}, nil
}
