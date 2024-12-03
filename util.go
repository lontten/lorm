package lorm

import "strings"

func gen(num int) string {
	var queryArr []string
	for i := 0; i < num; i++ {
		queryArr = append(queryArr, "?")
	}
	return strings.Join(queryArr, ",")
}

// todo 下面未重构--------------
