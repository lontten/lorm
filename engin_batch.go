package lorm

//批量操作
type EngineBatch struct {
	db        DBer
	ormConfig LormConf

	context OrmContext
}
