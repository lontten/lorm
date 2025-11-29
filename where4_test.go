package lorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestWhereBuilderFilterNot1(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, 2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse)
	as.ErrorIs(err, ErrNoPk)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderFilterNot2(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, 2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.Nil(err)
	as.Equal("NOT (NOT (id IN (?,?)))", query)
	as.Equal([]any{1, 2}, args)
}

func TestWhereBuilderFilterNot3(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.Nil(err)
	as.Equal("NOT (NOT (id IN (?)))", query)
	as.Equal([]any{1}, args)
}

func TestWhereBuilderFilterNot4(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, struct {
	}{}).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.ErrorIs(err, ErrTypePkArgs)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderFilterNot5(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(1, 2).Not()

	query, args, err := engine.ToWhereSQL(w1, "id", "name")
	as.ErrorIs(err, ErrNeedMultiPk)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderFilterNot6(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(struct {
		Id   int
		Name string
	}{
		Id:   1,
		Name: "name",
	}).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id", "name")
	as.Nil(err)
	as.Equal("NOT (NOT (id = ? AND name = ?))", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderFilterNot7(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().FilterPrimaryKey(struct {
		DocId   int
		DocName string
	}{
		DocId:   1,
		DocName: "name",
	}).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (NOT (doc_id = ? AND doc_name = ?))", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderFilterNot8(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	var m = make(map[string]any)
	m["doc_id"] = 1
	m["doc_name"] = "name"

	w1 := W().FilterPrimaryKey(m).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (NOT (doc_id = ? AND doc_name = ?))", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderFilterNot9(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	var m = make(map[string]any)
	m["doc_id"] = 1
	m["doc_name"] = "name"

	var m2 = make(map[string]any)
	m2["doc_id"] = 2
	m2["doc_name"] = "name2"

	w1 := W().FilterPrimaryKey(m, m2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (NOT ((doc_id = ? AND doc_name = ?) OR (doc_id = ? AND doc_name = ?)))", query)
	as.Equal([]any{1, "name", 2, "name2"}, args)
}
