package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"test/ldb"
)

func QueryOne() {
	var ka Ka
	num, err := lorm.QueryScan(ldb.DB, "select * from t_ka where id=?", 2).ScanOne(&ka)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	fmt.Println(*ka.Id)
	fmt.Println(*ka.Name)
}
func QueryOne1() {
	var n int
	num, err := lorm.QueryScan(ldb.DB, "select 1").ScanOne(&n)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	fmt.Println(n)
}

func QueryList() {
	var list []Ka
	num, err := lorm.QueryScan(ldb.DB, "select * from t_ka where id>1").ScanList(list)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)

	for _, ka := range list {
		fmt.Println(*ka.Id)
		fmt.Println(*ka.Name)
	}
}

func QueryList2() {
	var list []Ka
	num, err := lorm.QueryScan(ldb.DB, "select * from t_ka where id>1").ScanList(&list)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)

	for _, ka := range list {
		fmt.Println(*ka.Id)
		fmt.Println(*ka.Name)
	}
}
