package utils

import (
	"github.com/lontten/lorm/softdelete"
	"reflect"
	"testing"
)

type TestSoftDel1 struct {
	softdelete.DeleteGormMilli
}

type TestSoftDel2 struct {
	TestSoftDel1
}

func TestCheckSoftDelType(t *testing.T) {

	delType := GetSoftDelType(reflect.TypeOf(TestSoftDel2{}))
	t.Log(delType)
}

func TestErr(t *testing.T) {
	var kb = &Kb{}                  // 使用指针
	v := reflect.ValueOf(kb).Elem() // 获取指针指向的值

	f := v.Field(0)

	var num = 10

	f.Set(reflect.ValueOf(&num))

}

type Kb struct {
	Id *int
}
