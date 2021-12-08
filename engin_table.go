package lorm

import (
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
)

//------------------------------------TableCreate--------------------------------------------

// TableCreate
//v0.8
//1.ptr
//2.comp-struct
func (e Engine) TableCreate(v interface{}) (num int64, err error) {
	e.setTargetDest(v)
	e.initColumnsValue()
	if e.ctx.err != nil {
		return 0, e.ctx.err
	}
	sqlStr := e.ctx.tableCreateGen()

	sqlStr += " RETURNING id"
	return e.query(sqlStr, e.ctx.columnValues...)
}

// TableCreateOrUpdate
//v0.6
//1.ptr
//2.comp-struct
func (e Engine) TableCreateOrUpdate(v interface{}) OrmTableCreate {
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
//v0.6
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

//------------------------------------TableDelete--------------------------------------------

// TableDelete
//delete
func (e Engine) TableDelete(v interface{}) OrmTableDelete {
	e.setTargetDest2TableName(v)
	return OrmTableDelete{base: e}
}

// ByPrimaryKey
//v0.8
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
//v0.6
//ptr
//comp,只能一个comp-struct
func (orm OrmTableDelete) ByModel(v interface{}) (int64, error) {
	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doDel()
}

// ByWhere
//v0.6
func (orm OrmTableDelete) ByWhere(w *WhereBuilder) (int64, error) {
	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doDel()
}

//------------------------------------TableUpdate--------------------------------------------

// TableUpdate
//v0.6
func (e Engine) TableUpdate(v interface{}) OrmTableUpdate {
	e.setTargetDest(v)
	e.initColumnsValue()
	return OrmTableUpdate{base: e}
}

// ByPrimaryKey
//v0.8
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

//------------------------------------TableSelect--------------------------------------------

// TableSelect
//select
func (e Engine) TableSelect(v interface{}) OrmTableSelect {
	e.setTargetDest2TableName(v)
	return OrmTableSelect{base: e}
}

// ByPrimaryKey
//v0.8
func (orm OrmTableSelect) ByPrimaryKey(v ...interface{}) OrmTableSelectWhere {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	return OrmTableSelectWhere{base: orm.base}
}

// ByModel
//v0.6
//ptr-comp
func (orm OrmTableSelect) ByModel(v interface{}) OrmTableSelectWhere {
	orm.base.initByModel(v)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) OrmTableSelectWhere {
	orm.base.initByWhere(w)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) GetFirst(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect("limit 1")
}

func (orm OrmTableSelectWhere) GetOne(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect("")
}

func (orm OrmTableSelectWhere) GetList(v interface{}) (int64, error) {
	orm.base.ctx.initScanDestList(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect("")
}

type OrmTableCreate struct {
	base Engine
}

type OrmTableSelect struct {
	base Engine

	query string
	args  []interface{}
}

type OrmTableSelectWhere struct {
	base Engine
}

type OrmTableUpdate struct {
	base Engine
}

type OrmTableDelete struct {
	base Engine
}
