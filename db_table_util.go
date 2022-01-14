package lorm

import (
	"bytes"
	"database/sql"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

//update
func (db DB) doUpdate() (int64, error) {
	if err := db.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := db.ctx
	tableName := ctx.tableName
	cs := ctx.columns

	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
	bb.Write(db.genWhereSqlByToken())

	return db.dialect.exec(bb.String(), append(ctx.columnValues, db.args...)...)

}

//del
func (db DB) doDel() (int64, error) {
	if err := db.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer
	tableName := db.ctx.tableName
	where := db.genWhereSqlByToken()

	if db.ctx.conf.LogicDeleteSetSql == "" {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.Write(where)
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(db.ctx.conf.LogicDeleteSetSql)
		bb.Write(where)
	}

	return db.dialect.exec(bb.String(), db.args...)
}

//update
func (db DB) doSelect() (int64, error) {
	if err := db.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := db.ctx
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
	bb.Write(db.genWhereSqlByToken())

	return db.query(bb.String(), db.args...)
}

//has
func (db DB) doHas() (bool, error) {
	if err := db.ctx.err; err != nil {
		return false, err
	}
	var bb bytes.Buffer

	ctx := db.ctx
	tableName := ctx.tableName

	bb.WriteString("SELECT 1 FROM ")
	bb.WriteString(tableName)
	bb.Write(db.genWhereSqlByToken())
	bb.WriteString("LIMIT 1")
	rows, err := db.dialect.query(bb.String(), db.args...)
	if err != nil {
		return false, err
	}
	//update
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

//-------------------------------init------------------------

//根据 byModel 生成的where token
func (db *DB) initByPrimaryKey() {
	ctx := db.ctx
	if err := ctx.err; err != nil {
		return
	}
	pkNum := len(ctx.primaryKeyValues)
	db.whereTokens = append(db.whereTokens, utils.GenwhereTokenOfBatch(ctx.primaryKeyNames, pkNum))

	for _, value := range ctx.primaryKeyValues {
		db.args = append(db.args, value...)
	}
}

//根据 byModel 生成的where token
func (db *DB) initByModel(v interface{}) {
	if err := db.ctx.err; err != nil {
		return
	}
	if v == nil {
		db.ctx.err = errors.New("model is nil")
		return
	}

	columns, values, err := getCompCV(v, db.ctx.conf)
	if err != nil {
		db.ctx.err = err
		return
	}
	db.whereTokens = append(db.whereTokens, utils.GenwhereToken(columns)...)
	db.args = append(db.args, values...)
}

//根据 byWhere 生成的where token
func (db *DB) initByWhere(w *WhereBuilder) {
	if err := db.ctx.err; err != nil {
		return
	}
	if w == nil {
		db.ctx.err = errors.New("ByWhere is nil")
		return
	}

	args := w.context.args
	wheres := w.context.wheres

	db.whereTokens = append(db.whereTokens, wheres...)
	db.args = append(db.args, args...)
}

//init 逻辑删除、租户
func (db *DB) initExtra() {
	if err := db.ctx.err; err != nil {
		return
	}

	if db.ctx.conf.LogicDeleteYesSql != "" {
		db.whereTokens = append(db.whereTokens, db.ctx.conf.LogicDeleteYesSql)
	}

	if db.ctx.conf.TenantIdFieldName != "" {
		db.whereTokens = append(db.whereTokens, db.ctx.conf.TenantIdFieldName)
		db.args = append(db.args, db.ctx.conf.TenantIdValueFun())
	}

	var buf bytes.Buffer
	buf.Write(db.extraWhereSql)

	if len(db.orderByTokens) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(db.orderByTokens, ","))
	}
	if db.limit > 0 {
		buf.WriteString(" LIMIT ? ")
		db.args = append(db.args, db.limit)
	}
	if db.offset > 0 {
		buf.WriteString(" OFFSET ? ")
		db.args = append(db.args, db.offset)
	}
	db.extraWhereSql = buf.Bytes()

}

