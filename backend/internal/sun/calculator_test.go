package sun

import (
	"fmt"
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
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
	}
	if result.SunriseMinutesLocal < 0 {
		t.Errorf("sunrise_minutes_local = %d, want >= 0", result.SunriseMinutesLocal)
	}
	if result.SunsetMinutesLocal < 0 {
		t.Errorf("sunset_minutes_local = %d, want >= 0", result.SunsetMinutesLocal)
	}
	if result.SunriseMinutesLocal >= result.SunsetMinutesLocal {
		t.Errorf("sunrise_minutes_local (%d) should be < sunset_minutes_local (%d)", result.SunriseMinutesLocal, result.SunsetMinutesLocal)
	}
	t.Logf("Helsinki summer (June 21, Europe/Helsinki): SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, TZ=%s, Pol=%s, RiseMin=%d, SetMin=%d",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.Timezone, result.PolarType, result.SunriseMinutesLocal, result.SunsetMinutesLocal)
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
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
	}
	if result.SunriseMinutesLocal < 0 {
		t.Errorf("sunrise_minutes_local = %d, want >= 0", result.SunriseMinutesLocal)
	}
	if result.SunsetMinutesLocal < 0 {
		t.Errorf("sunset_minutes_local = %d, want >= 0", result.SunsetMinutesLocal)
	}
	t.Logf("Helsinki winter (Dec 21, Europe/Helsinki): SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, TZ=%s, Pol=%s, RiseMin=%d, SetMin=%d",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.Timezone, result.PolarType, result.SunriseMinutesLocal, result.SunsetMinutesLocal)
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
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
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
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
	}
	t.Logf("NYC December 21: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, Pol=%s",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.PolarType)
}

func TestCalculateSunTimes_Equator(t *testing.T) {
	lat := 0.0
	lon := 0.0
	date := time.Date(2024, time.March, 20, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("Equator March 20: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, Pol=%s, RiseMin=%d, SetMin=%d",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.PolarType, result.SunriseMinutesLocal, result.SunsetMinutesLocal)
	if result.DayLength == "0h 0m" {
		t.Error("day length should not be 0 at the equator")
	}
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
	}
}

func TestCalculateSunTimes_PolarNight(t *testing.T) {
	lat := 78.0
	lon := 16.0
	date := time.Date(2024, time.December, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Arctic/Longyearbyen")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DayLength != "Polar night" {
		t.Errorf("day_length = %q, want %q", result.DayLength, "Polar night")
	}
	if result.PolarType != "polar_night" {
		t.Errorf("polar_type = %q, want %q", result.PolarType, "polar_night")
	}
	if result.SunriseMinutesLocal != -1 {
		t.Errorf("sunrise_minutes_local = %d, want -1", result.SunriseMinutesLocal)
	}
	if result.SunsetMinutesLocal != -1 {
		t.Errorf("sunset_minutes_local = %d, want -1", result.SunsetMinutesLocal)
	}
	if result.SunriseLocal != "N/A" {
		t.Errorf("sunrise_local = %q, want N/A", result.SunriseLocal)
	}
	if result.SunsetLocal != "N/A" {
		t.Errorf("sunset_local = %q, want N/A", result.SunsetLocal)
	}
	t.Logf("Svalbard December 21: DayLength=%s, Pol=%s", result.DayLength, result.PolarType)
}

func TestCalculateSunTimes_PolarDay(t *testing.T) {
	lat := 78.0
	lon := 16.0
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Arctic/Longyearbyen")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DayLength != "Midnight sun" {
		t.Errorf("day_length = %q, want %q", result.DayLength, "Midnight sun")
	}
	if result.PolarType != "midnight_sun" {
		t.Errorf("polar_type = %q, want %q", result.PolarType, "midnight_sun")
	}
	if result.SunriseMinutesLocal != -1 {
		t.Errorf("sunrise_minutes_local = %d, want -1", result.SunriseMinutesLocal)
	}
	if result.SunsetMinutesLocal != -1 {
		t.Errorf("sunset_minutes_local = %d, want -1", result.SunsetMinutesLocal)
	}
	if result.SunriseLocal != "N/A" {
		t.Errorf("sunrise_local = %q, want N/A", result.SunriseLocal)
	}
	if result.SunsetLocal != "N/A" {
		t.Errorf("sunset_local = %q, want N/A", result.SunsetLocal)
	}
	t.Logf("Svalbard June 21: DayLength=%s, Pol=%s", result.DayLength, result.PolarType)
}

func TestCalculateSunTimes_Equinox(t *testing.T) {
	lat := 51.5074
	lon := -0.1278
	date := time.Date(2024, time.September, 22, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Europe/London")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
	}
	t.Logf("London September 22: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, Pol=%s, RiseMin=%d, SetMin=%d",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.PolarType, result.SunriseMinutesLocal, result.SunsetMinutesLocal)
}

func TestCalculateSunTimes_NewFieldsMatchFormattedTimes(t *testing.T) {
	lat := 60.1699
	lon := 24.9384
	date := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	result, err := CalculateSunTimes(lat, lon, date, "Europe/Helsinki")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// sunrise_local is "HH:MM", verify sunrise_minutes_local matches
	var sunriseHour, sunriseMin, sunsetHour, sunsetMin int
	n, err := fmt.Sscanf(result.SunriseLocal, "%d:%d", &sunriseHour, &sunriseMin)
	if n != 2 {
		t.Fatalf("failed to parse sunrise_local %q: %v (got %d fields)", result.SunriseLocal, err, n)
	}
	expectedMinutes := sunriseHour*60 + sunriseMin
	if result.SunriseMinutesLocal != expectedMinutes {
		t.Errorf("sunrise_minutes_local = %d, want %d (from %s)", result.SunriseMinutesLocal, expectedMinutes, result.SunriseLocal)
	}

	n, err = fmt.Sscanf(result.SunsetLocal, "%d:%d", &sunsetHour, &sunsetMin)
	if n != 2 {
		t.Fatalf("failed to parse sunset_local %q: %v (got %d fields)", result.SunsetLocal, err, n)
	}
	expectedMinutes = sunsetHour*60 + sunsetMin
	if result.SunsetMinutesLocal != expectedMinutes {
		t.Errorf("sunset_minutes_local = %d, want %d (from %s)", result.SunsetMinutesLocal, expectedMinutes, result.SunsetLocal)
	}
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
	t.Logf("Sydney June 21: SunriseUTC=%s, SunsetUTC=%s, SunriseLocal=%s, SunsetLocal=%s, DayLength=%s, Pol=%s, RiseMin=%d, SetMin=%d",
		result.SunriseUTC, result.SunsetUTC, result.SunriseLocal, result.SunsetLocal, result.DayLength, result.PolarType, result.SunriseMinutesLocal, result.SunsetMinutesLocal)
	if result.SunriseUTC == "N/A" && result.SunsetUTC == "N/A" {
		t.Error("Sydney should have sun data in June")
	}
	if result.PolarType != "" {
		t.Errorf("polar_type = %q, want empty string", result.PolarType)
	}
}
