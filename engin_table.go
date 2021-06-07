package lorm

import (
	"errors"
	"github.com/lontten/lorm/utils"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type EngineTable struct {
	db *DbPool

	idName string
	//当前表名
	tableName string
	//当前struct对象
	dest interface{}

	columns      []string
	columnValues []interface{}
}

func (e *EngineTable) query(query string, args ...interface{}) (int64, error) {
	log.Println(query, args)
	rows, err := e.db.db.Query(query, args...)
	if err != nil {
		return 0, err
	}

	return StructScan(rows, e.dest)
}

func (e *EngineTable) setDest(v interface{}) {
	e.dest = v
	e.initTableName()
}

type OrmTableCreate struct {
	base *EngineTable
}

type OrmTableSelect struct {
	base *EngineTable

	query string
	args  []interface{}
}

type OrmTableSelectWhere struct {
	base *EngineTable
}

type OrmTableUpdate struct {
	base *EngineTable
}

type OrmTableDelete struct {
	base *EngineTable
}


//create
func (engine *EngineTable) Create(v interface{}) (int64, error) {
	engine.setDest(v)
	engine.initColumnsValue()

	createSqlStr := tableCreateArgs2SqlStr(engine.columns)

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(engine.tableName + " ")
	sb.WriteString(createSqlStr)

	return engine.db.Exec(sb.String(), engine.columnValues)
}

func (engine *EngineTable) CreateOrUpdate(v interface{}) *OrmTableCreate {
	engine.setDest(v)
	engine.initColumnsValue()
	return &OrmTableCreate{
		base: engine,
	}
}

func (orm *OrmTableCreate) ByModel(v interface{}) (int64, error) {
	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues
	columns, values, err := getStructMappingColumnsValue(v, orm.base.db.ormConfig)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		panic(err)
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns)
	var sb strings.Builder
	sb.WriteString("SELECT 1 ")
	sb.WriteString(" FROM ")
	sb.WriteString(tableName)
	sb.WriteString(whereArgs2SqlStr)
	log.Println(sb.String(), values)
	rows, err := orm.base.db.db.Query(sb.String(), values...)
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
		cv = append(cv, values)

		return orm.base.db.Exec(sb.String(), cv...)
	}
	columnSqlStr := tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return orm.base.db.Exec(sb.String(), cv...)
}

func (orm *OrmTableCreate) ByWhere(w *WhereBuilder) (int64, error) {
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
	rows, err := orm.base.db.db.Query(sb.String(), args...)
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

		return orm.base.db.Exec(sb.String(), cv...)
	}
	columnSqlStr := tableCreateArgs2SqlStr(c)

	sb.Reset()
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	sb.WriteString(columnSqlStr)

	return orm.base.db.Exec(sb.String(), cv...)
}

//delete
func (engine *EngineTable) Delete(v interface{}) *OrmTableDelete {
	engine.setDest(v)
	return &OrmTableDelete{base: engine}
}

func (orm *OrmTableDelete) ById() (int64, error) {
	orm.base.initIdName()

	var sb strings.Builder

	sb.WriteString("DELETE FROM ")
	sb.WriteString(orm.base.tableName)
	sb.WriteString("WHERE ")
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ?")

	return orm.base.db.Exec(sb.String(), orm.base.dest)
}

func (orm *OrmTableDelete) ByModel() (int64, error) {
	columns, values, err := getStructMappingColumnsValue(orm.base.dest, orm.base.db.ormConfig)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		panic(err)
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns)
	var sb strings.Builder
	sb.WriteString("DELETE ")
	sb.WriteString(" FROM ")
	sb.WriteString(orm.base.tableName)
	sb.WriteString(whereArgs2SqlStr)

	return orm.base.db.Exec(sb.String(), values)
}

func (orm *OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
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

	return orm.base.db.Exec(sb.String(), args)
}

//update
func (engine *EngineTable) Update(v interface{}) *OrmTableUpdate {
	engine.setDest(v)
	engine.initColumnsValue()
	return &OrmTableUpdate{base: engine}
}

func (orm *OrmTableUpdate) ById(v interface{}) (int64, error) {
	orm.base.initIdName()

	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))
	sb.WriteString("WHERE ")
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ?")
	cv = append(cv, v)

	return orm.base.db.Exec(sb.String(), cv...)
}

