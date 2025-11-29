package main

import (
	"example/dbinit"
	"fmt"

	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm"
)

func Del() {
	var u = User{
		Id:   nil,
		Name: types.NewString("xxx"),
	}
	fmt.Println(u)
	var m = make(map[string]any)
	m["a"] = 1
	m["b"] = "bb"
	m["c"] = nil

	num, err := lorm.Delete[User](dbinit.DB, lorm.W().PrimaryKey(1), lorm.E().ShowSql())
	fmt.Println(num)
	fmt.Println(err)
}
