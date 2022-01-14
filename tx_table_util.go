package lorm

import (
	"bytes"
	"database/sql"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"strings"
)

//update
func (tx Tx) doUpdate() (int64, error) {
	if err := tx.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := tx.ctx
	tableName := ctx.tableName
	cs := ctx.columns

	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
	bb.Write(tx.genWhereSqlByToken())

	return tx.dialect.exec(bb.String(), append(ctx.columnValues, tx.args...)...)

}

//del
func (tx Tx) doDel() (int64, error) {
	if err := tx.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer
	tableName := tx.ctx.tableName
	where := tx.genWhereSqlByToken()

	if tx.ctx.ormConf.LogicDeleteSetSql == "" {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.Write(where)
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(tx.ctx.ormConf.LogicDeleteSetSql)
		bb.Write(where)
	}

	return tx.dialect.exec(bb.String(), tx.args...)
}

//update
func (tx Tx) doSelect() (int64, error) {
	if err := tx.ctx.err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := tx.ctx
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
	bb.Write(tx.genWhereSqlByToken())

	return tx.query(bb.String(), tx.args...)
}

//has
func (tx Tx) doHas() (bool, error) {
	if err := tx.ctx.err; err != nil {
		return false, err
	}
	var bb bytes.Buffer

	ctx := tx.ctx
	tableName := ctx.tableName

	bb.WriteString("SELECT 1 FROM ")
	bb.WriteString(tableName)
	bb.Write(tx.genWhereSqlByToken())
	bb.WriteString("LIMIT 1")
	rows, err := tx.dialect.query(bb.String(), tx.args...)
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
func (tx *Tx) initByPrimaryKey() {
	ctx := tx.ctx
	if err := ctx.err; err != nil {
		return
	}
	pkNum := len(ctx.primaryKeyValues)
	tx.whereTokens = append(tx.whereTokens, utils.GenwhereTokenOfBatch(ctx.primaryKeyNames, pkNum))

	for _, value := range ctx.primaryKeyValues {
		tx.args = append(tx.args, value...)
	}
}

//根据 byModel 生成的where token
func (tx *Tx) initByModel(v interface{}) {
	if err := tx.ctx.err; err != nil {
		return
	}
	if v == nil {
		tx.ctx.err = errors.New("model is nil")
		return
	}

	columns, values, err := getCompCV(v, tx.ctx.ormConf)
	if err != nil {
		tx.ctx.err = err
		return
	}
	tx.whereTokens = append(tx.whereTokens, utils.GenwhereToken(columns)...)
	tx.args = append(tx.args, values...)
}

//根据 byWhere 生成的where token
func (tx *Tx) initByWhere(w *WhereBuilder) {
	if err := tx.ctx.err; err != nil {
		return
	}
	if w == nil {
		tx.ctx.err = errors.New("ByWhere is nil")
		return
	}

	args := w.context.args
	wheres := w.context.wheres

	tx.whereTokens = append(tx.whereTokens, wheres...)
	tx.args = append(tx.args, args...)
}

//init 逻辑删除、租户
func (tx *Tx) initExtra() {
	if err := tx.ctx.err; err != nil {
		return
	}

	if tx.ctx.ormConf.LogicDeleteYesSql != "" {
		tx.whereTokens = append(tx.whereTokens, tx.ctx.ormConf.LogicDeleteYesSql)
	}

	if tx.ctx.ormConf.TenantIdFieldName != "" {
		tx.whereTokens = append(tx.whereTokens, tx.ctx.ormConf.TenantIdFieldName)
		tx.args = append(tx.args, tx.ctx.ormConf.TenantIdValueFun())
	}

	var buf bytes.Buffer
	buf.Write(tx.extraWhereSql)

	if len(tx.orderByTokens) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(tx.orderByTokens, ","))
	}
	if tx.limit > 0 {
		buf.WriteString(" LIMIT ? ")
		tx.args = append(tx.args, tx.limit)
	}
	if tx.offset > 0 {
		buf.WriteString(" OFFSET ? ")
		tx.args = append(tx.args, tx.offset)
	}
	tx.extraWhereSql = buf.Bytes()

}

//初始化逻辑删除
func (tx *Tx) initLgDel() {
	if err := tx.ctx.err; err != nil {
		return
	}
	if tx.ctx.ormConf.LogicDeleteYesSql != "" {
		tx.extraWhereSql = []byte(tx.ctx.ormConf.LogicDeleteYesSql)
	}
}

//-------------------------------target------------------------

//*.comp
//target scanDest 一个comp-struct
func (tx *Tx) setTargetDest(v interface{}) {
	if tx.ctx.err != nil {
		return
	}
	tx.ctx.initTargetDest(v)
	tx.ctx.checkTargetDestField()
	tx.initTableName()
}

func (tx *Tx) setTargetDest2TableName(v interface{}) {
	if tx.ctx.err != nil {
		return
	}
	tx.ctx.initTargetDest2TableName(v)
	tx.initTableName()
}

//初始化主键
func (tx *Tx) initPrimaryKeyName() {
	if tx.ctx.err != nil {
		return
	}
	tx.ctx.primaryKeyNames = tx.ctx.ormConf.primaryKeys(tx.ctx.tableName)
}

//初始化 表名
func (tx *Tx) initTableName() {
	if tx.ctx.err != nil {
		return
	}
	if tx.ctx.tableName != "" {
		return
	}
	tableName, err := tx.ctx.ormConf.tableName(tx.ctx.destBaseType)
	if err != nil {
		tx.ctx.err = err
		return
	}
	tx.ctx.tableName = tableName
}

//获取struct对应的字段名 和 其值，
//slice为全部，一个为非nil字段。
func (tx *Tx) initColumnsValue() {
	if tx.ctx.err != nil {
		return
	}
	columns, valuess, err := tx.ctx.ormConf.getCompColumnsValueNoNil(tx.ctx.destValue)
	if err != nil {
		tx.ctx.err = err
		return
	}
	tx.ctx.columns = columns
	tx.ctx.columnValues = valuess
	return
}

//获取struct对应的字段名 有效部分
func (tx *Tx) initColumns() {
	if tx.ctx.err != nil {
		return
	}

	columns, err := tx.ctx.ormConf.initColumns(tx.ctx.scanDestBaseType)
	if err != nil {
		tx.ctx.err = err
		return
	}
	tx.ctx.columns = columns
}

//------------------------query--------------------------
func (tx Tx) query(query string, args ...interface{}) (int64, error) {
	rows, err := tx.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	if tx.ctx.scanIsSlice {
		return tx.ctx.Scan(rows)
	}
	return tx.ctx.ScanLn(rows)
}

func (tx Tx) queryBatch(query string, args [][]interface{}) (int64, error) {
	stmt, err := tx.dialect.queryBatch(query)
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
	return tx.ctx.ScanBatch(rowss)
}

//------------------------gen-sql---------------------------

//根据whereTokens生成的where sql
func (tx Tx) genWhereSqlByToken() []byte {
	if len(tx.whereTokens) == 0 && tx.extraWhereSql == nil {
		return nil
	}
	var buf bytes.Buffer
	buf.WriteString(" WHERE ")
	for i, token := range tx.whereTokens {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(token)
	}
	buf.Write(tx.extraWhereSql)
	return buf.Bytes()
}
