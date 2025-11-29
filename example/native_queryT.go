package main

import (
	"encoding/json"
	"example/dbinit"
	"fmt"

	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm"
)

func QueryOneT() {
	ka, err := lorm.NativeQuery[User](dbinit.DB,
		"select * from t_user where id=?", 2222).One()
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(ka)
	fmt.Println(string(bytes))
}
func QueryOneT2() {
	ka, err := lorm.NativeQuery[types.StringList](dbinit.DB,
		"select img_list from public.user  where id=$1", 6).One()
	if err != nil {
		panic(err)
	}
	fmt.Println(*ka)
}
