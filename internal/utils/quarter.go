package utils

import (
	"fmt"
	"time"
)

type Quarter struct {
	Year    int
	Quarter int
}

func ParseQuarter(s string) (*Quarter, error) {
	var q Quarter
	_, err := fmt.Sscanf(s, "%d-T%d", &q.Year, &q.Quarter)
	if err != nil {
		return nil, err
	}
	if q.Quarter < 1 || q.Quarter > 4 {
		return nil, fmt.Errorf("invalid quarter: %d", q.Quarter)
	}
	return &q, nil
}

func (q *Quarter) From() time.Time {
	return time.Date(q.Year, time.Month((q.Quarter-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
}

func (q *Quarter) To() time.Time {
	return time.Date(q.Year, time.Month(q.Quarter*3), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
}
