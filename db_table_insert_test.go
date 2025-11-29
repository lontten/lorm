package lorm

type UserP struct {
	Id   *int
	Name *string
}

//func _TestCreateOrUpdate(t *testing.T) {
//	as := assert.New(t)
//	lorm, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(lorm, &PgConf{})
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
//	lorm, mock, err := sqlmock.New()
//	as.Nil(err, "new sqlmock error")
//	engine := MustConnectMock(lorm, &PgConf{})
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
