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

	"github.com/pkg/errors"
)

func (ctx *ormContext) initScanDestList(dest any) {
	if ctx.hasErr() {
		return
	}

	v := reflect.ValueOf(dest)
	isPtr, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if !isPtr {
		ctx.err = errors.New("scan must be a ptr ")
		return
	}

	isSlice, t := baseSliceType(v.Type())
	if !isSlice {
		ctx.err = errors.New("scanList must be a slice ")
		return
	}
	_, base := basePtrType(t)
	ctyp := checkAtomType(base)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = isPtr

	ctx.destBaseType = base
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = true
	ctx.destSliceItemIsPtr = t.Kind() == reflect.Ptr
}

func (ctx *ormContext) initScanDestListT(dest any, v, baseV reflect.Value, t reflect.Type, destSliceItemIsPtr bool) {
	if ctx.hasErr() {
		return
	}

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = true

	ctx.destBaseValue = baseV
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = true
	ctx.destSliceItemIsPtr = destSliceItemIsPtr
}

func (ctx *ormContext) initScanDestOne(dest any) {
	if ctx.hasErr() {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if !isPtr {
		ctx.err = errors.New("scan must be a ptr ")
		return
	}

	t := v.Type()

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = isPtr

	ctx.destBaseValue = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}

// dest 类型 struct 、所有 valuer 类型
func (ctx *ormContext) initScanDestOneT(dest any) {
	if ctx.hasErr() {
		return
	}

	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = true

	ctx.destBaseValue = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}
