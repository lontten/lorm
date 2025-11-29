package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDelete_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM t_user WHERE name = ?;")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().Eq("name", "tom"), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestDelete_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM t_user WHERE name = $1;")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().Eq("name", "tom"), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}
