package lorm

//批量操作
type EngineBatch struct {
	ormConf OrmConf
	dialect Dialect

	context OrmContext
}
