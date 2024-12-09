package lorm

//
//import (
//	"errors"
//	"reflect"
//	"strings"
//)
//
//// todo 下面未重构--------------
//func (db *lnDB) Page(size int, current int64) *SqlBuilder {
//	if size < 1 || current < 1 {
//		db.core.getCtx().err = errors.New("size,current must be greater than 0")
//	}
//	return &SqlBuilder{
//		core:        db.core,
//		selectQuery: &strings.SelectBuilder{},
//		otherQuery:  &strings.SelectBuilder{},
//		other: PageCnfig{
//			size:    size,
//			current: current,
//		},
//	}
//}
//
//type PageCnfig struct {
//	size    int
//	current int64
//}
//
//type Page struct {
//	Records any   `json:"records"`
//	Size    int   `json:"size"`
//	Current int64 `json:"current"`
//	Total   int64 `json:"total"`
//	Pages   int64 `json:"pages"`
//}
//
//// PageSelect 查询分页
//func (b *SqlBuilder) PageScan(dest any) (rowsNum int64, dto Page, err error) {
//	if err = b.core.getCtx().err; err != nil {
//		return
//	}
//	if b.other == nil {
//		err = errors.New("PageCnfig is nil")
//		return
//	}
//	var total int64
//	var size = b.other.(PageCnfig).size
//	var current = b.other.(PageCnfig).current
//
//	b.initSelectSql()
//	b.core.getCtx().initScanDestList(dest)
//	b.core.getCtx().checkScanDestField()
//	if err = b.core.getCtx().err; err != nil {
//		return
//	}
//	var countSql = "select count(*) " + b.otherQuery.String()
//
//	rows, err := b.core.doQuery(countSql, b.otherArgs...)
//
//	if err != nil {
//		return
//	}
//
//	defer rows.Close()
//	for rows.Next() {
//		box := reflect.ValueOf(&total).Interface()
//		err = rows.Scan(box)
//		if err != nil {
//			return
//		}
//	}
//	// 计算总页数
//
//	var selectSql = b.query + " limit ? offset ?"
//	var offset = (current - int64(1)) * int64(size)
//	args := append(b.args, size, offset)
//	listRows, err := b.core.doQuery(selectSql, args...)
//
//	if err != nil {
//		return
//	}
//	defer listRows.Close()
//	num, err := b.core.getCtx().ScanT(listRows)
//	if err != nil {
//		return
//	}
//
//	if num == 0 {
//		dest = make([]any, 0)
//	}
//	dto = Page{
//		Records: dest,
//		Size:    size,
//		Current: current,
//		Total:   total,
//		Pages:   total / int64(size),
//	}
//	return num, dto, nil
//}
