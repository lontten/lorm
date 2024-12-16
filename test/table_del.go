package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/types"
	"test/ldb"
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
	eq := lorm.Wb().Eq("abc", "xxx")

	num, err := lorm.Delete[User](ldb.DB, lorm.Wb().Or(eq).
		Model(u), lorm.Extra().ShowSql().SkipSoftDelete())
	fmt.Println(num)
	fmt.Println(err)
}
