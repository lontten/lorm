package lorm

// todo 下面未重构--------------
// 批量操作
type EngineBatch struct {
	dialect Dialecter

	context ormContext
}
