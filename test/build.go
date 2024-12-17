package main

import (
	"encoding/json"
	"fmt"
	"github.com/lontten/lorm"
	"test/ldb"
)

func Build1() {
	var u []User
	num, dto, err := lorm.QueryBuild(ldb.DB).ShowSql().
		Select("*").
		From("t_user u").
		Native(`
`).
		Where("id > 1").
		Page(1, 10).
		ScanPage(&u)
	fmt.Println(num, dto, err)
	bytes, err := json.Marshal(dto)
	fmt.Println(string(bytes))

	for _, user := range u {
		bytes2, _ := json.Marshal(user)
		fmt.Println(string(bytes2))
	}

}

func Build2() {
	neq := lorm.W().Neq("id", 7)
	var u User
	num, err := lorm.QueryBuild(ldb.DB).
		Select("id").Select("name").
		From("t_user").
		Where("id = 2").
		WhereBuilder(lorm.W().
			Eq("id", 3).
			Eq("id", 5).
			Or(neq)).
		Limit(2).ShowSql().
		ScanOne(&u)
	fmt.Println(num, err)
	bytes, err := json.Marshal(u)
	fmt.Println(string(bytes))
}
