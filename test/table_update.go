package main

import (
	"encoding/json"
	"fmt"
	"github.com/lontten/lorm"
	return_type "github.com/lontten/lorm/return-type"
	"github.com/lontten/lorm/types"
	"test/ldb"
)

func TableUpdate() {
	var user = User{
		Name: types.NewString("abc"),
	}
	num, err := lorm.Update(ldb.DB, lorm.W(), &user, lorm.E().
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
