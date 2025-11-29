package lorm

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createColBoxT(t *testing.T) {
	as := assert.New(t)
	as.True(true)

	rb := new(sql.RawBytes)
	fmt.Println(rb)
	fmt.Println(rb)
	fmt.Println(rb == nil)

	rb2 := new([]byte)
	fmt.Println(rb2)
	fmt.Println(rb2 == nil)

	var v *int = nil
	fmt.Println(v, v == nil)
	fmt.Println(v == nil || *v == 1)
}
