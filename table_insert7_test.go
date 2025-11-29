package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsert7_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO t_user (name) VALUES (?);")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(10, 1))

	var u = User{
		Id:   0,
		Name: "tom",
	}
	num, err := Insert(engine, &u, E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(int64(10), u.Id, "id error")
}

func TestInsert7_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO t_user (name) VALUES ($1) RETURNING id;")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(10),
		)

	var u = User{
		Id:   0,
		Name: "tom",
	}
	num, err := Insert(engine, &u, E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(int64(10), u.Id, "id error")
}
