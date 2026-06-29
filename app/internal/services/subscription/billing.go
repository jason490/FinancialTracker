package subscription

import "time"

// BillingPeriodStart returns the Unix start of the current monthly billing period for an anchor date.
func BillingPeriodStart(anchorUnix, nowUnix int64) int64 {
	anchor := time.Unix(anchorUnix, 0).UTC()
	now := time.Unix(nowUnix, 0).UTC()

	year, month, _ := now.Date()
	start := anniversaryInMonth(year, month, anchor)
	if now.Before(start) {
		year, month = previousMonth(year, month)
		start = anniversaryInMonth(year, month, anchor)
	}
	return start.Unix()
}

// BillingPeriodEnd returns the Unix end of the billing period that begins at periodStartUnix.
func BillingPeriodEnd(periodStartUnix, anchorUnix int64) int64 {
	start := time.Unix(periodStartUnix, 0).UTC()
	anchor := time.Unix(anchorUnix, 0).UTC()
	year, month, _ := start.Date()
	year, month = nextMonth(year, month)
	return anniversaryInMonth(year, month, anchor).Unix()
}

// anniversaryInMonth returns midnight UTC on the anchor's day-of-month, clamped to the month's length.
func anniversaryInMonth(year int, month time.Month, anchor time.Time) time.Time {
	day := anchor.Day()
	lastDay := daysInMonth(year, month)
	if day > lastDay {
		day = lastDay
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func previousMonth(year int, month time.Month) (int, time.Month) {
	if month == time.January {
		return year - 1, time.December
	}
	return year, month - 1
}

func nextMonth(year int, month time.Month) (int, time.Month) {
	if month == time.December {
		return year + 1, time.January
	}
	return year, month + 1
}
