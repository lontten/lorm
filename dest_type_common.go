package lorm

import (
	"database/sql/driver"
	"errors"
	"github.com/lontten/lorm/types"
	"reflect"
)

// v0.5
//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法
//Deprecated
func checkDestSingleSqlValuer(v reflect.Value) (bool, error) {

	is := baseBaseValue(v)
	if is {
		return true, nil
	}

	is = baseStructValue(v)
	if is {
		_, ok := v.Interface().(driver.Valuer)
		if !ok {
			return false, errors.New("struct field " + v.String() + " need imp sql Value")
		}
		return true, nil
	}

	return false, errors.New("need a struct or base type")
}

// v0.5
//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法
//Deprecated
func checkDestSingleNullerSqlValuer(v reflect.Value) (bool, error) {
	is, base := basePtrValue(v)
	if is && v.IsNil() { //数值无效，直接返回false，不再进行合法性检查
		return false, nil
	}

	if !is { //不是ptr，必须 Nuller
		is = baseStructValue(base)
		if is {
			_, ok := base.Interface().(driver.Valuer)
			if !ok {
				return false, errors.New("struct field " + base.String() + " need imp sql Value")
			}
			_, ok = base.Interface().(types.NullEr)
			if !ok {
				return false, errors.New("struct field " + base.String() + " need imp core NullEr ")
			}
			return true, nil
		}
		return false, errors.New("struct field " + base.String() + " need ptr or NullEr")
	}

	is = baseBaseValue(base)
	if is {
		return true, nil
	}

	is = baseStructValue(base)
	if is {
		_, ok := base.Interface().(driver.Valuer)
		if !ok {
			return false, errors.New("struct field " + base.String() + " need imp sql Value")
		}
		return true, nil
	}

	return false, errors.New("need a struct or base type")
}

// todo
// v0.5
// byid single
// slice filed 不进行 内部检查，只要是 slice就行
//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法
//Deprecated
func checkDestSingleContainSliceNullerSqlValuer(v reflect.Value) (bool, error) {
	is, base := basePtrValue(v)
	if is && v.IsNil() { //数值无效，直接返回false，不再进行合法性检查
		return false, nil
	}

	if !is { //不是ptr，必须 Nuller
		is = baseStructValue(base)
		if is {
			_, ok := base.Interface().(driver.Valuer)
			if !ok {
				return false, errors.New("struct field " + base.String() + " need imp sql Value")
			}
			_, ok = base.Interface().(types.NullEr)
			if !ok {
				return false, errors.New("struct field " + base.String() + " need imp core NullEr ")
			}
			return true, nil
		}

		is, _, _ = baseSliceValue(base)
		if is {
			return true, nil
		}

		return false, errors.New("struct field " + base.String() + " need ptr or NullEr")
	}

	is = baseBaseValue(base)
	if is {
		return true, nil
	}

	is, _, _ = baseSliceValue(base)
	if is {
		return true, nil
	}

	is = baseStructValue(base)
	if is {
		_, ok := base.Interface().(driver.Valuer)
		if !ok {
			return false, errors.New("struct field " + base.String() + " need imp sql Value")
		}
		return true, nil
	}

	return false, errors.New("need a struct or base type")
}

// v0.5
// slice - ptr				1
// map[string]intface - ptr		2
// struct - ptr					3
// base struct-base- ptr		4
// isPtr  是否指针
// hasContain 内容是否为空

// base 基础value map slice 返回自身

type ArgTyp struct {
	isPtr      bool
	typ        int
	base       reflect.Value
	hasContain bool
}

func checkDestTyp(v reflect.Value) (a ArgTyp, err error) {
	isPtr, base := basePtrValue(v)
	if isPtr && v.IsNil() { //数值无效，直接返回false，不再进行合法性检查
		err = ErrNil
		return
	}
	a.isPtr = isPtr

	//map
	is, has, key, _ := baseMapValue(base)
	if is {
		if key.Kind() != reflect.String {
			err = errors.New(" map type key err no string  ")
			return
		}
		a.typ = 2
		a.hasContain = has
		a.base = base
		return
	}

	//slice
	is, has, _ = baseSliceValue(base)
	if is {
		a.typ = 1
		a.base = base
		a.hasContain = has
		return
	}

	// base
	is = baseBaseValue(base)
	if is {
		a.typ = 4
		a.base = base
		return
	}

	is = baseStructValue(base)
	if is {
		//struct-base
		_, ok := base.Interface().(driver.Valuer)
		if ok {
			a.typ = 4
			a.base = base
			return
		}
		a.typ = 3
		a.base = base
		return
	}
	err = errors.New("  type err   " + base.Kind().String())
	return
}
