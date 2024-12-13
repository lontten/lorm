package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/softdelete"
	"github.com/lontten/lorm/types"
	"test/ldb"
)

type Ka struct {
	Id   *int    `ldb:"id"`
	Name *string `ldb:"name"`

	softdelete.DeleteGormMilli
}

func (k Ka) TableConf() *lorm.TableConf {
	return new(lorm.TableConf).
		Table("t_ka").AutoIncrements("id")
}

func (u User) TableConf() *lorm.TableConf {
	return new(lorm.TableConf).PrimaryKeys("id", "name").
		Table("t_user").AutoIncrements("id")
}

type User struct {
	Id   *int
	Name *string
	Age  *int

	softdelete.DeleteGormMilli
}

func main() {
	//QueryOneT2()
	//QueryListT2()
	//
	//QueryOneT()
	//QueryList2()
	//
	//Prepare4()
	//time.Sleep(1 * time.Hour)
	//TableInsert()
	//Build1()
	//Build2()

	var u = User{
		Id:   nil,
		Name: types.NewString("xxx"),
		Age:  types.NewInt(1),
	}
	fmt.Println(u)
	var m = make(map[string]any)
	m["a"] = 1
	m["b"] = "bb"
	m["c"] = nil

	lorm.Delete[User](ldb.DB, lorm.Wb().PrimaryKey(User{
		Id:   types.NewInt(1),
		Name: types.NewString("a"),
	}, 2, 3).FilterPrimaryKey(2, 2, 2))
}
