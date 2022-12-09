package typeparser

import (
	"fmt"
	"reflect"

	tagparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/graphql-tag-parser"
	jsontagparser "github.com/warpspeed-cloud/graphql-schema-generator/internal/json-tag-parser"
)

const (
	unnamedMapTemplate    = "Map%d"
	unnamedStructTemplate = "Struct%d"
)

type TypeDescriptor struct {
	Name            *string
	Type            string
	IsSlice         bool
	IsPointer       bool
	IsStruct        bool
	IsMap           bool
	IsEnum          bool
	IncludeInOutput bool
	ParsedTag       *tagparser.Tag
}

type Struct struct {
	Name   string
	Fields *[]TypeDescriptor
}

type Map struct {
	Name string
	Key  TypeDescriptor
	Val  TypeDescriptor
}

type EnumKeyPairOptions struct {
	Key         string
	Value       interface{}
	Description *string
}

type Enum struct {
	Name   string
	Values []EnumKeyPairOptions
}

type TypeParser struct {
	Structs *[]Struct
	Maps    *[]Map
	Enums   *[]Enum

	// Keep a list of types that are pending being added to the schema.
	// This is used to prevent infinite recursion when a struct has a field that is a pointer to itself
	// or a slice of itself or when structs have circular references.
	pendingStructTypeNames *[]string
}

type AddStructOptions struct {
	Name *string
}

func NewTypeParser(_ *interface{}) *TypeParser {
	return &TypeParser{}
}

// A function that returns a boolean whether a struct exists by this name or is pending.
func (t TypeParser) structExistsAndIsntPending(name string) bool {
	if t.pendingStructTypeNames != nil {
		for _, pendingName := range *t.pendingStructTypeNames {
			if pendingName == name {
				return true
			}
		}
	}

	if t.Structs != nil {
		for _, s := range *t.Structs {
			if s.Name == name {
				return true
			}
		}
	}

	return false
}

func (t *TypeParser) internalAddMap(name string, m reflect.Type, depth int) *TypeParser {
	var mapValueTypeName string

	key := TypeDescriptor{}
	val := TypeDescriptor{}

	if m.Kind() == reflect.Ptr {
		m = m.Elem()
	}

	if name == "" {
		name = fmt.Sprintf(unnamedMapTemplate, depth)
	}

	if m.Kind() != reflect.Map {
		panic(fmt.Sprintf("AddMap must be called with a map type, '%s' is a '%s'", name, m.Kind().String()))
	}

	mapKeyType := m.Key()
	mapValueType := m.Elem()

	// First we check if the key and value are pointers and if so, we unroll them.
	if mapValueType.Kind() == reflect.Ptr {
		val.IsPointer = true
		mapValueType = mapValueType.Elem()
	}

	if mapKeyType.Kind() == reflect.Ptr {
		key.IsPointer = true
		mapKeyType = mapKeyType.Elem()
	}

	// Now we check if the value is a slice and if so, we unroll it.
	if mapValueType.Kind() == reflect.Slice {
		val.IsSlice = true
		mapValueType = mapValueType.Elem()
	}

	// If the elem is a map then we need to add that map too
	// At this point we know that the type is a map so we also
	// increase the depth counter and generate a new name for the map.
	if mapValueType.Kind() == reflect.Map {
		mapName := fmt.Sprintf(unnamedMapTemplate, depth+1)
		mapValueTypeName = fmt.Sprintf("%s%s", name, mapName)

		t.internalAddMap(fmt.Sprintf("%s%s", name, mapName), mapValueType, depth+1)
	} else {
		mapValueTypeName = mapValueType.Kind().String()
	}

	key.Type = mapKeyType.Kind().String()
	val.Type = mapValueTypeName

	if t.Maps == nil {
		t.Maps = &[]Map{}
	}

	*t.Maps = append(*t.Maps, Map{
		Name: name,
		Key:  key,
		Val:  val,
	})

	return t
}

// getStructField returns a TypeDescriptor for a struct field.

