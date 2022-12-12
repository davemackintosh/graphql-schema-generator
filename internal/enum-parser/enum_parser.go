package enumparser

import "reflect"

func Parse(input any) *string {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	if input == nil {
		return nil
	}

	if _, ok := input.(string); !ok {
		return nil
	}

	name := inputType.Name()

	return &name
}
