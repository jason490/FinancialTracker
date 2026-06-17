package dashboard

import "time"

// CurrentMonthBounds returns unix start (inclusive) and end (exclusive) for the current calendar month.
func CurrentMonthBounds() (start, end int64) {
	now := time.Now()
	loc := now.Location()
	startTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	endTime := startTime.AddDate(0, 1, 0)
	return startTime.Unix(), endTime.Unix()
}
