package lorm

import (
	"database/sql"
	"fmt"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

// ScanLn
//接收一行结果
// v0.7
// 1.ptr single/comp
// 2.slice- single
func (ctx OrmContext) ScanLn(rows *sql.Rows) (num int64, err error) {
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)

	num = 1
	base := ctx.destValue
	t := ctx.destBaseType


	columns, err := rows.Columns()
	if err != nil {
		return
	}
	cfm, err := ormConfig.getColFieldIndexLinkMap(columns, t)
	if err != nil {
		return
	}
	if rows.Next() {
		box, _, v := createColBox(t, cfm)
		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return
		}
		base.Set(v)
	}

	if rows.Next() {
		return 0, errors.New("result to many for one")
	}
	return
}

// ScanLnBatch
//接收一行结果-批量
// v0.7
// 1.ptr single/comp
// 2.slice- single
func (ctx OrmContext) ScanLnBatch(rowss []*sql.Rows) (int64, error) {
	var num int64 = 0
	for _, rows := range rowss {
		n, err := ctx.ScanLn(rows)
		if err != nil {
			return num, err
		}
		num += n
	}
	return num, nil
}

//Scan
// v0.7
//接收多行结果
//1.[]- *
func (ctx OrmContext) Scan(rows *sql.Rows) (int64, error) {
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)

	arr := ctx.destValue
	t := ctx.destBaseType
	isPtr := ctx.sliceItemIsPtr

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	cfm, err := ormConfig.getColFieldIndexLinkMap(columns, t)
	if err != nil {
		return 0, err
	}
	var num int64 = 0
	for rows.Next() {
		box, vp, v := createColBox(t, cfm)

		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		if isPtr {
			arr.Set(reflect.Append(arr, vp))
		} else {
			arr.Set(reflect.Append(arr, v))
		}
		num++
	}
	return num, nil
}

//ScanBatch
// v0.6
//接收多行结果-批量
//1.[]- *
func (ctx OrmContext) ScanBatch(rowss []*sql.Rows) (num int64, err error) {
	for _, rows := range rowss {
		n, err := ctx.Scan(rows)
		if err != nil {
			return num, err
		}
		num += n
	}
	return
}



//检查sturct的filed是否合法，valuer，nuller
func (ctx *OrmContext) checkScanDestField() {
	if ctx.err != nil {
		return
	}
	if !ctx.destTypeIsComp {
		return
	}
	v := ctx.destBaseType
	err := checkCompFieldScan(v)
	ctx.err = err
}
