package lorm


type EngineExtra struct {
	db       DBer
	lormConf OrmCore
	dialect  Dialect

	context OrmContext
}
