package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"test/ldb"
)

func Exec() {
	num, err := lorm.Exec(ldb.DB, "delete from t_ka where id=?", 1)
	fmt.Println(num, err)
}
