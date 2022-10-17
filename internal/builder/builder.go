package builder

import (
	"fmt"
	"reflect"
	"strings"

	tagparser "github.com/warpspeedboilerplate/graphql-schema-generator/internal/tag-parser"
)

type Field struct {
	Name            string
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
	Fields *[]*Field
}

type Map struct {
	Name    string
	KeyType string
	Field   Field
}

type EnumKeyPairOptions struct {
	Key         string
	Value       interface{}
	Description *string
}

type Enum struct {
	Name   string
	Values []*EnumKeyPairOptions
}

type GraphQLSchemaBuilderWriter interface {
	WriteSchema(schema string)
}

type GraphQLSchemaBuilderOptions struct {
	// A callback that takes the type name and the generated schema and writes it to an ioWriter.
	Writer GraphQLSchemaBuilderWriter
}

type GraphQLSchemaBuilder struct {
	Structs []*Struct
	Maps    []*Map
	Enums   []*Enum
	Options *GraphQLSchemaBuilderOptions

	// Keep a list of types that are pending being added to the schema.
	// This is used to prevent infinite recursion when a struct has a field that is a pointer to itself
	// or a slice of itself or when structs have circular references.
	pendingStructTypeNames *[]string
}

func getGoBuiltInTypeNames() []string {
	return []string{
		"bool",
		"byte",
		"complex64",
		"complex128",
		"error",
		"float32",
		"float64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"rune",
		"string",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
	}
}

func NewGraphQLSchemaBuilder(options *GraphQLSchemaBuilderOptions) *GraphQLSchemaBuilder {
	return &GraphQLSchemaBuilder{
		Options: options,
	}
}

func (b *GraphQLSchemaBuilder) AddEnum(enum Enum) *GraphQLSchemaBuilder {
	b.Enums = append(b.Enums, &enum)

	return b
}

type AddStructOptions struct {
	Name *string
}

