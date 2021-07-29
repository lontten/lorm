package lorm

//批量操作
type EngineBatch struct {
	db       DBer
	lormConf OrmCore
	dialect  Dialect

	context OrmContext
}
