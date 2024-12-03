package lorm

import (
	"reflect"
)

type By struct {
	fieldNames []string

	// model 需要 ormContext 才能解析
	model any
	wb    *WhereBuilder

	err error
}

func NewBy() *By {
	return &By{
		wb: Wb(),
	}
}
func (b *By) Field(name ...string) *By {
	b.fieldNames = append(b.fieldNames, name...)
	return b
}

func (b *By) Map(v map[string]any) *By {
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

func (b *By) Model(v any) *By {
	b.model = v
	return b
}

func (b *By) Where(wb *WhereBuilder) *By {
	b.wb.And(wb)
	return b
}
