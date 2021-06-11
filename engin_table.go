package lorm

import (
	"database/sql/driver"
	"errors"
	"github.com/lontten/lorm/utils"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type EngineTable struct {
	context OrmContext

	db DBer

	idName string
	//当前表名
	tableName string
	//当前struct对象
	dest interface{}

	columns      []string
	columnValues []interface{}
}

func (engine EngineTable) queryLn(query string, args ...interface{}) (int64, error) {
	rows, err := engine.db.query(query, args...)
	if err != nil {
		return 0, err
	}

	return StructScanLn(rows, engine.dest)
}

func (engine *EngineTable) setDest(v interface{}) error {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return err
	}
	engine.dest = v
	engine.initTableName()
	return nil
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

	return e.db.exec(sb.String(), e.columnValues...)
}

func (e EngineTable) CreateOrUpdate(v interface{}) OrmTableCreate {
	err := e.setDest(v)
	if err != nil {
		panic(err)
	}
	err = e.initColumnsValue()
	if err != nil {
		panic(err)
	}
	return OrmTableCreate{
		base: e,
	}
}

func (orm OrmTableCreate) ById() (int64, error) {
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
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}
	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues
	config := orm.base.db.OrmConfig()
	columns, values, err := getStructMappingColumnsValue(v, config)
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
	rows, err := orm.base.db.query(sb.String(), values...)
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

		return orm.base.db.exec(sb.String(), cv...)
	}
	columnSqlStr := tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return orm.base.db.exec(sb.String(), cv...)
}

func (orm OrmTableCreate) ByWhere(w *WhereBuilder) (int64, error) {
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
func (engine EngineTable) Delete(v interface{}) OrmTableDelete {
	err := engine.setDest(v)
	if err != nil {
		panic(err)
	}
	return OrmTableDelete{base: engine}
}

func (orm OrmTableDelete) ById(v interface{}) (int64, error) {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}

	orm.base.initIdName()
	tableName := orm.base.tableName
	idName := orm.base.idName
	logicDeleteSetSql := orm.base.db.OrmConfig().LogicDeleteSetSql

	var sb strings.Builder
	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	if logicDeleteSetSql == lgSql {
		sb.WriteString("DELETE FROM ")
		sb.WriteString(tableName)
		sb.WriteString("WHERE ")
		sb.WriteString(idName)
		sb.WriteString(" = ? ")
	} else {
		sb.WriteString("UPDATE ")
		sb.WriteString(tableName)
		sb.WriteString(" SET ")
		sb.WriteString(lgSql)
		sb.WriteString("WHERE ")
		sb.WriteString(idName)
		sb.WriteString(" = ? ")
	}

	return orm.base.db.exec(sb.String(), v)
}

func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}
	config := orm.base.db.OrmConfig()

	columns, values, err := getStructMappingColumnsValue(v, config)
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
	sb.WriteString(orm.base.tableName)
	sb.WriteString(whereArgs2SqlStr)

	return orm.base.db.exec(sb.String(), values...)
}

func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
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
	err := e.setDest(v)
	if err != nil {
		panic(err)
	}
	err = e.initColumnsValue()
	if err != nil {
		panic(err)
	}
	return OrmTableUpdate{base: e}
}

func (orm OrmTableUpdate) ById(v interface{}) (int64, error) {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}

	orm.base.initIdName()

	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))
	sb.WriteString(" WHERE ")
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ? ")
	cv = append(cv, v)
	return orm.base.db.exec(sb.String(), cv...)
}

func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}

	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString(" UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))
	config := orm.base.db.OrmConfig()
	columns, values, err := getStructMappingColumnsValue(v, config)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		return 0, err
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns, config.LogicDeleteNoSql)
	sb.WriteString(whereArgs2SqlStr)

	cv = append(cv, values...)

	return orm.base.db.exec(sb.String(), cv...)
}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
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
		panic(err)
	}

	return OrmTableSelect{base: e}
}

func (orm OrmTableSelect) ById(v interface{}) (int64, error) {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}
	orm.base.initColumns()
	orm.base.initIdName()
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
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ? ")

	return orm.base.queryLn(sb.String(), v)
}

func (orm OrmTableSelectWhere) getOne() (int64, error) {
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
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ? ")

	return orm.base.queryLn(sb.String(), orm.base.dest)
}

func (orm OrmTableSelectWhere) getList() (int64, error) {
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
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ? ")

	return orm.base.queryLn(sb.String(), orm.base.dest)
}

func (orm OrmTableSelect) ByModel(v interface{}) (int64, error) {
	err := checkValidStruct(reflect.ValueOf(v))
	if err != nil {
		return 0, err
	}
	orm.base.initColumns()
	orm.base.initIdName()

	tableName := orm.base.tableName
	c := orm.base.columns
	config := orm.base.db.OrmConfig()
	columns, values, err := getStructMappingColumnsValue(v, config)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		panic(err)
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

	return orm.base.queryLn(sb.String(), values...)
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) (int64, error) {
	if w == nil {
		return 0, errors.New("table select where can't nil")
	}
	orm.base.initColumns()
	orm.base.initIdName()

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
func (e *EngineTable) initIdName() {
	idNameFun := e.db.OrmConfig().IdNameFun
	idName := e.db.OrmConfig().IdName
	if idNameFun != nil {
		e.idName = idNameFun(e.tableName, e.dest)
	}
	if idName != "" {
		e.idName = idName
	}
	e.idName = "id"
}

func (e *EngineTable) initTableName() {
	tableName, err := getStructTableName(e.dest, e.db.OrmConfig())
	if err != nil {
		Log.Println("tableName",err)
		panic(err)
	}
	e.tableName = tableName
}

//获取struct对应的字段名 和 其值   有效部分
func (e *EngineTable) initColumnsValue() error {
	dest := e.dest
	config := e.db.OrmConfig()

	t := reflect.TypeOf(dest)
	base, err := baseStructType(t)
	if err != nil {
		return err
	}

	mappingColumns, err := getStructMappingColumns(base, config)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(dest)
	structValue, err := baseStructValue(v)
	if err != nil {
		return err
	}

	for column, i := range mappingColumns {
		field := structValue.Field(i)
		indirect := reflect.Indirect(field)
		if !field.IsNil() {
			e.columns = append(e.columns, column)
			value, err := indirect.Interface().(driver.Valuer).Value()
			if err != nil {
				return err
			}
			e.columnValues = append(e.columnValues, value)
		}
	}
	return nil
}

//获取struct对应的字段名 有效部分
func (e *EngineTable) initColumns() {
	dest := e.dest
	typ := reflect.TypeOf(dest)
	base, err := baseStructType(typ)
	panicErr(err)

	config := e.db.OrmConfig()

	cMap := make(map[string]int)

	numField := base.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := base.Field(i)
		name := field.Name

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
				panic(errors.New("字段::" + "error"))
			}
			continue
		}

		fieldNamePrefix := config.FieldNamePrefix
		if fieldNamePrefix != "" {
			cMap[fieldNamePrefix+name] = i
			num++
			if len(cMap) < num {
				panic(errors.New("字段::" + "error"))
			}
			continue
		}

		cMap[name] = i
		num++
		if len(cMap) < num {
			panic(errors.New("字段::" + "error"))
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
