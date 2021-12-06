package lorm

import (
	"bytes"
	"database/sql"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

//update
func (e EngineTable) doUpdate() (int64, error) {
	if err := e.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := e.ctx
	tableName := ctx.tableName
	cs := ctx.columns

	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
	bb.Write(e.genWhereSqlByToken())

	return e.dialect.exec(bb.String(), append(ctx.columnValues, e.args...)...)

}

//del
func (e EngineTable) doDel() (int64, error) {
	if err := e.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer
	tableName := e.ctx.tableName
	where := e.genWhereSqlByToken()

	if ormConfig.LogicDeleteSetSql == "" {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.Write(where)
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(ormConfig.LogicDeleteSetSql)
		bb.Write(where)
	}

	return e.dialect.exec(bb.String(), e.args...)
}

//update
func (e EngineTable) doSelect(extra string) (int64, error) {
	if err := e.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := e.ctx
	tableName := ctx.tableName
	columns := ctx.columns

	bb.WriteString("SELECT ")
	for i, column := range columns {
		if i == 0 {
			bb.WriteString(column)
		} else {
			bb.WriteString(" , ")
			bb.WriteString(column)
		}
	}
	bb.WriteString(" FROM ")
	bb.WriteString(tableName)
	bb.Write(e.genWhereSqlByToken())
	bb.WriteString(extra)

	return e.query(bb.String(), e.args...)
}

//-------------------------------init------------------------

//根据 byModel 生成的where token
func (e *EngineTable) initByPrimaryKey() {
	ctx := e.ctx
	if err := ctx.err; err != nil {
		return
	}

	e.whereTokens = append(e.whereTokens, utils.GenwhereToken(ctx.primaryKeyNames)...)

	for _, value := range ctx.primaryKeyValues {
		e.args = append(e.args, value...)
	}
}

//根据 byModel 生成的where token
func (e *EngineTable) initByModel(v interface{}) {
	if err := e.ctx.err; err != nil {
		return
	}
	if v == nil {
		e.ctx.err = errors.New("model is nil")
		return
	}

	columns, values, err := getCompCV(v)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.whereTokens = append(e.whereTokens, utils.GenwhereToken(columns)...)
	e.args = append(e.args, values...)
}

//根据 byWhere 生成的where token
func (e *EngineTable) initByWhere(w *WhereBuilder) {
	if err := e.ctx.err; err != nil {
		return
	}
	if w == nil {
		e.ctx.err = errors.New("ByWhere is nil")
		return
	}

	args := w.context.args
	wheres := w.context.wheres

	e.whereTokens = append(e.whereTokens, wheres...)
	e.args = append(e.args, args...)
}

//init 逻辑删除、租户
func (e *EngineTable) initExtra() {
	if err := e.ctx.err; err != nil {
		return
	}

	if ormConfig.LogicDeleteYesSql != "" {
		e.extraWhereSql = []byte(ormConfig.LogicDeleteYesSql)
	}

	if ormConfig.TenantIdFieldName != "" {
		e.whereTokens = append(e.whereTokens, ormConfig.TenantIdFieldName)
		e.args = append(e.args, ormConfig.TenantIdValueFun())
	}
}

//初始化逻辑删除
func (e *EngineTable) initLgDel() {
	if err := e.ctx.err; err != nil {
		return
	}
	if ormConfig.LogicDeleteYesSql != "" {
		e.extraWhereSql = []byte(ormConfig.LogicDeleteYesSql)
	}
}

//-------------------------------target------------------------

//v0.6
//*.comp
//target scanDest 一个comp-struct
func (e *EngineTable) setTargetDest(v interface{}) {
	if e.ctx.err != nil {
		return
	}
	e.ctx.initTargetDest(v)
	e.ctx.checkTargetDestField()
	e.initTableName()
}

//v0.6
func (e *EngineTable) setTargetDest2TableName(v interface{}) {
	if e.ctx.err != nil {
		return
	}
	e.ctx.initTargetDest2TableName(v)
	e.initTableName()
}

//0.6
//初始化主键
func (e *EngineTable) initPrimaryKeyName() {
	if e.ctx.err != nil {
		return
	}
	e.ctx.primaryKeyNames = ormConfig.primaryKeys(e.ctx.tableName)
}

//0.6
//初始化 表名
func (e *EngineTable) initTableName() {
	if e.ctx.err != nil {
		return
	}
	if e.ctx.tableName != "" {
		return
	}
	tableName, err := ormConfig.tableName(e.ctx.destBaseType)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.ctx.tableName = tableName
}

//0.6
//获取struct对应的字段名 和 其值，
//slice为全部，一个为非nil字段。
func (e *EngineTable) initColumnsValue() {
	if e.ctx.err != nil {
		return
	}
	columns, valuess, err := ormConfig.getCompColumnsValueNoNil(e.ctx.destValue)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.ctx.columns = columns
	e.ctx.columnValues = valuess
	return
}

//v0.6
//获取struct对应的字段名 有效部分
func (e *EngineTable) initColumns() {
	if e.ctx.err != nil {
		return
	}

	columns, err := ormConfig.initColumns(e.ctx.scanDestBaseType)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.ctx.columns = columns
}

//-------------------------utils------------------------
//v0.6
//获取comp 的 cv
//排除 nil 字段
func getCompCV(v interface{}) ([]string, []interface{}, error) {
	value := reflect.ValueOf(v)
	_, value, err := basePtrDeepValue(value)
	if err != nil {
		return nil, nil, err
	}

	return getCompValueCV(value)
}

//v0.6
//排除 nil 字段
func getCompValueCV(v reflect.Value) ([]string, []interface{}, error) {
	if !isCompType(v.Type()) {
		return nil, nil, errors.New("getvcv not comp")
	}
	err := checkCompField(v)
	if err != nil {
		return nil, nil, err
	}

	columns, values, err := ormConfig.getCompColumnsValueNoNil(v)
	if err != nil {
		return nil, nil, err
	}
	if len(columns) < 1 {
		return nil, nil, errors.New("where model valid field need ")
	}
	return columns, values, nil
}

//------------------------query--------------------------
//v0.8
func (e EngineTable) query(query string, args ...interface{}) (int64, error) {
	rows, err := e.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	if e.ctx.scanIsSlice {
		return e.ctx.Scan(rows)
	}
	return e.ctx.ScanLn(rows)
}

//v0.8
func (e EngineTable) queryBatch(query string, args [][]interface{}) (int64, error) {
	stmt, err := e.dialect.queryBatch(query)
	if err != nil {
		return 0, err
	}

	rowss := make([]*sql.Rows, 0)
	for _, arg := range args {
		rows, err := stmt.Query(arg...)
		if err != nil {
			return 0, err
		}
		rowss = append(rowss, rows)
	}
	return e.ctx.ScanBatch(rowss)
}

//------------------------gen-sql---------------------------

//根据whereTokens生成的where sql
func (e EngineTable) genWhereSqlByToken() []byte {
	if len(e.whereTokens) == 0 && e.extraWhereSql == nil {
		return nil
	}
	var buf bytes.Buffer
	buf.WriteString(" WHERE ")
	for i, token := range e.whereTokens {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(token)
	}
	buf.Write(e.extraWhereSql)
	return buf.Bytes()
}
