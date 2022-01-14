package lorm

import (
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
)

//------------------------------------Insert--------------------------------------------

// Insert
//1.ptr
//2.comp-struct
func (tx Tx) Insert(v interface{}) (num int64, err error) {
	tx.setTargetDest(v)
	tx.initColumnsValue()
	if tx.ctx.err != nil {
		return 0, tx.ctx.err
	}
	sqlStr := tx.ctx.tableInsertGen()

	if tx.ctx.destIsPtr {
		sqlStr += " RETURNING id"
		return tx.query(sqlStr, tx.ctx.columnValues...)
	}

	return tx.dialect.exec(sqlStr, tx.ctx.columnValues...)
}

// InsertOrUpdate
//1.ptr
//2.comp-struct
func (tx Tx) InsertOrUpdate(v interface{}) OrmTableCreateTx {
	tx.setTargetDest(v)
	tx.initColumnsValue()
	return OrmTableCreateTx{base: tx}
}

// ByPrimaryKey
//ptr
//single / comp复合主键
func (orm OrmTableCreateTx) ByPrimaryKey() (int64, error) {
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
//ptr-comp
func (orm OrmTableCreateTx) ByUnique(fs types.Fields) (int64, error) {
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

//------------------------------------Delete--------------------------------------------

// Delete
//delete
func (tx Tx) Delete(v interface{}) OrmTableDeleteTx {
	tx.setTargetDest2TableName(v)
	return OrmTableDeleteTx{base: tx}
}

// ByPrimaryKey
//[]
//single -> 单主键
//comp -> 复合主键
func (orm OrmTableDeleteTx) ByPrimaryKey(v ...interface{}) (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doDel()
}

// ByModel
//ptr
//comp,只能一个comp-struct
func (orm OrmTableDeleteTx) ByModel(v interface{}) (int64, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doDel()
}

func (orm OrmTableDeleteTx) ByWhere(w *WhereBuilder) (int64, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doDel()
}

//------------------------------------Update--------------------------------------------

func (tx Tx) Update(v interface{}) OrmTableUpdateTx {
	tx.setTargetDest(v)
	tx.initColumnsValue()
	return OrmTableUpdateTx{base: tx}
}

func (orm OrmTableUpdateTx) ByPrimaryKey() (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initSelfPrimaryKeyValues()
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdateTx) ByModel(v interface{}) (int64, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdateTx) ByWhere(w *WhereBuilder) (int64, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doUpdate()
}

//------------------------------------Select--------------------------------------------

func (tx Tx) Select(v interface{}) OrmTableSelectTx {
	tx.setTargetDest2TableName(v)
	return OrmTableSelectTx{base: tx}
}

func (orm OrmTableSelectTx) ByPrimaryKey(v ...interface{}) OrmTableSelectWhereTx {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	return OrmTableSelectWhereTx{base: orm.base}
}

// ByModel
//ptr-comp
func (orm OrmTableSelectTx) ByModel(v interface{}) OrmTableSelectWhereTx {
	orm.base.initByModel(v)
	return OrmTableSelectWhereTx{base: orm.base}
}

func (orm OrmTableSelectTx) ByWhere(w *WhereBuilder) OrmTableSelectWhereTx {
	orm.base.initByWhere(w)
	return OrmTableSelectWhereTx{base: orm.base}
}

func (orm OrmTableSelectWhereTx) OrderBy(name string, condition ...bool) OrmTableSelectWhereTx {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhereTx{base: orm.base}
		}
	}

	orm.base.orderByTokens = append(orm.base.orderByTokens, name)
	return OrmTableSelectWhereTx{base: orm.base}
}

func (orm OrmTableSelectWhereTx) OrderDescBy(name string, condition ...bool) OrmTableSelectWhereTx {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhereTx{base: orm.base}
		}
	}

	orm.base.orderByTokens = append(orm.base.orderByTokens, name+" desc")
	return OrmTableSelectWhereTx{base: orm.base}
}

func (orm OrmTableSelectWhereTx) Limit(num int64, condition ...bool) OrmTableSelectWhereTx {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhereTx{base: orm.base}
		}
	}
	orm.base.limit = num
	return OrmTableSelectWhereTx{base: orm.base}
}

func (orm OrmTableSelectWhereTx) Offset(num int64, condition ...bool) OrmTableSelectWhereTx {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhereTx{base: orm.base}
		}
	}
	orm.base.offset = num
	return OrmTableSelectWhereTx{base: orm.base}
}

func (orm OrmTableSelectWhereTx) ScanFirst(v interface{}) (int64, error) {
	orm.Limit(1)
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelectWhereTx) ScanOne(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelectWhereTx) ScanList(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestList(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

//------------------------------------has--------------------------------------------

func (tx Tx) Has(v interface{}) OrmTableHasTx {
	tx.setTargetDest2TableName(v)
	return OrmTableHasTx{base: tx}
}

func (orm OrmTableHasTx) ByPrimaryKey(v ...interface{}) (bool, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doHas()
}

// ByModel
//v0.6
//ptr-comp
func (orm OrmTableHasTx) ByModel(v interface{}) (bool, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doHas()
}

func (orm OrmTableHasTx) ByWhere(w *WhereBuilder) (bool, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doHas()
}

//-----------------------------------------------------------

type OrmTableCreateTx struct {
	base Tx
}

type OrmTableSelectTx struct {
	base Tx

	query string
	args  []interface{}
}

type OrmTableHasTx struct {
	base Tx

	query string
	args  []interface{}
}

type OrmTableSelectWhereTx struct {
	base Tx
}

type OrmTableUpdateTx struct {
	base Tx
}

type OrmTableDeleteTx struct {
	base Tx
}
