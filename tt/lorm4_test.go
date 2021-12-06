package lorm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func create(engine lorm.Engine) {
	v4 := types.NewV4()
	task := Task{
		Name:      types.NewString("990"),
		Num:       types.NewInt(990),
		Info:      types.NewString("990"),
		Position:  types.NewInt(990),
		PatternId: &v4,
	}

	num, err := engine.Table.Create(&task)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(task.String())
}

func createOrUpdById(engine lorm.Engine) {
	v4 := types.NewV4()
	task := Task{
		Id:        &v4,
		Name:      types.NewString("990"),
		Num:       types.NewInt(990),
		Info:      types.NewString("990"),
		Position:  types.NewInt(990),
		PatternId: &v4,
	}

	num, err := engine.Table.CreateOrUpdate(&task).ByPrimaryKey()
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(task.String())
}

func createOrUpdByFields(engine lorm.Engine) {
	v4 := types.NewV4()
	task := Task{
		Id:        &v4,
		Name:      types.NewString("990"),
		Num:       types.NewInt(77),
		Info:      types.NewString("777"),
		Position:  types.NewInt(777),
		PatternId: &v4,
	}

	num, err := engine.Table.CreateOrUpdate(&task).ByUnique([]string{"name"})
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(task.String())
}

func updateByid(engine lorm.Engine) {
	v4 := types.NewV4()
	id := types.Str2UUIDMust("6db8b4a9-9b68-4f3c-953f-78ae69b5b780")
	task := Task{
		Id:        &id,
		Name:      types.NewString("66"),
		Num:       types.NewInt(66),
		Info:      types.NewString("66"),
		Position:  types.NewInt(66),
		PatternId: &v4,
	}

	num, err := engine.Table.Update(&task).ByPrimaryKey()
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(task.String())
}

func selectByid(engine lorm.Engine) {
	task := Task{}
	must := types.Str2UUIDMust("f59882000e474a54b8dd74c42d2a195d")

	num, err := engine.Table.Select(&task).ByPrimaryKey(must)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(task.String())
}

func selectByids(engine lorm.Engine) {
	task := Task{}
	tasks := make([]Task, 0)
	must := types.Str2UUIDMust("a83787f1-655f-4d07-b9a9-be154646534b")
	must2 := types.Str2UUIDMust("6db8b4a9-9b68-4f3c-953f-78ae69b5b780")

	num, err := engine.Table.Select(&tasks).ByPrimaryKey(must, must2)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(task.String())
	fmt.Println(tasks[0].String())
	fmt.Println(tasks[1].String())
}

func delByWhere(engine lorm.Engine) {
	num, err := engine.Table.Delete(Task{}).ByWhere(new(lorm.WhereBuilder).
		Eq("name", "asf").
		Eq("num", 0).
		Ne("name", "asf").
		Like("name", "asf"),
	)
	fmt.Println(num)
	fmt.Println(err)
}

func count(engine lorm.Engine) {
	n := 0

	num, err := engine.Classic.Query("select count(*) from common.t_task ").GetOne(&n)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(n)
}

func selectUUid(engine lorm.Engine) {
	uuid := types.UUID{}

	num, err := engine.Classic.Query("select gen_random_uuid() ").GetOne(&uuid)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(uuid)
}

func selectList(engine lorm.Engine) {
	list := types.StringList{}
	kk := make([]string, 0)

	num, err := engine.Classic.Query("select  ARRAY['os''dba', '123''456']  ").GetOne(&kk)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(list)
	fmt.Println(kk)
}

func TestDB434(t *testing.T) {

	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err)

	ormConfig := lorm.OrmConf{
		TableNamePrefix: "t_",
		PrimaryKeyNames: []string{"id"},
	}
	engine := lorm.MustConnectMock(db, &lorm.PgConf{}).Db(&ormConfig)

	mock.ExpectQuery("select 2").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(2),
		)

	//create(engine) //suc

	//delByid(engine) //suc
	//delByids(engine) //suc
	//updateByid(engine)		//suc

	//count(engine)			//suc

	//selectByid(engine) //suc
	//selectByids(engine)		//suc
	//selectUUid(engine)	//suc
	//selectList(engine)
	//createOrUpdById(engine)
	//createOrUpdByFields(engine)
	//delByModel(engine)
	delByWhere(engine)

}

type Task struct {
	Id        *types.UUID `tableName:"common.t_task"`
	Name      *string
	CreatedAt *types.DateTime
	UpdatedAt *types.DateTime
	Num       *int
	Info      *string
	PatternId *types.UUID
	Position  *int
}

func (conf *Task) String() string {
	b, err := json.Marshal(*conf)
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	return out.String()
}
