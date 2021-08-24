package lorm

import (
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type EngineTable struct {
	ormConf OrmConf
	dialect Dialect

	ctx OrmContext
}

func (e EngineTable) queryLn(query string, args ...interface{}) (int64, error) {
	rows, err := e.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return e.ormConf.ScanLn(rows, e.ctx.dest)
}

func (e *EngineTable) setTargetDest(v interface{}) {
	e.ctx.initTargetDest(v)
	if e.ctx.err != nil {
		return
	}
	err := checkValidFieldTypStruct(e.ctx.destBaseValue)
	if err != nil {
		e.ctx.err = err
		return
	}

	e.initTableName()
}

func (e *EngineTable) setTargetDestOnlyTableName(v interface{}) {
	value := reflect.ValueOf(v)
	_, base := basePtrValue(value)
	is, base := baseStructValue(base)
	if !is {
		e.ctx.err = errors.New("need a struct")
		return
	}
	e.ctx.destBaseValue = base
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

//create
func (e EngineTable) Create(v interface{}) (num int64, err error) {
	e.setTargetDest(v)
	if e.ctx.err != nil {
		return
	}
	e.initColumnsValue()
	if e.ctx.err != nil {
		return
	}

	createSqlStr := e.ctx.tableCreateArgs2SqlStr( )

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(e.ctx.tableName + " ")
	sb.WriteString(createSqlStr)

	return e.dialect.exec(sb.String(), e.ctx.columnValues...)
}

func (e EngineTable) CreateOrUpdate(v interface{}) OrmTableCreate {
	e.setTargetDest(v)
	if e.ctx.err != nil {
		return OrmTableCreate{base: e}
	}
	e.initColumnsValue()
	if e.ctx.err != nil {
		return OrmTableCreate{base: e}
	}
	return OrmTableCreate{base: e}
}

func (orm OrmTableCreate) ById() (int64, error) {
	base := orm.base
	if err := base.ctx.err; err != nil {
		return 0, err
	}
	tableName := base.ctx.tableName
	c := base.ctx.columns
	cv := base.ctx.columnValues
	var idValue interface{}
	for i, s := range c {
		if s == "id" {
			idValue = cv[i]
		}
	}

	var sb strings.Builder
	sb.WriteString("SELECT 1 ")
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(" WHERE id = ? ")
	rows, err := base.dialect.query(sb.String(), idValue)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
		sb.WriteString(" WHERE id = ? ")
		cv = append(cv, idValue)

		return base.dialect.exec(sb.String(), cv...)
	}
	columnSqlStr := base.ctx.tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return base.dialect.exec(sb.String(), cv...)
}

func (orm OrmTableCreate) ByModel(v interface{}) (int64, error) {
	base := orm.base
	if err := base.ctx.err; err != nil {
		return 0, err
	}

	va := reflect.ValueOf(v)
	_, va = basePtrValue(va)

	err := checkValidFieldTypStruct(va)
	if err != nil {
		return 0, err
	}
	tableName := base.ctx.tableName
	c := base.ctx.columns
	cv := base.ctx.columnValues

	columns, values, err := ormConf.getStructMappingColumnsValueNotNil(va)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		panic(err)
	}

	whereArgs2SqlStr := base.ctx.tableWhereArgs2SqlStr(columns)

	var sb strings.Builder
	sb.WriteString("SELECT 1 ")
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(whereArgs2SqlStr)
	rows, err := base.dialect.query(sb.String(), values...)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
		sb.WriteString(whereArgs2SqlStr)
		cv = append(cv, values...)

		return base.dialect.exec(sb.String(), cv...)
	}
	columnSqlStr := base.ctx.tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return base.dialect.exec(sb.String(), cv...)
}

