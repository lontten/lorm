package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lcore/v2/types"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Id       int64
	Name     string
	Name2    string
	Age      int
	Age2     int
	Birthday types.LocalDate
}

func (User) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}

func TestInsert_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO t_user (name) VALUES (?);")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = User{
		Id:   0,
		Name: "tom",
	}
	num, err := Insert(engine, u, E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestInsert_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO t_user (name) VALUES ($1);")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = User{
		Id:   0,
		Name: "tom",
	}
	num, err := Insert(engine, u, E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestInsert_2(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO t_user (name) VALUES (?);")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = struct {
		Id int64
	}{
		Id: 1,
	}
	num, err := Insert(engine, u, E().ShowSql())
	as.ErrorIs(err, ErrNoTableName)
	as.Equal(int64(0), num, "num error")
}
