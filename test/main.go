package main

import (
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/soft_delete"
)

type Ka struct {
	Id   *int    `ldb:"id"`
	Name *string `ldb:"name"`

	soft_delete.DeleteGormMilli
}

func (k Ka) TableConf() *lorm.TableConf {
	return new(lorm.TableConf).Table("t_ka")
}

type User struct {
	Id   *int
	Name *string
	Age  *int

	//UserState *int    `ldb:"user_state"`

	//Uuid uuid.UUID `ldb:"uuid"`

	soft_delete.DeleteGormMilli
}

func main() {
	//QueryOneT2()
	//QueryListT2()
	//
	//QueryOne1()
	//QueryList()
	//QueryList2()
	//
	//Prepare4()

	TableInsert()
}