func (orm OrmTableCreate) ByWhere(w *WhereBuilder) (int64, error) {
	base := orm.base

	if err := base.ctx.err; err != nil {
		return 0, err
	}
	tableName := base.ctx.tableName
	c := base.ctx.columns
	cv := base.ctx.columnValues

	if w == nil {
		return 0, nil
	}
	wheres := w.context.wheres
	args := w.context.args

	var sb strings.Builder
	sb.WriteString("WHERE ")
	for i, where := range wheres {
		if i == 0 {
			sb.WriteString(" WHERE " + where)
			continue
		}
		sb.WriteString(" AND " + where)
	}
	whereSql := sb.String()

	sb.Reset()
	sb.WriteString("SELECT 1 ")
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(whereSql)

	log.Println(sb.String(), args)
	rows, err := base.dialect.query(sb.String(), args...)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
		sb.WriteString(whereSql)
		cv = append(cv, args)

		return base.dialect.exec(sb.String(), cv...)
	}
	columnSqlStr := base.ctx.tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return base.dialect.exec(sb.String(), cv...)
}

//delete
func (e EngineTable) Delete(v interface{}) OrmTableDelete {
	e.setTargetDestOnlyTableName(v)
	if e.ctx.err != nil {
		return OrmTableDelete{base: e}
	}
	return OrmTableDelete{base: e}
}

func (orm OrmTableDelete) ByPrimaryKey(v ...interface{}) (int64, error) {
	base := orm.base

	if err := base.ctx.err; err != nil {
		return 0, err
	}

	pkLen := len(v)
	if pkLen == 0 {
		return 0, errors.New("ByPrimaryKey arg num 0")
	}
	targetDestLen := len(base.ctx.destValueArr)
	if targetDestLen > 1 && targetDestLen != pkLen {
		return 0, errors.New("need Pk num " + strconv.Itoa(targetDestLen) + "but Pk len is " + strconv.Itoa(pkLen))
	}

	base.initPrimaryKeyName()
	base.ctx.checkValidPrimaryKey(v)

	logicDeleteSetSql := ormConf.LogicDeleteSetSql
	logicDeleteYesSql := ormConf.LogicDeleteYesSql
	tableName := base.ctx.tableName
	whereSql := base.ctx.tableWherePrimaryKey2SqlStr(idNames, base.ormConf)

	var sb strings.Builder
	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	logicDeleteYesSql = strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	if logicDeleteSetSql == lgSql {
		sb.WriteString("DELETE FROM ")
		sb.WriteString(tableName)
		sb.WriteString("WHERE ")
		sb.WriteString(whereSql)
	} else {
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(lgSql)
		sb.WriteString("WHERE ")
		sb.WriteString(whereSql)
		sb.WriteString(" and ")
		sb.WriteString(logicDeleteYesSql)
	}

	return base.dialect.exec(sb.String(), v)
}

func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	base := orm.base
	if err := base.ctx.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkValidFieldTypStruct(va)
	if err != nil {
		return 0, err
	}

	columns, values, err := ormConf.getStructMappingColumnsValueNotNil(va)
	if err != nil {
		return 0, err
	}
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	whereArgs2SqlStr := base.ctx.tableWhereArgs2SqlStr(columns)
	var sb strings.Builder
	sb.WriteString("DELETE ")
	sb.WriteString(" FROM ")
	sb.WriteString(base.ctx.tableName)
	sb.WriteString(whereArgs2SqlStr)

	return base.dialect.exec(sb.String(), values...)
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
	e.setTargetDest(v)
	if e.ctx.err != nil {
		return OrmTableUpdate{base: e}
	}
	e.initColumnsValue()
	if e.ctx.err != nil {
		return OrmTableUpdate{base: e}
	}
	return OrmTableUpdate{base: e}
}

func (orm OrmTableUpdate) ByPrimaryKey(v ...interface{}) (int64, error) {
	base := orm.base
	if err := base.ctx.err; err != nil {
		return 0, err
	}

	err := checkValidFieldTypStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}

	base.initPrimaryKeyName()

	tableName := base.ctx.tableName
	c := base.ctx.columns
	cv := base.ctx.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
	sb.WriteString(" WHERE ")
	//sb.WriteString(orm.base.primaryKeyNames)
	sb.WriteString(" = ? ")
	cv = append(cv, v)
	return base.dialect.exec(sb.String(), cv...)
}

