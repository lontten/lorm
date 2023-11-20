package lorm

import (
	"testing"
)

func Test_commit(t *testing.T) {
	db := lnDB{}
	tx := db.BeginTx(nil, nil)
	tx.Commit()

}
