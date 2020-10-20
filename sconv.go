package utils

import (
	"errors"
	"reflect"
	"strings"
)

// Struct2Map 结构体转化为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Name
		}
		data[fieldName] = v.Field(i).Interface()
	}
	return data
}

// Contains s中是否包含sub
func Contains(s interface{}, sub interface{}) (bool, error) {
	targetValue := reflect.ValueOf(s)
	switch reflect.TypeOf(s).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == sub {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(sub)).IsValid() {
			return true, nil
		}
		return false, nil
	case reflect.String:
		return strings.Contains(s.(string), sub.(string)), nil
	}

	return false, errors.New("not support type match")
}
