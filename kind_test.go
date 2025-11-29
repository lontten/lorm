package lorm

import (
	"reflect"
	"testing"

	"github.com/lontten/lcore/v2/types"
	"github.com/stretchr/testify/assert"
)

func Test_basePtrValue(t *testing.T) {
	as := assert.New(t)
	var a *int
	v := reflect.ValueOf(a)
	is, v, err := basePtrValue(v)
	as.NotNil(err)
	as.False(is)
}

func Test_isValuerType(t *testing.T) {
	as := assert.New(t)

	t1 := reflect.TypeOf(types.LocalDateTime{})
	as.True(t1.Implements(ImpValuer))

	t2 := reflect.TypeOf(new(types.LocalDateTime))
	as.True(t2.Implements(ImpValuer))
}

func Test_isScannerType(t *testing.T) {
	as := assert.New(t)

	t1 := reflect.TypeOf(types.LocalDateTime{})
	as.False(t1.Implements(ImpScanner))

	t2 := reflect.TypeOf(new(types.LocalDateTime))
	as.True(t2.Implements(ImpScanner))
}
