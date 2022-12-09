package jsontagparser

import "strings"

type TargetType string

const (
	String  TargetType = "string"
	Int     TargetType = "int"
	Float   TargetType = "float"
	Boolean TargetType = "bool"
)

type JSONTag struct {
	Name       string
	Private    bool
	OmitEmpty  bool
	TargetType TargetType
}

func Parse(tag string) *JSONTag {
	var result JSONTag

	if tag == "" {
		return nil
	}

	parts := strings.Split(tag, ",")

	if len(parts) == 0 {
		return nil
	}

	// If the tag specifies that it should be omitted in the output
	// then we should return nothing
	if parts[0] == "-" {
		result.Private = true
	} else {
		// Otherwise...
		// Get the name of the field and pop it off the list.
		result.Name, parts = parts[0], parts[1:]
	}

	// If there are no more parts then we can return the result.
	if len(parts) == 0 {
		result.TargetType = String

		return &result
	}

	// For each part in the tag, work out what it is and set the
	// appropriate field. Most of the code here is to handle the
	// different types that can be specified.
	for _, rawPart := range parts {
		part := strings.TrimSpace(rawPart)

		switch part {
		case "omitempty":
			result.OmitEmpty = true
		case "string":
			result.TargetType = String
		case "int":
			result.TargetType = Int
		case "float":
			result.TargetType = Float
		case "bool":
			result.TargetType = Boolean
		}
	}

	if result.TargetType == "" {
		result.TargetType = String
	}

	return &result
}
