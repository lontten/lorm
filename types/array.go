package types

import (
	"database/sql/driver"
	"reflect"
)

func AsArray(v interface{}) Array {
	array := Array{}
	for _, i := range v.([]interface{}) {
		array = append(array, i.(driver.Valuer))
	}
	return array
}

type Array []driver.Valuer

func (a Array) Value() (driver.Value, error) {
	v := reflect.ValueOf(a)
	v = v.Elem()
	var str = "{"
	for _, e := range a {
		value, err := e.Value()
		if err != nil {
			return nil, err
		}
		str += value.(string) + ","
	}
	str += "}"
	return str, nil
}
