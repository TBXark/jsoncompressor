package jsoncompressor

import (
	"encoding/json"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	compressed, err := compressValue(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	return json.Marshal(compressed)
}

func compressStruct(val reflect.Value) ([]interface{}, error) {
	typ := val.Type()
	numFields := val.NumField()
	result := make([]interface{}, 0, numFields)

	for i := 0; i < numFields; i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		value, err := compressValue(field)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}

func compressValue(v reflect.Value) (interface{}, error) {
	switch v.Kind() {
	case reflect.Struct:
		return compressStruct(v)
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			val, err := compressValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case reflect.Pointer:
		if v.IsNil() {
			return nil, nil
		}
		return compressValue(v.Elem())
	default:
		return v.Interface(), nil
	}
}
