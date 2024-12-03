package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jackc/pgtype"
)

func Pg2Arr64(v pgtype.Int8Array) []int64 {
	var arr []int64
	for _, element := range v.Elements {
		if element.Status == pgtype.Present {
			arr = append(arr, element.Int)
		}
	}
	return arr
}

func Arr2Pg64(arr []int64) pgtype.Int8Array {
	var list = pgtype.Int8Array{}
	list.Set(arr)
	return list
}

func Pg2Arr32(v pgtype.Int4Array) []int32 {
	var arr []int32
	for _, element := range v.Elements {
		if element.Status == pgtype.Present {
			arr = append(arr, element.Int)
		}
	}
	return arr
}

func Arr2Pg32(arr []int32) pgtype.Int4Array {
	var list = pgtype.Int4Array{}
	list.Set(arr)
	return list
}

func Pg2Arr16(v pgtype.Int2Array) []int16 {
	var arr []int16
	for _, element := range v.Elements {
		if element.Status == pgtype.Present {
			arr = append(arr, element.Int)
		}
	}
	return arr
}

func Arr2Pg16(arr []int16) pgtype.Int2Array {
	var list = pgtype.Int2Array{}
	list.Set(arr)
	return list
}

type IntList []int

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p IntList) Value() (driver.Value, error) {
	var k []int
	k = p
	marshal, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	var s = string(marshal)
	if s != "null" {
		s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
	} else {
		s = "{}"
	}
	return s, nil
}

// Scan 实现方法
func (p *IntList) Scan(data any) error {
	array := pgtype.VarcharArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []string
	list = make([]string, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = element.String
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
