package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDelete3_mysql1(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM t_user WHERE id IN (?,?);")).
		WithArgs(1, 2).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().PrimaryKey(1, 2), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestDelete3_mysql2(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM t_user WHERE id IN (?);")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().PrimaryKey(1), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestDelete3_pg1(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM t_user WHERE id IN ($1,$2);")).
		WithArgs(1, 2).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().PrimaryKey(1, 2), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestDelete3_pg2(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM t_user WHERE id IN ($1);")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := Delete[User](engine, W().PrimaryKey(1), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}
