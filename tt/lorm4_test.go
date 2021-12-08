package lorm

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lontten/lorm"
	"github.com/lontten/lorm/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func selectList(engine lorm.Engine) {
	list := types.StringList{}
	kk := make([]string, 0)

	num, err := engine.Query("select  ARRAY['os''dba', '123''456']  ").GetOne(&kk)
	fmt.Println(num)
	fmt.Println(err)
	fmt.Println(list)
	fmt.Println(kk)
}

func TestDB434(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err)

	ormConfig := lorm.OrmConf{
		TableNamePrefix: "t_",
		PrimaryKeyNames: []string{"id"},
	}
	engine := lorm.MustConnectMock(db, &lorm.PgConf{}).Db(&ormConfig)

	mock.ExpectQuery("select 2").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(2),
		)

	selectList(engine)

}
