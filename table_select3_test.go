package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHasOrInsert_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT 1 FROM t_user WHERE id = ? ORDER BY name ASC;")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "lontten"),
		)

	ok, err := HasOrInsert(engine, W().Eq("id", 1), User{Name: "name"},
		E().ShowSql().OrderBy("name"))
	as.Nil(err)
	as.Equal(true, ok, "has error")
}

func TestHasOrInsert_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT 1 FROM t_user WHERE id = $1 ORDER BY name ASC;")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "lontten"),
		)

	ok, err := HasOrInsert(engine, W().Eq("id", 1), User{Name: "name"},
		E().ShowSql().OrderBy("name"))
	as.Nil(err)
	as.Equal(true, ok, "has error")
}

func TestHasOrInsert2_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT 1 FROM t_user WHERE id = ? ORDER BY name ASC;")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO t_user (name) VALUES (?);")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(2, 1))

	user := User{Name: "tom"}
	ok, err := HasOrInsert(engine, W().Eq("id", 1), &user,
		E().ShowSql().OrderBy("name"))
	as.Nil(err)
	as.Equal(false, ok, "ok error")
	as.Equal(int64(2), user.Id, "id error")
}

func TestHasOrInsert2_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT 1 FROM t_user WHERE id = $1 ORDER BY name ASC;")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"?column?"}))

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO t_user (name) VALUES ($1) RETURNING id;")).
		WithArgs("tom").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(10),
		)

	user := User{Name: "tom"}
	ok, err := HasOrInsert(engine, W().Eq("id", 1), &user,
		E().ShowSql().OrderBy("name"))
	as.Nil(err)
	as.Equal(false, ok, "ok error")
	as.Equal(int64(10), user.Id, "id error")
}
