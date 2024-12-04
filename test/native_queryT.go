package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/types"
	"test/ldb"
)

func QueryOneT() {
	ka, err := lorm.QueryOne[Ka](ldb.DB, "select * from t_ka where id=?", 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(*ka.Id)
	fmt.Println(*ka.Name)
}
func QueryOneT2() {
	ka, err := lorm.QueryOne[types.StringList](ldb.DB, "select img_list from public.user  where id=$1", 6)
	if err != nil {
		panic(err)
	}
	fmt.Println(*ka)
}

func QueryListT() {
	list, err := lorm.QueryList[Ka](ldb.DB, "select * from t_ka where id>1")
	if err != nil {
		panic(err)
	}
	for _, ka := range list {
		fmt.Println(*ka.Id)
		fmt.Println(*ka.Name)
	}

}

func QueryListT2() {
	list, err := lorm.QueryListP[Ka](ldb.DB, "select * from t_ka where id>1")
	if err != nil {
		panic(err)
	}
	for _, ka := range list {
		fmt.Println(*ka.Id)
		fmt.Println(*ka.Name)
	}
}