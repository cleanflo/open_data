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
	latS := r.FormValue("lat")
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse lat: %s", err.Error()), http.StatusBadRequest)
		return
	}

	lngS := r.FormValue("lng")
	lng, err := strconv.ParseFloat(lngS, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse lng: %s", err.Error()), http.StatusBadRequest)
		return
	}

	maxS := r.FormValue("max")
	max, err := strconv.ParseInt(maxS, 10, 16)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse max: %s", err.Error()), http.StatusBadRequest)
		return
	}

	s := StationInventory.Find(lat, lng, int(max))
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
