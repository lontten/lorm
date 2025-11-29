//  Copyright 2025 lontten lontten@163.com
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package lorm

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/lontten/lorm/utils"
)

// ListPage 查询分页
func (b *SqlBuilder[T]) ListPage() (dto PageResult[T], err error) {
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if err = ctx.err; err != nil {
		return
	}
	if b.pageConfig == nil {
		err = errors.New("no set pageConfig")
		return
	}
	var total int64
	var pageSize = b.pageConfig.pageSize
	var pageIndex = b.pageConfig.pageIndex

	var dest = &[]T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, false)
	if err = ctx.err; err != nil {
		return
	}

	b.initSelectSql()

	var countSql = b.countField
	if countSql == "" {
		countSql = "*"
	}

	countSql = "select count(" + countSql + ") " + b.otherSqlBuilder.String()

	dialect.getSql(countSql)
	ctx.originalArgs = b.otherSqlArgs
	ctx.printSql()

	if !ctx.noRun {
		if b.fakeTotalNum > 0 {
			total = b.fakeTotalNum
		} else {
			rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
			if err != nil {
				return dto, err
			}
			defer func(rows *sql.Rows) {
				utils.PanicErr(rows.Close())
			}(rows)
			for rows.Next() {
				box := reflect.ValueOf(&total).Interface()
				err = rows.Scan(box)
				if err != nil {
					return dto, err
				}
			}
		}
	}

	// 计算总页数
	var pageNum = total / pageSize
	if total%pageSize != 0 {
		pageNum++
	}

	var selectSql = b.query + " limit ? offset ?"
	var offset = (pageIndex - int64(1)) * pageSize
	args := append(b.args, pageSize, offset)

	dialect.getSql(selectSql)
	ctx.originalArgs = args
	ctx.printSql()

	if ctx.noRun {
		return dto, nil
	}
	if b.noGetList {
		dto = PageResult[T]{
			List:      *dest,
			PageSize:  pageSize,
			PageNum:   pageNum,
			PageIndex: pageIndex,
			Total:     total,
			HasMore:   total > pageSize*pageIndex,
		}
		return dto, nil
	}

	listRows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
	if err != nil {
		return
	}

	_, err = ctx.Scan(listRows)
	if err != nil {
		return
	}

	dto = PageResult[T]{
		List:      *dest,
		PageSize:  pageSize,
		PageNum:   pageNum,
		PageIndex: pageIndex,
		Total:     total,
		HasMore:   total > pageSize*pageIndex,
	}
	return dto, nil
}

// ListPageP 查询分页
func (b *SqlBuilder[T]) ListPageP() (dto PageResultP[T], err error) {
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if err = ctx.err; err != nil {
		return
	}
	if b.pageConfig == nil {
		err = errors.New("no set pageConfig")
		return
	}
	var total int64
	var pageSize = b.pageConfig.pageSize
	var pageIndex = b.pageConfig.pageIndex

	var dest = &[]*T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, true)
	if err = ctx.err; err != nil {
		return
	}

	b.initSelectSql()

	var countSql = b.countField
	if countSql == "" {
		countSql = "*"
	}

	countSql = "select count(" + countSql + ") " + b.otherSqlBuilder.String()

	dialect.getSql(countSql)
	ctx.originalArgs = b.otherSqlArgs
	ctx.printSql()

	if !ctx.noRun {
		if b.fakeTotalNum > 0 {
			total = b.fakeTotalNum
		} else {
			rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
			if err != nil {
				return dto, err
			}
			defer func(rows *sql.Rows) {
				utils.PanicErr(rows.Close())
			}(rows)
			for rows.Next() {
				box := reflect.ValueOf(&total).Interface()
				err = rows.Scan(box)
				if err != nil {
					return dto, err
				}
			}
		}
	}

	// 计算总页数
	var pageNum = total / pageSize
	if total%pageSize != 0 {
		pageNum++
	}

	var selectSql = b.query + " limit ? offset ?"
	var offset = (pageIndex - int64(1)) * pageSize
	args := append(b.args, pageSize, offset)

	dialect.getSql(selectSql)
	ctx.originalArgs = args
	ctx.printSql()

	if ctx.noRun {
		return dto, nil
	}
	if b.noGetList {
		dto = PageResultP[T]{
			List:      *dest,
			PageSize:  pageSize,
			PageNum:   pageNum,
			PageIndex: pageIndex,
			Total:     total,
			HasMore:   total > pageSize*pageIndex,
		}
		return dto, nil
	}

	listRows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
	if err != nil {
		return
	}

	_, err = ctx.Scan(listRows)
	if err != nil {
		return
	}

	dto = PageResultP[T]{
		List:      *dest,
		PageSize:  pageSize,
		PageNum:   pageNum,
		PageIndex: pageIndex,
		Total:     total,
		HasMore:   total > pageSize*pageIndex,
	}
	return dto, nil
}
