package lsql

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm/types"
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
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("DELETE FROM *").
		WithArgs(1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := engine.Delete(User{}).ByPrimaryKey(1)
	as.Nil(err)
	as.Equal(int64(1), num)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestDeleteByPrimaryKeys(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("DELETE FROM *").
		WithArgs(1, 2, 3).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 3))

	num, err := engine.Delete(User{}).ByPrimaryKey(1, 2, 3)
	as.Nil(err)
	as.Equal(int64(3), num, "num error")

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

type Whe struct {
	Name *string
	Age  *int
	Uid  *types.UUID
}

func TestDeleteByModel(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("DELETE FROM *").
		WithArgs("kk").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM *").
		WithArgs(233, "kk").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := engine.Delete(User{}).ByModel(Whe{
		Name: types.NewString("kk"),
		Age:  nil,
		Uid:  nil,
	})
	as.Nil(err)
	as.Equal(int64(1), num)

	num, err = engine.Delete(User{}).ByModel(Whe{
		Name: types.NewString("kk"),
		Age:  types.NewInt(233),
		Uid:  nil,
	})
	as.Nil(err)
	as.Equal(int64(1), num)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestDeleteByWhere(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("DELETE FROM *").
		WithArgs("kk", "%kk%").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := engine.Delete(User{}).ByWhere(new(WhereBuilder).
		Eq("name", "kk").
		Like("age", "kk"),
	)
	as.Nil(err)
	as.Equal(int64(1), num)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
