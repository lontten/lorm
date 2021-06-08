package lorm

//批量操作
type EngineBatch struct {
	context   OrmContext

	db DB
}
