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
		Page(10, 1).
		PageScan(&u)
	fmt.Println(num, dto, err)
	bytes, err := json.Marshal(dto)
	fmt.Println(string(bytes))

	for _, user := range u {
		bytes2, _ := json.Marshal(user)
		fmt.Println(string(bytes2))
	}

}

func Build2() {
	var u User
	num, err := lorm.QueryBuild(ldb.DB).
		Select("id").Select("name").
		From("t_user").Where("id = ?", 1 == 2).Arg(1, 1 == 2).
		Where("id = 2").Limit(2).
		ScanOne(&u)
	fmt.Println(num, err)
	fmt.Println(*u.Id)
	fmt.Println(*u.Name)
}
