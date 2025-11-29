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
	"reflect"
	"strings"

	"github.com/lontten/lcore/v2/lcutils"
	"github.com/pkg/errors"
)

/*
*
各个语句之间的and or关系和具体的数据库无关，直接在这里实现，parse。
每个语句的具体sql生成和数据库有关，但是不需要其他参数，例如orm_config  orm_context 等，
所以，生成具体sql的方法 toSql 直接接受 外界传过来的 parseFun 处理函数，代码结构比较简单，
不然，whereBuilder 里面要有 dialecter 的两种实现，代码结构复杂

primaryKeyFieldNames 主键字段名称列表
*/

func (w WhereBuilder) toSql(f parseFun, primaryKeyColumnNames ...string) (string, []any, error) {
	if w.err != nil {
		return "", nil, w.err
	}
	_, _, sql, args, err := w.parse(f, primaryKeyColumnNames...)
	if err != nil {
		return "", nil, err
	}
	return sql, args, err
}
func (w *WhereBuilder) parsePkClause(not bool, primaryKeyColumnNames ...string) error {
	if w.err != nil {
		return w.err
	}
	if len(primaryKeyColumnNames) == 0 {
		return ErrNoPk
	}
	var args = w.clause.args
	argsLen := len(args)
	if argsLen == 0 {
		return nil
	}

	if not {
		w.not = !w.not
	}

	// 0未设置；1struct复合主键;2map复合主键;3单主键
	var kindType = 0

	var nw = W()
	for _, arg := range args {
		var nw2 = W()

		var v = reflect.ValueOf(arg)
		_, v, err := basePtrValue(v)
		if err != nil {
			return err
		}
		kind := v.Kind()

		if kind == reflect.Struct {
			if kindType == 0 {
				kindType = 1
			} else {
				if kindType != 1 {
					return ErrTypePkArgs
				}
			}
			list := getStructCV(v)
			for _, cv := range list {
				if cv.isSoftDel || cv.isZero {
					continue
				}
				if lcutils.StrContainsAny(cv.columnName, primaryKeyColumnNames...) {
					nw2.fieldValue(cv.columnName, cv.value)
				}
			}

		} else if kind == reflect.Map {
			if kindType == 0 {
				kindType = 2
			} else {
				if kindType != 2 {
					return ErrTypePkArgs
				}
			}
			list := getMapCV(v)
			for _, cv := range list {
				if cv.isSoftDel || cv.isZero {
					continue
				}
				if lcutils.StrContainsAny(cv.columnName, primaryKeyColumnNames...) {
					nw2.fieldValue(cv.columnName, cv.value)
				}
			}
		} else {
			kindType = 3
		}

		nw.Or(nw2)
	}

	if kindType == 3 {
		if len(primaryKeyColumnNames) != 1 {
			return ErrNeedMultiPk
		}
		nw.In(primaryKeyColumnNames[0], args...)
	}

	w.clause = nil
	w.And(nw)
	return nil
}
func (w WhereBuilder) parse(f parseFun, primaryKeyFieldNames ...string) (hasOr bool, hasAnd bool, sql string, args []any, err error) {
	if w.err != nil {
		return false, false, "", nil, w.err
	}
	sb := strings.Builder{}
	var ors = w.wheres
	var ands = w.andWheres
	var allArgs []any

	var orLen = len(ors)
	var andLen = len(ands)

	var orNum = orLen - 1
	if andLen > 0 {
		orNum++
	}
	var localHasOr = orNum > 0
	var localHasAnd = andLen > 1

	if w.clause != nil {
		var c = *w.clause
		if c.Type == PrimaryKeys || c.Type == FilterPrimaryKeys {
			var _err = w.parsePkClause(c.Type == FilterPrimaryKeys, primaryKeyFieldNames...)
			if _err != nil {
				err = _err
				return
			}
			return w.parse(f, primaryKeyFieldNames...)
		} else {
			result, _err := f(c)
			if _err != nil {
				err = errors.Wrap(_err, "parse WhereBuilder")
				return
			}
			if w.not {
				sb.WriteString("NOT (")
			}
			sb.WriteString(result)
			if w.not {
				sb.WriteString(")")
			}
			return false, false, sb.String(), c.args, nil
		}
	}

	if !w.has() {
		return
	}

	if w.not {
		sb.WriteString("NOT (")
	}

	var needS = orLen > 0 && andLen > 1
	if needS {
		sb.WriteString("(")
	}

	for i, wt := range ands {
		var _hasOr, _hasAnd, _sql, _args, _err = wt.parse(f, primaryKeyFieldNames...)
		if _err != nil {
			err = errors.Wrap(_err, "parse WhereBuilder")
			return
		}
		if _hasOr {
			hasOr = true
		}
		if _hasAnd {
			hasAnd = true
		}
		allArgs = append(allArgs, _args...)
		if localHasAnd && _hasOr {
			sb.WriteString("(")
		}
		sb.WriteString(_sql)
		if localHasAnd && _hasOr {
			sb.WriteString(")")
		}
		if i < andLen-1 {
			sb.WriteString(" AND ")
			hasAnd = true
		}
	}

	if needS {
		sb.WriteString(")")
	}
	if andLen > 0 && orLen > 0 {
		sb.WriteString(" OR ")
		hasOr = true
	}

	for i, wt := range ors {
		var _hasOr, _hasAnd, _sql, _args, _err = wt.parse(f, primaryKeyFieldNames...)
		if _err != nil {
			err = errors.Wrap(_err, "parse WhereBuilder")
			return
		}
		if _hasOr {
			hasOr = true
		}
		if _hasAnd {
			hasAnd = true
		}
		allArgs = append(allArgs, _args...)
		if localHasOr && _hasAnd {
			sb.WriteString("(")
		}
		sb.WriteString(_sql)
		if localHasOr && _hasAnd {
			sb.WriteString(")")
		}
		if i < orLen-1 {
			sb.WriteString(" OR ")
			hasOr = true
		}
	}
	if w.not {
		sb.WriteString(")")
	}

	return hasOr, hasAnd, sb.String(), allArgs, nil
}
