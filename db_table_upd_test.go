package lsql

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lsql/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserP struct {
	Name *string
	Id   *int
}
type UserUuid struct {
	Name *string
	ID   *types.UUID ``
}

func TestUpdateByPrimaryKey(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	v4 := types.NewV4()

	mock.ExpectExec("UPDATE *").
		WithArgs(v4, "nn", v4).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := engine.Update(&UserUuid{
		ID:   &v4,
		Name: types.NewString("nn"),
	}).ByPrimaryKey()
	as.Nil(err)
	as.Equal(int64(1), num)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestUpdateByModel(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("UPDATE *").
		WithArgs(1, "nn", 22, "nmmn").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user := UserP{
		Id:   types.NewInt(1),
		Name: types.NewString("nn"),
	}
	num, err := engine.Update(&user).ByModel(struct {
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
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("UPDATE *").
		WithArgs(1, "nn", "name_name", 233).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user := UserP{
		Id:   types.NewInt(1),
		Name: types.NewString("nn"),
	}
	num, err := engine.Update(&user).ByWhere(new(WhereBuilder).
		Eq("name", "name_name").
		Eq("age", 233),
	)
	as.Nil(err)
	as.Equal(int64(1), num)
	as.Equal(1, *user.Id)
	as.Equal("nn", *user.Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
