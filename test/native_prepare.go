package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"test/ldb"
)

func Prepare() {
	stmt, err := lorm.Prepare(ldb.DB, "delete from t_ka where id= $1 ")
	if err != nil {
		panic(err)
	}
	num, err := stmt.Exec(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
}

func Prepare2() {
	stmt, err := lorm.Prepare(ldb.DB, "select id from t_ka where id=$1 ")
	if err != nil {
		panic(err)
	}
	var n int64
	num, err := stmt.QueryScan(2).ScanOne(&n)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
	fmt.Println(num)
}

func Prepare3() {
	stmt, err := lorm.Prepare(ldb.DB, "select id from t_ka where id>$1 ")
	if err != nil {
		panic(err)
	}
	var list []int64
	num, err := stmt.QueryScan(2).ScanList(&list)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
	fmt.Println(num)
}

func Prepare4() {
	stmt, err := lorm.Prepare(ldb.DB, "select * from t_ka where id>$1 ")
	if err != nil {
		panic(err)
	}
	n, err := lorm.StmtQueryOne[Ka](stmt, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(*n.Id)
	fmt.Println(*n.Name)
}
