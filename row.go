package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

func StructScan(rows *sql.Rows, dest interface{}, fieldNamePrefix string) (int64, error) {
	defer rows.Close()

	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return 0, errors.New("dest need a struct pointer")
	}
	arr := reflect.Indirect(value)

	typ := reflect.TypeOf(dest)
	slice, err := baseSliceTypePtr(typ)
	if err != nil {
		return 0, err
	}

	base := slice.Elem()
	var isPtr = base.Kind() == reflect.Ptr
	base, err = baseStructTypePtr(base)
	if err != nil {
		return 0, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	cfm, err := getColFieldIndexLinkMap(columns, base, fieldNamePrefix)
	if err != nil {
		return 0, err
	}
	var num int64 = 0
	for rows.Next() {
		num++
		box, vp, v := createColBox(base, cfm)

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

	}
	return num, nil
}

// StructScanLn 只有一个结果的row
func StructScanLn(rows *sql.Rows, dest interface{}, fieldNamePrefix string) (num int64, err error) {
	defer rows.Close()
	value := reflect.ValueOf(dest)
	code, base := basePtrStructBaseValue(value)
	if code == -1 {
		return 0, errors.New("dest need a  ptr")
	}
	if code == -2 {
		return 0, errors.New("need a ptr struct or base type")
	}

	num=1
	t := base.Type()

	if code==2 {
		box, v := getSignleRowFieldBox(t)
		if rows.Next() {
			err = rows.Scan(box)
			if err != nil {
				fmt.Println(err)
				return
			}
			value.Set(v)
		}
	}

	if code==1 {
		columns, err := rows.Columns()
		if err != nil {
			return
		}
		cfm, err := getColFieldIndexLinkMap(columns, t, fieldNamePrefix)
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
			value.Set(v)
		}
	}
	if rows.Next() {
		return 0, errors.New("result to many for one")
	}
	return
}

//创建用来存放row中值得 引用
func createColBox(base reflect.Type, cfLink ColFieldIndexLinkMap) (box []interface{}, vp, v reflect.Value) {
	vp = newStruct(base)
	v = reflect.Indirect(vp)
	box = make([]interface{}, len(cfLink))
	for c, f := range cfLink {
		if f < 0 { // -1 表示此列不接收
			box[c] = new([]uint8)
		} else {
			box[c] = v.Field(f).Addr().Interface()
		}
	}
	return
}

//用来存放row中值得 引用
func getRowBox(base reflect.Type, rsFM ColFieldIndexLinkMap) (num int, box []interface{}, vp, v reflect.Value) {
	vp = newStruct(base)
	v = reflect.Indirect(vp)
	fieldNum := len(rsFM)
	box = make([]interface{}, fieldNum)
	for r, s := range rsFM {
		if s < 0 {
			empt := new([]uint8)
			box[r] = empt
			continue
		}
		box[r] = v.Field(s).Addr().Interface()
	}
	return
}

//用来存放row中值得 引用
func getSignleRowFieldBox(base reflect.Type) (interface{}, reflect.Value) {
	vp := reflect.New(base)
	v := reflect.Indirect(vp)
	return v.Addr().Interface(), v
}

type ColFieldIndexLinkMap []int

func getColFieldIndexLinkMap(columns []string, typ reflect.Type, fieldNamePrefix string) (ColFieldIndexLinkMap, error) {
	cfm := make([]int, len(columns))
	fm, err := getFieldMap(typ, fieldNamePrefix)
	if err != nil {
		return nil, err
	}

	for i, column := range columns {
		index, ok := fm[column]
		if !ok {
			cfm[i] = -1
			continue
		}
		cfm[i] = index
	}
	return cfm, nil
}
