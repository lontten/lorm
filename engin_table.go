package lorm

import (
	"database/sql"
	"github.com/lontten/lorm/utils"
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
func (e EngineTable) queryLnBatch(query string, args [][]interface{}) (int64, error) {
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

	return ormConfig.ScanLnBatch(rowss, utils.ToSlice(e.ctx.destValue))
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
func (e *EngineTable) setTargetDestSlice(v interface{}) {
	if e.ctx.err != nil {
		return
	}
	e.ctx.initTargetDestSlice(v)
	e.ctx.checkTargetDestField()
	e.initTableName()
}

//v0.6
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
//v0.6
//1.ptr-comp-struct
//2.slice-comp-struct
func (e EngineTable) Create(v interface{}) (num int64, err error) {
	e.setTargetDestSlice(v)
	e.initColumnsValue()
	if e.ctx.err != nil {
		return 0, e.ctx.err
	}
	sqlStr := e.ctx.tableCreateGen()

	if e.ctx.isSlice {
		return e.dialect.execBatch(sqlStr, e.ctx.columnValues)
	}
	return e.dialect.exec(sqlStr, e.ctx.columnValues[0])
}

//v0.6
//只能一个
func (e EngineTable) CreateOrUpdate(v interface{}) OrmTableCreate {
	e.setTargetDest(v)
	e.initColumnsValue()
	return OrmTableCreate{base: e}
}

//v0.6
//ptr
//single / comp复合主键
func (orm OrmTableCreate) ById(v interface{}) (int64, error) {
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	keyNum := len(ctx.primaryKeyNames)

	tableName := ctx.tableName
	cs := ctx.columns
	cvs := ctx.columnValues[0]

	idValues := make([]interface{}, 0)
	//主键为nil，从dest中获取id value
	if v == nil {
		for _, key := range ctx.primaryKeyNames {
			for i, s := range cs {
				if s == key {
					idValues = append(idValues, cvs[i])
					continue
				}
			}
		}
	} else {
		columns, values, err := getCompCV(v)
		if err != nil {
			return 0, err
		}

		for _, key := range ctx.primaryKeyNames {
			for i, c := range columns {
				if c == key {
					idValues = append(idValues, values[i])
					continue
				}
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

	whereStr := ctx.tableWherePrimaryKey2SqlStr(ctx.primaryKeyNames)

	var sb strings.Builder
	sb.WriteString("SELECT 1 ")
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(whereStr)
	rows, err := base.dialect.query(sb.String(), idValues)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
		sb.WriteString(whereStr)
		cvs = append(cvs, idValues)

		return base.dialect.exec(sb.String(), cvs...)
	}

	columnSqlStr := ctx.tableCreateArgs2SqlStr()

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return base.dialect.exec(sb.String(), cvs...)
}

//func (orm OrmTableCreate) ByModel(v interface{}) (int64, error) {
//	base := orm.base
//	if err := base.ctx.err; err != nil {
//		return 0, err
//	}
//
//	va := reflect.ValueOf(v)
//	_, va = basePtrValue(va)
//
//	err := checkStructValidFieldNuller(va)
//	if err != nil {
//		return 0, err
//	}
//	tableName := base.ctx.tableName
//	c := base.ctx.columns
//	cv := base.ctx.columnValues
//
//	columns, values, err := ormConfig.getCompColumnsValueNoNil(va)
//	if len(columns) < 1 {
//		return 0, errors.New("where model valid field need ")
//	}
//	if err != nil {
//		panic(err)
//	}
//
//	whereArgs2SqlStr := base.ctx.tableWhereArgs2SqlStr(columns)
//
//	var sb strings.Builder
//	sb.WriteString("SELECT 1 ")
//	sb.WriteString(" FROM ")
//	sb.WriteString(tableName)
//	sb.WriteString(whereArgs2SqlStr)
//	rows, err := base.dialect.query(sb.String(), values...)
//	if err != nil {
//		return 0, err
//	}
//	//update
//	if rows.Next() {
//		sb.Reset()
//		sb.WriteString("UPDATE ")
//		sb.WriteString(tableName)
//		sb.WriteString(" SET ")
//		sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
//		sb.WriteString(whereArgs2SqlStr)
//		cv = append(cv, values...)
//
//		return base.dialect.exec(sb.String(), cv...)
//	}
//	columnSqlStr := base.ctx.tableCreateArgs2SqlStr(c)
//
//	sb.Reset()
//	sb.WriteString("INSERT INTO ")
//	sb.WriteString(tableName)
//	sb.WriteString(columnSqlStr)
//
//	return base.dialect.exec(sb.String(), cv...)
//}

//func (orm OrmTableCreate) ByWhere(w *WhereBuilder) (int64, error) {
//	base := orm.base
//
//	if err := base.ctx.err; err != nil {
//		return 0, err
//	}
//	tableName := base.ctx.tableName
//	c := base.ctx.columns
//	cv := base.ctx.columnValues
//
//	if w == nil {
//		return 0, nil
//	}
//	wheres := w.context.wheres
//	args := w.context.args
//
//	var sb strings.Builder
//	sb.WriteString("WHERE ")
//	for i, where := range wheres {
//		if i == 0 {
//			sb.WriteString(" WHERE " + where)
//			continue
//		}
//		sb.WriteString(" AND " + where)
//	}
//	whereSql := sb.String()
//
//	sb.Reset()
//	sb.WriteString("SELECT 1 ")
//	sb.WriteString(" FROM ")
//	sb.WriteString(tableName)
//	sb.WriteString(whereSql)
//
//	log.Println(sb.String(), args)
//	rows, err := base.dialect.query(sb.String(), args...)
//	if err != nil {
//		return 0, err
//	}
//	//update
//	if rows.Next() {
//		sb.Reset()
//		sb.WriteString("UPDATE ")
//		sb.WriteString(tableName)
//		sb.WriteString(" SET ")
//		sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
//		sb.WriteString(whereSql)
//		cv = append(cv, args)
//
//		return base.dialect.exec(sb.String(), cv...)
//	}
//	columnSqlStr := base.ctx.tableCreateArgs2SqlStr(c)
//
//	sb.Reset()
//	sb.WriteString("INSERT INTO ")
//	sb.WriteString(tableName)
//	sb.WriteString(columnSqlStr)
//
//	return base.dialect.exec(sb.String(), cv...)
//}

//delete
func (e EngineTable) Delete(v interface{}) OrmTableDelete {
	e.setTargetDestOnlyTableName(v)
	return OrmTableDelete{base: e}
}

//v0.6
//[]-single -> 单主键
//[]-comp -> 复合主键
//func (orm OrmTableDelete) ByPrimaryKey(v ...interface{}) (int64, error) {
//	base := orm.base
//
//	if err := base.ctx.err; err != nil {
//		return 0, err
//	}
//
//	idLen := len(v)
//	if idLen == 0 {
//		return 0, errors.New("ByPrimaryKey arg num 0")
//	}
//
//	pkLen := len(base.ctx.primaryKeyNames)
//	if pkLen==1 { //单主键
//		for _, i := range v {
//			value := reflect.ValueOf(i)
//			_, value = basePtrDeepValue(value)
//			ctyp := checkCompTypeValue(value, false)
//			if ctyp != Single {
//				return 0,  errors.New("ByPrimaryKey typ err")
//			}
//		}
//	}else {
//        for _, i := range v {
//            value := reflect.ValueOf(i)
//            _, value = basePtrDeepValue(value)
//            ctyp := checkCompTypeValue(value, false)
//            if ctyp != Composite {
//                return 0,  errors.New("ByPrimaryKey typ err")
//            }
//        }
//
//		//todo 多主键数量，和pkLen 是否相等
//	}
//
//
//
//	base.initPrimaryKeyName()
//	base.ctx.checkValidPrimaryKey(v)
//
//	logicDeleteSetSql := ormConfig.LogicDeleteSetSql
//	logicDeleteYesSql := ormConfig.LogicDeleteYesSql
//	tableName := base.ctx.tableName
//	whereSql := base.ctx.tableWherePrimaryKey2SqlStr(idNames)
//
//	var sb strings.Builder
//	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
//	logicDeleteYesSql = strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
//	if logicDeleteSetSql == lgSql {
//		sb.WriteString("DELETE FROM ")
//		sb.WriteString(tableName)
//		sb.WriteString("WHERE ")
//		sb.WriteString(whereSql)
//	} else {
//		sb.WriteString("UPDATE ")
//		sb.WriteString(tableName)
//		sb.WriteString(" SET ")
//		sb.WriteString(lgSql)
//		sb.WriteString("WHERE ")
//		sb.WriteString(whereSql)
//		sb.WriteString(" and ")
//		sb.WriteString(logicDeleteYesSql)
//	}
//
//	return base.dialect.exec(sb.String(), v)
//}

//v0.6
//comp,只能一个comp-struct
func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	base := orm.base
	ctx := base.ctx
	if err := ctx.err; err != nil {
		return 0, err
	}

	columns, values, err := getCompCV(v)
	if err != nil {
		return 0, err
	}

	whereArgs2SqlStr := ctx.tableWhereArgs2SqlStr(columns)
	var sb strings.Builder
	sb.WriteString("DELETE ")
	sb.WriteString(" FROM ")
	sb.WriteString(ctx.tableName)
	sb.WriteString(whereArgs2SqlStr)

	return base.dialect.exec(sb.String(), values...)
}

func getCompCV(v interface{}) ([]string, []interface{}, error) {
	value := reflect.ValueOf(v)
	_, value = basePtrDeepValue(value)
	ctyp := checkCompTypeValue(value, false)
	if ctyp != Composite {
		return nil, nil, errors.New("ByModel typ err")
	}
	err := checkCompValidFieldNuller(value)
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

func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}

	if w == nil {
		return 0, nil
	}
	wheres := w.context.wheres
	args := w.context.args

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	sb.WriteString(orm.base.ctx.tableName)
	sb.WriteString(" WHERE ")
	for i, where := range wheres {
		if i == 0 {
			sb.WriteString(where)
			continue
		}
		sb.WriteString(" AND " + where)
	}

	return orm.base.dialect.exec(sb.String(), args)
}

//update
func (e EngineTable) Update(v interface{}) OrmTableUpdate {
	if e.ctx.err != nil {
		return OrmTableUpdate{base: e}
	}
	e.setTargetDestSlice(v)
	if e.ctx.err != nil {
		return OrmTableUpdate{base: e}
	}
	e.initColumnsValue()
	if e.ctx.err != nil {
		return OrmTableUpdate{base: e}
	}
	return OrmTableUpdate{base: e}
}

//func (orm OrmTableUpdate) ByPrimaryKey(v ...interface{}) (int64, error) {
//	base := orm.base
//	if err := base.ctx.err; err != nil {
//		return 0, err
//	}
//
//	err := checkStructValidFieldNuller(reflect.ValueOf(v))
//	if err != nil {
//		return 0, err
//	}
//
//	base.initPrimaryKeyName()
//
//	tableName := base.ctx.tableName
//	c := base.ctx.columns
//	cv := base.ctx.columnValues
//
//	var sb strings.Builder
//	sb.WriteString(" UPDATE ")
//	sb.WriteString(tableName)
//	sb.WriteString(" SET ")
//	sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
//	sb.WriteString(" WHERE ")
//	//sb.WriteString(orm.base.primaryKeyNames)
//	sb.WriteString(" = ? ")
//	cv = append(cv, v)
//	return base.dialect.exec(sb.String(), cv...)
//}

//func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
//	base := orm.base
//
//	if err := base.ctx.err; err != nil {
//		return 0, err
//	}
//	va := reflect.ValueOf(v)
//	err := checkStructValidFieldNuller(va)
//	if err != nil {
//		return 0, err
//	}
//
//	tableName := base.ctx.tableName
//	c := base.ctx.columns
//	cv := base.ctx.columnValues
//
//	var sb strings.Builder
//	sb.WriteString(" UPDATE ")
//	sb.WriteString(tableName)
//	sb.WriteString(" SET ")
//	sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
//	columns, values, err := ormConfig.getCompColumnsValueNoNil(va)
//	if len(columns) < 1 {
//		return 0, errors.New("where model valid field need ")
//	}
//	if err != nil {
//		return 0, err
//	}
//	whereArgs2SqlStr := base.ctx.tableWhereArgs2SqlStr(columns)
//	sb.WriteString(whereArgs2SqlStr)
//
//	cv = append(cv, values...)
//
//	return base.dialect.exec(sb.String(), cv...)
//}

//func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
//	base := orm.base
//	if err := base.ctx.err; err != nil {
//		return 0, err
//	}
//
//	if w == nil {
//		return 0, nil
//	}
//	wheres := w.context.wheres
//	args := w.context.args
//
//	tableName := base.ctx.tableName
//	c := base.ctx.columns
//	cv := base.ctx.columnValues
//
//	var sb strings.Builder
//	sb.WriteString(" UPDATE ")
//	sb.WriteString(tableName)
//	sb.WriteString(" SET ")
//	sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
//	sb.WriteString(" WHERE ")
//	for i, where := range wheres {
//		if i == 0 {
//			sb.WriteString(where)
//			continue
//		}
//		sb.WriteString(" AND " + where)
//	}
//
//	cv = append(cv, args...)
//
//	return base.dialect.exec(sb.String(), cv...)
//}

//select
func (e EngineTable) Select(v interface{}) OrmTableSelect {
	e.setTargetDestOnlyTableName(v)
	if e.ctx.err != nil {
		return OrmTableSelect{base: e}
	}

	return OrmTableSelect{base: e}
}

//func (orm OrmTableSelect) ById(v ...interface{}) (int64, error) {
//	if err := orm.base.ctx.err; err != nil {
//		return 0, err
//	}
//
//	err := checkStructValidFieldNuller(reflect.ValueOf(v))
//	if err != nil {
//		return 0, err
//	}
//	err = orm.base.initColumns()
//	if err != nil {
//		return 0, err
//	}
//	orm.base.initPrimaryKeyName()
//	tableName := orm.base.ctx.tableName
//	c := orm.base.ctx.columns
//
//	var sb strings.Builder
//	sb.WriteString(" SELECT ")
//	for i, column := range c {
//		if i == 0 {
//			sb.WriteString(column)
//		} else {
//			sb.WriteString(" , ")
//			sb.WriteString(column)
//		}
//	}
//	sb.WriteString(" FROM ")
//	sb.WriteString(tableName)
//	sb.WriteString(" WHERE ")
//	//sb.WriteString(orm.base.primaryKeyNames)
//	sb.WriteString(" = ? ")
//
//	return orm.base.queryLn(sb.String(), v)
//}

func (orm OrmTableSelectWhere) getOne() (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}

	tableName := orm.base.ctx.tableName
	c := orm.base.ctx.columns

	var sb strings.Builder
	sb.WriteString(" SELECT ")
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
	//sb.WriteString(orm.base.primaryKeyNames)
	sb.WriteString(" = ? ")

	return orm.base.queryLn(sb.String(), orm.base.ctx.dest)
}

func (orm OrmTableSelectWhere) getList() (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}

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
	sb.WriteString("WHERE ")
	//sb.WriteString(orm.base.primaryKeyNames)
	sb.WriteString(" = ? ")

	return orm.base.queryLn(sb.String(), orm.base.ctx.dest)
}

func (orm OrmTableSelect) ByModel(v interface{}) (int64, error) {
	base := orm.base

	if err := base.ctx.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkStructValidFieldNuller(va)
	if err != nil {
		return 0, err
	}
	base.initColumns()
	if err != nil {
		return 0, err
	}
	base.initPrimaryKeyName()

	tableName := base.ctx.tableName
	c := base.ctx.columns
	columns, values, err := ormConfig.getCompColumnsValueNoNil(va)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		return 0, err
	}

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
	sb.WriteString(base.ctx.tableWhereArgs2SqlStr(columns))

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
//获取struct对应的字段名 和 其值   有效部分
func (e *EngineTable) initColumnsValue() {
	if e.ctx.err != nil {
		return
	}
	columns, valuess, err := ormConfig.initColumnsValue(e.ctx.destValueArr)
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
