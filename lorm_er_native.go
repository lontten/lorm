package lorm

import (
	"reflect"
)

/*
*
num,dto,err:=QueryOne[User](ldb,"",id)
num,dtos,err:=QueryList[User](ldb,"",id)

num,err:=Exec(ldb,"",id)
*/

func StmtQueryOne[T any](db Stmter, args ...any) (*T, error) {
	db.getDialect().initContext()
	ctx := db.getCtx()

	dest := new(T)

	ctx.initScanDestOneT(dest)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(args...)
	if err != nil {
		return nil, err
	}
	_, err = ctx.ScanLnT(rows)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func StmtQueryList[T any](db Stmter, args ...any) ([]T, error) {
	db.getDialect().initContext()
	ctx := db.getCtx()

	var dest = &[]T{}
	v := reflect.ValueOf(dest).Elem()
	t := reflect.TypeFor[T]()

	ctx.initScanDestListT(dest, v, t, false)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(args...)
	if err != nil {
		return nil, err
	}
	_, err = ctx.ScanT(rows)
	if err != nil {
		return nil, err
	}
	return *dest, nil
}

func StmtQueryListP[T any](db Stmter, args ...any) ([]*T, error) {
	db.getDialect().initContext()
	ctx := db.getCtx()

	var dest = &[]*T{}
	v := reflect.ValueOf(dest).Elem()
	t := reflect.TypeFor[T]()

	ctx.initScanDestListT(dest, v, t, true)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(args...)
	if err != nil {
		return nil, err
	}
	_, err = ctx.ScanT(rows)
	if err != nil {
		return nil, err
	}
	return *dest, nil
}

func QueryOne[T any](db Engine, query string, args ...any) (*T, error) {
	db.getDialect().initContext()
	ctx := db.getCtx()

	dest := new(T)

	ctx.initScanDestOneT(dest)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(query, args...)
	if err != nil {
		return nil, err
	}
	_, err = ctx.ScanLnT(rows)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func QueryList[T any](db Engine, query string, args ...any) ([]T, error) {
	db.getDialect().initContext()
	ctx := db.getCtx()

	var dest = &[]T{}
	v := reflect.ValueOf(dest).Elem()
	t := reflect.TypeFor[T]()

	ctx.initScanDestListT(dest, v, t, false)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(query, args...)
	if err != nil {
		return nil, err
	}
	_, err = ctx.ScanT(rows)
	if err != nil {
		return nil, err
	}
	return *dest, nil
}

func QueryListP[T any](db Engine, query string, args ...any) ([]*T, error) {
	db.getDialect().initContext()
	ctx := db.getCtx()

	var dest = &[]*T{}
	v := reflect.ValueOf(dest).Elem()
	t := reflect.TypeFor[T]()

	ctx.initScanDestListT(dest, v, t, true)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(query, args...)
	if err != nil {
		return nil, err
	}
	_, err = ctx.ScanT(rows)
	if err != nil {
		return nil, err
	}
	return *dest, nil
}

func Exec(db Engine, query string, args ...any) (int64, error) {
	db.getDialect().initContext()
	exec, err := db.exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

type NativeQuery struct {
	db    Engine
	query string
	args  []any
}

func QueryScan(db Engine, query string, args ...any) NativeQuery {
	return NativeQuery{
		db:    db,
		query: query,
		args:  args,
	}
}

func (q NativeQuery) ScanOne(dest any) (num int64, err error) {
	db := q.db
	ctx := db.getCtx()
	ctx.initScanDestOne(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}

	query := q.query
	args := q.args

	rows, err := db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLnT(rows)
}

// scanList 切片 必须 ptr ，才能赋值
// get操作必须 ptr，但是 insert 可是不是ptr，只是dest 不是 ptr，无法返回 自增id
func (q NativeQuery) ScanList(dest any) (num int64, err error) {
	db := q.db
	ctx := db.getCtx()
	ctx.initScanDestList(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}

	query := q.query
	args := q.args

	rows, err := db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanT(rows)
}
