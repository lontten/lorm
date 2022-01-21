package lsql

import (
	"fmt"
	"reflect"
	"testing"
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
		want interface{}
	}{
		{
			name: "234 text",
			args: args{v: v},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inter := getFieldInter(tt.args.v)
			fmt.Println(inter)
			fmt.Println(a)
			fmt.Println(reflect.ValueOf(inter).Kind())

		})
	}
}
