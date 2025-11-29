package lorm

func Prepare(db Engine, query string) (Stmter, error) {
	db = db.init()
	return db.prepare(query)
}

func StmtExec(db Stmter, args ...any) (int64, error) {
	exec, err := db.exec(args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}
