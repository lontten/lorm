package lorm

import (
	"bytes"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type EngineTable struct {
	dialect Dialect

	ctx OrmContext
}

//v0.6
func (e EngineTable) queryLn(query string, args ...interface{}) (int64, error) {
	rows, err := e.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ormConfig.ScanLn(rows, e.ctx.dest)
}

//v0.6
func (e EngineTable) query(query string, args ...interface{}) (int64, error) {
	rows, err := e.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ormConfig.Scan(rows, e.ctx.dest)
}

//v0.6
//func (e EngineTable) queryLnBatch(query string, args [][]interface{}) (int64, error) {
//	stmt, err := e.dialect.queryBatch(query)
//	if err != nil {
//		return 0, err
//	}
//
//	rowss := make([]*sql.Rows, 0)
//	for _, arg := range args {
//		rows, err := stmt.Query(arg...)
//		if err != nil {
//			return 0, err
//		}
//		rowss = append(rowss, rows)
//	}
//
//	return ormConfig.ScanLnBatch(rowss, utils.ToSlice(e.ctx.destValue))
//}
//v0.6
//func (e EngineTable) queryBatch(query string, args [][]interface{}) (int64, error) {
//	stmt, err := e.dialect.queryBatch(query)
//	if err != nil {
//		return 0, err
//	}
//
//	rowss := make([]*sql.Rows, 0)
//	for _, arg := range args {
//		rows, err := stmt.Query(arg...)
//		if err != nil {
//			return 0, err
//		}
//		rowss = append(rowss, rows)
//	}
//
//	return ormConfig.ScanBatch(rowss, utils.ToSlice(e.ctx.destValue))
//}

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
	return e.queryLn(sqlStr, e.ctx.columnValues...)
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
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	keyNum := len(ctx.primaryKeyNames)
	idValues := make([]interface{}, 0)
	columns, values, err := getCompCV(ctx.dest)
	if err != nil {
		return 0, err
	}
	//只要主键字段
	for _, key := range ctx.primaryKeyNames {
		for i, c := range columns {
			if c == key {
				idValues = append(idValues, values[i])
				continue
			}
		}
	}
	idLen := len(idValues)
	if idLen == 0 {
		return 0, errors.New("no pk")
	}
	if keyNum != idLen {
		return 0, errors.New("comp pk num err")
	}

	cs := ctx.columns
	cvs := ctx.columnValues
	tableName := ctx.tableName

	whereStr := ctx.genWhereByPrimaryKey()

	var bb bytes.Buffer
	bb.WriteString("SELECT 1 ")
	bb.WriteString(" FROM ")
	bb.WriteString(tableName)
	bb.Write(whereStr)
	bb.WriteString("limit 1")
	rows, err := base.dialect.query(bb.String(), idValues)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		bb.Reset()
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
		bb.Write(whereStr)
		cvs = append(cvs, idValues)

		return base.dialect.exec(bb.String(), cvs...)
	}

	columnSqlStr := ctx.tableCreateArgs2SqlStr()

	bb.Reset()
	bb.WriteString("INSERT INTO ")
	bb.WriteString(tableName)
	bb.WriteString(columnSqlStr)

	return base.dialect.exec(bb.String(), cvs...)
}

// ByModel
//v0.6
//ptr-comp
func (orm OrmTableCreate) ByModel(v interface{}) (int64, error) {
	if v == nil {
		return 0, errors.New("ByModel is nil")
	}
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	c := ctx.columns
	cv := ctx.columnValues
	tableName := ctx.tableName

	columns, values, err := getCompCV(v)
	if err != nil {
		return 0, err
	}
	where := ctx.genWhere(columns)

	var bb bytes.Buffer
	bb.WriteString("SELECT 1 ")
	bb.WriteString(" FROM ")
	bb.WriteString(tableName)
	bb.Write(where)
	bb.WriteString("limit 1")
	rows, err := base.dialect.query(bb.String(), values...)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		bb.Reset()
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(ctx.tableUpdateArgs2SqlStr(c))
		bb.Write(where)
		cv = append(cv, values...)

		return base.dialect.exec(bb.String(), cv...)
	}
	columnSqlStr := ctx.tableCreateArgs2SqlStr()

	bb.Reset()
	bb.WriteString("INSERT INTO ")
	bb.WriteString(tableName)
	bb.WriteString(columnSqlStr)

	return base.dialect.exec(bb.String(), cv...)
}

func (orm OrmTableCreate) ByWhere(w *WhereBuilder) (int64, error) {
	if w == nil {
		return 0, errors.New("ByWhere is nil")
	}
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}
	c := ctx.columns
	cv := ctx.columnValues
	tableName := ctx.tableName

	wheres := w.context.wheres
	args := w.context.args

	var bb bytes.Buffer
	bb.WriteString("WHERE ")
	for i, where := range wheres {
		if i == 0 {
			bb.WriteString(" WHERE " + where)
			continue
		}
		bb.WriteString(" AND " + where)
	}
	whereSql := bb.String()

	bb.Reset()
	bb.WriteString("SELECT 1 ")
	bb.WriteString(" FROM ")
	bb.WriteString(tableName)
	bb.WriteString(whereSql)
	bb.WriteString("limit 1")

	rows, err := base.dialect.query(bb.String(), args...)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		bb.Reset()
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(ctx.tableUpdateArgs2SqlStr(c))
		bb.WriteString(whereSql)
		cv = append(cv, args)

		return base.dialect.exec(bb.String(), cv...)
	}
	columnSqlStr := ctx.tableCreateArgs2SqlStr()

	bb.Reset()
	bb.WriteString("INSERT INTO ")
	bb.WriteString(tableName)
	bb.WriteString(columnSqlStr)

	return base.dialect.exec(bb.String(), cv...)
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
	orm.base.ctx.initSetPrimaryKey(v)

	base := orm.base
	ctx := orm.base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	delSql := ctx.genDelByPrimaryKey()
	return base.dialect.exec(string(delSql), v...)
}

