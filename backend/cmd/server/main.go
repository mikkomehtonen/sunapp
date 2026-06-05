package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/user/sunapp/backend/internal/sun"
)

func sunHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	latStr := query.Get("lat")
	lonStr := query.Get("lon")
	dateStr := query.Get("date")
	tz := query.Get("tz")

	if latStr == "" || lonStr == "" {
		http.Error(w, "lat and lon query parameters are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid lat parameter", http.StatusBadRequest)
		return
	}
	if lat < -90 || lat > 90 {
		http.Error(w, "lat must be between -90 and 90", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid lon parameter", http.StatusBadRequest)
		return
	}
	if lon < -180 || lon > 180 {
		http.Error(w, "lon must be between -180 and 180", http.StatusBadRequest)
		return
	}

	var date time.Time
	if dateStr == "" {
		date = time.Now().UTC()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date parameter. Use YYYY-MM-DD format", http.StatusBadRequest)
			return
		}
	}

	result, err := sun.CalculateSunTimes(lat, lon, date, tz)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating sunrise/sunset: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/sun", sunHandler)

	port := "8080"
	log.Printf("Starting server on http://localhost:%s", port)
	log.Printf("API endpoint: http://localhost:%s/api/sun?lat=51.5074&lon=-0.1278&date=2024-06-21", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
