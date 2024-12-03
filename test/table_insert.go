package main

import (
	"encoding/json"
	"fmt"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/types"
	"test/ldb"
)

func TableInsert() {
	var user = User{
		Name: types.NewString("abc"),
		Age:  types.NewInt(44),
	}
	num, err := lorm.Insert(ldb.DB, &user, new(lorm.Extra).
		TableName("t_user").
		WhenDuplicateKey("name").DoUpdate(lorm.Set()).
		ShowSql(),
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

func TableInsert2() {
	var user = User{
		Name: types.NewString("abc"),
	}
	num, err := lorm.Insert(ldb.DB, &user, new(lorm.Extra).
		TableName("t_user2").
		SetNull("uuid").
		WhenDuplicateKey().DoUpdate(lorm.Set().Set("user_state", 1).SetNull("name")).
		ShowSql(),
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
