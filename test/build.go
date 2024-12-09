package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"test/ldb"
)

func Build1() {
	var u User
	num, err := lorm.SelectBuilder(ldb.DB).
		Select("id").Select("name").
		From("t_user").Where("id = 1").Where("id = 2").Limit(2).
		ScanOne(&u)
	fmt.Println(num, err)
	fmt.Println(*u.Id)
	fmt.Println(*u.Name)
}

func Build2() {
	var u User
	num, err := lorm.SelectBuilder(ldb.DB).
		Select("id").Select("name").
		From("t_user").Where("id = ?", 1 == 2).Arg(1, 1 == 2).
		Where("id = 2").Limit(2).
		ScanOne(&u)
	fmt.Println(num, err)
	fmt.Println(*u.Id)
	fmt.Println(*u.Name)
}
