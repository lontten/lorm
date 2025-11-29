package lorm

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lcore/v2/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type User21 struct {
	Id int64
}

func (User21) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst21_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id FROM t_user WHERE id = ? ORDER BY id DESC LIMIT 1;")).
		WithArgs(11).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(11),
		)

	user, err := First[User21](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), user.Id, "id error")
}

type User22 struct {
	Id  *int64
	Val string
}

func (User22) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst22_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id ,val FROM t_user WHERE id = ? ORDER BY id DESC LIMIT 1;")).
		WithArgs(11).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "val"}).
			AddRow(11, nil),
		)

	user, err := First[User22](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.Equal("", user.Val, "val error")
}

type User23 struct {
	Id  *int64
	Val time.Time
}

func (User23) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst23_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id ,val FROM t_user WHERE id = ? ORDER BY id DESC LIMIT 1;")).
		WithArgs(11).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "val"}).
			AddRow(11, nil),
		)

	user, err := First[User23](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.True(user.Val.IsZero(), "val error")
}

type User24 struct {
	Id  *int64
	Val types.LocalDateTime
}

func (User24) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst24_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id ,val FROM t_user WHERE id = ? ORDER BY id DESC LIMIT 1;")).
		WithArgs(11).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "val"}).
			AddRow(11, nil),
		)

	user, err := First[User24](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.True(user.Val.IsZero(), "val error")
}

type User25 struct {
	Id  *int64
	Val decimal.Decimal
}

func (User25) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst25_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id ,val FROM t_user WHERE id = ? ORDER BY id DESC LIMIT 1;")).
		WithArgs(11).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "val"}).
			AddRow(11, nil),
		)

	user, err := First[User25](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.True(user.Val.IsZero(), "val error")
}
