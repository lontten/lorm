package lorm

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_checkValidSingleFiled(t *testing.T) {
	ka := Ka{}
	field := reflect.ValueOf(ka).Field(0)

	is, has, base := singleFieldBaseSliceValue(field)
	fmt.Println(is)
	fmt.Println(has)
	fmt.Println(base)

}

type Ka struct {
	Name [][]string
}
