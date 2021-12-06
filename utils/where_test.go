package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenwhereToken(t *testing.T) {
	as := assert.New(t)
	as.Equal([]string{"a = ? ", "b = ? ","c = ? "}, GenwhereToken([]string{"a", "b", "c"}))

}
