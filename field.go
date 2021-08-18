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



//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法

func checkValidFieldNullEr(v reflect.Value) error {
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

func checkValidFieldSqlValueEr(v reflect.Value) error {
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

func checkValidFieldNullErSqlValuer(v reflect.Value) error {
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
