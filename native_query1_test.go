package lorm

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lcore/v2/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type Ka struct {
	Id   int
	Name string
	Day1 types.LocalDate
	Day2 types.LocalDate
}

func (k Ka) TableConf() *TableConfContext {
	return TableConf("t_ka").
		PrimaryKeys("id").
		AutoColumn("id")
}

func TestQuery1_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectQuery("select q1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "day1", "day2"}).
			AddRow(12, "nil", nil, "2022-02-02"),
		)

	list, err := NativeQuery[Ka](engine, "select q1").List()
	as.Nil(err)
	as.NotNil(list)
	as.Equal(1, len(list), "list length error")
	ka := list[0]
	as.Equal(12, ka.Id, "id error")
	as.Equal("nil", ka.Name, "name error")
	as.True(ka.Day1.IsZero(), "day1 error")
	as.Equal("2022-02-02", ka.Day2.String(), "day2 error")
}

type UserNil2 struct {
	Id    int
	Name  string
	Money decimal.Decimal
	Day1  types.LocalDate
	Day2  time.Time
}

func (u UserNil2) TableConf() *TableConfContext {
	return TableConf("t_user").
		PrimaryKeys("id").
		AutoColumn("id")
}

func TestQuery2_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, "new sqlmock error")
	engine := MustConnectMock(db, &PgConf{})

	column1 := mock.NewColumn("id").OfType("int", 12).Nullable(false)
	column2 := mock.NewColumn("name").OfType("VARCHAR", nil).Nullable(true)
	column3 := mock.NewColumn("money").OfType("DECIMAL", nil).Nullable(true).WithPrecisionAndScale(10, 4)
	column4 := mock.NewColumn("day1").OfType("TIMESTAMP", nil).Nullable(true)
	column5 := mock.NewColumn("day2").OfType("TIMESTAMP", nil).Nullable(true)
	rows := mock.NewRowsWithColumnDefinition(column1, column2, column3, column4, column5)
	rows.AddRow(12, nil, nil, nil, nil)

	mock.ExpectQuery("select q1").
		WillReturnError(nil).
		WillReturnRows(rows)

	list, err := NativeQuery[UserNil2](engine, "select q1").List()
	as.Nil(err)
	as.NotNil(list)
	as.Equal(1, len(list), "list length error")
	ka := list[0]
	as.Equal(12, ka.Id, "id error")
	as.Equal("", ka.Name, "name error")
	as.True(ka.Day1.IsZero(), "day1 error")
	as.True(ka.Day2.IsZero(), "day2 error")
}
