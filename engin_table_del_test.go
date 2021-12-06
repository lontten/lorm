package lorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Id   int64 `lorm:"pk" tableName:"t_user"`
	Name string
}

func TestDeleteByPrimaryKey(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")

	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectPrepare("DELETE FROM *").ExpectExec().
		WithArgs(1, 2).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	num, err := engine.Table.Delete(User{}).ByPrimaryKey(1, 2)
	as.Nil(err)
	as.Equal(int64(2), num)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
