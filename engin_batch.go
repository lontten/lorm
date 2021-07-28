package lorm

//批量操作
type EngineBatch struct {
	db DBer
	lormConf Lorm
	dialect  Dialect

	context OrmContext
}
