package lorm

import (
	"errors"
	"reflect"
)

type ByContext struct {
	fieldNames []string

	primaryKeyValue []any // 主键值列表

	// model 需要 ormContext 才能解析
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
func By(v ...any) *ByContext {
	return &ByContext{
		primaryKeyValue: v,
		wb:              Wb(),
	}
}
func (b *ByContext) Field(name ...string) *ByContext {
	b.fieldNames = append(b.fieldNames, name...)
	return b
}

func (b *ByContext) Map(v map[string]any) *ByContext {
	cv, err := getMapCV(reflect.ValueOf(v))
	if err != nil {
		b.err = err
		return b
	}
	for i, column := range cv.columns {
		b.wb.Eq(column, cv.columnValues[i])
	}
	return b
}

func (b *ByContext) Model(v any) *ByContext {
	b.model = v
	return b
}

func (b *ByContext) Where(wb *WhereBuilder) *ByContext {
	b.wb.And(wb)
	return b
}
