package utils

import "reflect"

// IsNil
// nil 值在传到 interface{}、any 时，无法通过 == nil,来判断是否为nil
// ptr nil 为 true
// 非ptr，零值为 false
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return false
	}
	if reflect.ValueOf(v).IsNil() {
		return true
	}
	return false
}
