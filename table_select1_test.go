package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestFirst_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id ,name ,name2 ,age ,age2 ,birthday FROM t_user WHERE id = ? ORDER BY id DESC LIMIT 1;")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "lontten"),
		)

	user, err := First[User](engine, W().Eq("id", 1), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), user.Id, "id error")
	as.Equal("lontten", user.Name, "name error")
}

func TestFirst_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id ,name ,name2 ,age ,age2 ,birthday FROM t_user WHERE id = $1 ORDER BY id DESC LIMIT 1;")).
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "lontten"),
		)

	user, err := First[User](engine, W().Eq("id", 1), E().ShowSql())
	as.Nil(err)
	as.Equal(int64(1), user.Id, "id error")
	as.Equal("lontten", user.Name, "name error")
}
