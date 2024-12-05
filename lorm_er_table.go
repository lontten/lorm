package lorm

import (
	"fmt"
	"github.com/lontten/lorm/sql-type"
	"reflect"
)

// ------------------------------------Insert--------------------------------------------

// Insert
func Insert(db Engine, v any, extra ...*ExtraContext) (num int64, err error) {
	dialect := db.getDialect()
	ctx := dialect.initContext()
	ctx.initExtra(extra...)
	ctx.sqlType = sql_type.Insert
	ctx.sqlIsQuery = true
	dialect.appendBaseToken(baseToken{
		typ:  tInsert,
		dest: v,
	})

	ctx.setModelDest(v)
	dialect.tableInsertGen()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	sql := dialect.getSql()
	if ctx.showSql {
		fmt.Println(sql, ctx.args)
	}

	if ctx.sqlIsQuery {
		rows, err := db.query(sql, ctx.args...)
		if err != nil {
			return 0, err
		}
		return ctx.ScanLnT(rows)
	}

	exec, err := db.exec(sql, ctx.args...)
	if err != nil {
		return 0, err
	}
	if ctx.needLastInsertId {
		id, err := exec.LastInsertId()
		if err != nil {
			return 0, err
		}
		if id > 0 {
			// todo 将自增id 赋值到v
			numField := ctx.destV.NumField()
			fmt.Println(numField)
			for i := range numField {
				field := ctx.destV.Field(i)
				isPtr, v, err := basePtrValue(field)
				fmt.Println(isPtr, v.Type().Name(), err)
				fmt.Println(field.Interface())

				fmt.Println("------")
			}
			//field := ctx.destV.Field(0)
			//field.Set(reflect.ValueOf(id))
			fmt.Println(ctx.destV.Interface())
		}
	}
	return exec.RowsAffected()
}

// Insert
func InsertOrHas(db Engine, v any, extra ...*ExtraContext) (num int64, err error) {
	dialect := db.getDialect()
	ctx := dialect.initContext()
	ctx.initExtra(extra...)
	ctx.sqlType = sql_type.Insert
	ctx.sqlIsQuery = true
	dialect.appendBaseToken(baseToken{
		typ:  tInsert,
		dest: v,
	})

	ctx.setModelDest(v)
	dialect.tableInsertGen()

	sql := dialect.getSql()
	if ctx.showSql {
		fmt.Println(sql, ctx.args)
	}

	if ctx.sqlIsQuery {
		rows, err := db.query(sql, ctx.args...)
		if err != nil {
			return 0, err
		}
		return ctx.ScanLnT(rows)
	}

	exec, err := db.exec(sql, ctx.args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

// InsertOrUpdate
// 1.ptr
// 2.comp-struct
func (db *lnDB) InsertOrUpdate(v any) OrmTableCreate {
	db.core.appendBaseToken(baseToken{
		typ:  tInsert,
		dest: v,
	})

	//db.setModelDest(v)
	//ldb.initColumnsValue()
	return OrmTableCreate{base: db.core}
}

// ByPrimaryKey
// ptr
// single / comp复合主键
func (orm OrmTableCreate) ByPrimaryKey() (int64, error) {
	//orm.base.initPrimaryKeyName()
	//orm.base.ctx.initSelfPrimaryKeyValues()
	//base := orm.base
	//ctx := base.ctx
	//if err := ctx.err; err != nil {
	//	return 0, err
	//}
	//
	//cs := ctx.columns
	//cvs := ctx.columnValues
	//tableName := ctx.tableName
	//idNames := ctx.primaryKeyNames
	//query, args := base.dialect.insertOrUpdateByPrimaryKey(tableName, idNames, cs, cvs...)
	//return base.doExec(query, args...)
	return 0, nil
}

// ByUnique
// ptr-comp
func (orm OrmTableCreate) ByUnique(fs ...string) (int64, error) {
	//if len(fs) == 0 {
	//	return 0, errors.New("ByUnique is empty")
	//}
	//
	//base := orm.base
	//ctx := base.ctx
	//if err := ctx.err; err != nil {
	//	return 0, err
	//}
	//
	//cs := ctx.columns
	//cvs := ctx.columnValues
	//tableName := ctx.tableName
	//query, args := base.dialect.insertOrUpdateByUnique(tableName, fs, cs, cvs...)
	//return base.doExec(query, args...)
	return 0, nil
}

//------------------------------------Delete--------------------------------------------

func (db *lnDB) Delete(v any) OrmTableDelete {
	core := db.core
	core.getCtx().tableSqlType = dDelete
	core.appendBaseToken(baseToken{
		typ: tTableName,
		t:   reflect.TypeOf(v),
	})
	return OrmTableDelete{base: core}
}

func (orm OrmTableDelete) ByPrimaryKey(v ...any) OrmTableDelete {
	orm.base.appendBaseToken(baseToken{
		typ: tPrimaryKey,
		pk:  v,
	})
	return orm
}

func (orm OrmTableDelete) ByModel(v any) OrmTableDelete {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})
	return orm
}

