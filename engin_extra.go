package lorm


type EngineExtra struct {
	db        DBer
	ormConfig LormConf

	context OrmContext
}
