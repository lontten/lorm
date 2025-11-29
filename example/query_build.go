package main

import (
	"example/dbinit"
	"fmt"

	"github.com/lontten/lcore/v2/logutil"
	"github.com/lontten/lorm"
)

func QueryBuild() {

	list, err := lorm.QueryBuild[User](dbinit.DB).ShowSql().
		Select("u.*").
		From("t_user u").
		Convert(lorm.ConvertRegister("age", func(v *int) any {
			fmt.Println(v, v == nil)
			if v == nil {
				return "kk"
			}
			if *v == 1 {
				return "one"
			}
			return "abc"
		})).
		List()
	if err != nil {
		panic(err)
	}

	for _, v := range list {
		logutil.Log(v)
	}
}
