package lorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// a successful case
func TestClassic(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err)

	engine := lorm.MustConnectMock(db, &lorm.PgConf{}).Db(nil)

	//set return values
	mock.ExpectQuery("select 2").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(2),
		)

	var name = 0
	num, err := engine.Classic.Query("select 2 ").GetOne(&name)
	as.Nil(err)
	as.Equal(int64(1), num)
	as.Equal(2, name)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
