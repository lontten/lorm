package lorm


type EngineExtra struct {
	ormConf OrmConf
	dialect Dialect

	context OrmContext
}
