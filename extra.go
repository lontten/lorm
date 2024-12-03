package lorm

import (
	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/insert_type"
	"github.com/lontten/lorm/return_type"
)

// ExtraContext 扩展参数
type ExtraContext struct {
	insertType     insert_type.InsertType
	returnType     return_type.ReturnType
	showSql        bool
	skipSoftDelete bool
	tableName      string

	columns      []string
	columnValues []field.Value

	// 唯一索引字段名列表
	duplicateKeyNames []string
	set               *SetContext
	err               error
}

func Extra() *ExtraContext {
	return &ExtraContext{}
}

func (e *ExtraContext) GetErr() error {
	if e.err != nil {
		return e.err
	}
	if e.set != nil {
		return e.set.err
	}
	return nil
}

func (e *ExtraContext) ShowSql() *ExtraContext {
	e.showSql = true
	return e
}

func (e *ExtraContext) SkipSoftDelete() *ExtraContext {
	e.skipSoftDelete = true
	return e
}

func (e *ExtraContext) TableName(name string) *ExtraContext {
	e.tableName = name
	return e
}

func (e *ExtraContext) ReturnType(typ return_type.ReturnType) *ExtraContext {
	e.returnType = typ
	return e
}

func (e *ExtraContext) SetNull(name string) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type: field.Null,
	})
	return e
}

func (e *ExtraContext) SetNow(name string) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type: field.Now,
	})
	return e
}

func (e *ExtraContext) Set(name string, value any) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Val,
		Value: value,
	})
	return e
}

// 自增，自减
func (e *ExtraContext) SetIncrement(name string, num any) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Increment,
		Value: num,
	})
	return e
}

// 自定义表达式
// SetExpression("name", "substr(time('now'), 12)") // sqlite 设置时分秒
func (e *ExtraContext) SetExpression(name string, expression string) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Expression,
		Value: expression,
	})
	return e
}

type DuplicateKey struct {
	e *ExtraContext
}

//.whenDuplicateKey(name ...string, )
//.do(nothing, nil)
//.do(update, all, .set(), select ("name", "age"))
//.do(replace, all, .set(), select ("name", "age"))

// 唯一索引冲突
func (e *ExtraContext) WhenDuplicateKey(name ...string) *DuplicateKey {
	e.duplicateKeyNames = name
	return &DuplicateKey{
		e: e,
	}
}

func (dk *DuplicateKey) DoNothing() *ExtraContext {
	dk.e.insertType = insert_type.Ignore
	return dk.e
}

func (dk *DuplicateKey) update(insertType insert_type.InsertType, set ...*SetContext) *ExtraContext {
	var sc *SetContext
	if len(set) == 0 {
		sc = &SetContext{}
	} else {
		if set[0] == nil {
			sc = &SetContext{}
		} else {
			sc = set[0]
		}
	}
	dk.e.insertType = insertType
	dk.e.set = sc
	return dk.e
}

func (dk *DuplicateKey) DoUpdate(set ...*SetContext) *ExtraContext {
	return dk.update(insert_type.Update, set...)
}

func (dk *DuplicateKey) DoReplace(set ...*SetContext) *ExtraContext {
	return dk.update(insert_type.Replace, set...)
}
