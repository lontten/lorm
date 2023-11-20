package lorm

import (
	"testing"
)

func Test_commit(t *testing.T) {
	db := DB{}
	tx := db.BeginTx(nil, nil)
	tx.Commit()

}
