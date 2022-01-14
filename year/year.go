package year

import (
	"strconv"
	"time"
)

// Year Since yyyy-01-01 ~ (yyyy+1)-01-01
type Year struct {
	Since time.Time
	Until time.Time
}

const datetimeFormat = "2006-01-02"

func (u *Year) String() string {
	return "Since: " + u.Since.String() + ", Until: " + u.Until.String()
}

// New is return new instance
func New(sinceYear int) (*Year, error) {
	nowYearStart, err := time.Parse(datetimeFormat, strconv.Itoa(sinceYear)+"-01-01")
	if err != nil {
		return nil, err
	}
	until := nowYearStart.AddDate(1, 0, 0)
	return &Year{
		Since: nowYearStart,
		Until: until,
	}, nil
}
