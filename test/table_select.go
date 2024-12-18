package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/types"
	"test/ldb"
)

func First() {
	var u = User{
		Id:   nil,
		Name: types.NewString("xxx"),
	}
	fmt.Println(u)
	var m = make(map[string]any)
	m["a"] = 1
	m["b"] = "bb"
	m["c"] = nil
	eq := lorm.W().Eq("abc", "xxx")

	num, err := lorm.First[User](ldb.DB, lorm.W().Or(eq).
		Model(u), lorm.E().ShowSql().SkipSoftDelete())
	fmt.Println(num)
	fmt.Println(err)
}

func List() {
	num, err := lorm.List[User](ldb.DB, lorm.W().Eq("id", 1), lorm.E().ShowSql())
	fmt.Println(num)
	fmt.Println(err)
}
