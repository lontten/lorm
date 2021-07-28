package lorm


type EngineExtra struct {
	db DBer
	lormConf Lorm
	dialect  Dialect

	context OrmContext
}
