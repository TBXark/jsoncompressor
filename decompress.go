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
	dataIndex := 0
	jsonMap := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag, ok := getJsonKey(&fieldType)
		if !ok {
			continue
		}
		if dataIndex >= len(data) {
			break
		}
		value, err := decompressValue(data[dataIndex], field)
		if err != nil {
			return nil, err
		}
		jsonMap[jsonTag] = value
		dataIndex++
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
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return "", false
	}
	tagParts := strings.Split(tag, ",")
	return tagParts[0], true
}
