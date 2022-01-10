package ayear

import (
	"strconv"
	"time"
)

// AYear Since yyyy-01-01 ~ (yyyy+1)-01-01
type AYear struct {
	Since time.Time
	Until time.Time
}

const datetimeFormat = "2006-01-02"

func (u *AYear) String() string {
	return "Since: " + u.Since.String() + ", Until: " + u.Until.String()
}

// New is return new instance
func New(sinceYear int) (*AYear, error) {
	nowYearStart, err := time.Parse(datetimeFormat, strconv.Itoa(sinceYear)+"-01-01")
	if err != nil {
		return nil, err
	}
	until := nowYearStart.AddDate(1, 0, 0)
	return &AYear{
		Since: nowYearStart,
		Until: until,
	}, nil
}
