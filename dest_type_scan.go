package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)


//struct 只检查 struct是否合格，不检查 filed
func checkScanTypeLn(t reflect.Type) (reflect.Type, error) {
	is, base := basePtrType(t)
	if !is {
		return t, errors.New("need a ptr")
	}
	code, base := baseStructBaseType(t)
	if code<0 {
		return t, errors.New("need a ptr struct or base type")
	}
	return base, nil
}
// slice 只检查 struct是否合格，不检查 filed
func checkScanType(t reflect.Type) (reflect.Type, error) {
	_, base := basePtrType(t)
	is, base := baseSliceType(base)
	if !is {
		return t, errors.New("need a slice type")
	}

	baseType, _ := baseStructBaseSliceType(base)

	if baseType < 0 {
		return t, errors.New("need a slice struct or base type")
	}
	return base, nil
}


