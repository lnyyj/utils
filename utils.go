package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

// MD5 以16进制字符串输出
func MD5(v []byte) string {
	return fmt.Sprintf("%x", md5.Sum(v))
}

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

// Float2Int 浮点->整型
func Float2Int(v float64, rate int64) int64 {
	val := decimal.NewFromFloat(v)
	if rate > 0 {
		return val.Mul(decimal.NewFromInt(rate)).IntPart()
	} else if rate < 0 {
		return val.Div(decimal.NewFromInt(-rate)).IntPart()
	}
	return val.IntPart()
}

// Int2Float 整型->浮点
func Int2Float(v int64, rate int64) (r float64) {
	val := decimal.NewFromInt(v)
	if rate > 0 {
		r, _ = val.Mul(decimal.NewFromInt(rate)).Float64()
	} else if rate < 0 {
		r, _ = val.Div(decimal.NewFromInt(-rate)).Float64()
	} else {
		r, _ = val.Float64()
	}
	return
}

// Mul 乘法
func Mul(a, b float64) float64 {
	r, _ := decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(b)).Float64()
	return r
}

// Add 加法
func Add(a, b float64) float64 {
	r, _ := decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)).Float64()
	return r
}

// Sub 减法
func Sub(a, b float64) float64 {
	r, _ := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).Float64()
	return r
}
