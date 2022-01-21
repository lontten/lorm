package lsql

import (
	"github.com/lontten/lsql/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test__isBaseType(t *testing.T) {
	as := assert.New(t)

	is := _isBaseType(reflect.TypeOf(types.NewV4()))
	as.False(is)

	is = _isBaseType(reflect.TypeOf(12))
	as.True(is)

	is = _isBaseType(reflect.TypeOf("jfaskf"))
	as.True(is)

	is = _isBaseType(reflect.TypeOf(types.NewString("fjakls")))
	as.False(is)

}
