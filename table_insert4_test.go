package lorm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// 插入时，唯一索引冲突，处理策略-忽略
func TestInsert4_mysql(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &MysqlConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT IGNORE t_user (id, name) VALUES (?, ?);")).
		WithArgs(1, "tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = User{
		Id:   1,
		Name: "tom",
	}
	num, err := Insert(engine, u, E().ShowSql().WhenDuplicateKey().DoNothing())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}

func TestInsert4_pg(t *testing.T) {
	as := assert.New(t)
	db, mock, err := sqlmock.New()
	as.Nil(err, fmt.Sprintf("failed to open sqlmock database: %s", err))
	engine := MustConnectMock(db, &PgConf{})

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO t_user (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING;")).
		WithArgs(1, "tom").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	var u = User{
		Id:   1,
		Name: "tom",
	}
	num, err := Insert(engine, u, E().ShowSql().WhenDuplicateKey().DoNothing())
	as.Nil(err)
	as.Equal(int64(1), num, "num error")
}
