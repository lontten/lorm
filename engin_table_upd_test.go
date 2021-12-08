package lorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserP struct {
	Name *string
	Id   *int
}

func TestUpdateByPrimaryKey(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectExec("UPDATE *").
		WithArgs(1, "nn", 1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user := UserP{
		Id:   types.NewInt(1),
		Name: types.NewString("nn"),
	}
	num, err := engine.Table.Update(&user).ByPrimaryKey()
	as.Nil(err)
	as.Equal(int64(1), num)
	as.Equal(1, *user.Id)
	as.Equal("nn", *user.Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestUpdateByModel(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectExec("UPDATE *").
		WithArgs(1, "nn", 22, "nmmn").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user := UserP{
		Id:   types.NewInt(1),
		Name: types.NewString("nn"),
	}
	num, err := engine.Table.Update(&user).ByModel(struct {
		Age  *int
		Name *string
	}{
		Age:  types.NewInt(22),
		Name: types.NewString("nmmn"),
	})
	as.Nil(err)
	as.Equal(int64(1), num)
	as.Equal(1, *user.Id)
	as.Equal("nn", *user.Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestUpdateByWhere(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectExec("UPDATE *").
		WithArgs(1, "nn", "name_name", 233).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user := UserP{
		Id:   types.NewInt(1),
		Name: types.NewString("nn"),
	}
	num, err := engine.Table.Update(&user).ByWhere(new(WhereBuilder).
		Eq("name", "name_name").
		Eq("age", 233),
	)
	as.Nil(err)
	as.Equal(int64(1), num)
	as.Equal(1, *user.Id)
	as.Equal("nn", *user.Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
