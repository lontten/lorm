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
	"reflect"
)

// 获取：参数，表名,scan
func (ctx *ormContext) initModelDest(dest any) {
	if ctx.hasErr() {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	ctx.scanIsPtr = isPtr
	t := v.Type()

	if ctx.checkParam {
		if _isStructType(t) {
			ctx.err = errors.New("dest need is struct")
			return
		}
		err = checkCompFieldVS(v)
		if err != nil {
			ctx.err = err
			return
		}
	}

	ctx.paramModelBaseV = v

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = isPtr

	ctx.destBaseValue = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = true

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}

// 只作为参数，主要用在insert,update时，获取更新的数据
func (ctx *ormContext) initTargetDestOne(dest any) {
	if ctx.hasErr() {
		return
	}
	v := reflect.ValueOf(dest)
	_, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}

	t := v.Type()

	if ctx.checkParam {
		if _isStructType(t) {
			ctx.err = errors.New("dest need is struct")
			return
		}
		err = checkCompFieldVS(v)
		if err != nil {
			ctx.err = err
			return
		}
	}
	ctx.paramModelBaseV = v

	ctx.destBaseValue = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = true
}
