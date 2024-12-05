package lorm

import (
	"reflect"
)

type LormTabler interface {
	TableConf() *TableConf
}

// 表配置,会缓存，不可设置动态属性
type TableConf struct {
	tableName       *string  // 表名
	primaryKeyNames []string // 主键字段列表
	indexs          []Index  // 索引列表
	AutoIncrements  []string // 自增字段列表
}
type Index struct {
	Name      string   // 索引名称
	Unique    bool     // 是否唯一
	Columns   []string // 索引列
	IndexType string   // 索引类型
	Comment   string   // 索引注释
}

func (c *TableConf) Table(name string) *TableConf {
	c.tableName = &name
	return c
}

func (c *TableConf) PrimaryKeys(name ...string) *TableConf {
	c.primaryKeyNames = name
	return c
}

var TableConfCache = map[reflect.Type]TableConf{}

func GetTableConf(v reflect.Value) *TableConf {
	n, has := TableConfCache[v.Type()]
	if has {
		return &n
	}
	method := v.MethodByName("TableConf")
	if !method.IsValid() || method.IsZero() {
		return nil
	}

	values := method.Call(nil)

	if len(values) != 1 {
		return nil
	}
	value := values[0]
	if value.IsNil() {
		return nil
	}
	tc, ok := value.Interface().(*TableConf)
	if !ok {
		return nil
	}
	return tc
}

func GetTableName(v reflect.Value) *string {
	tc := GetTableConf(v)
	if tc == nil {
		return nil
	}
	return tc.tableName
}

func GetPrimaryKeyNames(v reflect.Value) []string {
	tc := GetTableConf(v)
	if tc == nil {
		return nil
	}
	return tc.primaryKeyNames
}
