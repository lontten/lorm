package jsonint

import (
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
