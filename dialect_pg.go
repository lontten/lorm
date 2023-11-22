package lorm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lontten/lorm/utils"
	"strconv"
	"strings"
)

type PgDialect struct {
	ctx ormContext
}

func (m PgDialect) query(query string, args ...interface{}) (string, []interface{}) {
	query = toPgSql(query)
	m.ctx.log.Println(query, args)
	//return m.db.Query(query, args...)
	return query, args
}

func (m PgDialect) queryBatch(query string) string {
	query = toPgSql(query)

	//return m.db.Prepare(query)
	return query
}

func (m PgDialect) insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (string, []interface{}) {
	return m.insertOrUpdateByUnique(table, fields, columns, args...)
}

func (m PgDialect) insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (string, []interface{}) {
	cs := make([]string, 0)
	vs := make([]interface{}, 0)

	for i, column := range columns {
		if utils.Contains(fields, column) {
			continue
		}
		cs = append(cs, column)
		vs = append(vs, args[i])
	}

	var query = "INSERT INTO " + table + "(" + strings.Join(columns, ",") +
		") VALUES (" + strings.Repeat(" ? ,", len(args)-1) +
		" ? ) ON CONFLICT (" + strings.Join(fields, ",") + ") DO"
	if len(vs) == 0 {
		query += "NOTHING"
	} else {
		query += " UPDATE SET " + strings.Join(cs, "= ? , ") + "= ? "
	}
	args = append(args, vs...)
	query = toPgSql(query)
	m.ctx.log.Println(query, args)
	//exec, err := m.db.Exec(query, args...)
	//if err != nil {
	//	if errors.As(err, &ErrNoPkOrUnique) {
	//		return 0, errors.New("insertOrUpdateByUnique fields need to be unique or primary key:" + strings.Join(fields, ",") + err.Error())
	//	}
	//	return 0, err
	//}
	return query, args
}

func (m PgDialect) exec(query string, args ...interface{}) (string, []interface{}) {
	query = toPgSql(query)
	m.ctx.log.Println(query, args)

	//exec, err := m.db.Exec(query, args...)
	//if err != nil {
	//	return 0, err
	//}
	return query, args
}

func (m PgDialect) execBatch(query string, args [][]interface{}) (string, [][]interface{}) {
	query = toPgSql(query)
	//var num int64 = 0
	//stmt, err := m.db.Prepare(query)
	//defer stmt.Close()
	//if err != nil {
	//	return 0, err
	//}
	//for _, arg := range args {
	//	exec, err := stmt.Exec(arg...)
	//
	//	m.log.Println(query, arg)
	//	if err != nil {
	//		return num, err
	//	}
	//	rowsAffected, err := exec.RowsAffected()
	//	if err != nil {
	//		return num, err
	//	}
	//	num += rowsAffected
	//}
	return query, args
}

func toPgSql(sql string) string {
	var i = 1
	for {
		t := strings.Replace(sql, "?", " $"+strconv.Itoa(i)+" ", 1)
		if t == sql {
			break
		}
		i++
		sql = t
	}
	return sql
}

func (m PgDialect) parse(c Clause) (string, error) {
	sb := strings.Builder{}
	switch c.Type {
	case Eq:
		sb.WriteString(c.query + " = ?")
	case Neq:
		sb.WriteString(c.query + " <> ?")
	case Less:
		sb.WriteString(c.query + " < ?")
	case LessEq:
		sb.WriteString(c.query + " <= ?")
	case Greater:
		sb.WriteString(c.query + " > ?")
	case GreaterEq:
		sb.WriteString(c.query + " >= ?")
	case Like:
		sb.WriteString(c.query + " LIKE ?")
	case NotLike:
		sb.WriteString(c.query + " NOT LIKE ?")
	case In:
		sb.WriteString(c.query + " IN (")
		sb.WriteString(gen(c.argsNum))
		sb.WriteString(")")
	case NotIn:
		sb.WriteString(c.query + " NOT IN (")
		sb.WriteString(gen(c.argsNum))
		sb.WriteString(")")
	case Between:
		sb.WriteString(c.query + " BETWEEN ? AND ?")
	case NotBetween:
		sb.WriteString(c.query + " NOT BETWEEN ? AND ?")
	case IsNull:
		sb.WriteString(c.query + " IS NULL")
	case IsNotNull:
		sb.WriteString(c.query + " IS NOT NULL")
	case IsFalse:
		sb.WriteString(c.query + " IS FALSE")
	default:
		return "", errors.New("unknown where token type")
	}

	return sb.String(), nil
}

func (m PgDialect) prepare(query string) (Stmt, error) {
	query = toPgSql(query)
	stmt, err := m.db.Prepare(query)
	return Stmt{stmt: stmt}, err
}

//todo 下面未重构--------------

func (m PgDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) TXer {
	tx := m.core.beginTx(ctx, opts)
	return lnDB{
		core: tx,
		ctx:  db.ctx,
	}
}

func (m PgDialect) Rollback() error {
	err := m.core.rollback()
	if err != nil {
		return err
	}
	db.ctx.log.Println("rollback")
	return nil
}

func (m PgDialect) Commit() error {
	err := m.core.commit()
	if err != nil {
		return err
	}
	m.ctx.log.Println("commit")
	return nil
}
func (m PgDialect) C() {
}
func (m PgDialect) R() {
}

func (m PgDialect) U() {
}
func (m PgDialect) D() {
}
func (m PgDialect) Query(query string, args ...interface{}) *NativeQuery {
	return m.core.query(query, args...)
}
func (m PgDialect) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	//query, args = db.dialect.exec(query, args...)
	//return tx.doExec(query, args...)
	return 0, nil
}
