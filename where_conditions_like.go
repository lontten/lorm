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
	"fmt"
)

type likeType int

const (
	LikeAnywhere   likeType = iota // 两边加通配符，匹配任意位置
	LikeEndsWith                   // 左边加通配符，匹配结尾
	LikeStartsWith                 // 右边加通配符，匹配开头
)

// columns 多个字段，任意一个字段满足 LIKE ?
func (w *WhereBuilder) _like(key *string, likeType likeType, noLike bool, columns ...string) *WhereBuilder {
	if w.err != nil {
		return w
	}
	if key == nil {
		w.err = fmt.Errorf("invalid use of like: key  is nil.")
		return w
	}
	if *key == "" {
		return w
	}

	var k = ""
	switch likeType {
	case LikeAnywhere:
		k = "%" + *key + "%"
	case LikeEndsWith:
		k = "%" + *key
	case LikeStartsWith:
		k = *key + "%"
	}

	var likeTokenType = Like
	if noLike {
		likeTokenType = NotLike
	}

	likeW := W()
	for _, field := range columns {
		likeW.wheres = append(likeW.wheres, WhereBuilder{
			clause: &Clause{
				Type:  likeTokenType,
				query: field,
				args:  []any{k},
			},
		})
	}
	w.And(likeW)
	return w
}

// LikeP
// LIKE %?%
func (w *WhereBuilder) LikeP(query string, arg *string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w._like(arg, LikeAnywhere, false, query)
	return w
}

// LikeLeftP
// LIKE %?
func (w *WhereBuilder) LikeLeftP(query string, arg *string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w._like(arg, LikeEndsWith, false, query)
	return w
}

// LikeRightP
// LIKE ?%
func (w *WhereBuilder) LikeRightP(query string, arg *string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w._like(arg, LikeStartsWith, false, query)
	return w
}

// Like
// LIKE ?
func (w *WhereBuilder) Like(query string, arg string, condition ...bool) *WhereBuilder {
	w.LikeP(query, &arg, condition...)
	return w
}

func (w *WhereBuilder) LikeLeft(query string, arg string, condition ...bool) *WhereBuilder {
	w.LikeLeftP(query, &arg, condition...)
	return w
}
func (w *WhereBuilder) LikeRight(query string, arg string, condition ...bool) *WhereBuilder {
	w.LikeRightP(query, &arg, condition...)
	return w
}

// NoLikeP
// NOT LIKE %?%
func (w *WhereBuilder) NoLikeP(query string, arg *string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w._like(arg, LikeAnywhere, true, query)
	return w
}
func (w *WhereBuilder) NoLikeLeftP(query string, arg *string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w._like(arg, LikeEndsWith, true, query)
	return w
}
func (w *WhereBuilder) NoLikeRightP(query string, arg *string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w._like(arg, LikeStartsWith, true, query)
	return w
}
func (w *WhereBuilder) NoLike(query string, arg string, condition ...bool) *WhereBuilder {
	w.NoLikeP(query, &arg, condition...)
	return w
}
func (w *WhereBuilder) NoLikeLeft(query string, arg string, condition ...bool) *WhereBuilder {
	w.NoLikeLeftP(query, &arg, condition...)
	return w
}
func (w *WhereBuilder) NoLikeRight(query string, arg string, condition ...bool) *WhereBuilder {
	w.NoLikeRightP(query, &arg, condition...)
	return w
}

// LikeAny
// 多个字段，任意一个字段满足 LIKE ?
func (w *WhereBuilder) LikeAny(key *string, columns ...string) *WhereBuilder {
	w._like(key, LikeAnywhere, false, columns...)
	return w
}

// BoolLikeAny
// 多个字段，任意一个字段满足 LIKE ?
func (w *WhereBuilder) BoolLikeAny(condition bool, key *string, columns ...string) *WhereBuilder {
	if !condition {
		return w
	}
	w.LikeAny(key, columns...)
	return w
}

// LikeLeftAny
// 多个字段，任意一个字段满足 LIKE ?
func (w *WhereBuilder) LikeLeftAny(key *string, columns ...string) *WhereBuilder {
	w._like(key, LikeEndsWith, false, columns...)
	return w
}
func (w *WhereBuilder) BoolLikeLeftAny(condition bool, key *string, columns ...string) *WhereBuilder {
	if !condition {
		return w
	}
	w.LikeLeftAny(key, columns...)
	return w
}

// LikeRightAny
// 多个字段，任意一个字段满足 LIKE ?
func (w *WhereBuilder) LikeRightAny(key *string, columns ...string) *WhereBuilder {
	w._like(key, LikeStartsWith, false, columns...)
	return w
}
func (w *WhereBuilder) BoolLikeRightAny(condition bool, key *string, columns ...string) *WhereBuilder {
	if !condition {
		return w
	}
	w.LikeRightAny(key, columns...)
	return w
}
