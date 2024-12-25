package main

import (
	"encoding/json"
	"fmt"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/softdelete"
)

func Log(v any) {
	bytes, _ := json.Marshal(v)
	fmt.Println(string(bytes))
}

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
	return new(lorm.TableConf).Table("t_user").
		PrimaryKeys("id").
		AutoIncrements("id")
}

type Kaaa struct {
}
type User struct {
	Id   *int
	Name *string
	Age  *int
	Kaaa
	softdelete.DeleteTimeNil
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
	Del()
	//First()
	//First2()
	//QueryOneT()
	//List()
	//TableUpdate()
	//GetOrInsert()
	//InsertOrHas()
}
