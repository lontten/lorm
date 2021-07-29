package lorm

import (
	"fmt"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type EngineTable struct {
	db      DBer
	lorm    Lorm
	dialect Dialect

	context OrmContext

	primaryKeyNames []string

	//当前表名
	tableName string
	//当前struct对象
	dest          interface{}
	destBaseValue reflect.Value
	destIsSlice   bool

	columns      []string
	columnValues []interface{}
}

func (e EngineTable) queryLn(query string, args ...interface{}) (int64, error) {
	rows, err := e.db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return e.lorm.ScanLn(rows, e.dest)
}

func (e *EngineTable) setDest(v interface{}) error {
	value := reflect.ValueOf(v)
	is, OrmCore, err := destBaseValueCheckSlice(value)
	if err != nil {
		return err
	}

	err = checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return err
	}
	e.dest = v
	e.destBaseValue = OrmCore
	e.destIsSlice = is
	return e.initTableName()
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
	err = e.setDest(v)
	if err != nil {
		return
	}
	err = e.initColumnsValue()
	if err != nil {
		return
	}
	createSqlStr := tableCreateArgs2SqlStr(e.columns)

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(e.tableName + " ")
	sb.WriteString(createSqlStr)

	sql := e.dialect.ToDialectSql(sb.String())

	return e.db.exec(sql, e.columnValues...)
}

func (e EngineTable) CreateOrUpdate(v interface{}) OrmTableCreate {
	err := e.setDest(v)
	if err != nil {
		e.context.err = err
		return OrmTableCreate{base: e}
	}
	err = e.initColumnsValue()
	if err != nil {
		e.context.err = err
		return OrmTableCreate{base: e}
	}
	return OrmTableCreate{base: e}
}

func (orm OrmTableCreate) ById() (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}
	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues
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
	rows, err := orm.base.db.query(sb.String(), idValue)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(tableUpdateArgs2SqlStr(c))
		sb.WriteString(" WHERE id = ? ")
		cv = append(cv, idValue)

		return orm.base.db.exec(sb.String(), cv...)
	}
	columnSqlStr := tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return orm.base.db.exec(sb.String(), cv...)
}

func (orm OrmTableCreate) ByModel(v interface{}) (int64, error) {
	base := orm.base

	if err := base.context.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkValidStruct(va)
	if err != nil {
		return 0, err
	}
	tableName := base.tableName
	c := base.columns
	cv := base.columnValues
	config := base.lorm
	columns, values, err := getStructMappingColumnsValueNotNull(va, config)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		panic(err)
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns, config.LogicDeleteNoSql)
	var sb strings.Builder
	sb.WriteString("SELECT 1 ")
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(whereArgs2SqlStr)
	rows, err := base.db.query(sb.String(), values...)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(tableUpdateArgs2SqlStr(c))
		sb.WriteString(whereArgs2SqlStr)
		cv = append(cv, values...)

		return base.db.exec(sb.String(), cv...)
	}
	columnSqlStr := tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return base.db.exec(sb.String(), cv...)
}

func (orm OrmTableCreate) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}
	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

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
	rows, err := orm.base.db.query(sb.String(), args...)
	if err != nil {
		return 0, err
	}
	//update
	if rows.Next() {
		sb.Reset()
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(tableUpdateArgs2SqlStr(c))
		sb.WriteString(whereSql)
		cv = append(cv, args)

		return orm.base.db.exec(sb.String(), cv...)
	}
	columnSqlStr := tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return orm.base.db.exec(sb.String(), cv...)
}

//delete
func (e EngineTable) Delete(v interface{}) OrmTableDelete {
	err := e.setDest(v)
	if err != nil {
		e.context.err = err
		return OrmTableDelete{base: e}
	}
	return OrmTableDelete{base: e}
}

func (orm OrmTableDelete) ById(v ...interface{}) (int64, error) {
	base := orm.base

	if err := base.context.err; err != nil {
		return 0, err
	}

	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}

	base.initPrimaryKeyName()
	tableName := base.tableName
	idNames := base.primaryKeyNames
	fmt.Println(idNames)
	logicDeleteSetSql := base.lorm.LogicDeleteSetSql

	var sb strings.Builder
	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	if logicDeleteSetSql == lgSql {
		sb.WriteString("DELETE FROM ")
		sb.WriteString(tableName)
		sb.WriteString("WHERE ")
		//sb.WriteString(idName)
		sb.WriteString(" = ? ")
	} else {
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(lgSql)
		sb.WriteString("WHERE ")
		//sb.WriteString(idName)
		sb.WriteString(" = ? ")
	}

	return base.db.exec(sb.String(), v)
}

func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	base := orm.base
	if err := base.context.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkValidStruct(va)
	if err != nil {
		return 0, err
	}
	config := base.lorm

	columns, values, err := getStructMappingColumnsValueNotNull(va, config)
	if err != nil {
		return 0, err
	}
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns, config.LogicDeleteNoSql)
	var sb strings.Builder
	sb.WriteString("DELETE ")
	sb.WriteString(" FROM ")
	sb.WriteString(base.tableName)
	sb.WriteString(whereArgs2SqlStr)

	return base.db.exec(sb.String(), values...)
}

func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	if w == nil {
		return 0, nil
	}
	wheres := w.context.wheres
	args := w.context.args

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	sb.WriteString(orm.base.tableName)
	sb.WriteString(" WHERE ")
	for i, where := range wheres {
		if i == 0 {
			sb.WriteString(where)
			continue
		}
		sb.WriteString(" AND " + where)
	}

	return orm.base.db.exec(sb.String(), args)
}

