package lorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{}).Db(nil)

	mock.ExpectQuery("a.t_num*").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(4),
		)

	n := 0
	num, err := engine.Classic.Query("select count(*) from a.t_num ").GetOne(&n)
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
	as.Equal(4, n, "n error")

	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
}
