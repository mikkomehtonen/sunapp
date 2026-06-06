package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/user/sunapp/backend/internal/sun"
	"github.com/user/sunapp/backend/internal/web"
)

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func sunHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	latStr := query.Get("lat")
	lonStr := query.Get("lon")
	dateStr := query.Get("date")
	tz := query.Get("tz")

	if latStr == "" || lonStr == "" {
		writeJSONError(w, "lat and lon query parameters are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		writeJSONError(w, "Invalid lat parameter", http.StatusBadRequest)
		return
	}
	if lat < -90 || lat > 90 {
		writeJSONError(w, "lat must be between -90 and 90", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		writeJSONError(w, "Invalid lon parameter", http.StatusBadRequest)
		return
	}
	if lon < -180 || lon > 180 {
		writeJSONError(w, "lon must be between -180 and 180", http.StatusBadRequest)
		return
	}

	var date time.Time
	if dateStr == "" {
		date = time.Now().UTC()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			writeJSONError(w, "Invalid date parameter. Use YYYY-MM-DD format", http.StatusBadRequest)
			return
		}
	}

	result, err := sun.CalculateSunTimes(lat, lon, date, tz)
	if err != nil {
		writeJSONError(w, fmt.Sprintf("Error calculating sunrise/sunset: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/sun", sunHandler)
	mux.Handle("/", web.NewHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on http://localhost:%s", port)
		log.Printf("API endpoint: http://localhost:%s/api/sun?lat=51.5074&lon=-0.1278&date=2024-06-21", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
