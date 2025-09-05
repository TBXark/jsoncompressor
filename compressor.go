package jsoncompressor

import (
	"encoding/json"
	"reflect"
)

func Marshal(v any) ([]byte, error) {
	compressed, err := compressValue(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	return json.Marshal(compressed)
}

func compressStruct(val reflect.Value) ([]any, error) {
	meta := getStructMeta(val.Type())
	if meta == nil {
		// no fields or not a struct (defensive)
		return []any{}, nil
	}
	result := make([]any, 0, len(meta.fields))
	for _, f := range meta.fields {
		v, err := compressValue(val.Field(f.index))
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func compressValue(v reflect.Value) (any, error) {
	switch v.Kind() {
	case reflect.Struct:
		return compressStruct(v)
	case reflect.Slice, reflect.Array:
		result := make([]any, v.Len())
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
