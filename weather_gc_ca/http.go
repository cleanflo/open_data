package weather_gc_ca

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// SearchHandler processes a standard search request and returns a JSON response
// corresponding to []StationMetadata
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	latS := q.Get("lat")
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse lat: %s", err.Error()), http.StatusBadRequest)
		return
	}

	lngS := q.Get("lng")
	lng, err := strconv.ParseFloat(lngS, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse lng: %s", err.Error()), http.StatusBadRequest)
		return
	}

	maxS := q.Get("max")
	max, err := strconv.ParseInt(maxS, 10, 16)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse max: %s", err.Error()), http.StatusBadRequest)
		return
	}

	interval := Daily
	intervalS := q.Get("interval")
	if intervalS != "" {
		intParsed, err := strconv.ParseInt(intervalS, 10, 16)
		if err != nil {
			fmt.Println("failed to parse interval:", err)
		}
		interval = Interval(intParsed)
	}

	s := StationInventory.FindWithInterval(lat, lng, int(max), interval)
	if s == nil || len(s) == 0 {
		http.Error(w, "No stations found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

}
