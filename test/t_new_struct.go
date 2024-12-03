package main

import (
	"github.com/lontten/lorm/types"
	"reflect"
)

type Kk struct {
	ID         types.UUID
	CreateTime types.DateTime
	UpdateTime types.DateTime

	KkName string
	KkInfo string
	KkImg  string
}

func get[T any]() T {
	t := new(T)
	zero := reflect.New(reflect.TypeOf(*t)).Elem()
	k := zero.FieldByName("KkName")
	k.SetString("b")
	defaultValue := zero.Interface()
	return defaultValue.(T)
}

// 创建T类型引用
func getOne3[T any]() T {
	return *new(T)
}

// 创建T类型的数组，长度为3
func getOne[T any]() [3]T {
	return [3]T{}
}

// 创建T类型的切片
func getOne2[T any]() []T {
	return []T{}
}
func main7888() {
	var arr = &[]Ka{}
	v := reflect.ValueOf(arr).Elem()

	ka1 := Ka{
		Id:   nil,
		Name: nil,
	}
	ka1v := reflect.ValueOf(ka1)

	v.Set(reflect.Append(v, ka1v))
}
