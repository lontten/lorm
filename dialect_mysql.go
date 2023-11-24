package lorm

import (
	"errors"
	"github.com/lontten/lorm/utils"
	"strings"
)

type MysqlDialect struct {
	ctx *ormContext
}

func (m MysqlDialect) getCtx() *ormContext {
	return m.ctx
}

func (m MysqlDialect) appendBaseToken(token baseToken) {
	m.ctx.baseTokens = append(m.ctx.baseTokens, token)
}
func (m MysqlDialect) query(query string, args ...interface{}) (string, []interface{}) {
	m.ctx.log.Println(query, args)
	return query, args
}

func (m MysqlDialect) queryBatch(query string) string {
	return query
}
func (m MysqlDialect) insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (string, []interface{}) {
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
		") VALUES (" + strings.Repeat("?", len(args)) +
		") ON duplicate key UPDATE " + strings.Join(cs, "=?, ") + "=?"

	args = append(args, vs...)
	m.ctx.log.Println(query, args)

	//exec, err := m.db.Exec(query, args...)
	//if err != nil {
	//	return 0, err
	//}
	return query, args
}

func (m MysqlDialect) insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (string, []interface{}) {
	//return 0, errors.New("MySQL insertOrUpdateByUnique not support, please use insertOrUpdateByPrimaryKey")
	return "", nil
}

func (m MysqlDialect) exec(query string, args ...interface{}) (string, []interface{}) {
	m.ctx.log.Println(query, args)

	//exec, err := m.db.Exec(query, args...)
	//if err != nil {
	//	return 0, err
	//}
	return query, args
}

func (m MysqlDialect) execBatch(query string, args [][]interface{}) (string, [][]interface{}) {
	m.ctx.log.Println(query, args)

	//var num int64 = 0
	//stmt, err := m.db.Prepare(query)
	//if err != nil {
	//	return 0, err
	//}
	//for _, arg := range args {
	//	exec, err := stmt.Exec(arg...)
	//	m.log.Println(query, args)
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

func (m MysqlDialect) parse(c Clause) (string, error) {
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

//todo 下面未重构--------------
