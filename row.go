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
		return 0, errors.New("need a ptr type")
	}
	arr := value.Elem()
	if arr.Kind() != reflect.Slice {
		return 0, errors.New("need a slice type")
	}

	slice := arr.Type()

	base := slice.Elem()
	isPtr := base.Kind() == reflect.Ptr
	code, base := baseStructBaseType(base)
	if code == -2 {
		return 0, errors.New("need a struct or base type in  slice")
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
		num++
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

	num = 1
	t := base.Type()

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
		base.Set(v)
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
	length := len(cfLink)
	box = make([]interface{}, 1)
	if length == 0 {
		box[0] = v.Addr().Interface()
		return
	}
	box = make([]interface{}, length)
	for c, f := range cfLink {
		if f < 0 { // -1 表示此列不接收
			box[c] = new([]uint8)
		} else {
			box[c] = v.Field(f).Addr().Interface()
		}

	}
	return
}

type ColFieldIndexLinkMap []int

func getColFieldIndexLinkMap(columns []string, typ reflect.Type, fieldNamePrefix string) (ColFieldIndexLinkMap, error) {
	is := baseBaseType(typ)
	if is {
		return ColFieldIndexLinkMap{}, nil
	}

	colNum := len(columns)
	cfm := make([]int, colNum)
	fm, err := getFieldMap(typ, fieldNamePrefix)
	if err != nil {
		return nil, err
	}

	validNum := 0
	for i, column := range columns {
		index, ok := fm[column]
		if !ok {
			cfm[i] = -1
			continue
		}
		cfm[i] = index
		validNum++
	}

	if colNum == 1 && validNum == 0 {
		return ColFieldIndexLinkMap{}, nil
	}
	return cfm, nil
}
