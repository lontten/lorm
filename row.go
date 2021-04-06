package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type Rows struct {
	*sql.Rows
	started bool
	base    reflect.Type
	fields  []int
	values  []interface{}
}

type Row struct {
	rows *sql.Rows
	err  error
}

func (r *Row) Scan(dest ...interface{}) error {
	return nil
}

func (r *Row) StructScan(dest interface{}) error {
	return nil
}

func (r *Rows) StructScan(dest interface{}) (int64, error) {
	return StructScan(r.Rows, dest)
}

func StructScan(rows *sql.Rows, dest interface{}) (int64, error) {
	defer rows.Close()
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return 0, errors.New("dest need a struct pointer")
	}
	arr := reflect.Indirect(value)

	typ := reflect.TypeOf(dest)
	slice, err := baseSliceType(typ)
	if err != nil {
		return 0, err
	}

	base := slice.Elem()
	var isPtr = base.Kind() == reflect.Ptr
	base, err = baseStructType(base)
	if err != nil {
		return 0, err
	}
	isNullable := checkNullStruct(base)
	if !isNullable {
		return 0, errors.New("struct need fields all pointer")
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

func StructScanLn(rows *sql.Rows, dest interface{}) (int64, error) {
	defer rows.Close()
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return 0, errors.New("dest need a struct pointer")
	}

	typ := reflect.TypeOf(dest)
	base, err := baseStructType(typ)
	if err != nil {
		return 0, err
	}
	isNullable := checkNullStruct(base)
	if !isNullable {
		return 0, errors.New("struct need fields all pointer")
	}

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	rsFM, err := getRowStructFieldMap(columns, base)
	if err != nil {
		return 0, err
	}
	var num int64=0
	if rows.Next() {
		num++
		box, _, p := getRowFieldBox(columns, base, rsFM)
		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		value.Elem().Set(p)
	}
	if rows.Next() {
		return 0, errors.New("result to many for one")
	}

	return num, nil
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
