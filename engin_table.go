package lorm

import (
	"bytes"
	"database/sql"
	"github.com/lontten/lorm/types"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

type EngineTable struct {
	dialect Dialect
	ctx     OrmContext

	//where tokens
	whereTokens []string

	extraWhereSql []byte

	//where values
	args []interface{}
}

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

//update
func (e EngineTable) doUpdate() (int64, error) {
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

//根据 byModel 生成的where token
func (e EngineTable) initByModel(v interface{}) {
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
func (e EngineTable) initByWhere(w *WhereBuilder) {
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
func (e EngineTable) initExtra() {
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
func (e EngineTable) initLgDel() {
	if err := e.ctx.err; err != nil {
		return
	}
	if ormConfig.LogicDeleteYesSql != "" {
		e.extraWhereSql = []byte(ormConfig.LogicDeleteYesSql)
	}
}

// Create
//v0.8
//1.ptr
//2.comp-struct
func (e EngineTable) Create(v interface{}) (num int64, err error) {
	e.setTargetDest(v)
	e.initColumnsValue()
	if e.ctx.err != nil {
		return 0, e.ctx.err
	}
	sqlStr := e.ctx.tableCreateGen()

	sqlStr += " RETURNING id"
	return e.query(sqlStr, e.ctx.columnValues...)
}

// CreateOrUpdate
//v0.6
//1.ptr
//2.comp-struct
func (e EngineTable) CreateOrUpdate(v interface{}) OrmTableCreate {
	e.setTargetDest(v)
	e.initColumnsValue()
	return OrmTableCreate{base: e}
}

// ByPrimaryKey
//v0.6
//ptr
//single / comp复合主键
func (orm OrmTableCreate) ByPrimaryKey() (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initSelfPrimaryKeyValues()
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	cs := ctx.columns
	cvs := ctx.columnValues
	tableName := ctx.tableName
	idNames := ctx.primaryKeyNames
	return base.dialect.insertOrUpdateByPrimaryKey(tableName, idNames, cs, cvs...)
}

// ByUnique
//v0.6
//ptr-comp
func (orm OrmTableCreate) ByUnique(fs types.Fields) (int64, error) {
	if fs == nil {
		return 0, errors.New("ByUnique is nil")
	}
	if len(fs) == 0 {
		return 0, errors.New("ByUnique is empty")
	}

	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	cs := ctx.columns
	cvs := ctx.columnValues
	tableName := ctx.tableName
	return base.dialect.insertOrUpdateByUnique(tableName, fs, cs, cvs...)
}

// Delete
//delete
func (e EngineTable) Delete(v interface{}) OrmTableDelete {
	e.setTargetDestOnlyTableName(v)
	return OrmTableDelete{base: e}
}

// ByPrimaryKey
//v0.8
//[]
//single -> 单主键
//comp -> 复合主键
func (orm OrmTableDelete) ByPrimaryKey(v ...interface{}) (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)

	base := orm.base
	ctx := orm.base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	delSql := ctx.genDelByPrimaryKey()
	idValues := orm.base.ctx.primaryKeyValues

	if len(v) == 1 {
		return base.dialect.exec(string(delSql), idValues[0]...)
	}
	return base.dialect.execBatch(string(delSql), idValues)
}

// ByModel
//v0.6
//ptr
//comp,只能一个comp-struct
func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initByModel(v)
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initExtra()
	return orm.base.doDel()
}

// ByWhere
//v0.6
func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initByWhere(w)
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initExtra()
	return orm.base.doDel()
}

// Update
//v0.6
func (e EngineTable) Update(v interface{}) OrmTableUpdate {
	e.setTargetDest(v)
	e.initColumnsValue()
	return OrmTableUpdate{base: e}
}

// ByPrimaryKey
//v0.8
func (orm OrmTableUpdate) ByPrimaryKey() (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initSelfPrimaryKeyValues()

	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	tableName := ctx.tableName
	cs := ctx.columns
	cvs := ctx.columnValues
	idValues := ctx.primaryKeyValues[0]

	whereStr := ctx.genWhereByPrimaryKey()

	var bb bytes.Buffer

	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
	bb.Write(whereStr)
	cvs = append(cvs, idValues...)

	return base.dialect.exec(bb.String(), cvs...)
}

func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initByModel(v)
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initByWhere(w)
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initExtra()
	return orm.base.doUpdate()
}

// Select
//select
func (e EngineTable) Select(v interface{}) OrmTableSelect {
	e.setScanDestSlice(v)
	e.initColumns()
	return OrmTableSelect{base: e}
}

// ByPrimaryKey
//v0.8
func (orm OrmTableSelect) ByPrimaryKey(v ...interface{}) (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)

	ctx := orm.base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	selSql := ctx.genSelectByPrimaryKey()
	idValues := ctx.primaryKeyValues
	if len(v) == 1 {
		return orm.base.query(string(selSql), idValues[0]...)
	}
	return orm.base.queryBatch(string(selSql), idValues)
}

//update
func (e EngineTable) doSelect() (int64, error) {
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

	return e.query(bb.String(), e.args...)
}

// ByModel
//v0.6
//ptr-comp
func (orm OrmTableSelect) ByModel(v interface{}) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initByModel(v)
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initByWhere(w)
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}
	orm.base.initExtra()
	return orm.base.doSelect()
}













//0.6
//初始化主键
func (e *EngineTable) initPrimaryKeyName() {
	if e.ctx.err != nil {
		return
	}
	e.ctx.primaryKeyNames = ormConfig.primaryKeys(e.ctx.tableName, e.ctx.destBaseValue)
}

//0.6
//初始化 表名
func (e *EngineTable) initTableName() {
	if e.ctx.err != nil {
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

	columns, err := ormConfig.initColumns(e.ctx.destBaseType)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.ctx.columns = columns
}

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

//v0.8
func (e EngineTable) query(query string, args ...interface{}) (int64, error) {
	rows, err := e.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	if e.ctx.isSlice {
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

//v0.7
// *.comp / slice.comp
//scan dest 一个comp-struct，或者一个slice-comp-struct
func (e *EngineTable) setScanDestSlice(v interface{}) {
	if e.ctx.err != nil {
		return
	}
	e.ctx.initScanDestSlice(v)
	e.ctx.checkScanDestField()
	e.initTableName()
}

//v0.6
//*.comp
//target dest 一个comp-struct
func (e *EngineTable) setTargetDest(v interface{}) {
	if e.ctx.err != nil {
		return
	}
	e.ctx.initTargetDest(v)
	e.ctx.checkTargetDestField()
	e.initTableName()
}

//v0.6
func (e *EngineTable) setTargetDestOnlyTableName(v interface{}) {
	if e.ctx.err != nil {
		return
	}
	e.ctx.initTargetDestOnlyBaseValue(v)
	e.ctx.checkTargetDestField()
	e.initTableName()
}

type OrmTableCreate struct {
	base EngineTable
}

type OrmTableSelect struct {
	base EngineTable

	query string
	args  []interface{}
}

type OrmTableSelectWhere struct {
	base EngineTable
}

type OrmTableUpdate struct {
	base EngineTable
}

type OrmTableDelete struct {
	base EngineTable
}
