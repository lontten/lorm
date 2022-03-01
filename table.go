package lorm

import (
	"github.com/pkg/errors"
	"reflect"
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

	query, args := db.dialect.exec(sqlStr, db.ctx.columnValues...)
	return db.doExec(query, args...)
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
	query, args := base.dialect.insertOrUpdateByPrimaryKey(tableName, idNames, cs, cvs...)
	return base.doExec(query, args...)
}

// ByUnique
//ptr-comp
func (orm OrmTableCreate) ByUnique(fs ...string) (int64, error) {
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
	query, args := base.dialect.insertOrUpdateByUnique(tableName, fs, cs, cvs...)
	return base.doExec(query, args...)
}

//------------------------------------Delete--------------------------------------------

func (db DB) Delete(v interface{}) OrmTableDelete {
	db.typ = dDelete
	db.baseTokens = append(db.baseTokens, baseToken{
		typ:  tDelete,
		dest: v,
		v:    reflect.ValueOf(v),
	})
	return OrmTableDelete{base: db}
}

func (orm OrmTableDelete) ByPrimaryKey(v ...interface{}) Resulter {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tPrimaryKey,
		pk:  v,
	})
	return orm.base.Do()
}

func (orm OrmTableDelete) ByModel(v interface{}) Resulter {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})
	return orm.base.Do()
}

func (orm OrmTableDelete) ByWhere(w *WhereBuilder) Resulter {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ:   tWhereBuilder,
		where: w,
	})
	return orm.base.Do()
}

//------------------------------------Update--------------------------------------------

func (db DB) Update(v interface{}) OrmTableUpdate {
	db.typ = dUpdate
	db.baseTokens = append(db.baseTokens, baseToken{
		typ:  tUpdate,
		dest: v,
		v:    reflect.ValueOf(v),
	})

	db.setTargetDest(v)
	db.initColumnsValue()
	return OrmTableUpdate{base: db}
}

func (orm OrmTableUpdate) ByPrimaryKey() (int64, error) {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tPrimaryKey,
	})

	orm.base.initPrimaryKeyName()
	orm.base.ctx.initSelfPrimaryKeyValues()
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdate) ByModel(v interface{}) (int64, error) {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})

	orm.base.initByModel(v)
	orm.base.initExtra()
	return orm.base.doUpdate()
}

func (orm OrmTableUpdate) ByWhere(w *WhereBuilder) (int64, error) {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ:   tWhereBuilder,
		where: w,
	})

	orm.base.initByWhere(w)
	orm.base.initExtra()
	return orm.base.doUpdate()
}

//------------------------------------Select--------------------------------------------

func (db DB) Select(v interface{}) OrmTableSelect {
	db.baseTokens = append(db.baseTokens, baseToken{
		typ:  tSelect,
		dest: v,
		v:    reflect.ValueOf(v),
	})

	db.setTargetDest2TableName(v)
	return OrmTableSelect{base: db}
}

func (orm OrmTableSelect) ByPrimaryKey(v ...interface{}) OrmTableSelectWhere {
	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	return OrmTableSelectWhere{base: orm.base}
}

//ptr-comp
func (orm OrmTableSelect) ByModel(v interface{}) OrmTableSelectWhere {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})

	orm.base.initByModel(v)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelect) ByWhere(w *WhereBuilder) OrmTableSelectWhere {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ:   tWhereBuilder,
		where: w,
	})

	if w == nil {
		return OrmTableSelectWhere{base: orm.base}
	}
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
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tScanFirst,
		v:   reflect.ValueOf(v),
	})

	orm.Limit(1)
	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelectWhere) ScanOne(v interface{}) (int64, error) {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tScanOne,
		v:   reflect.ValueOf(v),
	})

	orm.base.ctx.initScanDestOne(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

func (orm OrmTableSelectWhere) ScanList(v interface{}) (int64, error) {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tScanList,
		v:   reflect.ValueOf(v),
	})

	orm.base.ctx.initScanDestList(v)
	orm.base.ctx.checkScanDestField()
	orm.base.initColumns()

	orm.base.initExtra()
	return orm.base.doSelect()
}

//------------------------------------has--------------------------------------------

func (db DB) Has(v interface{}) OrmTableHas {
	db.typ = dHas

	db.setTargetDest2TableName(v)
	return OrmTableHas{base: db}
}

// ByPrimaryKey
//v0.8
func (orm OrmTableHas) ByPrimaryKey(v ...interface{}) (bool, error) {
	orm.base.baseTokens = append(orm.base.baseTokens, baseToken{
		typ: tPrimaryKey,
		pk:  v,
	})

	orm.base.initPrimaryKeyName()
	orm.base.ctx.initPrimaryKeyValues(v)
	orm.base.initByPrimaryKey()
	orm.base.initExtra()
	return orm.base.doHas()
}

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
