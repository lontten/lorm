// Copyright (c) 2024 lontten
// lorm is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lorm

// 参数分为struct、map ## map只能作为纯字段参数，struct还可以作为表名
// 接受是comp、atom都可以

//分类：
// 1. 表名							--nameDest
// 2. 参数struct、map				--paramDest
// 3. 参数sturct + 表名				--tableDest
// 4. scan接收						--scanDest
// 5. 参数sturct + 表名 + scan接收	--targetDest

// scan需要v，参数需要vn

// * struct
func (db *lnDB) setNameDest(v any) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initNameDest(v) //初始化参数
	db.core.getCtx().initConf()      //初始化表名
}

// * struct/map
func (db *lnDB) setParamDest(v any) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initModelDest(v)     //初始化参数
	db.core.getCtx().checkDestParamScan() //检查dest参数和接收
	db.core.getCtx().initColumnsValue()   //初始化cv
}

// * struct/map
func (db *lnDB) setTableDest(v any) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initModelDest(v)     //初始化参数
	db.core.getCtx().initConf()           //初始化表名
	db.core.getCtx().checkDestParamScan() //检查dest参数和接收
	db.core.getCtx().initColumnsValue()   //初始化cv
}

// * struct/map
func (db *lnDB) setScanDest(v any) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initModelDest(v)   //初始化参数
	db.core.getCtx().checkDestScan()    //检查dest参数和接收
	db.core.getCtx().initColumnsValue() //初始化cv
}

// * struct
// model 设置 model,表名
// dest  设置 dest
func (ctx *ormContext) setModelDest(v any) {
	if ctx.hasErr() {
		return
	}
	ctx.initModelDest(v)   //初始化参数
	ctx.initConf()         //初始化表名，主键
	ctx.initColumnsValue() //初始化cv
}
