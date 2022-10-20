package typeparser

import (
	"fmt"
	"reflect"

	"github.com/warpspeedboilerplate/graphql-schema-generator/internal/ptr"
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
	Structs *[]Struct
	Maps    *[]Map
	Enums   *[]Enum

	// Keep a list of types that are pending being added to the schema.
	// This is used to prevent infinite recursion when a struct has a field that is a pointer to itself
	// or a slice of itself or when structs have circular references.
}

func NewTypeParser(_ *interface{}) *TypeParser {
	return &TypeParser{}
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

func (t *TypeParser) internalAddMap(name string, m reflect.Type, depth *int) *TypeParser {
	if depth == nil {
		depth = ptr.Of(0)
	}
	// if this initial type is a pointer, dig deeper.
	if m.Kind() == reflect.Ptr {
		m = m.Elem()
	}

	if name == "" {
		name = fmt.Sprintf("Map%d", *depth)
	}

	if m.Kind() != reflect.Map {
		panic(fmt.Sprintf("AddMap must be called with a map type, '%s' is a '%s'", name, m.Kind().String()))
	}

	typeType := m.Elem()
	typeName := typeType.Name()

	// If the elem is a map then we need to add that map too
	// and increase the depth counter.
	// First check if it's a pointer.
	if typeType.Kind() == reflect.Ptr {
		typeType = typeType.Elem()
	}

	if typeType.Kind() == reflect.Map {
		mapName := fmt.Sprintf("Map%d", *depth+1)
		typeName = fmt.Sprintf("%s%s", name, mapName)
		// Send the interface of the type to the AddMap function
		t.internalAddMap(fmt.Sprintf("%s%s", name, mapName), typeType, ptr.Of(*depth+1))
	} else {
		typeName = typeType.Kind().String()
	}

	if t.Maps == nil {
		t.Maps = &[]Map{}
	}

	*t.Maps = append(*t.Maps, Map{
		Name: name,
		Key: TypeDescriptor{
			Type:      m.Key().Name(),
			IsPointer: m.Key().Kind() == reflect.Ptr,
		},
		Val: TypeDescriptor{
			Type:      typeName,
			IsPointer: m.Elem().Kind() == reflect.Ptr,
		},
	})

	return t
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

	return t.internalAddMap(name, mapType, nil)
}