func (orm OrmTableDelete) ByWhere(wb *WhereBuilder) OrmTableDelete {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereBuilder,
		wb:  wb,
	})
	return orm
}

func (orm OrmTableDelete) Exec(w *WhereBuilder) Resulter {
	//return orm.base.getCtx()
	return nil
}

//------------------------------------Update--------------------------------------------

func (db *lnDB) Update(v any) OrmTableUpdate {
	core := db.core
	core.getCtx().tableSqlType = dUpdate
	core.appendBaseToken(baseToken{
		typ:  tTableNameDestValue,
		t:    reflect.TypeOf(v),
		dest: v,
		v:    reflect.ValueOf(v),
	})

	//ldb.setModelDest(v)
	//ldb.initColumnsValue()

	return OrmTableUpdate{base: core}
}

func (orm OrmTableUpdate) ByPrimaryKey(v ...any) OrmTableUpdate {
	orm.base.appendBaseToken(baseToken{
		typ: tPrimaryKey,
		pk:  v,
	})

	//orm.base.initPrimaryKeyName()
	//orm.base.ctx.initSelfPrimaryKeyValues()
	//orm.base.initByPrimaryKey()
	//orm.base.initExtra()

	return orm
}

func (orm OrmTableUpdate) ByModel(v any) OrmTableUpdate {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})

	//orm.base.initByModel(v)
	//orm.base.initExtra()
	return orm
}

func (orm OrmTableUpdate) ByWhere(wb *WhereBuilder) OrmTableUpdate {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereBuilder,
		wb:  wb,
	})

	//orm.base.initByWhere(w)
	//orm.base.initExtra()

	return orm
}

//------------------------------------Select--------------------------------------------

func (db *lnDB) Select(v any) OrmTableSelect {
	core := db.core
	core.getCtx().tableSqlType = dSelect
	core.appendBaseToken(baseToken{
		typ:  tTableNameDestValue,
		t:    reflect.TypeOf(v),
		dest: v,
		v:    reflect.ValueOf(v),
	})

	db.setNameDest(v)
	return OrmTableSelect{base: core}
}

func (orm OrmTableSelect) ByPrimaryKey(v ...any) OrmTableSelect {
	orm.base.appendBaseToken(baseToken{
		typ: tPrimaryKey,
		pk:  v,
	})
	//orm.base.initPrimaryKeyName()
	//orm.base.ctx.initPrimaryKeyValues(v)
	//orm.base.initByPrimaryKey()
	return orm
}

// ptr-comp
func (orm OrmTableSelect) ByModel(v any) OrmTableSelect {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})
	//orm.base.initByModel(v)
	return orm
}

func (orm OrmTableSelect) ByWhere(wb *WhereBuilder) OrmTableSelect {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereBuilder,
		wb:  wb,
	})

	//if w == nil {
	//	return OrmTableSelectWhere{base: orm.base}
	//}
	//orm.base.initByWhere(w)
	return orm
}

