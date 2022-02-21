package hook

import (
	"github.com/jackc/pgtype"
	"github.com/lontten/lorm/types"
)

/**



 */
func StringEn(v string) interface{} {

	return v[1:]
}

type Hello struct {
	_db  struct{} `lorm:"column(name)"`
	Name string   `lorm:"hook.after:StringEn" json:"name"`
}

func ToArr(src pgtype.DateArray) []types.Date {
	var arr []types.Date
	err := src.AssignTo(&arr)
	if err != nil {
		panic(err)
	}
	return arr
}

func ToPgDateArr(src []types.Date) pgtype.DateArray {
	var arr pgtype.DateArray
	err := arr.Set(src)
	if err != nil {
		panic(err)
	}
	return arr
}
