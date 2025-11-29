package lorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestWhereBuilderFilter1(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, 2)

	query, args, err := w1.toSql(engine.getDialect().parse)
	as.ErrorIs(err, ErrNoPk)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderFilter2(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, 2)

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.Nil(err)
	as.Equal("NOT (id IN (?,?))", query)
	as.Equal([]any{1, 2}, args)
}

func TestWhereBuilderFilter3(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1)

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.Nil(err)
	as.Equal("NOT (id IN (?))", query)
	as.Equal([]any{1}, args)
}

func TestWhereBuilderFilter4(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, struct {
	}{})

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.ErrorIs(err, ErrTypePkArgs)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderFilter5(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, 2)

	query, args, err := w1.toSql(engine.getDialect().parse, "id", "name")
	as.ErrorIs(err, ErrNeedMultiPk)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderFilter6(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(struct {
		Id   int
		Name string
	}{
		Id:   1,
		Name: "name",
	})

	query, args, err := w1.toSql(engine.getDialect().parse, "id", "name")
	as.Nil(err)
	as.Equal("NOT (id = ? AND name = ?)", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderFilter7(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(struct {
		DocId   int
		DocName string
	}{
		DocId:   1,
		DocName: "name",
	})

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (doc_id = ? AND doc_name = ?)", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderFilter8(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	var m = make(map[string]any)
	m["doc_id"] = 1
	m["doc_name"] = "name"

	w1 := W().FilterPrimaryKey(m)

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (doc_id = ? AND doc_name = ?)", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderFilter9(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	var m = make(map[string]any)
	m["doc_id"] = 1
	m["doc_name"] = "name"

	var m2 = make(map[string]any)
	m2["doc_id"] = 2
	m2["doc_name"] = "name2"

	w1 := W().FilterPrimaryKey(m, m2)

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT ((doc_id = ? AND doc_name = ?) OR (doc_id = ? AND doc_name = ?))", query)
	as.Equal([]any{1, "name", 2, "name2"}, args)
}