//初始化逻辑删除
func (db *DB) initLgDel() {
	if err := db.ctx.err; err != nil {
		return
	}
	if db.ctx.conf.LogicDeleteYesSql != "" {
		db.extraWhereSql = []byte(db.ctx.conf.LogicDeleteYesSql)
	}
}

//-------------------------------target------------------------

//*.comp
//target scanDest 一个comp-struct
func (db *DB) setTargetDest(v interface{}) {
	if db.ctx.err != nil {
		return
	}
	db.ctx.initTargetDest(v)
	db.ctx.checkTargetDestField()
	db.initTableName()
}

func (db *DB) setTargetDest2TableName(v interface{}) {
	if db.ctx.err != nil {
		return
	}
	db.ctx.initTargetDest2TableName(v)
	db.initTableName()
}

//初始化主键
func (db *DB) initPrimaryKeyName() {
	if db.ctx.err != nil {
		return
	}
	db.ctx.primaryKeyNames = db.ctx.conf.primaryKeys(db.ctx.tableName)
}

//初始化 表名
func (db *DB) initTableName() {
	if db.ctx.err != nil {
		return
	}
	if db.ctx.tableName != "" {
		return
	}
	tableName, err := db.ctx.conf.tableName(db.ctx.destBaseType)
	if err != nil {
		db.ctx.err = err
		return
	}
	db.ctx.tableName = tableName
}

//获取struct对应的字段名 和 其值，
//slice为全部，一个为非nil字段。
func (db *DB) initColumnsValue() {
	if db.ctx.err != nil {
		return
	}
	columns, valuess, err := db.ctx.conf.getCompColumnsValueNoNil(db.ctx.destValue)
	if err != nil {
		db.ctx.err = err
		return
	}
	db.ctx.columns = columns
	db.ctx.columnValues = valuess
	return
}

//获取struct对应的字段名 有效部分
func (db *DB) initColumns() {
	if db.ctx.err != nil {
		return
	}

	columns, err := db.ctx.conf.initColumns(db.ctx.scanDestBaseType)
	if err != nil {
		db.ctx.err = err
		return
	}
	db.ctx.columns = columns
}

//-------------------------utils------------------------
//获取comp 的 cv
//排除 nil 字段
func getCompCV(v interface{}, c OrmConf) ([]string, []interface{}, error) {
	value := reflect.ValueOf(v)
	_, value, err := basePtrDeepValue(value)
	if err != nil {
		return nil, nil, err
	}

	return getCompValueCV(value, c)
}

//排除 nil 字段
func getCompValueCV(v reflect.Value, c OrmConf) ([]string, []interface{}, error) {
	if !isCompType(v.Type()) {
		return nil, nil, errors.New("getvcv not comp")
	}
	err := checkCompField(v)
	if err != nil {
		return nil, nil, err
	}

	columns, values, err := c.getCompColumnsValueNoNil(v)
	if err != nil {
		return nil, nil, err
	}
	if len(columns) < 1 {
		return nil, nil, errors.New("where model valid field need ")
	}
	return columns, values, nil
}

//------------------------query--------------------------
func (db DB) query(query string, args ...interface{}) (int64, error) {
	rows, err := db.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	if db.ctx.scanIsSlice {
		return db.ctx.Scan(rows)
	}
	return db.ctx.ScanLn(rows)
}

func (db DB) queryBatch(query string, args [][]interface{}) (int64, error) {
	stmt, err := db.dialect.queryBatch(query)
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
	return db.ctx.ScanBatch(rowss)
}

//------------------------gen-sql---------------------------

//根据whereTokens生成的where sql
func (db DB) genWhereSqlByToken() []byte {
	if len(db.whereTokens) == 0 && db.extraWhereSql == nil {
		return nil
	}
	var buf bytes.Buffer
	buf.WriteString(" WHERE ")
	for i, token := range db.whereTokens {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(token)
	}
	buf.Write(db.extraWhereSql)
	return buf.Bytes()
}
