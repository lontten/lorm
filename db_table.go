package lsql

import (
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
)

//------------------------------------Insert--------------------------------------------

// Insert
//1.ptr
//2.comp-struct
func (db DB) Insert(v interface{}) (num int64, err error) {
	db.setTargetDest(v)
	db.initColumnsValue()
	if db.ctx.err != nil {
		return 0, db.ctx.err
	}
	sqlStr := db.ctx.tableInsertGen()

	if db.ctx.destIsPtr {
		sqlStr += " RETURNING id"
		return db.query(sqlStr, db.ctx.columnValues...)
	}

	return db.dialect.exec(sqlStr, db.ctx.columnValues...)
}

// InsertOrUpdate
//1.ptr
//2.comp-struct
func (db DB) InsertOrUpdate(v interface{}) OrmTableCreate {
	db.setTargetDest(v)
	db.initColumnsValue()
	return OrmTableCreate{base: db}
}

// ByPrimaryKey
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

//------------------------------------Delete--------------------------------------------

// Delete
//delete
func (db DB) Delete(v interface{}) OrmTableDelete {
	db.setTargetDest2TableName(v)
	return OrmTableDelete{base: db}
}

// ByPrimaryKey
//[]
//single -> 单主键
//comp -> 复合主键
func (orm OrmTableDelete) ByPrimaryKey(v ...interface{}) (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doDel()
}

// ByModel
//ptr
//comp,只能一个comp-struct
func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doDel()
}

func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doDel()
}

//------------------------------------Update--------------------------------------------

func (db DB) Update(v interface{}) OrmTableUpdate {
	db.setTargetDest(v)
	db.initColumnsValue()
	return OrmTableUpdate{base: db}
}

func (orm OrmTableUpdate) ByPrimaryKey() (int64, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initSelfPrimaryKeyValues()
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doUpdate()
}

//------------------------------------Select--------------------------------------------

func (db DB) Select(v interface{}) OrmTableSelect {
	db.setTargetDest2TableName(v)
	return OrmTableSelect{base: db}
}

func (orm OrmTableSelect) ByPrimaryKey(v ...interface{}) OrmTableSelectWhere {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	return OrmTableSelectWhere{base: orm.base}
}

// ByModel
//ptr-comp
func (orm OrmTableSelect) ByModel(v interface{}) OrmTableSelectWhere {
	orm.base.initByModel(v)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) OrmTableSelectWhere {
	orm.base.initByWhere(w)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) OrderBy(name string, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}

	orm.base.orderByTokens = append(orm.base.orderByTokens, name)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) OrderDescBy(name string, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}

	orm.base.orderByTokens = append(orm.base.orderByTokens, name+" desc")
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) Limit(num int64, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}
	orm.base.limit = num
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) Offset(num int64, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}
	orm.base.offset = num
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) ScanFirst(v interface{}) (int64, error) {
	orm.Limit(1)
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelectWhere) ScanOne(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelectWhere) ScanList(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestList(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

//------------------------------------has--------------------------------------------

// Select
//select
func (db DB) Has(v interface{}) OrmTableHas {
	db.setTargetDest2TableName(v)
	return OrmTableHas{base: db}
}

// ByPrimaryKey
//v0.8
func (orm OrmTableHas) ByPrimaryKey(v ...interface{}) (bool, error) {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doHas()
}

// ByModel
//v0.6
//ptr-comp
func (orm OrmTableHas) ByModel(v interface{}) (bool, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doHas()
}

func (orm OrmTableHas) ByWhere(w *WhereBuilder) (bool, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doHas()
}

//-----------------------------------------------------------

type OrmTableCreate struct {
	base DB
}

type OrmTableSelect struct {
	base DB

	query string
	args  []interface{}
}

type OrmTableHas struct {
	base DB

	query string
	args  []interface{}
}

type OrmTableSelectWhere struct {
	base DB
}

type OrmTableUpdate struct {
	base DB
}

type OrmTableDelete struct {
	base DB
}
