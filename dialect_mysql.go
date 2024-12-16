package lorm

import (
	"errors"
	"github.com/lontten/lorm/insert-type"
	"github.com/lontten/lorm/softdelete"
	"github.com/lontten/lorm/utils"
	"strings"
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
func (d *MysqlDialect) initContext() Dialecter {
	return &MysqlDialect{ctx: &ormContext{
		ormConf:                 d.ctx.ormConf,
		query:                   &strings.Builder{},
		wb:                      Wb(),
		insertType:              insert_type.Err,
		dialectNeedLastInsertId: d.ctx.dialectNeedLastInsertId,
	}}
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
	// mysql insert 时，无法直接返回数据，只能借助 last_insert_id
	ctx.sqlIsQuery = false

	extra := ctx.extra
	set := extra.set

	columns := ctx.columns
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

	query.WriteString("(")
	query.WriteString(strings.Join(columns, ","))
	query.WriteString(") ")
	query.WriteString("VALUES")
	query.WriteString("(")
	ctx.genInsertValuesSqlBycolumnValues()
	query.WriteString(" ) ")

	switch ctx.insertType {
	case insert_type.Update:
		query.WriteString(" AS new ON DUPLICATE KEY UPDATE ")
		// 当未设置更新字段时，默认为所有字段
		if len(set.columns) == 0 && len(set.fieldNames) == 0 {
			list := append(ctx.columns, extra.columns...)

			for _, name := range list {
				find := utils.Find(extra.duplicateKeyNames, name)
				if find < 0 {
					set.fieldNames = append(set.fieldNames, name)
				}
			}
		}

		// 当 软删除 字段 未删除状态 为 0 时，这里fieldNames 会有 软删除字段，
		// DUPLICATE KEY UPDATE 时，软删除字段是否应该更新，问题分析：当更新时：
		// 1.假设 唯一索引 字段 为 name (因为有软删除，逻辑上这样加唯一索引是错误的。) ，更新值为 abc，则数据库中，一定是只有一条 name为 abc的字段，
		// 更新时，会把 软删除设为未删除状态。
		// 1.1 旧数据 未删除，更新数据，符合预期
		// 1.2 旧数据 已删除，更新数据并变成未删除状态，相当于替换，数据insert成功，只是 id 是 原来的数据，之前的旧数据没有了，勉强算是 符合预期
		// 2. 将 name 和 软删除字段（唯一索引设置正确的情况），设为 符合唯一索引，则 数据库中已有数据为 name=abc,del=0; name=abc,del=大于0的数 （已软删除数据众多）
		// 2.1 只有已删除数据，直接插入，符合预期
		// 2.2 同时有未删除，已删除数据，只对未删除数据 进行更新，成功更新，符合预期

		// DUPLICATE KEY UPDATE 时，软删除字段是否应该更新，问题分析：当不更新时：
		// 1.假设 唯一索引 字段 为 name (因为有软删除，逻辑上这样加唯一索引是错误的。) ，更新值为 abc，则数据库中，一定是只有一条 name为 abc的字段，更新时，不会修改软删除状态。
		// 1.1 旧数据 未删除，更新数据，符合预期
		// 1.2 旧数据 已删除，更新数据，已删除被更新，数据还是被删除状态，无法查询到添加的数据，不符合预期！
		// 2. 将 name 和 软删除字段（唯一索引设置正确的情况），设为 符合唯一索引，则 数据库中已有数据为 name=abc,del=0; name=abc,del=大于0的数 （已软删除数据众多）
		// 2.1 只有已删除数据，直接插入，符合预期
		// 2.2 同时有未删除，已删除数据，同时对未删除数据和已删除 进行更新，勉强算是 符合预期

		// 从上面分析可知，DUPLICATE KEY UPDATE 时，软删除字段不进行更新是最差方案，会出现 不符合预期情况。
		// 软删除字段进行更新 时，如果 唯一索引设置正确，是完美执行；如果 唯一索引 错误，也可以达到 基本复合预期的效果。

		for i, name := range set.fieldNames {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString(name + " = new." + name)
		}
		for i, column := range set.columns {
			query.WriteString(column + " = ? , ")
			ctx.args = append(ctx.args, set.columnValues[i].Value)
		}
		break
	default:
		break
	}

	query.WriteString(";")
}

// del 生成
func (d *MysqlDialect) tableDelGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}

	var query = d.ctx.query

	tableName := ctx.tableName

	//  没有软删除 或者 跳过软删除 ，执行物理删除
	if ctx.softDeleteType == softdelete.None || ctx.skipSoftDelete {
		query.WriteString("DELETE FROM ")
		query.WriteString(tableName)

		query.WriteString(" WHERE ")
		query.WriteString(ctx.extraWhereSql)
	} else {
		query.WriteString("UPDATE ")
		query.WriteString(tableName)

		query.WriteString(" SET ")
		ctx.genSetSqlBycolumnValues()

		query.WriteString(" WHERE ")
		query.WriteString(ctx.extraWhereSql)
	}

	query.WriteString(";")
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
