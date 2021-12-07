package utils

import (
	"bytes"
	"strings"
)

func GenwhereToken(columns []string) []string {
	tokens := make([]string, 0)
	for _, column := range columns {
		tokens = append(tokens, column+" = ? ")
	}
	return tokens
}

func GenwhereTokenOfBatch(columns []string, num int) string {
	where := _genWhere(columns)
	tokens := make([]string, 0)
	for i := 0; i < num; i++ {
		tokens = append(tokens, where)
	}
	sql := strings.Join(tokens, " or ")
	if num == 1 {
		return sql
	}

	var bb bytes.Buffer
	bb.WriteString("(")
	bb.WriteString(sql)
	bb.WriteString(")")
	return bb.String()
}

// columns:["id","name"] return: "id = ? and name = ?"
func _genWhere(columns []string) string {
	tokens := GenwhereToken(columns)
	sql := strings.Join(tokens, "and ")
	if len(columns) == 1 {
		return sql
	}

	var bb bytes.Buffer
	bb.WriteString("(")
	bb.WriteString(sql)
	bb.WriteString(")")
	return bb.String()
}
