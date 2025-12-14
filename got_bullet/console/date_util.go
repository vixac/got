package console

import (
	"errors"
	"fmt"
	"math"
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

const (
	PastMany   = -2
	Yesterday  = -1
	Today      = 0
	Tomorrow   = 1
	FutureMany = 2
)

type SpaceTime struct {
	TimeType int
}
type RFC3339Time time.Time

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).UTC().Format(time.RFC3339) + `"`), nil
}

func (t *RFC3339Time) UnmarshalJSON(b []byte) error {
	parsed, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}
	*t = RFC3339Time(parsed)
	return nil
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

// HumanizeDate converts target into a relative human-readable string
// compared to the reference time "now".
func HumanizeDate(target, now time.Time) (string, SpaceTime) {
	// Normalize both times to midnight to avoid hour drift
	targetDay := normalizeDate(target)
	nowDay := normalizeDate(now)

	diffDays := int(targetDay.Sub(nowDay).Hours() / 24)

	switch diffDays {
	case -1:
		return "yesterday", SpaceTime{TimeType: Yesterday}
	case 0:
		return "today", SpaceTime{TimeType: Today}
	case 1:
		return "tomorrow", SpaceTime{Tomorrow}
	}

	absDays := int(math.Abs(float64(diffDays)))

	// Switch to weeks if far away
	if absDays > 100 {
		weeks := absDays / 7
		if diffDays < 0 {
			return fmt.Sprintf("%d weeks ago", weeks), SpaceTime{TimeType: PastMany}
		}
		return fmt.Sprintf("in %d weeks", weeks), SpaceTime{TimeType: FutureMany}
	}

	if diffDays < 0 {
		return fmt.Sprintf("%d days ago", absDays), SpaceTime{TimeType: PastMany}
	}
	return fmt.Sprintf("in %d days", absDays), SpaceTime{TimeType: FutureMany}
}

func normalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
