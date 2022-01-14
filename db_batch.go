package lorm

//批量操作
type EngineBatch struct {
	dialect Dialect

	context OrmContext
}
