package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateByPrimaryKey_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("UPDATE t_user SET name = ? WHERE id = ?;")).
		WithArgs("tom", 10).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := UpdateByPrimaryKey(engine, User{Name: "tom", Id: 10}, E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestUpdateByPrimaryKey_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("UPDATE t_user SET name = $1 WHERE id = $2;")).
		WithArgs("tom", 10).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := UpdateByPrimaryKey(engine, User{Name: "tom", Id: 10}, E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}
