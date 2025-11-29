package lorm

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lcore/v2/types"
	"github.com/stretchr/testify/assert"
)

func TestSelectByPrimaryKey(t *testing.T) {
	//as := assert.New(t)
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("failed to open sqlmock database: %s", err)
	//}
	//defer db.Close()
	//engine := MustConnectMock(db, &MysqlConf{})
	//
	//mock.ExpectQuery("SELECT ").
	//	WithArgs(1).
	//	WillReturnError(nil).
	//	WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))
	//
	//user, err := First[User](engine, W().PrimaryKey(1), E().ShowSql())
	//as.Nil(err)
	//as.NotNil(user)
	//as.Equal(int64(1), user.Id)
	//as.Equal("test", user.Name)
	//
	//as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestSelectByPrimaryKeys(t *testing.T) {
	//as := assert.New(t)
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("failed to open sqlmock database: %s", err)
	//}
	//defer db.Close()
	//engine := MustConnectMock(db, &MysqlConf{})
	//
	//mock.ExpectQuery("SELECT ").
	//	WithArgs(1, 2).
	//	WillReturnError(nil).
	//	WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
	//		AddRow(1, "test").
	//		AddRow(2, "test2"),
	//	)
	//
	//users, err := List[User](engine, W().PrimaryKey(1, 2), E().ShowSql())
	//as.Nil(err)
	//as.Equal(2, len(users))
	//as.Equal(int64(1), users[0].Id)
	//as.Equal("test", users[0].Name)
	//as.Equal(int64(2), users[1].Id)
	//as.Equal("test2", users[1].Name)
	//
	//as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestSelectByModel(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	defer db.Close()
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectQuery("SELECT *").
		WithArgs("123").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test"),
		)

	mock.ExpectQuery("SELECT *").
		WithArgs("kk", 233).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test").
			AddRow(2, "test2"),
		)

	user, err := First[User](engine, W().Model(Whe{
		Name: types.NewString("123"),
		Age:  nil,
		Uid:  nil,
	}), E().ShowSql())
	as.Nil(err)
	as.NotNil(user)
	as.Equal(int64(1), user.Id)
	as.Equal("test", user.Name)

	users, err := List[User](engine, W().Model(Whe{
		Name: types.NewString("kk"),
		Age:  types.NewInt(233),
		Uid:  nil,
	}), E().ShowSql())
	as.Nil(err)
	as.Equal(2, len(users))
	as.Equal(int64(1), users[0].Id)
	as.Equal("test", users[0].Name)
	as.Equal(int64(2), users[1].Id)
	as.Equal("test2", users[1].Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestSelectByWhere(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()
	engine := MustConnectMock(db, &MysqlConf{})

	//---------------------------scan one------------------------

	mock.ExpectQuery("SELECT *").
		WithArgs("kk").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test"),
		)

	user, err := First[User](engine, W().Eq("name", "kk"), E().ShowSql())
	as.Nil(err)
	as.NotNil(user)
	as.Equal(int64(1), user.Id)
	as.Equal("test", user.Name)

	//---------------------------scan list------------------------
	mock.ExpectQuery("SELECT *").
		WithArgs("kk", 233).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test").
			AddRow(2, "test2"),
		)
	users, err := List[User](engine, W().
		Eq("name", "kk").
		Eq("age", 233),
		E().ShowSql())
	as.Nil(err)
	as.Equal(2, len(users))
	as.Equal(int64(1), users[0].Id)
	as.Equal("test", users[0].Name)
	as.Equal(int64(2), users[1].Id)
	as.Equal("test2", users[1].Name)

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
