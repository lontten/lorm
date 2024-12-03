package lorm

import (
	"github.com/lontten/lorm/field"
	"reflect"
)

type SetContext struct {
	// 更新的字段
	fieldNames []string

	// model 需要 ormContext 才能解析
	model    any
	hasModel bool

	// 其他可以直接解析
	columns      []string
	columnValues []field.Value

	err error
}

func Set() *SetContext {
	return &SetContext{}
}

func (s *SetContext) Field(name ...string) *SetContext {
	s.fieldNames = append(s.fieldNames, name...)
	return s
}

func (s *SetContext) Set(name string, value any) *SetContext {
	s.columns = append(s.columns, name)
	s.columnValues = append(s.columnValues, field.Value{
		Type:  field.Val,
		Value: value,
	})
	return s
}

func (s *SetContext) SetNull(name string) *SetContext {
	s.columns = append(s.columns, name)
	s.columnValues = append(s.columnValues, field.Value{
		Type: field.Null,
	})
	return s
}

// 自增，自减
func (s *SetContext) SetIncrement(name string, num any) *SetContext {
	s.columns = append(s.columns, name)
	s.columnValues = append(s.columnValues, field.Value{
		Type:  field.Increment,
		Value: num,
	})
	return s
}

// 自定义表达式
// SetExpression("name", "substr(time('now'), 12)") // sqlite 设置时分秒
func (s *SetContext) SetExpression(name string, expression string) *SetContext {
	s.columns = append(s.columns, name)
	s.columnValues = append(s.columnValues, field.Value{
		Type:  field.Expression,
		Value: expression,
	})
	return s
}

func (s *SetContext) Map(v map[string]any) *SetContext {
	cv, err := getMapCV(reflect.ValueOf(v))
	if err != nil {
		s.err = err
		return s
	}
	s.columns = append(s.columns, cv.columns...)
	s.columnValues = append(s.columnValues, cv.columnValues...)
	return s
}

// 这里不对 model 进行解析
// 在 initColumnsValueSet 中解析
func (s *SetContext) Model(v any) *SetContext {
	s.model = v
	return s
}
