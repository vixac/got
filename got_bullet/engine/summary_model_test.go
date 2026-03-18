package engine

import (
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

// fixedTime returns a time with nanoseconds to exercise sub-second precision paths.
func fixedTime() time.Time {
	return time.Date(2026, 3, 16, 10, 30, 45, 123456789, time.UTC)
}

// fixedTimeNoNanos returns a time with zero sub-second precision.
func fixedTimeNoNanos() time.Time {
	return time.Date(2026, 3, 16, 10, 30, 45, 0, time.UTC)
}

// TestNewDateTime_StoresQuotedDate checks that DateTime.Date contains the
// JSON-quoted date string produced by time.MarshalJSON.
func TestNewDateTime_StoresQuotedDate(t *testing.T) {
	dt, err := NewDateTime(fixedTimeNoNanos())
	assert.NilError(t, err)

	// time.MarshalJSON wraps the date in double-quotes because it produces a JSON string.
	if !strings.HasPrefix(dt.Date, `"`) || !strings.HasSuffix(dt.Date, `"`) {
		t.Errorf("expected DateTime.Date to be a JSON-quoted string, got: %s", dt.Date)
	}
}

// TestNewDateTime_MillisSet checks that Millis is populated correctly.
func TestNewDateTime_MillisSet(t *testing.T) {
	ts := fixedTimeNoNanos()
	dt, err := NewDateTime(ts)
	assert.NilError(t, err)
	assert.Equal(t, dt.Millis, ts.UnixMilli())
}

// TestToDate_RoundTrip checks that NewDateTime followed by ToDate recovers the original time
// (truncated to second precision, since RFC3339 has no sub-second resolution).
func TestToDate_RoundTrip_NoNanos(t *testing.T) {
	ts := fixedTimeNoNanos()
	dt, err := NewDateTime(ts)
	assert.NilError(t, err)

	recovered, err := dt.ToDate()
	assert.NilError(t, err)
	assert.Assert(t, recovered != nil)

	got := time.Time(*recovered).UTC()
	assert.Equal(t, got, ts.UTC())
}

// TestToDate_RoundTrip_WithNanos exercises the sub-second precision path.
// Go 1.20+ time.Parse accepts RFC3339 strings with fractional seconds, so the
// full nanosecond precision is preserved through the round-trip.
func TestToDate_RoundTrip_WithNanos(t *testing.T) {
	ts := fixedTime() // has 123456789 nanoseconds
	dt, err := NewDateTime(ts)
	assert.NilError(t, err)

	recovered, err := dt.ToDate()
	assert.NilError(t, err)
	assert.Assert(t, recovered != nil)

	got := time.Time(*recovered).UTC()
	assert.Equal(t, got, ts.UTC())
}

// TestToDate_NilReceiver checks that a nil *DateTime returns nil, nil.
func TestToDate_NilReceiver(t *testing.T) {
	var dt *DateTime
	result, err := dt.ToDate()
	assert.NilError(t, err)
	assert.Assert(t, result == nil)
}

// TestToDate_SpecialValue checks that a special DateTime (e.g. "now") returns nil, nil.
func TestToDate_SpecialValue(t *testing.T) {
	dt := &DateTime{Special: "n"}
	result, err := dt.ToDate()
	assert.NilError(t, err)
	assert.Assert(t, result == nil)
}

// TestEpochMillis_Normal checks EpochMillis for a regular DateTime.
func TestEpochMillis_Normal(t *testing.T) {
	ts := fixedTimeNoNanos()
	dt, err := NewDateTime(ts)
	assert.NilError(t, err)
	assert.Equal(t, dt.EpochMillis(), ts.UnixMilli())
}

// TestEpochMillis_Nil checks that a nil *DateTime returns 0.
func TestEpochMillis_Nil(t *testing.T) {
	var dt *DateTime
	assert.Equal(t, dt.EpochMillis(), int64(0))
}

// TestEpochMillis_Special checks that a special DateTime returns 0.
func TestEpochMillis_Special(t *testing.T) {
	dt := &DateTime{Special: "n"}
	assert.Equal(t, dt.EpochMillis(), int64(0))
}

// TestJsonDateToReadable checks the human-readable date string output.
func TestJsonDateToReadable(t *testing.T) {
	ts := time.Date(2026, 3, 16, 10, 30, 0, 0, time.UTC)
	dt, err := NewDateTime(ts)
	assert.NilError(t, err)

	readable, err := dt.JsonDateToReadable()
	assert.NilError(t, err)
	assert.Equal(t, readable, "2026-03-16")
}

// TestJsonDateToReadable_Nil checks nil *DateTime returns empty string.
func TestJsonDateToReadable_Nil(t *testing.T) {
	var dt *DateTime
	readable, err := dt.JsonDateToReadable()
	assert.NilError(t, err)
	assert.Equal(t, readable, "")
}

// TestNewDeadlineFromDateLookup_Now checks that "<now>" produces a special DateTime.
func TestNewDeadlineFromDateLookup_Now(t *testing.T) {
	now := time.Now()
	dt, err := NewDeadlineFromDateLookup("<now>", now)
	assert.NilError(t, err)
	assert.Equal(t, dt.Special, "n")
}

// TestNewDeadlineFromDateLookup_RelativeDay checks a day-offset string.
func TestNewDeadlineFromDateLookup_RelativeDay(t *testing.T) {
	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	dt, err := NewDeadlineFromDateLookup("3d", now)
	assert.NilError(t, err)

	recovered, err := dt.ToDate()
	assert.NilError(t, err)
	assert.Assert(t, recovered != nil)

	got := time.Time(*recovered).UTC()
	want := time.Date(2026, 3, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, got, want)
}

// TestNewDeadlineFromDateLookup_InvalidInput checks that an unknown string returns an error.
func TestNewDeadlineFromDateLookup_InvalidInput(t *testing.T) {
	_, err := NewDeadlineFromDateLookup("notadate", time.Now())
	assert.ErrorContains(t, err, "")
}

// TestNewTimeFromString_Valid checks that a valid layout string is parsed correctly.
func TestNewTimeFromString_Valid(t *testing.T) {
	// Layout used inside NewTimeFromString: "2006-01-02T15:04:05.999999Z"
	input := "2026-03-16T10:30:45.000000Z"
	result, err := NewTimeFromString(input)
	assert.NilError(t, err)
	assert.Assert(t, result != nil)

	got := time.Time(*result)
	want := time.Date(2026, 3, 16, 10, 30, 45, 0, time.UTC)
	assert.Equal(t, got, want)

	//now we test WITH QUOTES
	input = "\"2026-03-14T18:47:22.465879Z\""
	result, err = NewTimeFromString(input)
	assert.NilError(t, err)
	assert.Assert(t, result != nil)

	got = time.Time(*result)
	want = time.Date(2026, 3, 14, 18, 47, 22, 465879000, time.UTC)
	assert.Equal(t, got, want)

	//handle timezones.
	input = "2026-01-20T21:11:47.642531+07:00"
	result, err = NewTimeFromString(input)
	assert.NilError(t, err)
	assert.Assert(t, result != nil)

	got = time.Time(*result)
	// Compare as UTC instants: 21:11:47 +07:00 == 14:11:47 UTC
	want = time.Date(2026, 1, 20, 14, 11, 47, 642531000, time.UTC)
	assert.Equal(t, got.UTC(), want)
}

// TestNewTimeFromString_Invalid checks that an unparseable string returns an error.
func TestNewTimeFromString_Invalid(t *testing.T) {
	_, err := NewTimeFromString("not-a-date")
	assert.Assert(t, err != nil)
}

// TestDateTimeDate_ContainsNoUnescapedQuoteProblem documents that DateTime.Date
// holds a JSON-encoded string (with surrounding double-quotes). When this value
// is passed directly to json.Unmarshal it is interpreted as a JSON string and
// UnmarshalJSON receives the raw bytes including the quotes – matching the
// parse pattern in RFC3339Time.UnmarshalJSON.
func TestDateTimeDate_QuotedStringIsValidForUnmarshal(t *testing.T) {
	ts := fixedTimeNoNanos()
	dt, err := NewDateTime(ts)
	assert.NilError(t, err)

	// Sanity: Date must start and end with a double-quote.
	if len(dt.Date) < 2 || dt.Date[0] != '"' || dt.Date[len(dt.Date)-1] != '"' {
		t.Fatalf("DateTime.Date is not a JSON-quoted string: %q", dt.Date)
	}

	// ToDate must succeed — it relies on the quoted form.
	recovered, err := dt.ToDate()
	assert.NilError(t, err)
	assert.Assert(t, recovered != nil)
}
