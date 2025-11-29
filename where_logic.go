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

func (w *WhereBuilder) And(wb *WhereBuilder, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}
	if wb == nil {
		return w
	}
	if wb.Invalid() {
		return w
	}
	if wb.err != nil {
		w.err = wb.err
		return w
	}
	w.andWheres = append(w.andWheres, *wb)
	return w
}

func (w *WhereBuilder) Or(wb *WhereBuilder, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}
	if wb == nil {
		return w
	}
	if wb.Invalid() {
		return w
	}
	if wb.err != nil {
		w.err = wb.err
		return w
	}
	w.wheres = append(w.wheres, *wb)
	return w
}

func (w *WhereBuilder) Not(condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}
	if w.Invalid() {
		return w
	}
	w.not = true
	return w
}
