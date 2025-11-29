package lorm

import (
	"reflect"
	"testing"

	"github.com/lontten/lcore/v2/types"
	"github.com/stretchr/testify/assert"
)

func Test_checkHandleNull(t *testing.T) {
	as := assert.New(t)

	t1 := reflect.TypeOf(new(types.LocalDateTime))
	canNull, isScanner := checkHandleNull(t1)
	as.True(canNull)
	as.True(isScanner)

	t2 := reflect.TypeOf(types.LocalDateTime{})
	canNull, isScanner = checkHandleNull(t2)
	as.True(canNull)
	as.True(isScanner)
}

func Test_checkHandleNull2(t *testing.T) {
	as := assert.New(t)

	t1 := reflect.TypeOf(new(types.LocalDateTime))
	as.True(t1.Implements(ImpScanner))

	t2 := reflect.TypeOf(types.LocalDateTime{})
	as.False(t2.Implements(ImpScanner))
}
