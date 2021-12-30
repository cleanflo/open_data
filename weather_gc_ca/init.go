package weather_gc_ca

import (
	"embed"
	"encoding/json"
	"fmt"
)

var (
	//go:embed "station-inventory.json"
	stationInventoryFS embed.FS
	StationInventory   = RawStations{}
)

func init() {
	err := StationInventory.loadData()
	if err != nil {
		panic(fmt.Errorf("failed init: %s", err))
	}
}

func (r *RawStations) loadData() error {
	file := "station-inventory.json"
	b, err := stationInventoryFS.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to load data: %s", err)
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return fmt.Errorf("failed to load data: %s", err)
	}

	return nil
}
