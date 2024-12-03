package lorm

import (
	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/insert_type"
	"github.com/lontten/lorm/return_type"
)

// Extra
//
//	InsertType(InsertType.Overvide) // 插入执行逻辑
//	ShowSql()  // 打印sql
//	SkipSoftDelete() //跳过软删除逻辑
//	ReturnType(Field.All,Field.Nil,Field.Pk,Field.None) //返回所有字段，；只返回nil字段；只返回主键字段
//	tableName("")  //覆盖表名
type Extra struct {
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

func (e *Extra) GetErr() error {
	if e.err != nil {
		return e.err
	}
	if e.set != nil {
		return e.set.err
	}
	return nil
}

func (e *Extra) ShowSql() *Extra {
	e.showSql = true
	return e
}

func (e *Extra) SkipSoftDelete() *Extra {
	e.skipSoftDelete = true
	return e
}

func (e *Extra) TableName(name string) *Extra {
	e.tableName = name
	return e
}

func (e *Extra) ReturnType(typ return_type.ReturnType) *Extra {
	e.returnType = typ
	return e
}

func (e *Extra) SetNull(name string) *Extra {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type: field.Null,
	})
	return e
}

func (e *Extra) SetNow(name string) *Extra {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type: field.Now,
	})
	return e
}

func (e *Extra) Set(name string, value any) *Extra {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Val,
		Value: value,
	})
	return e
}

// 自增，自减
func (e *Extra) SetIncrement(name string, num any) *Extra {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Increment,
		Value: num,
	})
	return e
}

// 自定义表达式
// SetExpression("name", "substr(time('now'), 12)") // sqlite 设置时分秒
func (e *Extra) SetExpression(name string, expression string) *Extra {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Expression,
		Value: expression,
	})
	return e
}

type DuplicateKey struct {
	e *Extra
}

//.whenDuplicateKey(name ...string, )
//.do(nothing, nil)
//.do(update, all, .set(), select ("name", "age"))
//.do(replace, all, .set(), select ("name", "age"))

// 唯一索引冲突
func (e *Extra) WhenDuplicateKey(name ...string) *DuplicateKey {
	e.duplicateKeyNames = name
	return &DuplicateKey{
		e: e,
	}
}

func (dk *DuplicateKey) DoNothing() *Extra {
	dk.e.insertType = insert_type.Ignore
	return dk.e
}

func (dk *DuplicateKey) update(insertType insert_type.InsertType, set ...*SetContext) *Extra {
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

func (dk *DuplicateKey) DoUpdate(set ...*SetContext) *Extra {
	return dk.update(insert_type.Update, set...)
}

func (dk *DuplicateKey) DoReplace(set ...*SetContext) *Extra {
	return dk.update(insert_type.Replace, set...)
}
