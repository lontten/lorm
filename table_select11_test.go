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

type User11 struct {
	Id *int64
}

func (User11) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst11_mysql(t *testing.T) {
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

	user, err := First[User11](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
}

type User12 struct {
	Id  *int64
	Val *string
}

func (User12) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst12_mysql(t *testing.T) {
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

	user, err := First[User12](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.Nil(user.Val, "Val error")
}

type User13 struct {
	Id  *int64
	Val *time.Time
}

func (User13) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst13_mysql(t *testing.T) {
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

	user, err := First[User13](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.Nil(user.Val, "Val error")
}

type User14 struct {
	Id  *int64
	Val *types.LocalDateTime
}

func (User14) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst14_mysql(t *testing.T) {
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

	user, err := First[User14](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.Nil(user.Val, "Val error")
}

type User15 struct {
	Id  *int64
	Val *decimal.Decimal
}

func (User15) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}
func TestFirst15_mysql(t *testing.T) {
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

	user, err := First[User15](engine, W().Eq("id", 11), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(11), *user.Id, "id error")
	as.Nil(user.Val, "Val error")
}
