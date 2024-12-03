package utils

import "reflect"

func ToSliceValue(v reflect.Value) []reflect.Value {
	l := v.Len()
	arr := make([]reflect.Value, l)
	for i := 0; i < l; i++ {
		arr[i] = v.Index(i)
	}
	return arr
}

func ToSlice(v reflect.Value) []any {
	l := v.Len()
	arr := make([]any, l)
	for i := 0; i < l; i++ {
		arr[i] = v.Index(i).Interface()
	}
	return arr
}

func Contains(list []string, s string) bool {
	for _, a := range list {
		if a == s {
			return true
		}
	}
	return false
}

func Find(list []string, s string) int {
	for i, a := range list {
		if a == s {
			return i
		}
	}
	return -1
}
