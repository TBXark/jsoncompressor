package jsoncompressor

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	ErrorUnmarshalInputNotStruct  = fmt.Errorf("input must be a struct or pointer to struct")
	ErrorUnmarshalNotEnoughData   = fmt.Errorf("not enough data for field")
	ErrorUnmarshalTargetNotNilPtr = fmt.Errorf("target must be a non-nil pointer")
)

func Unmarshal(data []byte, target interface{}) error {
	var compressed []interface{}
	if err := json.Unmarshal(data, &compressed); err != nil {
		return err
	}
	return decompress(compressed, target)
}

func decompress(data []interface{}, target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrorUnmarshalTargetNotNilPtr
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrorUnmarshalInputNotStruct
	}

	return decompressIntoStruct(data, val)
}

func decompressIntoStruct(data []interface{}, val reflect.Value) error {
	typ := val.Type()
	dataIndex := 0

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		if dataIndex >= len(data) {
			return ErrorUnmarshalNotEnoughData
		}
		if err := decompressValue(data[dataIndex], field); err != nil {
			return err
		}
		dataIndex++
	}
	return nil
}

func decompressValue(data interface{}, field reflect.Value) error {
	if data == nil {
		return nil
	}
	switch field.Kind() {
	case reflect.Struct:
		dataSlice, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("expected array for struct field")
		}
		return decompressIntoStruct(dataSlice, field)
	case reflect.Slice:
		dataSlice, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("expected array for slice field")
		}
		slice := reflect.MakeSlice(field.Type(), len(dataSlice), len(dataSlice))
		for i := 0; i < len(dataSlice); i++ {
			if err := decompressValue(dataSlice[i], slice.Index(i)); err != nil {
				return err
			}
		}
		field.Set(slice)
	default:
		v := reflect.ValueOf(data)
		if v.Type().ConvertibleTo(field.Type()) {
			field.Set(v.Convert(field.Type()))
		} else {
			return fmt.Errorf("cannot convert %v to %v", v.Type(), field.Type())
		}
	}
	return nil
}
