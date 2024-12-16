package lorm

import (
	"fmt"
	"github.com/lontten/lorm/sqltype"
	"reflect"
)

// ------------------------------------Insert--------------------------------------------

// Insert 插入或者根据主键冲突更新
func Insert(db Engine, v any, extra ...*ExtraContext) (num int64, err error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Insert
	ctx.sqlIsQuery = true

	ctx.initModelDest(v)   //初始化参数
	ctx.initConf()         //初始化表名，主键，自增id
	ctx.initColumnsValue() //初始化cv

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
			ctx.setLastInsertId(id)
			if ctx.hasErr() {
				return 0, ctx.err
			}
		}
	}
	return exec.RowsAffected()
}

// InsertOrHas 根据条件查询是否已存在，不存在则直接插入
// 应用场景：例如添加 后台管理员 时，如果名字已存在，返回名字重复，否者正常添加。
func InsertOrHas(db Engine, v any, extra ...*ExtraContext) (num int64, err error) {
	return 0, err
}

//------------------------------------Delete--------------------------------------------

func Delete[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (int64, error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Delete
	ctx.sqlIsQuery = false

	dest := new(T)
	ctx.initScanDestOneT(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}

	ctx.initConf() //初始化表名，主键，自增id
	ctx.initColumnsValueSoftDel()

	ctx.initPrimaryKeyByWhere(wb)
	ctx.wb.And(wb)

	whereStr, args, err := ctx.wb.toSql(dialect.parse)
	if err != nil {
		return 0, err
	}
	ctx.whereSql = whereStr
	ctx.args = append(ctx.args, args...)

	dialect.tableDelGen()
	if ctx.hasErr() {
		return 0, ctx.err
	}
	sql := dialect.getSql()
	if ctx.showSql {
		fmt.Println(sql, ctx.args)
	}
	exec, err := db.exec(sql, ctx.args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
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

	return OrmTableUpdate{base: core}
}

//------------------------------------Select--------------------------------------------

// First 根据条件获取第一个
func First[T any](db Engine, by *ByContext, extra ...*ExtraContext) (t *T, err error) {
	return nil, err
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
