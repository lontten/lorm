package types

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

func ArrayOf(v ...any) Array {
	array := Array{}
	for _, i := range v {
		value := reflect.ValueOf(i)

		if value.Kind() == reflect.Struct {
			array.ints = append(array.ints, i.(driver.Valuer))
		}
		array.bases = append(array.bases, i)
	}
	return array
}

type Array struct {
	ints  []driver.Valuer
	bases []any
}

func (a Array) Value() (driver.Value, error) {
	var str = "{"
	for _, e := range a.ints {
		value, err := e.Value()
		if err != nil {
			return nil, err
		}
		str += value.(string) + ","
	}
	for _, e := range a.bases {
		switch e.(type) {
		case int:
			str += fmt.Sprintf("%v", e) + ","
		case int8:
			str += fmt.Sprintf("%v", e) + ","
		case int16:
			str += fmt.Sprintf("%v", e) + ","
		case int32:
			str += fmt.Sprintf("%v", e) + ","
		case float32:
			str += fmt.Sprintf("%v", e) + ","
		case float64:
			str += fmt.Sprintf("%v", e) + ","
		}
		str += fmt.Sprintf("\"%v\"", e) + ","
	}
	str = str[:len(str)-1]
	str += "}"
	return str, nil
}
