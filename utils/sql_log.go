package utils

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ANSI颜色常量，用于终端输出高亮
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// PrintSql 打印SQL语句及其参数，第一个参数控制是否启用颜色输出
func PrintSql(disableColor bool, originalSql, dialectSql string, args ...any) {
	// 打印分隔线
	fmt.Println(strings.Repeat("-", 80))

	// 打印SQL语句
	if disableColor {
		fmt.Println("完整 SQL 语句:")
		fmt.Println(combineSqlArgs(originalSql, args...))
		fmt.Println("预处理 SQL 语句:")
		fmt.Println(dialectSql)
	} else {
		fmt.Println(colorCyan + "完整 SQL 语句:" + colorReset)
		fmt.Println(HighlightSQL(combineSqlArgs(originalSql, args...)))
		fmt.Println(colorCyan + "预处理 SQL 语句:" + colorReset)
		fmt.Println(HighlightSQL(dialectSql))
	}

	// 打印参数
	if len(args) == 0 {
		if disableColor {
			fmt.Println("参数: 无")
		} else {
			fmt.Println(colorYellow + "参数: 无" + colorReset)
		}
	} else {
		if disableColor {
			fmt.Println("参数列表:")
		} else {
			fmt.Println(colorGreen + "参数列表:" + colorReset)
		}

		// 格式化打印每个参数，带索引
		for i, arg := range args {
			argStr := formatArg(arg)
			fmt.Printf("  参数 #%d: %s\n", i, argStr)
		}
	}

	// 打印结束分隔线
	fmt.Println(strings.Repeat("-", 80) + "\n")
}

// 格式化参数为 SQL 字符串
func formatArg(arg any) string {
	// 如果实现了 driver.Valuer
	if valuer, ok := arg.(driver.Valuer); ok {
		v, err := valuer.Value()
		if err != nil {
			return fmt.Sprintf("<invalid valuer: %v>", err)
		}
		return formatBasic(v)
	}
	return formatBasic(arg)
}

// 格式化基本类型
func formatBasic(arg any) string {
	switch v := arg.(type) {
	case nil:
		return "NULL"
	case string:
		return fmt.Sprintf("'%s'", v)
	case []byte:
		// 默认当作 BLOB 打印
		hexStr := ""
		for i := 0; i < len(v) && i < 8; i++ {
			hexStr += fmt.Sprintf("%02X ", v[i])
		}
		if len(v) > 8 {
			hexStr += "... "
		}
		return fmt.Sprintf("<BLOB:%s(%d bytes)>", hexStr, len(v))
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("'%v'", arg)
	}
}

// 用参数替换 SQL 的 ?
func combineSqlArgs(sql string, args ...any) string {
	var sb strings.Builder
	argIdx := 0
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' && argIdx < len(args) {
			sb.WriteString(formatArg(args[argIdx]))
			argIdx++
		} else {
			sb.WriteByte(sql[i])
		}
	}
	return sb.String()
}

func HighlightSQL(sql string) string {
	// 匹配常见关键字并加粗蓝色
	keywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE",
		"FROM", "WHERE", "VALUES", "AND", "OR", "SET", "ORDER BY",
		"GROUP BY", "HAVING", "JOIN", "ON", "AS", "NOT", "NULL",
		"TRUE", "FALSE", "LIKE", "IN", "BETWEEN", "IS", "UNION",
		"CREATE", "DROP", "TABLE", "VIEW", "INDEX", "PRIMARY KEY",
		"FOREIGN KEY", "UNIQUE", "CHECK", "DEFAULT", "AUTO_INCREMENT",
		"COMMENT", "TRIGGER", "PROCEDURE", "FUNCTION", "BEGIN", "END",
		"CALL", "RETURN", "CASE", "WHEN", "THEN", "ELSE",
		"IF", "ELSEIF", "DECLARE",
		"DESC", "ASC", "LIMIT", "OFFSET", "DISTINCT", "ALL", "EXISTS",
		"INNER JOIN", "LEFT JOIN", "RIGHT JOIN", "FULL JOIN",
		"CROSS JOIN", "NATURAL JOIN", "UNION ALL", "INTERSECT", "EXCEPT",
		"EXPLAIN", "ANALYZE", "OPTIMIZE",
	}

	// 预编译所有正则表达式
	regexps := make([]*regexp.Regexp, len(keywords))
	for i, kw := range keywords {
		regexps[i] = regexp.MustCompile(`(\b` + kw + `\b)`)
	}

	// 依次应用所有正则表达式
	for _, re := range regexps {
		sql = re.ReplaceAllString(sql, "\033[1;34m${1}\033[0m")
	}

	return sql
}
