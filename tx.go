package lorm

type Tx struct {
	Base    *EngineBase
	Extra   *EngineExtra
	Table   *EngineTable
	Classic *EngineClassic
}

func (tx *Tx) DriverName() string {
	return tx.Base.db.dbConfig.DriverName()
}

