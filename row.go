package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

func StructScan(rows *sql.Rows, dest interface{}) (int64, error) {
	defer rows.Close()
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return 0, errors.New("dest need a struct pointer")
	}
	arr := reflect.Indirect(value)
	err := checkValidStruct(arr.Elem())
	if err != nil {
		return 0, err
	}

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
	rsFM, err := getRowStructFieldMap(columns, base)
	if err != nil {
		return 0, err
	}
	var num int64 = 0
	for rows.Next() {
		num++
		box, vp, v := getRowFieldBox(columns, base, rsFM)

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

//只有一个结果的row
func StructScanLn(rows *sql.Rows, dest interface{}) (num int64, err error) {

	defer rows.Close()
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return 0, errors.New("dest need a  ptr")
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		box, v := getSignleRowFieldBox(value.Type())
		if rows.Next() {
			num++
			err = rows.Scan(box)
			if err != nil {
				fmt.Println(err)
				return
			}
			value.Set(v)
		}
		if rows.Next() {
			return 0, errors.New("result to many for one")
		}
		return
	}

	//err = checkValidStruct(value)
	//if err != nil {
	//	return 0, err
	//}

	typ := reflect.TypeOf(dest)
	base, err := baseStructTypePtr(typ)
	if err != nil {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}
	rsFM, err := getRowStructFieldMap(columns, base)
	if err != nil {
		return
	}
	if rows.Next() {
		num++
		box, _, p := getRowFieldBox(columns, base, rsFM)
		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return
		}
		value.Elem().Set(p)
	}
	if rows.Next() {
		return 0, errors.New("result to many for one")
	}

	return
}

//用来存放row中值得 引用
func getRowFieldBox(columns []string, base reflect.Type, rsFM RowStructFieldMap) (box []interface{}, vp, v reflect.Value) {
	vp = newStruct(base)
	v = reflect.Indirect(vp)
	fieldNum := len(columns)
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

type RowStructFieldMap []int

func getRowStructFieldMap(columns []string, typ reflect.Type) (RowStructFieldMap, error) {
	rfm := make([]int, len(columns))
	fm, err := getFieldMap(typ)
	if err != nil {
		return nil, err
	}

	for i, column := range columns {
		index, ok := fm[column]
		if !ok {
			rfm[i] = -1
			continue
		}
		rfm[i] = index
	}
	return rfm, nil
}
