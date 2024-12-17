package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lontten/lorm/utils"
	"reflect"
)

func (b *SqlBuilder) Page(current int64, pageSize int64) *SqlBuilder {
	if pageSize < 1 || current < 1 {
		b.db.getCtx().err = errors.New("pageSize,current must be greater than 0")
	}
	b.other = PageConfig{
		pageSize: pageSize,
		current:  current,
	}
	return b
}

type PageConfig struct {
	pageSize int64
	current  int64
}

type PageResult struct {
	List     any   `json:"list"`
	PageSize int64 `json:"pageSize"`
	Current  int64 `json:"current"`
	Total    int64 `json:"total"`
}

// ScanPage 查询分页
func (b *SqlBuilder) ScanPage(dest any) (rowsNum int64, dto PageResult, err error) {
	db := b.db
	ctx := db.getCtx()
	if err = ctx.err; err != nil {
		return
	}
	if b.other == nil {
		err = errors.New("no set pageConfig")
		return
	}
	var total int64
	var size = b.other.(PageConfig).pageSize
	var current = b.other.(PageConfig).current

	ctx.initScanDestList(dest)
	if err = ctx.err; err != nil {
		return
	}

	b.initSelectSql()

	var countSql = "select count(*) " + b.otherSqlBuilder.String()

	if ctx.showSql {
		fmt.Println(countSql, b.otherSqlArgs)
	}
	rows, err := db.query(countSql, b.otherSqlArgs...)
	if err != nil {
		return
	}
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)
	for rows.Next() {
		box := reflect.ValueOf(&total).Interface()
		err = rows.Scan(box)
		if err != nil {
			return
		}
	}
	// 计算总页数

	var selectSql = b.query + " limit ? offset ?"
	var offset = (current - int64(1)) * size
	args := append(b.args, size, offset)

	if ctx.showSql {
		fmt.Println(selectSql, args)
	}
	listRows, err := db.query(selectSql, args...)
	if err != nil {
		return
	}

	num, err := ctx.ScanT(listRows)
	if err != nil {
		return
	}

	if num == 0 {
		dest = make([]any, 0)
	}
	dto = PageResult{
		List:     dest,
		PageSize: size,
		Current:  current,
		Total:    total,
	}
	return num, dto, nil
}
