package lorm

//
//import (
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/lontten/lorm/types"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestCreate(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	mock.ExpectQuery("user_p *").
//		WithArgs(2, "up").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(88))
//
//	user := UserP{
//		Name: types.NewString("up"),
//		Id:   types.NewInt(2),
//	}
//	num, err := engine.Insert(&user)
//	as.Nil(err)
//	as.Equal(int64(1), num)
//	as.Equal(88, *user.Id)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
//
//func _TestCreateOrUpdate(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	//-----------------create--------------------------
//	mock.ExpectQuery("INSERT *").
//		WithArgs(2, "add", "add").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).
//			AddRow(1),
//		)
//
//	user := UserP{
//		Name: types.NewString("add"),
//		Id:   types.NewInt(2),
//	}
//	num, err := engine.InsertOrUpdate(&user).ByPrimaryKey()
//	as.Nil(err)
//	as.Equal(int64(1), num)
//	as.Equal(1, *user.Id)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
//
//func _TestCreateOrUpdate2(t *testing.T) {
//	as := assert.New(t)
//	ldb, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(ldb, &PgConf{})
//
//	//-----------------update--------------------------
//	mock.ExpectQuery("INSERT*").
//		WithArgs(1, "upd", "upd").
//		WillReturnError(nil).
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).
//			AddRow(2),
//		)
//
//	user := UserP{
//		Name: types.NewString("upd"),
//		Id:   types.NewInt(2),
//	}
//	num, err := engine.InsertOrUpdate(&user).ByPrimaryKey()
//	as.Nil(err)
//	as.Equal(int64(2), num)
//	as.Equal(2, *user.Id)
//
//	as.Nil(mock.ExpectationsWereMet(), "we make sure that all expectations were met")
//}
