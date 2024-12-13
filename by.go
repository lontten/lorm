package lorm

import (
	"errors"
	"reflect"
)

type ByContext struct {
	fieldNames []string

	primaryKeyValue       []any // 主键值列表
	filterPrimaryKeyValue []any // 主键值列表,过滤

	model any
	wb    *WhereBuilder

	err error
}

func (ctx ormContext) name(v ...any) {
	idLen := len(v)
	if idLen == 0 {
		return
	}
	pkLen := len(ctx.primaryKeyNames)

	idValuess := make([][]interface{}, 0)

	if pkLen == 1 {
		//单主键
		for _, i := range v {
			value := reflect.ValueOf(i)
			_, value, err := basePtrValue(value)
			if err != nil {
				ctx.err = err
				return
			}

			if ctx.checkParam {
				if !isValuerType(value.Type()) {
					ctx.err = errors.New("ByPrimaryKey typ err,not valuer")
					return
				}
			}

			idValues := make([]any, 1)
			idValues[0] = value.Interface()
			idValuess = append(idValuess, idValues)
		}

	} else {
		for _, i := range v {
			value := reflect.ValueOf(i)
			_, value, err := basePtrValue(value)
			if err != nil {
				ctx.err = err
				return
			}
			if !isCompType(value.Type()) {
				ctx.err = errors.New("ByPrimaryKey typ err,not comp")
				return
			}

			columns, values, err := getCompValueCV(value)
			if err != nil {
				ctx.err = err
				return
			}
			if len(columns) != pkLen {
				ctx.err = errors.New("复合主键，filed数量 len err")
				return
			}

			idValues := make([]any, 0)
			for _, f := range values {
				idValues = append(idValues, f)
			}
			idValuess = append(idValuess, idValues)
		}
	}
}