// ByModel
//v0.6
//ptr
//comp,只能一个comp-struct
func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	if v == nil {
		return 0, errors.New("ByModel is nil")
	}
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	columns, values, err := getCompCV(v)
	if err != nil {
		return 0, err
	}

	delSql := ctx.genDel(columns)
	return base.dialect.exec(string(delSql), values...)
}

// ByWhere
//v0.6
func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
	if w == nil {
		return 0, nil
	}
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	wheres := w.context.wheres
	args := w.context.args

	delSql := ctx.genDel(wheres)
	return base.dialect.exec(string(delSql), args...)
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
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	keyNum := len(ctx.primaryKeyNames)
	idValues := make([]interface{}, 0)
	columns, values, err := getCompCV(ctx.dest)
	if err != nil {
		return 0, err
	}
	//只要主键字段
	for _, key := range ctx.primaryKeyNames {
		for i, c := range columns {
			if c == key {
				idValues = append(idValues, values[i])
				continue
			}
		}
	}
	idLen := len(idValues)
	if idLen == 0 {
		return 0, errors.New("no pk")
	}
	if keyNum != idLen {
		return 0, errors.New("comp pk num err")
	}

	tableName := ctx.tableName
	cs := ctx.columns
	cvs := ctx.columnValues

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
	if v == nil {
		return 0, errors.New("ByModel is nil")
	}
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	c := ctx.columns
	cv := ctx.columnValues
	tableName := ctx.tableName

	columns, values, err := getCompCV(v)
	if err != nil {
		return 0, err
	}
	where := ctx.genWhere(columns)

	var bb bytes.Buffer
	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(c))
	bb.Write(where)
	cv = append(cv, values...)

	return base.dialect.exec(bb.String(), cv...)

}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	if w == nil {
		return 0, errors.New("ByWhere is nil")
	}
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}
	c := ctx.columns
	cv := ctx.columnValues
	tableName := ctx.tableName

	wheres := w.context.wheres
	args := w.context.args

	var bb bytes.Buffer
	bb.WriteString("WHERE ")
	for i, where := range wheres {
		if i == 0 {
			bb.WriteString(" WHERE " + where)
			continue
		}
		bb.WriteString(" AND " + where)
	}
	whereSql := bb.String()

	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(c))
	bb.WriteString(whereSql)
	cv = append(cv, args)

	return base.dialect.exec(bb.String(), cv...)

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
	orm.base.ctx.initSetPrimaryKey(v)

	ctx := orm.base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	selSql := ctx.genSelectByPrimaryKey()
	return orm.base.queryLn(string(selSql), v...)
}

// ByModel
//v0.6
//ptr-comp
func (orm OrmTableSelect) ByModel(v interface{}) (int64, error) {
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	columns, values, err := getCompCV(v)
	if err != nil {
		return 0, err
	}

	tableName := ctx.tableName
	c := ctx.columns

	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, column := range c {
		if i == 0 {
			sb.WriteString(column)
		} else {
			sb.WriteString(" , ")
			sb.WriteString(column)
		}
	}
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(ctx.tableWhereArgs2SqlStr(columns))

	return base.queryLn(sb.String(), values...)
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}

	if w == nil {
		return 0, errors.New("table select where can't nil")
	}
	orm.base.initColumns()
	orm.base.initPrimaryKeyName()

	wheres := w.context.wheres
	args := w.context.args

	tableName := orm.base.ctx.tableName
	c := orm.base.ctx.columns

	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, column := range c {
		if i == 0 {
			sb.WriteString(column)
		} else {
			sb.WriteString(" , ")
			sb.WriteString(column)
		}
	}
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(" WHERE ")
	for i, where := range wheres {
		if i == 0 {
			sb.WriteString(where)
			continue
		}
		sb.WriteString(" AND " + where)
	}

	return orm.base.queryLn(sb.String(), args...)
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
	tableName, err := ormConfig.tableName(e.ctx.destBaseValue)
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

	columns, err := ormConfig.initColumns(e.ctx.destBaseValue)
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

	ctyp := checkCompValue(value)
	if ctyp != Composite {
		return nil, nil, errors.New("getcv not comp")
	}
	err = checkCompField(value)
	if err != nil {
		return nil, nil, err
	}

	columns, values, err := ormConfig.getCompColumnsValueNoNil(value)
	if err != nil {
		return nil, nil, err
	}
	if len(columns) < 1 {
		return nil, nil, errors.New("where model valid field need ")
	}
	return columns, values, nil
}

//v0.6
//排除 nil 字段
func getCompValueCV(v reflect.Value) ([]string, []interface{}, error) {
	ctyp := checkCompValue(v)
	if ctyp != Composite {
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
