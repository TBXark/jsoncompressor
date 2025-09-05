package jsoncompressor

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Unmarshal(data []byte, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}
	var jsonData any
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}
	jsonObj, err := decompressValue(jsonData, val)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, target)
}

func decompressIntoStruct(data []any, val reflect.Value) (any, error) {
	meta := getStructMeta(val.Type())
	if meta == nil {
		if len(data) != 0 {
			return nil, fmt.Errorf("field count mismatch: have %d values, want %d", len(data), 0)
		}
		return map[string]any{}, nil
	}
	if len(data) != len(meta.fields) {
		return nil, fmt.Errorf("field count mismatch: have %d values, want %d", len(data), len(meta.fields))
	}
	jsonMap := make(map[string]any, len(meta.fields))
	for i, f := range meta.fields {
		field := val.Field(f.index)
		value, err := decompressValue(data[i], field)
		if err != nil {
			return nil, err
		}
		jsonMap[f.jsonName] = value
	}
	return jsonMap, nil
}

func decompressValue(data any, field reflect.Value) (any, error) {
	if data == nil {
		return nil, nil
	}
	for field.Kind() == reflect.Pointer {
		if field.IsNil() {
			field = reflect.New(field.Type().Elem())
		} else {
			field = field.Elem()
		}
	}
	switch field.Kind() {
	case reflect.Struct:
		dataSlice, ok := data.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array for struct field")
		}
		return decompressIntoStruct(dataSlice, field)
	case reflect.Slice, reflect.Array:
		dataSlice, ok := data.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array for slice field")
		}
		slice := reflect.MakeSlice(field.Type(), len(dataSlice), len(dataSlice))
		jsonSlice := make([]any, len(dataSlice))
		for i := 0; i < len(dataSlice); i++ {
			item, err := decompressValue(dataSlice[i], slice.Index(i))
			if err != nil {
				return nil, err
			}
			jsonSlice[i] = item
		}
		return jsonSlice, nil
	default:
		return data, nil
	}
}
