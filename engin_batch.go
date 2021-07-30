package lorm

//批量操作
type EngineBatch struct {
	core    OrmCore
	dialect Dialect

	context OrmContext
}
