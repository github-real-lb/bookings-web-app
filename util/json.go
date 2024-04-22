package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// CopyData uses json package to copy data from src to target using the json package.
// target needs to be a pointer to the type.
func CopyDataUsingJSON(src any, target any) error {
	// marshal src to json data
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, target)
	return err
}

// StructToMapUsingJSON converts a struct to a map[string]any.
// map keys are named using the JSON tags
func StructToMapUsingJSON(obj any) (map[string]any, error) {
	result := make(map[string]any)

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Use the json tag as the key in the map, or default to the field name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name // Default to field name if there is no json tag
		} else {
			// Handle cases where the json tag includes options, like omitempty
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Add to the map
		result[jsonTag] = v.Field(i).Interface()
	}

	return result, nil
}
