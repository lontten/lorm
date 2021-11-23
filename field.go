package lorm

import (
	"database/sql/driver"
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"reflect"
)

//只校验struct 的 field 是否合法

// v0.5
//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法
func checkFieldNullEr(v reflect.Value) error {
	isPtr, base := basePtrValue(v)
	if !isPtr {
		//必须 nuller struct   or slice
		is, base := baseStructValue(base)
		if is {
			_, ok := base.Interface().(types.NullEr)
			if !ok {
				return errors.New("struct field " + base.String() + " need imp core NullEr ")
			}
			return nil
		}

		is, _ = baseSliceDeepValue(base)
		if is {
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

	is, _ = baseSliceDeepValue(base)
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

	is, _ = baseSliceDeepValue(base)
	if is {
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

		is, _= baseSliceDeepValue(base)
		if is {
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

	is, _= baseSliceDeepValue(base)
	if is {
		return nil
	}

	return errors.New("struct field " + base.String() + " need  NullEr sql valuer")
}

// v0.5
func getFieldInter(v reflect.Value) interface{} {

	isPtr, base := basePtrValue(v)
	if !isPtr {
		//必须 nuller struct
		is, base := baseStructValue(base)
		if is {
			vv := base.Interface().(types.NullEr)
			if !vv.IsNull() {
				value, _ := base.Interface().(driver.Valuer).Value()
				return value
			}
		}

		_, has, _ := baseSliceDeepValue(base)
		if has {
			return base.Interface()
		}

	}

	code, base := baseStructBaseValue(base)
	if code > 0 {
		return base.Interface()
	}

	_, has, _ := baseSliceDeepValue(base)
	if has {
		return base.Interface()
	}

	return nil
}
