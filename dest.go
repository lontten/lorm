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
func (db *lnDB) setNameDest(v interface{}) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initNameDest(v) //初始化参数
	db.core.getCtx().initTableName() //初始化表名
}

// * struct/map
func (db *lnDB) setParamDest(v interface{}) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initTargetDest(v)    //初始化参数
	db.core.getCtx().checkDestParamScan() //检查dest参数和接收
	db.core.getCtx().initColumnsValue()   //初始化cv
}

// * struct/map
func (db *lnDB) setTableDest(v interface{}) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initTargetDest(v)    //初始化参数
	db.core.getCtx().initTableName()      //初始化表名
	db.core.getCtx().checkDestParamScan() //检查dest参数和接收
	db.core.getCtx().initColumnsValue()   //初始化cv
}

// * struct/map
func (db *lnDB) setScanDest(v interface{}) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initTargetDest(v)  //初始化参数
	db.core.getCtx().checkDestScan()    //检查dest参数和接收
	db.core.getCtx().initColumnsValue() //初始化cv
}

// * struct
// target scanDest 一个struct
func (db *lnDB) setTargetDest(v interface{}) {
	if db.core.hasErr() {
		return
	}
	db.core.getCtx().initTargetDest(v)    //初始化参数
	db.core.getCtx().initTableName()      //初始化表名
	db.core.getCtx().checkDestParamScan() //检查dest参数和接收
	db.core.getCtx().initColumnsValue()   //初始化cv
}
