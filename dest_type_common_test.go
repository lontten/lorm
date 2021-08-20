package lorm

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_checkValidSingleFiled(t *testing.T) {
	var n=make([]Ka,0)

	field := reflect.ValueOf(n)
	fmt.Println(field.Len())
	fmt.Println(field.IsNil())
	fmt.Println(field.IsZero())
	fmt.Println(field.IsValid())


}

type Ka struct {
	Name [][]string
}