// A function that returns a boolean whether a struct exists by this name or is pending.
func (b GraphQLSchemaBuilder) structExistsAndIsntPending(name string) bool {
	if b.pendingStructTypeNames != nil {
		for _, pendingName := range *b.pendingStructTypeNames {
			if pendingName == name {
				return true
			}
		}
	}

	for _, s := range b.Structs {
		if s.Name == name {
			return true
		}
	}

	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func (b *GraphQLSchemaBuilder) AddMap(name string, t interface{}) *GraphQLSchemaBuilder {
	// Get the type the map is to.
	baseType := reflect.TypeOf(t)
	if baseType.Kind() == reflect.Ptr {
		baseType = baseType.Elem()
	}

	mapType := baseType.Elem()

	// if the map type is to a pointer then dig deeper.
	if mapType.Kind() == reflect.Ptr {
		mapType = mapType.Elem()
	}

	// If it's a slice, then dig deeper.
	if mapType.Kind() == reflect.Slice {
		mapType = mapType.Elem()
	}

	// If it's a map, then dig deeper.
	if mapType.Kind() == reflect.Map {
		mapType = mapType.Elem()
	}

	// If the map type is to a struct then add it to the list of structs.
	if mapType.Kind() == reflect.Struct {
		if !b.structExistsAndIsntPending(mapType.Name()) {
			name := mapType.Name()
			b.AddStruct(reflect.Zero(mapType).Interface(), &AddStructOptions{
				Name: &name,
			})
		}
	}

	// Create a new struct object (which we'll add to the maps list because technically, they're the same.) by looping
	// through the fields of the incoming map type interface.
	keyType := baseType.Key().Name()

	// if the key type isn't built in, it's going to be an enum type string.
	if !stringInSlice(keyType, getGoBuiltInTypeNames()) {
		b.AddEnum(Enum{
			Name:   fmt.Sprintf("%s%s", name, keyType),
			Values: []*EnumKeyPairOptions{},
		})
	}

	s := &Map{
		Name:    name,
		KeyType: keyType,
		Field: Field{
			Name:            "",
			Type:            mapType.Name(),
			IsSlice:         baseType.Elem().Kind() == reflect.Slice,
			IsPointer:       baseType.Elem().Kind() == reflect.Ptr,
			IsStruct:        baseType.Elem().Kind() == reflect.Struct,
			IsMap:           baseType.Elem().Kind() == reflect.Map,
			IsEnum:          baseType.Elem().Kind() == reflect.String && !stringInSlice(baseType.Elem().Name(), getGoBuiltInTypeNames()),
			IncludeInOutput: true,
			// Maps don't have tags.
			ParsedTag: nil,
		},
	}

	b.Maps = append(b.Maps, s)

	return b
}

func (b *GraphQLSchemaBuilder) AddStruct(t interface{}, options *AddStructOptions) *GraphQLSchemaBuilder { //nolint: cyclop
	structType := reflect.ValueOf(t)
	structName := structType.Type().Name()

	if options != nil && options.Name != nil {
		structName = *options.Name
	}

	// If there is already a struct with this name, return.
	if b.structExistsAndIsntPending(structName) {
		return b
	}

	// Check if pending types has been initialized.
	if b.pendingStructTypeNames == nil {
		b.pendingStructTypeNames = &[]string{}
	}

	// Add this struct to the pending list.
	*b.pendingStructTypeNames = append(*b.pendingStructTypeNames, structName)

	// Loop over the struct's fields and add them to the list of fields.
	var fields []*Field

	// if the struct is a pointer, then dig deeper.
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Type().Field(i)
		fieldType := field.Type
		fieldTypeName := fieldType.String()

		// If the field is a pointer, get the type it points to.
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			fieldTypeName = fieldType.String()
		}

		// If the field is a slice of a struct, we need to get the struct and add it.
		if structType.Field(i).Kind() == reflect.Slice && structType.Field(i).Type().Elem().Kind() == reflect.Struct {
			target := reflect.New(structType.Field(i).Type().Elem()).Elem()
			embeddedName := fmt.Sprintf("%s_%s", structName, field.Name)

			// If this is an embedded struct with no name, we need to get the name of the struct it's embedded in
			// and use that as the name of the struct concatenated with the name of the field.
			targetName := target.Type().Name()
			if targetName == "" {
				targetName = embeddedName
				fieldTypeName = embeddedName
			}

			// If the struct doesn't already exist, add it.
			if !b.structExistsAndIsntPending(targetName) {
				b.AddStruct(target.Interface(), &AddStructOptions{
					Name: &targetName,
				})
			}
		} else if structType.Field(i).Kind() == reflect.Struct {
			target := structType.Field(i)
			embeddedName := fmt.Sprintf("%s_%s", structName, field.Name)

			// If this is an embedded struct with no name, we need to get the name of the struct it's embedded in
			// and use that as the name of the struct concatenated with the name of the field.
			targetName := target.Type().Name()
			if targetName == "" {
				targetName = embeddedName
				fieldTypeName = embeddedName
			}

			// If the struct doesn't already exist, add it.
			if !b.structExistsAndIsntPending(targetName) {
				b.AddStruct(target.Interface(), &AddStructOptions{
					Name: &targetName,
				})
			}
		} else if fieldType.Kind() == reflect.Map {
			// If the field is a map add it to the list of maps.
			b.AddMap(fmt.Sprintf("%s%s", structName, field.Name), structType.Field(i).Interface())
		}

		fieldName := field.Name

		// if the fieldTypeName starts with map then we need to get the type of the map.
		if strings.HasPrefix(fieldTypeName, "map") {
			fieldTypeName = fmt.Sprintf("%s%s", structName, field.Name)
		}

		// If the fieldTypeName has a period in it, it's a package name.Type and we only want the type name.
		if strings.Contains(fieldTypeName, ".") {
			fieldTypeNameParts := strings.Split(fieldTypeName, ".")
			fieldTypeName = fieldTypeNameParts[len(fieldTypeNameParts)-1]
		}

		// Get the field name from the json tag and fallback to the field name.
		jsonTag := field.Tag.Get("json")
		jsonTagParts := strings.Split(jsonTag, ",")
		jsonTagName := jsonTagParts[0]

		if jsonTagName != "" && jsonTagName != "-" {
			fieldName = jsonTagName
		}

		fields = append(fields, &Field{
			Name:            fieldName,
			Type:            fieldTypeName,
			IsPointer:       field.Type.Kind() == reflect.Ptr,
			IsSlice:         fieldType.Kind() == reflect.Slice,
			IsStruct:        fieldType.Kind() == reflect.Struct,
			IsMap:           fieldType.Kind() == reflect.Map,
			ParsedTag:       tagparser.ParseTag(field.Tag.Get("graphql"), field.Name),
			IncludeInOutput: field.Tag.Get("graphql") != "-" && field.Tag.Get("json") != "-" && field.IsExported(),
		})
	}

	b.Structs = append(b.Structs, &Struct{
		Name:   structName,
		Fields: &fields,
	})

	// Remove the struct name from the pending list.
	for i, pendingName := range *b.pendingStructTypeNames {
		if pendingName == structName {
			*b.pendingStructTypeNames = append((*b.pendingStructTypeNames)[:i], (*b.pendingStructTypeNames)[i+1:]...)
		}

		if len(*b.pendingStructTypeNames) == 0 {
			b.pendingStructTypeNames = nil
		}
	}

	return b
}

func (b *GraphQLSchemaBuilder) Build() string {
	return ""
}
