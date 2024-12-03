package lorm

import (
	"database/sql"
	"fmt"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

// ScanLn
// 接收一行结果
// 1.ptr single/comp
// 2.slice- single
//func (ctx ormContext) ScanLn(rows *sql.Rows) (num int64, err error) {
//	defer func(rows *sql.Rows) {
//		utils.PanicErr(rows.Close())
//	}(rows)
//
//	num = 0
//	base := ctx.destV
//	t := ctx.destBaseType
//
//	columns, err := rows.Columns()
//	if err != nil {
//		return
//	}
//	cfm, err := ctx.ormConf.getColFieldIndexLinkMap(columns, t)
//	if err != nil {
//		return
//	}
//	if rows.Next() {
//		box, _, v := createColBox(t, cfm)
//		err = rows.Scan(box...)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		base.SetContext(v)
//		num++
//	}
//
//	if rows.Next() {
//		return 0, errors.New("result to many for one")
//	}
//	return
//}

// ScanLn
// 接收一行结果
// 1.ptr single/comp
// 2.slice- single
func (ctx ormContext) ScanLnT(rows *sql.Rows) (num int64, err error) {
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)

	num = 0
	t := ctx.destBaseType
	v := ctx.destV
	tP := ctx.scanDest

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	cfm := ColFieldIndexLinkMap{}
	if ctx.destBaseTypeIsComp {
		cfm, err = ctx.ormConf.getColFieldIndexLinkMap(columns, t)
		if err != nil {
			return
		}
	}

	if rows.Next() {
		box := createColBoxT(v, tP, cfm)
		err = rows.Scan(box...)
		if err != nil {
			return
		}
		num++
	}

	if rows.Next() {
		return 0, errors.New("result to many for one")
	}
	return
}

// ScanBatch
// 批量
// 1.ptr single/comp
// 2.slice- single
func (ctx ormContext) ScanBatch(rowss []*sql.Rows) (int64, error) {
	var nums int64 = 0
	arr := ctx.destV
	t := ctx.destBaseType
	isPtr := ctx.destSliceItemIsPtr

	for _, rows := range rowss {

		defer func(rows *sql.Rows) {
			utils.PanicErr(rows.Close())
		}(rows)

		columns, err := rows.Columns()
		if err != nil {
			return 0, err
		}
		cfm, err := ctx.ormConf.getColFieldIndexLinkMap(columns, t)
		if err != nil {
			return 0, err
		}
		var num int64 = 0
		for rows.Next() {
			box, vp, v := createColBox(t, cfm)
			err = rows.Scan(box...)
			if err != nil {
				return 0, err
			}
			if isPtr {
				arr.Set(reflect.Append(arr, vp))
			} else {
				arr.Set(reflect.Append(arr, v))
			}
			num++
		}

		nums += num
	}
	return nums, nil
}

// Scan
// 接收多行结果
// 1.[]- *
//func (ctx ormContext) Scan(rows *sql.Rows) (int64, error) {
//	defer func(rows *sql.Rows) {
//		utils.PanicErr(rows.Close())
//	}(rows)
//
//	var num int64 = 0
//	t := ctx.destBaseType
//	arr := ctx.destV
//	isPtr := ctx.destSliceItemIsPtr
//
//	columns, err := rows.Columns()
//	if err != nil {
//		return 0, err
//	}
//	cfm, err := ctx.ormConf.getColFieldIndexLinkMap(columns, t)
//	if err != nil {
//		return 0, err
//	}
//	for rows.Next() {
//		box, vp, v := createColBox(t, cfm)
//
//		err = rows.Scan(box...)
//		if err != nil {
//			fmt.Println(err)
//			return 0, err
//		}
//		if isPtr {
//			arr.SetContext(reflect.Append(arr, vp))
//		} else {
//			arr.SetContext(reflect.Append(arr, v))
//		}
//		num++
//	}
//	return num, nil
//}

// Scan
// 接收多行结果
// 1.[]- *
func (ctx ormContext) ScanT(rows *sql.Rows) (int64, error) {
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)

	var num int64 = 0
	t := ctx.destBaseType
	arr := ctx.destV
	isPtr := ctx.destSliceItemIsPtr

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	cfm, err := ctx.ormConf.getColFieldIndexLinkMap(columns, t)
	if err != nil {
		return 0, err
	}
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

// 检查scan dest 的filed是否合法，valuer，nuller
func (ctx *ormContext) checkScanDestField() {
	if ctx.err != nil {
		return
	}
	if !ctx.destBaseTypeIsComp {
		return
	}
	ctx.err = checkCompFieldV(ctx.destBaseType)
}