//update
func (e EngineTable) Update(v interface{}) OrmTableUpdate {
	if e.context.err != nil {
		return OrmTableUpdate{base: e}
	}
	err := e.setDest(v)
	if err != nil {
		e.context.err = err
		return OrmTableUpdate{base: e}
	}
	err = e.initColumnsValue()
	if err != nil {
		e.context.err = err
		return OrmTableUpdate{base: e}
	}
	return OrmTableUpdate{base: e}
}

func (orm OrmTableUpdate) ById(v interface{}) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}

	orm.base.initPrimaryKeyName()

	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))
	sb.WriteString(" WHERE ")
	//sb.WriteString(orm.base.primaryKeyNames)
	sb.WriteString(" = ? ")
	cv = append(cv, v)
	return orm.base.db.exec(sb.String(), cv...)
}

func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	base := orm.base

	if err := base.context.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkValidStruct(va)
	if err != nil {
		return 0, err
	}

	tableName := base.tableName
	c := base.columns
	cv := base.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))
	config := base.lorm
	columns, values, err := getStructMappingColumnsValueNotNull(va, config)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		return 0, err
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns, config.LogicDeleteNoSql)
	sb.WriteString(whereArgs2SqlStr)

	cv = append(cv, values...)

	return base.db.exec(sb.String(), cv...)
}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	if w == nil {
		return 0, nil
	}
	wheres := w.context.wheres
	args := w.context.args

	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))
	sb.WriteString(" WHERE ")
	for i, where := range wheres {
		if i == 0 {
			sb.WriteString(where)
			continue
		}
		sb.WriteString(" AND " + where)
	}

	cv = append(cv, args...)

	return orm.base.db.exec(sb.String(), cv...)
}

//select
func (e EngineTable) Select(v interface{}) OrmTableSelect {
	err := e.setDest(v)
	if err != nil {
		e.context.err = err
		return OrmTableSelect{base: e}
	}

	return OrmTableSelect{base: e}
}

func (orm OrmTableSelect) ById(v ...interface{}) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}
	orm.base.initColumns()
	orm.base.initPrimaryKeyName()
	tableName := orm.base.tableName
	c := orm.base.columns

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
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	tableName := orm.base.tableName
	c := orm.base.columns

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

	return orm.base.queryLn(sb.String(), orm.base.dest)
}

func (orm OrmTableSelectWhere) getList() (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	tableName := orm.base.tableName
	c := orm.base.columns

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

	return orm.base.queryLn(sb.String(), orm.base.dest)
}

func (orm OrmTableSelect) ByModel(v interface{}) (int64, error) {
	base := orm.base

	if err := base.context.err; err != nil {
		return 0, err
	}
	va := reflect.ValueOf(v)
	err := checkValidStruct(va)
	if err != nil {
		return 0, err
	}
	base.initColumns()
	base.initPrimaryKeyName()

	tableName := base.tableName
	c := base.columns
	config := base.lorm
	columns, values, err := getStructMappingColumnsValueNotNull(va, config)
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
	sb.WriteString(tableWhereArgs2SqlStr(columns, config.LogicDeleteNoSql))

	return base.queryLn(sb.String(), values...)
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	if w == nil {
		return 0, errors.New("table select where can't nil")
	}
	orm.base.initColumns()
	orm.base.initPrimaryKeyName()

	wheres := w.context.wheres
	args := w.context.args

	tableName := orm.base.tableName
	c := orm.base.columns

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
	e.primaryKeyNames = e.lorm.primaryKeys(e.tableName, e.destBaseValue)
}

func (e *EngineTable) initTableName() error {
	tableName, err := e.lorm.tableName(e.destBaseValue)
	if err != nil {
		return err
	}
	e.tableName = tableName
	return nil
}


//获取struct对应的字段名 和 其值   有效部分
func (e *EngineTable) initColumnsValue() error {
	config := e.ormConfig

	columns, values, err := getStructMappingColumnsValueNotNull(e.destBaseValue, config)
	if err != nil {
		return err
	}
	e.columns = columns
	e.columnValues = values
	return nil
}

//获取struct对应的字段名 有效部分
func (e *EngineTable) initColumns() {
	dest := e.dest
	typ := reflect.TypeOf(dest)
	OrmCore, err := baseStructTypePtr(typ)
	if err != nil {
		e.context.err = err
		return
	}

	config := e.lorm

	cMap := make(map[string]int)

	numField := OrmCore.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := OrmCore.Field(i)
		name := field.Name
		if name == "ID" {
			cMap["id"] = i
			num++
			if len(cMap) < num {
				e.context.err = errors.New("字段:: id  error")
				return
			}
			continue
		}

		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}
		name = utils.Camel2Case(name)

		if tag := field.Tag.Get("lorm"); tag == "-" {
			continue
		}

		if tag := field.Tag.Get("db"); tag != "" {
			name = tag
			cMap[name] = i
			num++
			if len(cMap) < num {
				e.context.err = errors.New("字段::" + "error")
				return
			}
			continue
		}

		fieldNamePrefix := config.FieldNamePrefix
		if fieldNamePrefix != "" {
			cMap[fieldNamePrefix+name] = i
			num++
			if len(cMap) < num {
				e.context.err = errors.New("字段::" + "error")
				return
			}
			continue
		}

		cMap[name] = i
		num++
		if len(cMap) < num {
			e.context.err = errors.New("字段::" + "error")
			return
		}
	}
	arr := make([]string, len(cMap))

	var i = 0
	for s := range cMap {
		arr[i] = s
		i++
	}
	e.columns = arr
}
