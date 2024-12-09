package lorm

import (
	"reflect"
)

type ByContext struct {
	fieldNames []string

	// model 需要 ormContext 才能解析
	model any
	wb    *WhereBuilder

	err error
}

func By() *ByContext {
	return &ByContext{
		wb: Wb(),
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
