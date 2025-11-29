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

// 句子类型，用于whereBuilder
type clauseType int

const (
	Eq clauseType = iota
	Neq
	Less
	LessEq
	Greater
	GreaterEq
	Like
	NotLike
	In
	NotIn
	Between
	NotBetween
	IsNull
	IsNotNull
	IsFalse

	PrimaryKeys       // 主键
	FilterPrimaryKeys // 过滤主键

	// Contains 包含
	// pg 独有
	// [1] @< [1,2]
	Contains
)
