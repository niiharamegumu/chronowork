package timeutil

import (
	"testing"
	"time"
)

func TestIsToday(t *testing.T) {
	today := time.Now()
	want := true
	if IsToday(today) != want {
		t.Errorf("IsToday(%v) = false, want true", today)
	}
}

func TestIsNotToday(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1)
	want := false
	if IsToday(yesterday) != want {
		t.Errorf("IsToday(%v) = true, want false", yesterday)
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "00:00:00"},
		{30, "00:00:30"},
		{60, "00:01:00"},
		{90, "00:01:30"},
		{3600, "01:00:00"},
		{3661, "01:01:01"},
		{7265, "02:01:05"},
	}

	for _, tc := range tests {
		result := FormatTime(tc.seconds)
		if result != tc.expected {
			t.Errorf("FormatTime(%d) = %s, want %s", tc.seconds, result, tc.expected)
		}
	}
}

func TestFormatWithPersonDay(t *testing.T) {
	// 8 hours = 28800 seconds = 1 person day (with 8-hour workday)
	tests := []struct {
		seconds   int
		personDay uint
		display   bool
		expected  string
	}{
		{0, 8, true, "00:00:00(0.00)"},
		{3600, 8, true, "01:00:00(0.12)"},
		{28800, 8, true, "08:00:00(1.00)"},
		{3600, 8, false, "01:00:00"},
		{3600, 0, true, "01:00:00"},
	}

	for _, tc := range tests {
		result := FormatWithPersonDay(tc.seconds, tc.personDay, tc.display)
		if result != tc.expected {
			t.Errorf("FormatWithPersonDay(%d, %d, %v) = %s, want %s",
				tc.seconds, tc.personDay, tc.display, result, tc.expected)
		}
	}
}

func TestSecondsToHourAndMinute(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "00:00"},
		{3600, "01:00"},
		{5400, "01:30"},
	}

	for _, tc := range tests {
		result := SecondsToHourAndMinute(tc.seconds)
		if result != tc.expected {
			t.Errorf("SecondsToHourAndMinute(%d) = %s, want %s", tc.seconds, result, tc.expected)
		}
	}
}

func TestTodayEndTime(t *testing.T) {
	endTime := TodayEndTime()
	now := time.Now()

	if endTime.Year() != now.Year() || endTime.Month() != now.Month() || endTime.Day() != now.Day() {
		t.Error("TodayEndTime should return today's date")
	}
	if endTime.Hour() != 23 || endTime.Minute() != 59 || endTime.Second() != 59 {
		t.Errorf("TodayEndTime should return 23:59:59, got %02d:%02d:%02d",
			endTime.Hour(), endTime.Minute(), endTime.Second())
	}
}

func TestRelativeStartTimeWithDays(t *testing.T) {
	now := time.Now()

	// 0 days = today
	result := RelativeStartTimeWithDays(0)
	if result.Year() != now.Year() || result.Month() != now.Month() || result.Day() != now.Day() {
		t.Error("RelativeStartTimeWithDays(0) should return today")
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Error("RelativeStartTimeWithDays should return start of day (00:00:00)")
	}

	// 7 days ago
	result = RelativeStartTimeWithDays(7)
	expected := now.AddDate(0, 0, -7)
	if result.Year() != expected.Year() || result.Month() != expected.Month() || result.Day() != expected.Day() {
		t.Error("RelativeStartTimeWithDays(7) should return 7 days ago")
	}
}
