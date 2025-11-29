package main

import (
	"github.com/lontten/lorm"
)

type User struct {
	Id   *int `db:"id"`
	Name *string
	Age  string
}

func (u User) TableConf() *lorm.TableConfContext {
	return new(lorm.TableConfContext).Table("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}

func main() {
	QueryBuild()
	//TableInsert()
}
