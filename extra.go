package lorm

import (
	"errors"
	"reflect"
	"strings"
)

func (db DB) Page(size int, current int64) *SqlBuilder {
	if size < 1 || current < 1 {
		db.ctx.err = errors.New("size,current must be greater than 0")
	}
	return &SqlBuilder{
		db:          db,
		selectQuery: &strings.Builder{},
		otherQuery:  &strings.Builder{},
		other: PageCnfig{
			size:    size,
			current: current,
		},
	}
}

type PageCnfig struct {
	size    int
	current int64
}

type Page struct {
	Records interface{} `json:"records"`
	Size    int         `json:"size"`
	Current int64       `json:"current"`
	Total   int64       `json:"total"`
	Pages   int64       `json:"pages"`
}

// PageSelect 查询分页
func (b *SqlBuilder) PageScan(dest interface{}) (rowsNum int64, dto Page, err error) {
	if err = b.db.ctx.err; err != nil {
		return
	}
	if b.other == nil {
		err = errors.New("PageCnfig is nil")
		return
	}
	var total int64
	var size = b.other.(PageCnfig).size
	var current = b.other.(PageCnfig).current

	b.initSelectSql()
	b.db.ctx.initScanDestList(dest)
	b.db.ctx.checkScanDestField()
	if err = b.db.ctx.err; err != nil {
		return
	}
	var countSql = "select count(*) " + b.otherQuery.String()

	rows, err := b.db.dialect.query(countSql, b.otherArgs...)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		box := reflect.ValueOf(&total).Interface()
		err = rows.Scan(box)
		if err != nil {
			return
		}
	}
	// 计算总页数

	var selectSql = b.query + " limit ? offset ?"
	var offset = (current - int64(1)) * int64(size)
	var args = append(b.args, size, offset)
	listRows, err := b.db.dialect.query(selectSql, args...)
	if err != nil {
		return
	}
	defer listRows.Close()
	num, err := b.db.ctx.Scan(listRows)
	if err != nil {
		return
	}

	if num == 0 {
		dest = make([]interface{}, 0)
	}
	dto = Page{
		Records: dest,
		Size:    size,
		Current: current,
		Total:   total,
		Pages:   total / int64(size),
	}
	return num, dto, nil
}
