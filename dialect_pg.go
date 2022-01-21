package lsql

import (
	"database/sql"
	"errors"
	"github.com/lontten/lsql/utils"
	"strconv"
	"strings"
)

type PgDialect struct {
	db  DBer
	log Logger
}

func (m PgDialect) Copy(db DBer) Dialect {
	return &PgDialect{
		db:  db,
		log: m.log,
	}
}

func (m PgDialect) query(query string, args ...interface{}) (*sql.Rows, error) {
	query = toPgSql(query)
	m.log.Println(query, args)
	return m.db.Query(query, args...)
}

func (m PgDialect) queryBatch(query string) (*sql.Stmt, error) {
	query = toPgSql(query)

	return m.db.Prepare(query)
}

func (m PgDialect) insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (int64, error) {
	return m.insertOrUpdateByUnique(table, fields, columns, args...)
}

func (m PgDialect) insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (int64, error) {
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
	m.log.Println(query, args)
	exec, err := m.db.Exec(query, args...)
	if err != nil {
		if errors.As(err, &ErrNoPkOrUnique) {
			return 0, errors.New("insertOrUpdateByUnique fields need to be unique or primary key:" + strings.Join(fields, ",") + err.Error())
		}
		return 0, err
	}
	return exec.RowsAffected()
}

func (m PgDialect) exec(query string, args ...interface{}) (int64, error) {
	query = toPgSql(query)
	m.log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m PgDialect) execBatch(query string, args [][]interface{}) (int64, error) {
	query = toPgSql(query)
	var num int64 = 0
	stmt, err := m.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	for _, arg := range args {
		exec, err := stmt.Exec(arg...)

		m.log.Println(query, arg)
		if err != nil {
			return num, err
		}
		rowsAffected, err := exec.RowsAffected()
		if err != nil {
			return num, err
		}
		num += rowsAffected
	}
	return num, nil
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
