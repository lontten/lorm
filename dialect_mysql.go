package lorm

import (
	"errors"
	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/insert_type"
	"github.com/lontten/lorm/return_type"
	"github.com/lontten/lorm/utils"
	"strconv"
	"strings"
	"time"
)

type MysqlDialect struct {
	ctx *ormContext
}

// ===----------------------------------------------------------------------===//
// 获取上下文
// ===----------------------------------------------------------------------===//

func (d *MysqlDialect) getCtx() *ormContext {
	return d.ctx
}
func (d *MysqlDialect) initContext() *ormContext {
	d.ctx = &ormContext{
		ormConf:    d.ctx.ormConf,
		query:      &strings.Builder{},
		insertType: insert_type.Err,
	}
	return d.ctx
}
func (d *MysqlDialect) hasErr() bool {
	return d.ctx.err != nil
}

func (d *MysqlDialect) getErr() error {
	return d.ctx.err
}

// ===----------------------------------------------------------------------===//
// sql 方言化
// ===----------------------------------------------------------------------===//
func (d *MysqlDialect) query(query string, args ...any) (string, []any) {
	d.ctx.log.Println(query, args)
	return query, args
}

func (d *MysqlDialect) queryBatch(query string) string {
	return query
}

func (d *MysqlDialect) prepare(query string) string {
	return query
}

func (d *MysqlDialect) exec(query string, args ...any) (string, []any) {
	return query, args
}

func (d *MysqlDialect) insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...any) (string, []any) {
	cs := make([]string, 0)
	vs := make([]any, 0)

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
	d.ctx.log.Println(query, args)
	return query, args
}

func (d *MysqlDialect) insertOrUpdateByUnique(table string, fields []string, columns []string, args ...any) (string, []any) {
	d.ctx.err = errors.New("MySQL insertOrUpdateByUnique not support, please use insertOrUpdateByPrimaryKey")
	return "", nil
}

func (d *MysqlDialect) execBatch(query string, args [][]any) (string, [][]any) {
	d.ctx.log.Println(query, args)

	//var num int64 = 0
	//stmt, err := d.ldb.Prepare(query)
	//if err != nil {
	//	return 0, err
	//}
	//for _, arg := range args {
	//	exec, err := stmt.Exec(arg...)
	//	d.log.Println(query, args)
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

// ===----------------------------------------------------------------------===//
// 工具
// ===----------------------------------------------------------------------===//
func (d *MysqlDialect) appendBaseToken(token baseToken) {
	d.ctx.baseTokens = append(d.ctx.baseTokens, token)
}

// ===----------------------------------------------------------------------===//
// 中间服务
// ===----------------------------------------------------------------------===//
// 初始化主键
func (d *MysqlDialect) initPrimaryKeyName() {
	if d.ctx.err != nil {
		return
	}
	v := d.ctx.destV
	dest := d.ctx.scanDest
	d.ctx.primaryKeyNames = d.ctx.ormConf.primaryKeys(v, dest)
}
func (d *MysqlDialect) getSql() string {
	s := d.ctx.query.String()
	return s
}

// insert 生成
func (d *MysqlDialect) tableInsertGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}
	columns := ctx.columns
	values := ctx.columnValues
	var query = d.ctx.query

	switch ctx.insertType {
	case insert_type.Err:
		query.WriteString("INSERT INTO ")
		break
	case insert_type.Ignore:
		query.WriteString("INSERT IGNORE ")
		break
	case insert_type.Update:
		query.WriteString("INSERT INTO ")
		break
	case insert_type.Replace:
		query.WriteString("REPLACE INTO ")
		break
	}
	query.WriteString(ctx.tableName + " ")

	query.WriteString(" ( ")
	for i, v := range columns {
		if i == 0 {
			query.WriteString(v)
		} else {
			query.WriteString(" , " + v)
		}
	}
	query.WriteString(" ) ")
	query.WriteString(" VALUES ")
	query.WriteString("( ")
	for i, v := range values {
		if i > 0 {
			query.WriteString(" , ")
		}
		switch v.Type {
		case field.None:
			break
		case field.Null:
			query.WriteString("NULL")
			break
		case field.Now:
			query.WriteString("NOW()")
			break
		case field.UnixSecond:
			query.WriteString(strconv.Itoa(time.Now().Second()))
			break
		case field.UnixMilli:
			query.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))
			break
		case field.UnixNano:
			query.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
			break
		case field.Val:
			query.WriteString(" ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Increment:
			query.WriteString(columns[i] + " + ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Expression:
			query.WriteString(v.Value.(string))
			break
		case field.ID:
			if len(ctx.primaryKeyNames) > 0 {
				ctx.err = errors.New("软删除标记为主键id，需要单主键")
				return
			}
			query.WriteString(ctx.primaryKeyNames[0])
			break
		}
	}
	query.WriteString(" ) ")

	switch ctx.insertType {
	case insert_type.Update:
		set := ctx.extra.set
		query.WriteString(" AS new ON DUPLICATE KEY UPDATE ")
		// 当未设置更新字段时，默认为所有字段
		if len(set.columns) == 0 && len(set.fieldNames) == 0 {
			set.fieldNames = ctx.modelAllFieldNames
		}
		for _, name := range set.fieldNames {
			query.WriteString(name + " = new." + name + ", ")
		}
		for i, column := range set.columns {
			query.WriteString(column + " = ? , ")
			ctx.args = append(ctx.args, set.columnValues[i].Value)
		}

		break
	default:
		break
	}

	// INSERT IGNORE 无法和 RETURNING 共存，当 INSERT IGNORE 时，不返回
	if ctx.scanIsPtr && ctx.insertType != insert_type.Ignore {
		switch expr := ctx.returnType; expr {
		case return_type.None:
			ctx.sqlIsQuery = true
			break
		case return_type.PrimaryKey:
			query.WriteString(" RETURNING " + strings.Join(ctx.primaryKeyNames, ","))
		case return_type.ZeroField:
			query.WriteString(" RETURNING " + strings.Join(ctx.modelZeroFieldNames, ","))
		case return_type.AllField:
			query.WriteString(" RETURNING " + strings.Join(ctx.modelAllFieldNames, ","))
		}
	}
	query.WriteString(";")
}

// 获取struct对应的字段名 有效部分
func (d *MysqlDialect) initColumns() {
	if d.ctx.err != nil {
		return
	}

	columns, err := d.ctx.ormConf.initColumns(d.ctx.destBaseType)
	if err != nil {
		d.ctx.err = err
		return
	}
	d.ctx.columns = columns
}
func (d *MysqlDialect) parse(c Clause) (string, error) {
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
