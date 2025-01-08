package types

import (
	"fmt"
	"testing"
	"time"
)

func TestNowTime(t *testing.T) {
	now := time.Now()

	dateEnd := DateOf(now)
	date := dateEnd.Time.AddDate(0, 0, 1)
	of := DateOf(date)
	var dateTimeEnd = of.ToDateTimeP()

	fmt.Println(dateTimeEnd.String())
}
