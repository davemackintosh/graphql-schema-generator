package typeparser

import (
	"fmt"
	"reflect"

	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
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
	Fields *[]*TypeDescriptor
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
	Structs []Struct
	Maps    []Map
	Enums   []Enum

	// Keep a list of types that are pending being added to the schema.
	// This is used to prevent infinite recursion when a struct has a field that is a pointer to itself
	// or a slice of itself or when structs have circular references.
	pendingStructTypeNames []string

	// Keep a counter for the depth of the discovered maps and structs.
	// This is useful so we can name nested maps. We can also "lint" this
	// when we finish parsing and tell the user that perhaps the structure
	// is too deep and should be made flatter or more concrete.
	mapDepthCounter int
}

func NewTypeParser(_ *interface{}) *TypeParser {
	return &TypeParser{
		Structs:                []Struct{},
		Maps:                   []Map{},
		Enums:                  []Enum{},
		pendingStructTypeNames: []string{},
	}
}

// A function that returns a boolean whether a struct exists by this name or is pending.
/*func (t TypeParser) structExistsAndIsntPending(name string) bool {
	if t.pendingStructTypeNames != nil {
		for _, pendingName := range t.pendingStructTypeNames {
			if pendingName == name {
				return true
			}
		}
	}

	for _, s := range t.Structs {
		if s.Name == name {
			return true
		}
	}

	return false
}*/

func (t *TypeParser) AddMap(name string, m interface{}) *TypeParser {
	mapType := reflect.TypeOf(m)

	// if this initial type is a pointer, dig deeper.
	if mapType.Kind() == reflect.Ptr {
		mapType = mapType.Elem()
	}

	if name == "" {
		name = fmt.Sprintf("Map%d", t.mapDepthCounter)
		t.mapDepthCounter++
	}

	if mapType.Kind() != reflect.Map {
		panic(fmt.Sprintf("AddMap must be called with a map type, '%s' is a '%s'", name, mapType.Kind().String()))
	}

	typeType := mapType.Elem()

	// If the elem is a map then we need to add that map too
	// and increase the depth counter.
	// First check if it's a pointer.
	if typeType.Kind() == reflect.Ptr {
		typeType = typeType.Elem()
	}

	if typeType.Kind() == reflect.Map {
		mapName := fmt.Sprintf("Map%d", t.mapDepthCounter+1)
		t.mapDepthCounter++
		t.AddMap(fmt.Sprintf("%s%s", name, mapName), typeType)
	}

	t.Maps = append(t.Maps, Map{
		Name: name,
		Key: TypeDescriptor{
			Type: mapType.Key().Name(),
		},
		Val: TypeDescriptor{
			Type: mapType.Elem().Name(),
		},
	})

	return t
}
