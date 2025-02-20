package jsoncompressor

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	ErrMarshalInputNotStruct = fmt.Errorf("input must be a struct or pointer to struct")
)

func Marshal(v interface{}) ([]byte, error) {
	compressed, err := compress(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(compressed)
}

func compress(v interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, ErrMarshalInputNotStruct
	}
	return compressStruct(val)
}

func compressStruct(val reflect.Value) ([]interface{}, error) {
	typ := val.Type()
	result := make([]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {
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
	case reflect.Ptr:
		if v.IsNil() {
			return nil, nil
		}
		return compressValue(v.Elem())
	default:
		return v.Interface(), nil
	}
}
