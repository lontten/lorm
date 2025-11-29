package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDelete2_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM x_user WHERE name = ?;")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().Eq("name", "tom"), E().ShowSql().TableName("x_user"))
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestDelete2_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM x_user WHERE name = $1;")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().Eq("name", "tom"), E().ShowSql().TableName("x_user"))
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}
