package types

import (
	"github.com/jackc/pgtype"
)

func Pg2Arr(v pgtype.BoolArray) []bool {
	var arr []bool
	for _, element := range v.Elements {
		if element.Status == pgtype.Present {
			arr = append(arr, element.Bool)
		}
	}
	return arr
}

func Arr2Pg(arr []bool) pgtype.BoolArray {
	var list = pgtype.BoolArray{}
	list.Set(arr)
	return list
}
