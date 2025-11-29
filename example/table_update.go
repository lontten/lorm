package main

import (
	"encoding/json"
	"example/dbinit"
	"fmt"

	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm"
	return_type "github.com/lontten/lorm/return-type"
)

func TableUpdate() {
	var user = User{
		Name: types.NewString("abc"),
	}
	num, err := lorm.Update(dbinit.DB, &user, lorm.W(), lorm.E().
		SetNull("abc").
		TableName("t_user").
		ReturnType(return_type.Auto).
		WhenDuplicateKey("name").DoUpdate().
		ShowSql().NoRun(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	bytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func TableUpdate2() {
	var user = User{
		Name: types.NewString("abc"),
	}
	num, err := lorm.Update(dbinit.DB, &user, lorm.W().
		Eq("id", 1).
		In("id", 1, 2).
		Gt("id", 1).
		IsNull("name").
		Like("name", "abc"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	bytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
