package main

import (
	"encoding/json"
	"fmt"
	"github.com/lontten/lorm"
	return_type "github.com/lontten/lorm/return-type"
	"github.com/lontten/lorm/types"
	"test/ldb"
	"time"
)

func TableInsert() {
	var user = User{
		Name: types.NewString(time.Now().String()),
	}
	num, err := lorm.Insert(ldb.DB, &user, lorm.E().
		SetNull("abc").
		TableName("t_user").
		ReturnType(return_type.PrimaryKey).
		WhenDuplicateKey("name").DoUpdate().
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
	num, err := lorm.Insert(ldb.DB, &user, new(lorm.ExtraContext).
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
