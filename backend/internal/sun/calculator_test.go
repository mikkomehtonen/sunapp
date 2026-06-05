package sun

import (
	"math"
	"testing"
	"time"
)

func TestCalculateSunTimes_HelsinkiSummer(t *testing.T) {
	lat := 60.1699
	lon := 24.9384
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Europe/Helsinki")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SunriseUTC == "N/A" {
		t.Error("sunrise UTC should not be N/A for Helsinki on June 21")
	}
	if result.SunsetUTC == "N/A" {
		t.Error("sunset UTC should not be N/A for Helsinki on June 21")
	}
	if result.SunriseLocal == "N/A" {
		t.Error("sunrise local should not be N/A for Helsinki on June 21")
	}
	if result.SunsetLocal == "N/A" {
		t.Error("sunset local should not be N/A for Helsinki on June 21")
	}
	t.Logf("Helsinki summer (June 21, Europe/Helsinki): SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, TZ=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.Timezone)
}

func TestCalculateSunTimes_HelsinkiWinter(t *testing.T) {
	lat := 60.1699
	lon := 24.9384
	date := time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Europe/Helsinki")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SunriseUTC == "N/A" {
		t.Error("sunrise UTC should not be N/A for Helsinki on December 21")
	}
	if result.SunsetUTC == "N/A" {
		t.Error("sunset UTC should not be N/A for Helsinki on December 21")
	}
	if result.SunriseLocal == "N/A" {
		t.Error("sunrise local should not be N/A for Helsinki on December 21")
	}
	if result.SunsetLocal == "N/A" {
		t.Error("sunset local should not be N/A for Helsinki on December 21")
	}
	expectedTZ := "Europe/Helsinki"
	if result.Timezone != expectedTZ {
		t.Errorf("timezone = %q, want %q", result.Timezone, expectedTZ)
	}
	t.Logf("Helsinki winter (Dec 21, Europe/Helsinki): SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, TZ=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.Timezone)
}

func TestCalculateSunTimes_DefaultToUTC(t *testing.T) {
	lat := 51.5074
	lon := -0.1278
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTZ := "UTC"
	if result.Timezone != expectedTZ {
		t.Errorf("timezone = %q, want %q", result.Timezone, expectedTZ)
	}
}

func TestCalculateSunTimes_NYCDecember21(t *testing.T) {
	lat := 40.7128
	lon := -74.0060
	date := time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "America/New_York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SunriseUTC == "N/A" {
		t.Error("sunrise should not be N/A for NYC on December 21")
	}
	if result.SunsetUTC == "N/A" {
		t.Error("sunset should not be N/A for NYC on December 21")
	}
	t.Logf("NYC December 21: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength)
}

func TestCalculateSunTimes_Equator(t *testing.T) {
	lat := 0.0
	lon := 0.0
	date := time.Date(2024, time.March, 20, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("Equator March 20: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength)
	if result.DayLength == "0h 0m" {
		t.Error("day length should not be 0 at the equator")
	}
}

func TestCalculateSunTimes_PolarNight(t *testing.T) {
	lat := 78.0
	lon := 16.0
	date := time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Arctic/Longyearbyen")
	if err != nil {
		t.Logf("Expected error for invalid or polar timezone: %v", err)
	}
	t.Logf("Svalbard December 21: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength)
}

func TestCalculateSunTimes_Equinox(t *testing.T) {
	lat := 51.5074
	lon := -0.1278
	date := time.Date(2024, time.September, 22, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Europe/London")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("London September 22: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength)
}

func TestCalculateSunTimes_InvalidTimezone(t *testing.T) {
	lat := 51.5074
	lon := -0.1278
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	_, err := CalculateSunTimes(lat, lon, date, "Foo/Bar")
	if err == nil {
		t.Error("expected error for invalid timezone")
	}
}

func TestMinutesToTimeOfDay(t *testing.T) {
	tests := []struct {
		minutes  float64
		expected string
	}{
		{60, "01:00"},
		{720, "12:00"},
		{1080, "18:00"},
		{1380, "23:00"},
	}

	for _, tt := range tests {
		result := minutesToTimeOfDay(tt.minutes)
		if result != tt.expected {
			t.Errorf("minutesToTimeOfDay(%v) = %v, want %v", tt.minutes, result, tt.expected)
		}
	}
}

func TestFormatDayLength(t *testing.T) {
	tests := []struct {
		minutes  float64
		expected string
	}{
		{60, "1h 0m"},
		{123, "2h 3m"},
		{720, "12h 0m"},
		{-1, "0h 0m"},
	}

	for _, tt := range tests {
		result := formatDayLength(tt.minutes)
		if result != tt.expected {
			t.Errorf("formatDayLength(%v) = %v, want %v", tt.minutes, result, tt.expected)
		}
	}
}

func TestNormalizeMinutes(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{-60, 1380},
		{1500, 60},
		{720, 720},
	}

	for _, tt := range tests {
		result := normalizeMinutes(tt.input)
		if math.Abs(result-tt.expected) > 0.0001 {
			t.Errorf("normalizeMinutes(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2024, true},
		{2023, false},
		{2000, true},
		{1900, false},
	}

	for _, tt := range tests {
		result := isLeapYear(tt.year)
		if result != tt.expected {
			t.Errorf("isLeapYear(%v) = %v, want %v", tt.year, result, tt.expected)
		}
	}
}

func TestCalculateSunTimes_SouthernHemisphere(t *testing.T) {
	lat := -33.8688
	lon := 151.2093
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Australia/Sydney")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("Sydney June 21: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength)
	if result.SunriseUTC == "N/A" && result.SunsetUTC == "N/A" {
		t.Error("Sydney should have sun data in June")
	}
}
