package main

import (
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/softdelete"
)

type Ka struct {
	Id   *int    `ldb:"id"`
	Name *string `ldb:"name"`

	softdelete.DeleteGormMilli
}

func (k Ka) TableConf() *lorm.TableConf {
	return new(lorm.TableConf).
		Table("t_ka").
		AutoIncrements("id")
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
	//Del()
	//First()
	//List()
	TableUpdate()

}
