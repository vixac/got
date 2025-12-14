package console

import (
	"testing"
	"time"
)

func mustDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func TestParseRelativeDate_Nd(t *testing.T) {
	now := mustDate("2025-01-10")

	tests := []struct {
		in   string
		want string
	}{
		{"5d", "2025-01-15"},
		{"0d", "2025-01-10"},
		{"-3d", "2025-01-07"},
	}

	for _, tc := range tests {
		got, err := ParseRelativeDate(tc.in, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != mustDate(tc.want) {
			t.Errorf("input %q => got %v want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseRelativeDate_EOW(t *testing.T) {
	// Friday
	now := mustDate("2025-01-10")  // Friday
	want := mustDate("2025-01-12") // Sunday

	got, err := ParseRelativeDate("eow", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("eow => got %v want %v", got, want)
	}
}

func TestParseRelativeDate_EOM(t *testing.T) {
	now := mustDate("2025-01-10")
	want := mustDate("2025-01-31")

	got, err := ParseRelativeDate("eom", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("eom => got %v want %v", got, want)
	}
}

func TestParseRelativeDate_Weekdays(t *testing.T) {
	now := mustDate("2025-01-10") // Friday

	tests := []struct {
		in   string
		want string
	}{
		{"fri", "2025-01-17"}, // next Friday (7 days later)
		{"sat", "2025-01-11"},
		{"sun", "2025-01-12"},
		{"mon", "2025-01-13"},
		{"tue", "2025-01-14"},
		{"wed", "2025-01-15"},
		{"thu", "2025-01-16"},
	}

	for _, tc := range tests {
		got, err := ParseRelativeDate(tc.in, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != mustDate(tc.want) {
			t.Errorf("%s => got %v want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseRelativeDate_Invalid(t *testing.T) {
	now := mustDate("2025-01-10")

	invalid := []string{"", "xx", "99x", "tomorrow", "d5"}

	for _, in := range invalid {
		_, err := ParseRelativeDate(in, now)
		if err == nil {
			t.Errorf("expected error for input %q", in)
		}
	}
}

func TestHumanizeDate(t *testing.T) {
	ref := time.Date(2025, 12, 6, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		target time.Time
		want   string
	}{
		{
			name:   "yesterday",
			target: ref.AddDate(0, 0, -1),
			want:   "yesterday",
		},
		{
			name:   "today",
			target: ref,
			want:   "today",
		},
		{
			name:   "tomorrow",
			target: ref.AddDate(0, 0, 1),
			want:   "tomorrow",
		},
		{
			name:   "5 days ago",
			target: ref.AddDate(0, 0, -5),
			want:   "5 days ago",
		},
		{
			name:   "in 10 days",
			target: ref.AddDate(0, 0, 10),
			want:   "in 10 days",
		},
		{
			name:   "over 100 days past uses weeks",
			target: ref.AddDate(0, 0, -140),
			want:   "20 weeks ago",
		},
		{
			name:   "over 100 days future uses weeks",
			target: ref.AddDate(0, 0, 154),
			want:   "in 22 weeks",
		},
		{
			name:   "time of day ignored",
			target: time.Date(2025, 12, 6, 23, 59, 59, 0, time.UTC),
			want:   "today",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//VX:Note spacetime is not tested.
			got, _ := HumanizeDate(tt.target, ref)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
