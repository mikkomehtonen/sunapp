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
	SunriseUTC   string `json:"sunrise_utc"`
	SunsetUTC    string `json:"sunset_utc"`
	SunriseLocal string `json:"sunrise_local"`
	SunsetLocal  string `json:"sunset_local"`
	DayLength    string `json:"day_length"`
	Timezone     string `json:"timezone"`
}

// CalculateSunTimes computes sunrise and sunset times for a given
// latitude, longitude, date (UTC), and IANA timezone string.
// If tz is empty, it defaults to UTC.
func CalculateSunTimes(lat, lon float64, date time.Time, tz string) (*Result, error) {
	loc, err := resolveLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %q: %w", tz, err)
	}

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	year := date.Year()
	dayOfYear := date.YearDay()
	yearLen := 365.0
	if isLeapYear(year) {
		yearLen = 366.0
	}

	b := 2 * math.Pi * float64(dayOfYear-1) / yearLen

	eqtime := 229.18 * (0.000075 + 0.001868*math.Cos(b) - 0.032077*math.Sin(b) -
		0.014615*math.Cos(2*b) - 0.040849*math.Sin(2*b))

	decl := 0.006918 - 0.399912*math.Cos(b) + 0.070257*math.Sin(b) -
		0.006758*math.Cos(2*b) + 0.000907*math.Sin(2*b) -
		0.002697*math.Cos(3*b) + 0.00148*math.Sin(3*b)

	latRad := lat * deg2rad
	zenith := 90.833 * deg2rad

	cosHA := math.Cos(zenith)/(math.Cos(latRad)*math.Cos(decl)) - math.Tan(latRad)*math.Tan(decl)

	polar := cosHA < -1 || cosHA > 1
	if polar {
		return &Result{
			SunriseUTC:   "N/A",
			SunsetUTC:    "N/A",
			SunriseLocal: "N/A",
			SunsetLocal:  "N/A",
			DayLength:    "N/A (polar day/night)",
			Timezone:     loc.String(),
		}, nil
	}

	ha := math.Acos(cosHA) * rad2deg

	sunriseUTCMin := 720 - 4*(lon+ha) - eqtime
	sunsetUTCMin := 720 - 4*(lon-ha) - eqtime

	sunriseUTC := date.Add(time.Duration(math.Round(sunriseUTCMin)) * time.Minute).UTC()
	sunsetUTC := date.Add(time.Duration(math.Round(sunsetUTCMin)) * time.Minute).UTC()

	sunriseLocal := sunriseUTC.In(loc)
	sunsetLocal := sunsetUTC.In(loc)

	tzName := loc.String()
	if tzName == "" {
		tzName = "UTC"
	}

	return &Result{
		SunriseUTC:   FormatTime(sunriseUTC),
		SunsetUTC:    FormatTime(sunsetUTC),
		SunriseLocal: FormatTime(sunriseLocal),
		SunsetLocal:  FormatTime(sunsetLocal),
		DayLength:    formatDayLengthFromTimes(sunriseLocal, sunsetLocal),
		Timezone:     tzName,
	}, nil
}

func resolveLocation(tz string) (*time.Location, error) {
	if tz == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(tz)
}

func formatDayLengthFromTimes(start, end time.Time) string {
	dur := end.Sub(start)
	if dur < 0 {
		dur += 24 * time.Hour
	}
	totalMinutes := int(dur.Minutes())
	hours := totalMinutes / 60
	mins := totalMinutes % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// FormatTime formats a time for display as "HH:MM".
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("15:04")
}
