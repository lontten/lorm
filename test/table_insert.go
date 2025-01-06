package main

import (
	"encoding/json"
	"errors"
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

func TableInsert3() {
	tx, err := ldb.DB.Begin()
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
	Log(user)

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
	Log(user2)
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
