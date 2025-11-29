package utils

import (
	"reflect"
	"testing"

	"github.com/lontten/lorm/softdelete"
	"github.com/stretchr/testify/assert"
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
	as := assert.New(t)

	delType := GetSoftDelType(reflect.TypeOf(TestSoftDel2{}))
	as.Equal(softdelete.DelTimeGormMilli, delType)
}

func TestIsSoftDelFieldType(t *testing.T) {
	as := assert.New(t)

	delType := IsSoftDelFieldType(reflect.TypeOf(TestSoftDel2{}))
	as.False(delType)

}
