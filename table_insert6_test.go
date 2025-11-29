package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsert6_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO x_user (name) VALUES (?);")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = User{
		Id:   0,
		Name: "tom",
	}
	num, err := Insert(engine, u, E().ShowSql().TableName("x_user"))
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestInsert6_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO x_user (id) VALUES ($1);")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = struct {
		Id int64
	}{
		Id: 1,
	}
	num, err := Insert(engine, u, E().ShowSql().TableName("x_user"))
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}
