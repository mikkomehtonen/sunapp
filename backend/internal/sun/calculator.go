package sun

import (
	"fmt"
	"math"
	"time"
)

const (
	deg2rad = math.Pi / 180.0
	rad2deg = 180.0 / math.Pi
)

// Result holds sunrise/sunset calculation results.
type Result struct {
	Sunrise   string `json:"sunrise"`
	Sunset    string `json:"sunset"`
	DayLength string `json:"day_length"`
}

// CalculateSunTimes computes sunrise and sunset times for a given
// latitude, longitude, and date (expected UTC).
func CalculateSunTimes(lat, lon float64, date time.Time) (*Result, error) {
	year := date.Year()
	dayOfYear := date.YearDay()
	yearLen := 365.0
	if isLeapYear(year) {
		yearLen = 366.0
	}

	// Fractional year in radians
	b := 2 * math.Pi * float64(dayOfYear-1) / yearLen

	// Equation of time in minutes
	eqtime := 229.18 * (0.000075 + 0.001868*math.Cos(b) - 0.032077*math.Sin(b) -
		0.014615*math.Cos(2*b) - 0.040849*math.Sin(2*b))

	// Solar declination in radians
	decl := 0.006918 - 0.399912*math.Cos(b) + 0.070257*math.Sin(b) -
		0.006758*math.Cos(2*b) + 0.000907*math.Sin(2*b) -
		0.002697*math.Cos(3*b) + 0.00148*math.Sin(3*b)

	// Hour angle calculation
	latRad := lat * deg2rad
	zenith := 90.833 * deg2rad

	cosHA := math.Cos(zenith)/(math.Cos(latRad)*math.Cos(decl)) - math.Tan(latRad)*math.Tan(decl)

	if cosHA < -1 || cosHA > 1 {
		return &Result{
			Sunrise:   "N/A",
			Sunset:    "N/A",
			DayLength: "N/A (polar day/night)",
		}, nil
	}

	ha := math.Acos(cosHA) * rad2deg // Convert to degrees

	// Calculate sunrise and sunset times in UTC minutes from midnight
	sunriseUTC := 720 - 4*(lon+ha) - eqtime
	sunsetUTC := 720 - 4*(lon-ha) - eqtime

	// Calculate timezone offset in minutes based on longitude
	tzOffset := int(math.Round(lon/15.0)) * 60

	// Convert to local time
	sunriseLocal := sunriseUTC + float64(tzOffset)
	sunsetLocal := sunsetUTC + float64(tzOffset)

	// Normalize to 0-1440 range
	sunriseLocal = normalizeMinutes(sunriseLocal)
	sunsetLocal = normalizeMinutes(sunsetLocal)

	return &Result{
		Sunrise:   minutesToTimeOfDay(sunriseLocal),
		Sunset:    minutesToTimeOfDay(sunsetLocal),
		DayLength: formatDayLength(sunsetLocal - sunriseLocal),
	}, nil
}

func normalizeMinutes(m float64) float64 {
	for m < 0 {
		m += 1440
	}
	for m >= 1440 {
		m -= 1440
	}
	return m
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func minutesToTimeOfDay(minutes float64) string {
	totalMinutes := int(math.Round(minutes))
	hours := totalMinutes / 60
	mins := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hours, mins)
}

func formatDayLength(minutes float64) string {
	if minutes < 0 {
		return "0h 0m"
	}
	hours := int(minutes / 60)
	mins := int(minutes) % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}

// FormatTime formats a time for display as "HH:MM".
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("15:04")
}

// FormatTimeDetailed formats a time for display as "HH:MM:SS".
func FormatTimeDetailed(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("15:04:05")
}
