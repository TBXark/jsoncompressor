package jsoncompressor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}
	var jsonData interface{}
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

func decompressIntoStruct(data []interface{}, val reflect.Value) (interface{}, error) {
	typ := val.Type()
	type fld struct {
		name  string
		index int
	}
	fields := make([]fld, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		name, ok := getJsonKey(&fieldType)
		if ok {
			fields = append(fields, fld{name: name, index: i})
		}
	}

	if len(data) != len(fields) {
		return nil, fmt.Errorf("field count mismatch: have %d values, want %d", len(data), len(fields))
	}

	jsonMap := make(map[string]interface{}, len(fields))
	for i, f := range fields {
		field := val.Field(f.index)
		value, err := decompressValue(data[i], field)
		if err != nil {
			return nil, err
		}
		jsonMap[f.name] = value
	}
	return jsonMap, nil
}

func decompressValue(data interface{}, field reflect.Value) (interface{}, error) {
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
		dataSlice, ok := data.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected array for struct field")
		}
		return decompressIntoStruct(dataSlice, field)
	case reflect.Slice, reflect.Array:
		dataSlice, ok := data.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected array for slice field")
		}
		slice := reflect.MakeSlice(field.Type(), len(dataSlice), len(dataSlice))
		jsonSlice := make([]interface{}, len(dataSlice))
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

func getJsonKey(field *reflect.StructField) (string, bool) {
	if field.PkgPath != "" { // unexported
		return "", false
	}
	tag := field.Tag.Get("json")
	if tag == "-" || tag == "" {
		return "", false
	}
	name := strings.Split(tag, ",")[0]
	if name == "" {
		name = field.Name
	}
	return name, true
}
