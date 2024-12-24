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

func First2() {
	one, err := lorm.First[User](ldb.DB, lorm.W().Eq("name", "fjakdsf").
		IsNotNull("name"), lorm.E().ShowSql())
	fmt.Println(one)
	fmt.Println(err)
}

func List() {
	num, err := lorm.List[User](ldb.DB, lorm.W().Eq("id", 1), lorm.E().ShowSql())
	fmt.Println(num)
	fmt.Println(err)
}

func GetOrInsert() {
	var u = User{
		Name: types.NewString("kb"),
		Age:  types.NewInt(33),
	}
	d, err := lorm.GetOrInsert[User](ldb.DB, lorm.W().Eq("name", "kb"), &u, lorm.E().ShowSql())
	Log(d)
	Log(u)
	fmt.Println(err)

}

func InsertOrHas() {
	var u = User{
		Name: types.NewString("kc"),
		Age:  types.NewInt(33),
	}
	has, err := lorm.InsertOrHas(ldb.DB, lorm.W().Eq("name", "kc"), &u, lorm.E().ShowSql())
	fmt.Println(has)
	Log(u)
	fmt.Println(err)

}
