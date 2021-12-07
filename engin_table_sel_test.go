package lorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectByPrimaryKey(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectQuery("SELECT *").
		WithArgs(1).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))

	user := User{}
	num, err := engine.Table.Select(user).ByPrimaryKey(1).ScanOne(&user)
	as.Nil(err)
	as.Equal(int64(1), num)
	as.Equal(int64(1), user.Id)
	as.Equal("test", user.Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestSelectByPrimaryKeys(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectQuery("SELECT *").
		WithArgs(1, 2).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test").
			AddRow(2, "test2"),
		)

	users := make([]User, 0)
	num, err := engine.Table.Select(User{}).ByPrimaryKey(1, 2).ScanList(&users)
	as.Nil(err)
	as.Equal(int64(2), num)
	as.Equal(2, len(users))
	as.Equal(int64(1), users[0].Id)
	as.Equal("test", users[0].Name)
	as.Equal(int64(2), users[1].Id)
	as.Equal("test2", users[1].Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestSelectByModel(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")

	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectExec("SELECT *").
		WithArgs("kk").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("SELECT *").
		WithArgs("kk", 233).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := engine.Table.Delete(User{}).ByModel(Whe{
		Name: types.NewString("kk"),
		Age:  nil,
		Uid:  nil,
	})
	as.Nil(err)
	as.Equal(int64(1), num)

	num, err = engine.Table.Delete(User{}).ByModel(Whe{
		Name: types.NewString("kk"),
		Age:  types.NewInt(233),
		Uid:  nil,
	})
	as.Nil(err)
	as.Equal(1, num)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
