package lsql

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	//-------------base------------

	mock.ExpectQuery("select 1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(4),
		)

	n := 0
	num, err := engine.Query("select 1").GetOne(&n)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(4, n, "n error")

	mock.ExpectQuery("select 'kk' ").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow("kk"),
		)

	name := ""
	num, err = engine.Query("select 'kk' ").GetOne(&name)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal("kk", name, "name error")

	//-------------------uuid---------------

	v4 := types.NewV4()
	mock.ExpectQuery("select gen_random_uuid() ").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(v4),
		)

	uid := types.UUID{}
	num, err = engine.Query("select gen_random_uuid() ").GetOne(&uid)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(v4, uid, "uuid error")

	//-------------------date---------------
	date := types.Date{time.Now()}
	mock.ExpectQuery("select gen_random_uuid() ").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(date),
		)

	d := types.Date{}
	num, err = engine.Query("select gen_random_uuid() ").GetOne(&d)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(date, d, "uuid error")
	as.NotEqual(types.Date{}, d, "date error")

	//-------------------user---------------
	user := User{Id: 1, Name: "lontten"}
	mock.ExpectQuery("select id,name from user limit 1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "lontten"),
		)

	u := User{}
	num, err = engine.Query("select id,name from user limit 1").GetOne(&u)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(user, u, "user error")

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}

func TestExec(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec("delete from user where id = ? ").
		WithArgs(1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err := engine.Exec("delete from user where id = ?", 1)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")

	mock.ExpectExec("update user set name = 'kk' where id = ? ").
		WithArgs(1).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	num, err = engine.Exec("update user set name = 'kk' where id = ? ", 1)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
