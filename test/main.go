package main

import (
	"github.com/lontten/lorm"
	soft_del "github.com/lontten/lorm/soft-delete"
)

type Ka struct {
	Id   *int    `ldb:"id"`
	Name *string `ldb:"name"`

	soft_del.DeleteGormMilli
}

func (k Ka) TableConf() *lorm.TableConf {
	return new(lorm.TableConf).
		Table("t_ka").AutoIncrements("id")
}

func (u User) TableConf() *lorm.TableConf {
	return new(lorm.TableConf).
		Table("t_user").AutoIncrements("id")
}

type User struct {
	Id   *int
	Name *string
	Age  *int

	soft_del.DeleteGormMilli
}

func main() {
	//QueryOneT2()
	//QueryListT2()
	//
	QueryOne1()
	//QueryList()
	//QueryList2()
	//
	//Prepare4()

	//TableInsert()
	//Build1()
	//Build2()
}
