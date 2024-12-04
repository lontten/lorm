package utils

import (
	"github.com/lontten/lorm/soft-delete"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

type TestSoftDel1 struct {
	soft_delete.DeleteGormMilli
}

type TestSoftDel2 struct {
	TestSoftDel1
}

func TestCheckSoftDelType(t *testing.T) {

	delType := GetSoftDelType(reflect.TypeOf(TestSoftDel2{}))
	t.Log(delType)
}

func TestErr(t *testing.T) {
	err := kk()
	t.Log(err)
	t.Log(err == nil)
}

var gg error

func kk() (err error) {
	gg = errors.New("jfda")
	return gg
}
