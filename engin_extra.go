package lorm


type EngineExtra struct {
	core    OrmCore
	dialect Dialect

	context OrmContext
}