func (orm OrmTableSelectWhere) OrderBy(name string, condition ...bool) OrmTableSelectWhere {

	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}

	//orm.base.orderByTokens = append(orm.base.orderByTokens, name)
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) OrderDescBy(name string, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}

	//orm.base.orderByTokens = append(orm.base.orderByTokens, name+" desc")
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) Limit(num int64, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}
	//orm.base.limit = num
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) Offset(num int64, condition ...bool) OrmTableSelectWhere {
	for _, b := range condition {
		if !b {
			return OrmTableSelectWhere{base: orm.base}
		}
	}
	//orm.base.offset = num
	return OrmTableSelectWhere{base: orm.base}
}

func (orm OrmTableSelectWhere) ScanFirst(v any) (int64, error) {
	orm.base.appendBaseToken(baseToken{
		typ: tScanFirst,
		v:   reflect.ValueOf(v),
	})

	//orm.Limit(1)
	//orm.base.ctx.initScanDestOne(v)
	//orm.base.ctx.checkScanDestField()
	//orm.base.initColumns()
	//
	//orm.base.initExtra()
	//return orm.base.doSelect()
	return 0, nil
}

func (orm OrmTableSelectWhere) ScanOne(v any) (int64, error) {
	orm.base.appendBaseToken(baseToken{
		typ: tScanOne,
		v:   reflect.ValueOf(v),
	})

	//orm.base.ctx.initScanDestOne(v)
	//orm.base.ctx.checkScanDestField()
	//orm.base.initColumns()
	//
	//orm.base.initExtra()
	//return orm.base.doSelect()
	return 0, nil
}

func (orm OrmTableSelectWhere) ScanList(v any) (int64, error) {
	orm.base.appendBaseToken(baseToken{
		typ: tScanList,
		v:   reflect.ValueOf(v),
	})

	//orm.base.ctx.initScanDestList(v)
	//orm.base.ctx.checkScanDestField()
	//orm.base.initColumns()
	//
	//orm.base.initExtra()
	//return orm.base.doSelect()
	return 0, nil
}

//------------------------------------has--------------------------------------------

func (db *lnDB) Has(v any) OrmTableHas {
	core := db.core
	core.getCtx().tableSqlType = dHas
	core.appendBaseToken(baseToken{
		typ:  tTableNameDestValue,
		t:    reflect.TypeOf(v),
		dest: v,
		v:    reflect.ValueOf(v),
	})

	db.setNameDest(v)
	return OrmTableHas{base: core}
}

// ByPrimaryKey
// v0.8
func (orm OrmTableHas) ByPrimaryKey(v ...any) OrmTableHas {
	orm.base.appendBaseToken(baseToken{
		typ: tPrimaryKey,
		pk:  v,
	})

	//orm.base.initPrimaryKeyName()
	//orm.base.ctx.initPrimaryKeyValues(v)
	//orm.base.initByPrimaryKey()
	//orm.base.initExtra()
	return orm
}

// ptr-comp
func (orm OrmTableHas) ByModel(v any) OrmTableHas {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereModel,
		v:   reflect.ValueOf(v),
	})

	//orm.base.initByModel(v)
	//orm.base.initExtra()
	return orm
}

func (orm OrmTableHas) ByWhere(wb *WhereBuilder) OrmTableHas {
	orm.base.appendBaseToken(baseToken{
		typ: tWhereBuilder,
		wb:  wb,
	})

	//orm.base.initByWhere(wb)
	//orm.base.initExtra()
	return orm
}

//-----------------------------------------------------------

type OrmTableCreate struct {
	base corer
}

type OrmTableSelect struct {
	base corer

	query string
	args  []any
}

type OrmTableHas struct {
	base corer

	query string
	args  []any
}

type OrmTableSelectWhere struct {
	base corer
}

type OrmTableUpdate struct {
	base corer
}

type OrmTableDelete struct {
	base corer
}
