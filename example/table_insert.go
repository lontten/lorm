package main

import (
	"encoding/json"
	"errors"
	"example/dbinit"
	"fmt"

	"github.com/lontten/lcore/v2/logutil"
	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm"
	return_type "github.com/lontten/lorm/return-type"
)

func TableInsert() {
	var user = dbinit.TestModel{
		Id:   types.NewInt(1),
		Name: types.NewString("cc"),
	}
	num, err := lorm.Insert(dbinit.DB, &user, lorm.E().
		TableName("t_test").
		ReturnType(return_type.Auto).
		WhenDuplicateKey().DoUpdate().
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
	num, err := lorm.Insert(dbinit.DB, &user, new(lorm.ExtraContext).
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

func TableInsert3() {
	tx, err := dbinit.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return
		}
	}()
	var user = User{
		Name: types.NewString("99"),
	}
	num, err := lorm.Insert(tx, &user, lorm.E().
		ShowSql(),
	)
	err = errors.New("db tx")
	if err != nil {
		panic(err)
	}

	fmt.Println(num)
	logutil.Log(user)

	var user2 = User{
		Name: types.NewString("000"),
	}
	num, err = lorm.Insert(tx, &user2, lorm.E().
		ShowSql(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	logutil.Log(user2)
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
