package lorm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lontten/lorm/utils"
)

func Test_getFieldInter(t *testing.T) {
	var a = "a234"
	v := reflect.ValueOf(a)

	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "234 text",
			args: args{v: v},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inter := getFieldInterZero(tt.args.v)
			fmt.Println(inter)
			fmt.Println(a)
			fmt.Println(reflect.ValueOf(inter).Kind())

		})
	}
}

func Test_isFieldNull(t *testing.T) {
	var a1 = "aa"
	v1 := reflect.ValueOf(a1)
	is1 := isFieldNull(v1)
	t.Log(is1)

	var a2 string
	v2 := reflect.ValueOf(a2)
	is2 := isFieldNull(v2)
	t.Log(is2)

	var a3 *string
	v3 := reflect.ValueOf(a3)
	is3 := isFieldNull(v3)
	t.Log(is3)

	var a4 *string = new(string)
	v4 := reflect.ValueOf(a4)
	is4 := isFieldNull(v4)
	t.Log(is4)
	t.Log(*a4)

	var a5 *string = nil
	var is5 = toNoNil(a5)
	t.Log(is5)
}
func toNoNil(v any) bool {
	return utils.IsNil(v)
}