func (orm *OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(tableName)
	sb.WriteString(" SET ")
	sb.WriteString(tableUpdateArgs2SqlStr(c))

	columns, values, err := getStructMappingColumnsValue(v, orm.base.db.ormConfig)
	if len(columns) < 1 {
		return 0, errors.New("where model valid field need ")
	}
	if err != nil {
		panic(err)
	}
	whereArgs2SqlStr := tableWhereArgs2SqlStr(columns)
	sb.WriteString(whereArgs2SqlStr)

	cv = append(cv, values)

	return orm.base.db.Exec(sb.String(), cv...)
}

func (orm *OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	if w == nil {
		return 0, nil
	}
	wheres := w.context.wheres
	args := w.context.args

	tableName := orm.base.tableName
	c := orm.base.columns
	cv := orm.base.columnValues

	var sb strings.Builder
	sb.WriteString("UPDATE ")
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

	cv = append(cv, args)

	return orm.base.db.Exec(sb.String(), cv...)
}

//select
func (engine *EngineTable) Select(v interface{}) *OrmTableSelect {
	engine.setDest(v)

	return &OrmTableSelect{base: engine}
}

func (orm *OrmTableSelect) ById(v interface{}) (int64, error) {

	orm.base.initColumnsValue()
	orm.base.initIdName()
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
	sb.WriteString(orm.base.idName)
	sb.WriteString(" = ?")

	log.Println(sb.String(), v)
	rows, err := orm.base.db.db.Query(sb.String(), v)
	if err != nil {
		return 0, err
	}
	return StructScanLn(rows, orm.base.dest)
}

func (orm *OrmTableSelectWhere) getOne() (int64, error) {
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
	sb.WriteString(" = ?")

	return orm.base.query(sb.String(), orm.base.dest)
}

func (orm *OrmTableSelectWhere) getList() (int64, error) {
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
	sb.WriteString(" = ?")

	return orm.base.query(sb.String(), orm.base.dest)
}

func (orm *OrmTableSelect) ByModel(v interface{}) (int64, error) {
	tableName := orm.base.tableName
	c := orm.base.columns

	columns, values, err := getStructMappingColumnsValue(v, orm.base.db.ormConfig)
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
	sb.WriteString(tableWhereArgs2SqlStr(columns))

	return orm.base.query(sb.String(), values)
}

func (orm *OrmTableSelect) ByWhere(w *WhereBuilder) (int64, error) {
	if w == nil {
		return 0, errors.New("table select where can't nil")
	}
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

	return orm.base.query(sb.String(), args)
}

//init
func (e *EngineTable) initIdName() {
	idNameFun := e.db.ormConfig.IdNameFun
	idName := e.db.ormConfig.IdName
	if idNameFun != nil {
		e.idName = idNameFun(e.tableName, e.dest)
	}
	if idName != "" {
		e.idName = idName
	}
	e.idName = "id"
}

func (e *EngineTable) initTableName() {
	typ := reflect.TypeOf(e.dest)
	base, err := baseStructType(typ)
	if err != nil {
		panic(err)
	}
	name := base.String()
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[index+1:]
	}
	name = utils.Camel2Case(name)

	tableNameFun := e.db.ormConfig.TableNameFun
	if tableNameFun != nil {
		e.tableName = tableNameFun(name, base)
	}
	tableNamePrefix := e.db.ormConfig.TableNamePrefix
	e.tableName = tableNamePrefix + name
}

//获取struct对应的字段名 和 其值   有效部分
func (e *EngineTable) initColumnsValue() {
	dest := e.dest
	config := e.db.ormConfig

	t := reflect.TypeOf(dest)
	base, err := baseStructType(t)
	if err != nil {
		return
	}

	mappingColumns, err := getStructMappingColumns(base, config)
	if err != nil {
		return
	}

	v := reflect.ValueOf(dest)
	structValue, err := baseStructValue(v)
	if err != nil {
		return
	}

	for column, i := range mappingColumns {
		field := structValue.Field(i)
		indirect := reflect.Indirect(field)
		if !field.IsNil() {
			e.columns = append(e.columns, column)
			e.columnValues = append(e.columnValues, indirect.Interface())
		}
	}
	return
}

//获取struct对应的字段名 有效部分
func (e *EngineTable) initColumns() {
	dest := e.dest
	typ := reflect.TypeOf(dest)
	base, err := baseStructType(typ)
	panicErr(err)

	config := e.db.ormConfig

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
