package hook

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/lontten/lorm/types"
	"os/user"
	"reflect"
)

/*
*
 */
func StringEn(v string) any {

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

func cc(src string) user.Group {
	var u user.Group
	err := json.Unmarshal([]byte(src), &u)
	if err != nil {
		panic(err)
	}
	return u
}

func bind(src any, name string, fun func()) {
	value := reflect.ValueOf(fun)
	method := value.MethodByName("")
	method.Call([]reflect.Value{})
}

func h() {
	//lorm.DB{}.Builder().
	//	Select("").
	//	SelectOneModel("select * from user u where u.id = k.uid","user_info_dto")
	//	SelectListModel("select * from user u where u.id = k.uid","user_info_list")
	//     Select("").

}
