package utils

import "reflect"

func ToSliceValue(v reflect.Value) []reflect.Value {
	l := v.Len()
	arr := make([]reflect.Value, l)
	for i := 0; i < l; i++ {
		arr[i] =v.Index(i)
	}
	return arr
}

func ToSlice(v reflect.Value) []interface{} {
	l := v.Len()
	arr := make([]interface{}, l)
	for i := 0; i < l; i++ {
		arr[i] = v.Index(i).Interface()
	}
	return arr
}
