package lorm

//import (
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/lontten/lorm/types"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestHasByPrimaryKey(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	mock.ExpectQuery("t_user *").
//		WithArgs(1).
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(1))
//
//	has, err := engine.Has(User{}).ByPrimaryKey(1)
//	as.Nil(err)
//	as.Equal(true, has)
//
//	mock.ExpectQuery("t_false *").
//		WithArgs(1).
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}))
//
//	has, err = engine.Has("t_false").ByPrimaryKey(1)
//	as.Nil(err)
//	as.Equal(false, has)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
//
//func TestHasByPrimaryKeys(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	mock.ExpectQuery("t_user *").
//		WithArgs(1, 2).
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(1))
//
//	has, err := engine.Has(User{}).ByPrimaryKey(1, 2)
//	as.Nil(err)
//	as.Equal(true, has)
//
//	mock.ExpectQuery("t_false *").
//		WithArgs(1, 2).
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}))
//
//	has, err = engine.Has("t_false").ByPrimaryKey(1, 2)
//	as.Nil(err)
//	as.Equal(false, has)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
//
//func TestHasByModel(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	mock.ExpectQuery("t_user *").
//		WithArgs("kk").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}).
//			AddRow(1),
//		)
//
//	has, err := engine.Has(User{}).ByModel(Whe{
//		Name: types.NewString("kk"),
//		Age:  nil,
//		Uid:  nil,
//	})
//	as.Nil(err)
//	as.Equal(true, has)
//
//	mock.ExpectQuery("t_false *").
//		WithArgs("kk").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}))
//
//	has, err = engine.Has("t_false").ByModel(Whe{
//		Name: types.NewString("kk"),
//		Age:  nil,
//		Uid:  nil,
//	})
//	as.Nil(err)
//	as.Equal(false, has)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
//
//func TestHasByWhere(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	mock.ExpectQuery("t_user *").
//		WithArgs("kk").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}).
//			AddRow(1),
//		)
//	num, err := engine.Has(User{}).ByWhere(new(WhereBuilder).
//		Eq("name", "kk"),
//	)
//	as.Nil(err)
//	as.Equal(true, num)
//
//	mock.ExpectQuery("t_false *").
//		WithArgs("kk").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{""}))
//
//	num, err = engine.Has("t_false").ByWhere(new(WhereBuilder).
//		Eq("name", "kk"),
//	)
//	as.Nil(err)
//	as.Equal(false, num)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
