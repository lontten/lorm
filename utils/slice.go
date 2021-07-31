package utils

import "reflect"

func ToSlice(v reflect.Value) []interface{} {
	l := v.Len()
	arr := make([]interface{}, l)
	for i := 0; i < l; i++ {
		arr[i] = v.Index(i).Interface()
	}
	return arr
}
