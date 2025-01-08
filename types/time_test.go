package types

import (
	"fmt"
	"testing"
	"time"
)

func TestNowTime(t *testing.T) {
	now := time.Now()

	date := DateOf(now)
	fmt.Println(date.Time.String())

	time := TimeOf(now)
	fmt.Println(time.Time.String())
}
