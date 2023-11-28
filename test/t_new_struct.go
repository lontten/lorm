package main

import (
	"github.com/lontten/lorm/types"
	"reflect"
)

func main() {
	k := get[Kk]()
	println(k.KkName)

}

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
