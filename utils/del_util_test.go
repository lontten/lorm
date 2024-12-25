package utils

import (
	"github.com/lontten/lorm/softdelete"
	"reflect"
	"testing"
)

type TestSoftDel1 struct {
	softdelete.DeleteGormMilli
}
type T1 struct {
}
type TestSoftDel2 struct {
	T1
	TestSoftDel1
}

func TestCheckSoftDelType(t *testing.T) {

	delType := GetSoftDelType(reflect.TypeOf(TestSoftDel2{}))
	t.Log(delType)
}

func TestErr(t *testing.T) {
	var kb = &TestSoftDel2{} // 使用指针
	v := reflect.ValueOf(kb)

	f := v.Field(0)

	var num = 10

	f.Set(reflect.ValueOf(&num))

}

type Kb struct {
	Id *int
}

func TestIsSoftDelFieldType(t *testing.T) {

	delType := IsSoftDelFieldType(reflect.TypeOf(TestSoftDel2{}))
	t.Log(delType)
}
