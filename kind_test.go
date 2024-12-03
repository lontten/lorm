package lorm

import (
	"database/sql"
	"reflect"
	"testing"
)

func Test_basePtrValue(t *testing.T) {
	var a *int
	v := reflect.ValueOf(a)
	t.Log(v.Kind())
	t.Log(v.IsValid())
	t.Log(v.IsZero())
	t.Log(v.IsNil())
	t.Log(v.Elem())
	is, v, err := basePtrValue(v)
	t.Log(is, v, err)
}

func Test_vn(ts *testing.T) {
	var d sql.NullTime

	v := reflect.ValueOf(d)
	t := v.Type()
	ts.Log(t.Kind())
	ts.Log(t.String())
	ts.Log(t.Name())
}
