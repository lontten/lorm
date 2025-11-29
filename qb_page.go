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
	"errors"
)

// CountField 自定义count字段
func (b *SqlBuilder[T]) CountField(field string, conditions ...bool) *SqlBuilder[T] {
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	b.countField = field
	return b
}

// FakerTotalNum 分页时，直接使用 fakeTotalNum，不再查询实际总数
func (b *SqlBuilder[T]) FakerTotalNum(num int64, conditions ...bool) *SqlBuilder[T] {
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	b.fakeTotalNum = num
	return b
}

// NoGetList 分页时，只查询数量，不返回数据列表
func (b *SqlBuilder[T]) NoGetList(conditions ...bool) *SqlBuilder[T] {
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	b.noGetList = true
	return b
}

func (b *SqlBuilder[T]) Page(pageIndex int64, pageSize int64) *SqlBuilder[T] {
	if pageSize < 1 || pageIndex < 1 {
		b.db.getCtx().err = errors.New("pageSize,pageIndex must be greater than 0")
	}
	b.pageConfig = &PageConfig{
		pageSize:  pageSize,
		pageIndex: pageIndex,
	}
	return b
}

type PageConfig struct {
	pageSize  int64
	pageIndex int64
}

type PageResult[T any] struct {
	List      []T   `json:"list"`      // 结果
	PageSize  int64 `json:"pageSize"`  // 每页大小
	PageIndex int64 `json:"pageIndex"` // 当前页码
	Total     int64 `json:"total"`     // 总数
	PageNum   int64 `json:"totalPage"` // 总页数
	HasMore   bool  `json:"hasMore"`   // 是否有更多
}

type PageResultP[T any] struct {
	List      []*T  `json:"list"`      // 结果
	PageSize  int64 `json:"pageSize"`  // 每页大小
	PageIndex int64 `json:"pageIndex"` // 当前页码
	Total     int64 `json:"total"`     // 总数
	PageNum   int64 `json:"totalPage"` // 总页数
	HasMore   bool  `json:"hasMore"`   // 是否有更多
}
