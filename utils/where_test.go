package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenwhereToken(t *testing.T) {
	as := assert.New(t)
	as.Equal([]string{"a = ? ", "b = ? ", "c = ? "}, GenwhereToken([]string{"a", "b", "c"}))

}

func Test_genBatch(t *testing.T) {
	as := assert.New(t)
	as.Equal("(a = ? and b = ? and c = ? )", _genWhere([]string{"a", "b", "c"}))

}