// internalAddStruct loops over each field in the struct and add it to the schema
// recursively. It will unroll pointers and slices to find the underlying type
// automatically.
func (t *TypeParser) internalAddStruct(m reflect.Type, depth int) {
	newStruct := Struct{}

	if m.Kind() == reflect.Ptr {
		m = m.Elem()
	}

	if m.Kind() != reflect.Struct {
		panic(fmt.Sprintf("AddStruct must be called with a struct type, '%s' is a '%s'", m.Name(), m.Kind().String()))
	}

	if m.Name() == "" {
		newStruct.Name = fmt.Sprintf(unnamedStructTemplate, depth)
	} else {
		newStruct.Name = m.Name()
	}

	// If the struct is already in the schema or is pending then we don't need to add it again.
	if t.structExistsAndIsntPending(newStruct.Name) {
		return
	}

	// Add the struct name to the pending list so that we don't add it again.
	*t.pendingStructTypeNames = append(*t.pendingStructTypeNames, newStruct.Name)

	// Create a new slice to hold the fields for this struct.
	var fields []TypeDescriptor

	// Loop over each field in the struct and add it to the schema.
	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		jsonTag := jsontagparser.Parse(field.Tag.Get("json"))

		var fieldName string
		if jsonTag == nil || jsonTag.Name == "" {
			fieldName = field.Name
		} else {
			fieldName = jsonTag.Name
		}

		graphqlTag := tagparser.ParseTag(m.Field(i).Tag.Get("graphql"), fieldName)
		newField := TypeDescriptor{
			Name:      &fieldName,
			ParsedTag: graphqlTag,
		}

		fieldType := field.Type

		// First we check if the field is a pointer and if so, we unroll it.
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			newField.IsPointer = true
		}

		// Now we check if the field is a slice and if so, we unroll it.
		if fieldType.Kind() == reflect.Slice {
			fieldType = fieldType.Elem()
			newField.IsSlice = true
		}

		// If the field is a struct then we need to add that struct too
		// At this point we know that the type is a struct so we also
		// increase the depth counter and generate a new name for the struct.
		if fieldType.Kind() == reflect.Struct {
			if fieldType.Name() == "" {
				newField.Type = fmt.Sprintf("%s%s", newStruct.Name, fmt.Sprintf(unnamedStructTemplate, depth+1))
			} else {
				newField.Type = fieldType.Name()
			}
			t.internalAddStruct(fieldType, depth+1)

			if !newField.IsSlice {
				newField.IsStruct = true
			}
		} else if fieldType.Kind() == reflect.Map {
			t.internalAddMap(fmt.Sprintf("%s%s", newStruct.Name, fieldType.Name()), fieldType, 0)

			newField.Type = fmt.Sprintf("%s%s", newStruct.Name, fieldType.Name())
			newField.IsMap = true
		} else {
			newField.Type = fieldType.Kind().String()
		}

		newField.IncludeInOutput = field.IsExported() && (jsonTag == nil || !jsonTag.Private)

		fields = append(fields, newField)
	}

	if t.Structs == nil {
		t.Structs = &[]Struct{}
	}

	*t.Structs = append(*t.Structs, Struct{
		Name:   newStruct.Name,
		Fields: &fields,
	})

	// Remove the struct name from the pending list.
	for i, pendingName := range *t.pendingStructTypeNames {
		if pendingName == newStruct.Name {
			*t.pendingStructTypeNames = append((*t.pendingStructTypeNames)[:i], (*t.pendingStructTypeNames)[i+1:]...)
		}
	}
}

// AddMap adds a map to the schema and recursively adds any discovered types
// to the schema. It will unroll pointers and slices to find the underlying type
// automatically.
//
// You must pass in a map type and a name for the map. If you do not supply a name;
// i.e. you pass in an empty string, then the name will be generated automatically as
// Map1, Map2, Map3, etc. and depth is calculated automatically.
//
// If you don't pass a map type in (say a struct, reflect.Type, etc.) then it will panic.
func (t *TypeParser) AddMap(name string, m interface{}) *TypeParser {
	mapType := reflect.TypeOf(m)

	return t.internalAddMap(name, mapType, 0)
}

// AddStruct adds a struct to the schema and recursively adds any discovered types
// to the schema. It will unroll pointers and slices to find the underlying type
// automatically.
//
// You must pass in a struct type and you can optionally pass an AddStructOptions struct
// to specify a name for the struct. If you do not supply a name; i.e. you pass in an empty
// string, then the name will be generated automatically as Struct1, Struct2, Struct3, etc and
// depth is calculated automatically.
//
// If you don't pass a struct type in (say a map, reflect.Type, etc.) then it will panic.
func (t *TypeParser) AddStruct(s interface{}, options *AddStructOptions) *TypeParser {
	if options == nil {
		options = &AddStructOptions{}
	}

	structType := reflect.TypeOf(s)

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("AddStruct must be called with a struct type, '%s' is a '%s'", *options.Name, structType.Kind().String()))
	}

	if t.pendingStructTypeNames == nil {
		t.pendingStructTypeNames = &[]string{}
	}

	t.internalAddStruct(structType, 0)

	// Check there are no pending structs left and if so, panic because that means
	// something went wrong and not all structs were added to the schema.
	// If it's empty then we're good and we can clear the pending structs list entirely.
	if len(*t.pendingStructTypeNames) > 0 {
		panic(fmt.Sprintf("There are still pending structs to be added to the schema: %s", *t.pendingStructTypeNames))
	} else {
		t.pendingStructTypeNames = nil
	}

	return t
}
