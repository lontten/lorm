package lorm

import (
	"context"
	"fmt"
	"github.com/lontten/lorm/types"
	"reflect"
	"testing"
)

// ptr struct
func doStruct(v interface{}) {
	va := reflect.ValueOf(v)
	va = va.Elem()
	va.FieldByName("Name").SetString("waaaaaao")
}

func Test_struct(t *testing.T) {
	ha := Haa{types.NewV4()}

	fmt.Println(ha)
	doStruct(&ha)
	fmt.Println(ha)
	fmt.Println(ha.Name)
}

// ptr struct ptr
func doStruct_p(v interface{}) {
	va := reflect.ValueOf(v)
	va = va.Elem()    //1. struct 必须为 ptr

	fieldByName := va.FieldByName("Name")
	if fieldByName.Kind() == reflect.Ptr { // 2. 判断 filed 是否为 ptr ，是则
		newString := "wosadfaf"
		fieldByName.Set(reflect.ValueOf(&newString))  // 使用 &
	}
	fmt.Println(fieldByName.Kind().String())
	//base
	if fieldByName.Kind() == reflect.Struct {   //3. struct int 根据不同类型创建不同数据
		newString := types.NewV4()
		fieldByName.Set(reflect.ValueOf(newString))
	}
if fieldByName.Kind() == reflect.Int {
		newString := 23
		fieldByName.Set(reflect.ValueOf(newString))
	}

}

func Test_struct_p(t *testing.T) {
	ha := HaaPtr{ 2}

	fmt.Println(ha)
	doStruct_p(&ha)
	fmt.Println(ha)
	fmt.Println(ha.Name)
}

// slice type
func Test_slice(t *testing.T) {
	haas := make([]Haa, 0)
	ha := Haa{types.NewV4()}
	haas = append(haas, ha)

	fmt.Println(haas[0])
	do_slice(haas)
	fmt.Println(haas[0])
	fmt.Println(haas[0].Name)
context.Background()
}

type Haa struct {
	Name types.UUID
}

func do_slice(v interface{}) {
	va := reflect.ValueOf(v)
	va = va.Index(0)
	va.FieldByName("Name").SetString("waaaaaao")
}

func Test_slice_ptr(t *testing.T) {
	haas := make([]Haa, 0)
	ha := Haa{types.NewV4()}
	haas = append(haas, ha)

	fmt.Println(haas[0])
	do_slice_ptr(&haas)
	fmt.Println(haas[0])
	fmt.Println(haas[0].Name)

}

type HaaPtr struct {
	Name int
}

func do_slice_ptr(v interface{}) {
	va := reflect.ValueOf(v)
	va=reflect.Indirect(va)  // 若 slice ptr 取base
	va = va.Index(0)  //取 第n个， 按照 单个执行
	newString := types.NewV4()
	va.FieldByName("Name").Set(reflect.ValueOf(newString))
	//va.FieldByName("Name").SetString("waaaaaao")
}

func TestBasePtrDeepType(t *testing.T) {
	v4 := types.NewV4()
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 reflect.Type
	}{
		{
			name: "deep",
			args: args{t: reflect.TypeOf(&v4)},
			want: true,
			want1: reflect.TypeOf(v4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := basePtrDeepType(tt.args.t)
			if got != tt.want {
				t.Errorf("basePtrDeepType() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("basePtrDeepType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}