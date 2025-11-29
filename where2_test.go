package lorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestWhereBuilderNot1(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Eq("a", 1).Not()

	query, args, err := w1.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a = ?)", query)
	as.Equal([]any{1}, args)
}

func TestWhereBuilderNot2(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Eq("a", 1).Eq("b", 2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a = ? AND b = ?)", query)
	as.Equal([]any{1, 2}, args)
}

func TestWhereBuilderNot3(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Eq("a", 1).Not()
	w2 := W().Eq("b", 2).Not()

	w0 := w1.And(w2)

	query, args, err := w0.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a = ? AND NOT (b = ?))", query)
	as.Equal([]any{1, 2}, args)
}

func TestWhereBuilderNot4(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Eq("a", 1).Not()
	w2 := W().Eq("b", 2).Not()

	w0 := w1.Or(w2)

	query, args, err := w0.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a = ? OR NOT (b = ?))", query)
	as.Equal([]any{1, 2}, args)
}

func TestWhereBuilderNot5(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w01 := W().Eq("a1", "a1").Or(W().Eq("a2", "a2")).Not()
	w02 := W().Eq("b1", "b1").Or(W().Eq("b2", "b2")).Not()

	w0 := W().And(w01).And(w02)

	query, args, err := w0.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("(NOT (a1 = ? OR a2 = ?)) AND (NOT (b1 = ? OR b2 = ?))", query)
	as.Equal([]any{"a1", "a2", "b1", "b2"}, args)
}

func TestWhereBuilderNot6(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w01 := W().Eq("a1", "a1").Or(W().Eq("a2", "a2")).Not()

	w0 := W().Or(w01).Or(w01)

	w00 := W().Or(w0)

	query, args, err := w00.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a1 = ? OR a2 = ?) OR NOT (a1 = ? OR a2 = ?)", query)
	as.Equal([]any{"a1", "a2", "a1", "a2"}, args)
}

func TestWhereBuilderNot7(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Eq("a1", "a1").Not()
	w2 := W().Eq("a2", "a2").Not()
	w3 := W().Eq("a3", "a3").Not()

	w0 := W().And(w1).And(w2).And(w3)

	query, args, err := w0.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a1 = ?) AND NOT (a2 = ?) AND NOT (a3 = ?)", query)
	as.Equal([]any{"a1", "a2", "a3"}, args)
}

func TestWhereBuilderNot8(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Or(W()).Not()
	w2 := W().And(W()).Not()
	w3 := W().Or(W().Or(W())).Not()

	w0 := W().And(w1).And(w2).And(w3)

	query, args, err := w0.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderNot9(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(1, 2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse)
	as.ErrorIs(err, ErrNoPk)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderNot10(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(1, 2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.Nil(err)
	as.Equal("NOT (id IN (?,?))", query)
	as.Equal([]any{1, 2}, args)
}

func TestWhereBuilderNot11(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(1).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.Nil(err)
	as.Equal("NOT (id IN (?))", query)
	as.Equal([]any{1}, args)
}

func TestWhereBuilderNot12(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(1, struct {
	}{}).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id")
	as.ErrorIs(err, ErrTypePkArgs)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderNot13(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(1, 2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id", "name")
	as.ErrorIs(err, ErrNeedMultiPk)
	as.Equal("", query)
	as.Nil(args)
}

func TestWhereBuilderNot14(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(struct {
		Id   int
		Name string
	}{
		Id:   1,
		Name: "name",
	}).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "id", "name")
	as.Nil(err)
	as.Equal("NOT (id = ? AND name = ?)", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderNot15(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().PrimaryKey(struct {
		DocId   int
		DocName string
	}{
		DocId:   1,
		DocName: "name",
	}).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (doc_id = ? AND doc_name = ?)", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderNot16(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	var m = make(map[string]any)
	m["doc_id"] = 1
	m["doc_name"] = "name"

	w1 := W().PrimaryKey(m).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT (doc_id = ? AND doc_name = ?)", query)
	as.Equal([]any{1, "name"}, args)
}

func TestWhereBuilderNot17(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	var m = make(map[string]any)
	m["doc_id"] = 1
	m["doc_name"] = "name"

	var m2 = make(map[string]any)
	m2["doc_id"] = 2
	m2["doc_name"] = "name2"

	w1 := W().PrimaryKey(m, m2).Not()

	query, args, err := w1.toSql(engine.getDialect().parse, "doc_id", "doc_name")
	as.Nil(err)
	as.Equal("NOT ((doc_id = ? AND doc_name = ?) OR (doc_id = ? AND doc_name = ?))", query)
	as.Equal([]any{1, "name", 2, "name2"}, args)
}

func TestWhereBuilderNot18(t *testing.T) {
	as := assert.New(t)
	engine := getMockDB(PgConf{})

	w1 := W().Or(W().Eq("a", "a")).Not()

	w1.And(W().Or(W()))
	w1.And(W().Or(W()))

	query, args, err := w1.toSql(engine.getDialect().parse)
	as.Nil(err)
	as.Equal("NOT (a = ?)", query)
	as.Equal([]any{"a"}, args)
}
