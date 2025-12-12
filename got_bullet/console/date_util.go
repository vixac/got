package console

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var weekdayMap = map[string]time.Weekday{
	"sun": time.Sunday,
	"mon": time.Monday,
	"tue": time.Tuesday,
	"wed": time.Wednesday,
	"thu": time.Thursday,
	"fri": time.Friday,
	"sat": time.Saturday,
}

func ParseRelativeDate(input string, now time.Time) (time.Time, error) {
	in := strings.ToLower(strings.TrimSpace(input))

	// Case 1: Nd (number of days offset)
	if strings.HasSuffix(in, "d") && in != "wed" {
		numStr := strings.TrimSuffix(in, "d")
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return time.Time{}, errors.New("invalid Nd format")
		}
		return normalizeDate(now.AddDate(0, 0, n)), nil
	}

	// Case 2: End of week
	if in == "eow" {
		weekday := int(now.Weekday())
		daysToSunday := (7 - weekday) % 7
		return normalizeDate(now.AddDate(0, 0, daysToSunday)), nil
	}

	// Case 3: End of month
	if in == "eom" {
		y, m, _ := now.Date()
		loc := now.Location()
		firstNextMonth := time.Date(y, m+1, 1, 0, 0, 0, 0, loc)
		endOfMonth := firstNextMonth.AddDate(0, 0, -1)
		return normalizeDate(endOfMonth), nil
	}

	// Case 4: Next weekday
	if wd, ok := weekdayMap[in]; ok {
		todayW := now.Weekday()
		delta := (int(wd) - int(todayW) + 7) % 7
		if delta == 0 {
			delta = 7 // always next weekday, not today
		}
		return normalizeDate(now.AddDate(0, 0, delta)), nil
	}

	return time.Time{}, errors.New("unrecognized date expression")
}

func normalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