func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	base := orm.base

	if err := base.ctx.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkValidFieldTypStruct(va)
	if err != nil {
		return 0, err
	}

	tableName := base.ctx.tableName
	c := base.ctx.columns
	cv := base.ctx.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
	columns, values, err := ormConf.getStructMappingColumnsValueNotNil(va)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		return 0, err
	}
	whereArgs2SqlStr := base.ctx.tableWhereArgs2SqlStr(columns)
	sb.WriteString(whereArgs2SqlStr)

	cv = append(cv, values...)

	return base.dialect.exec(sb.String(), cv...)
}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	base := orm.base
	if err := base.ctx.err; err != nil {
		return 0, err
	}

	if w == nil {
		return 0, nil
	}
	wheres := w.context.wheres
	args := w.context.args

	tableName := base.ctx.tableName
	c := base.ctx.columns
	cv := base.ctx.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(base.ctx.tableUpdateArgs2SqlStr(c))
	sb.WriteString(" WHERE ")
	for i, where := range wheres {
		if i == 0 {
			sb.WriteString(where)
			continue
		}
		sb.WriteString(" AND " + where)
	}

	cv = append(cv, args...)

	return base.dialect.exec(sb.String(), cv...)
}

//select
func (e EngineTable) Select(v interface{}) OrmTableSelect {
	e.setTargetDest(v)
	if e.ctx.err != nil {
		return OrmTableSelect{base: e}
	}

	return OrmTableSelect{base: e}
}

func (orm OrmTableSelect) ById(v ...interface{}) (int64, error) {
	if err := orm.base.ctx.err; err != nil {
		return 0, err
	}

	err := checkValidFieldTypStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}
	err = orm.base.initColumns()
	if err != nil {
		return 0, err
	}
	orm.base.initPrimaryKeyName()
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

	return orm.base.queryLn(sb.String(), v)
}

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
	err := checkValidFieldTypStruct(va)
	if err != nil {
		return 0, err
	}
	err = base.initColumns()
	if err != nil {
		return 0, err
	}
	base.initPrimaryKeyName()

	tableName := base.ctx.tableName
	c := base.ctx.columns
	columns, values, err := ormConf.getStructMappingColumnsValueNotNil(va)
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
	err := orm.base.initColumns()
	if err != nil {
		return 0, err
	}
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

//init
func (e *EngineTable) initPrimaryKeyName() {
	e.ctx.primaryKeyNames = ormConf.primaryKeys(e.ctx.tableName, e.ctx.destBaseValue)
}

func (e *EngineTable) initTableName() {
	tableName, err := ormConf.tableName(e.ctx.destBaseValue)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.ctx.tableName = tableName
}

//获取struct对应的字段名 和 其值   有效部分
func (e *EngineTable) initColumnsValue() {
	if !e.ctx.isSlice {
		columns, values, err := ormConf.getStructMappingColumnsValueNotNil(e.ctx.destValueArr[0])
		if err != nil {
			e.ctx.err = err
			return
		}
		param := ModelParam{
			columns:      columns,
			columnValues: values,
		}
		e.ctx.modelParams = append(e.ctx.modelParams, param)
		return
	}

	modelParam, err := ormConf.getStructMappingColumnsValueList(e.ctx.destValueArr)
	if err != nil {
		e.ctx.err = err
		return
	}
	e.ctx.modelParams = append(e.ctx.modelParams, modelParam)
	return
}

//获取struct对应的字段名 有效部分
func (e *EngineTable) initColumns() error {
	columns, err := e.ormConf.initColumns(e.ctx.destBaseValue)
	if err != nil {
		return err
	}
	e.ctx.columns = columns
	return nil
}
