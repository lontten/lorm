package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func TestP(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(cap(list))
	fmt.Println(len(list))

	f2(list)
	fmt.Println(list)

}

func f1(list []int) {
	list[0] = 2
}

func f2(list []int) {
	list = append(list, 3)
}

func TestP2(t *testing.T) {
	var a sql.NullTime

	fmt.Println(a)
	v := reflect.ValueOf(a)
	fmt.Println(v.IsValid())
	fmt.Println(v.IsZero())

}
