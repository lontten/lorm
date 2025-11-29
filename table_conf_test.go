package lorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Demo01 struct {
}

func (u Demo01) TableConf() *TableConfContext {
	return TableConf("demo01").
		PrimaryKeys("id").
		AutoColumn("id")
}

func TestTableConfContext_PrimaryKeys(t *testing.T) {
	as := assert.New(t)
	var d = Demo01{}
	as.Equal("demo01", d.TableConf().tableName)
	as.Equal([]string{"id"}, d.TableConf().primaryKeyColumnNames)
	as.Equal([]string{"id"}, d.TableConf().allAutoColumnName)
	as.Equal("id", d.TableConf().autoPrimaryKeyColumnName)
	as.Equal([]string{}, d.TableConf().otherAutoColumnName)
}
