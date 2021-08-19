package lorm

import (
	"database/sql/driver"
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"reflect"
)

func checkTargetDestField(v reflect.Value) error {

	return nil
}

func checkArgsDestField(v reflect.Value) error {
	return nil
}

func checkScanDestField(v reflect.Value) error {
	return nil
}


// slice single filed
func singleFieldBaseSliceValue(value reflect.Value) (is, has bool, base reflect.Value) {
	_, base = basePtrValue(value)
	is, has, base = baseSliceValue(base)
	if !is { //不是 slice
		return
	}
	if !has { //是slice 内容空
		return
	}

base:
	is, has, base = baseSliceValue(base)
	if is {
		if !has { //是slice 内容空
			return
		}
		goto base
	}

	is, err := checkDestSingleNullerSqlValuer(base)
	if err != nil {
		return false, false, base
	}
	if is {
		return true, true, base
	}
	return false, false, base
}

// v0.5
//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法
func checkFieldNullEr(v reflect.Value) error {
	isPtr, base := basePtrValue(v)
	if !isPtr {
		//必须 nuller struct
		is, base := baseStructValue(base)
		if is {
			_, ok := base.Interface().(types.NullEr)
			if !ok {
				return errors.New("struct field " + base.String() + " need imp core NullEr ")
			}
			return nil
		}
		return errors.New("struct field " + base.String() + " need ptr or NullEr")
	}

	is := baseBaseValue(base)
	if is {
		return nil
	}

	is, base = baseStructValue(base)
	if is {
		return nil
	}
	return errors.New("struct field " + base.String() + " need ptr or NullEr")
}

// v0.5
func checkFieldSqlValueEr(v reflect.Value) error {
	_, base := basePtrValue(v)
	is := baseBaseValue(base)
	if is {
		return nil
	}

	is, base = baseStructValue(base)
	if is {
		_, ok := base.Interface().(driver.Valuer)
		if !ok {
			return errors.New("struct field " + base.String() + " need imp sql Value")
		}
		return nil
	}
	return errors.New("struct field " + base.String() + " need ptr or NullEr")
}

// v0.5
func checkFieldNullErSqlValuer(v reflect.Value) error {
	isPtr, base := basePtrValue(v)
	if !isPtr {
		//必须 nuller struct
		is, base := baseStructValue(base)
		if is {
			_, ok := base.Interface().(driver.Valuer)
			if !ok {
				return errors.New("struct field " + base.String() + " need imp sql Value")
			}
			_, ok = base.Interface().(types.NullEr)
			if !ok {
				return errors.New("struct field " + base.String() + " need imp core NullEr ")
			}
			return nil
		}
		return errors.New("struct field " + base.String() + " need ptr or NullEr")
	}

	is := baseBaseValue(base)
	if is {
		return nil
	}

	is, base = baseStructValue(base)
	if is {
		_, ok := base.Interface().(driver.Valuer)
		if !ok {
			return errors.New("struct field " + base.String() + " need imp sql Value")
		}
		return nil
	}
	return errors.New("struct field " + base.String() + " need  NullEr sql valuer")
}
