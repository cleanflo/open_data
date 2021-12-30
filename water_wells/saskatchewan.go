package water_wells

import (
	db "github.com/cleanflo/open_data/water_wells/db_layer"
)

var (
	SaskatchewanR = db.Retreiver{
		DBName: "saskatchewan",
		Table:  "tblWells",
		Data:   &[]Saskatchewan{},
		Query: db.ClientRequestData{
			Completed: db.TimeOption{
				Column: "completed",
				Layout: "2006.01.02", // 1973.10.28
			},
			Abandoned: db.TimeOption{
				Column: "date_decommisioned",
				Layout: "2006.01.02", // 1899.12.30
			},
			/*
				well_use			Count
									815
				Mineral Test Hole	175
				Observation			3110
				Quality Monitoring	311
				Recharge			11
				Seismic Test Hole	1130
				Soil Test Hole		806
				Waste Disposal		28
				Water Test Hole		40659
				Withdrawal			87871
			*/
			Status: db.NumberListOption{
				Column: "well_use",
				Items: map[string][]interface{}{
					"supply":     {nil, "Withdrawal"},
					"research":   {"Observation", "Quality Monitoring", "Seismic Test Hole", "Soil Test Hole", "Water Test Hole"},
					"geothermal": {},
					"abandoned":  {},
					"other":      {"Waste Disposal", "Recharge"},
					"unknown":    {},
				},
				Multiple: true,
			},
			/*
				water_use			Count
									1586
				Domestic			112266
				Drainage			21
				Industrial			3849
				Irrigation			319
				Mineral Recovery	141
				Mineral Water		3
				Multi-purpose		321
				Municipal			10977
				Other				45
				Recreation			193
				Research			5195
			*/
			Use: db.NumberListOption{
				Column: "water_use",
				Items: map[string][]interface{}{
					"domestic":    {nil, "Domestic"},
					"commercial":  {"Multi-purpose"},
					"industial":   {"Industrial", "Mineral Recovery", "Mineral Water"},
					"municipal":   {"Municipal", "Recreation"},
					"irrigation":  {"Irrigation", "Drainage"},
					"agriculture": {},
					"research":    {"Research"},
					"other":       {"Other"},
					"unknown":     {},
				},
				Multiple: true,
			},
			Colour: db.NumberListOption{},
			Taste:  db.NumberListOption{},
			Odour:  db.NumberListOption{},
			/*
				Rate
				NULL: Count(134916)
				Min: 0,
				Max: 3046
			*/
			Rate: db.NumberOption{
				Column: "recommended_pumping_rate",
			},
			/*
				Depth
				NULL: Count(134916)
				Min: 0,
				Max: 7452
			*/
			Depth: db.NumberOption{
				Column: "TotalOrFinishedDepth",
			},
			Bedrock: db.NumberOption{},
		},
		ResponseChunks: map[int]int{
			10000:  1,
			50000:  2,
			150000: 4,
			500000: 5,
		},
		Coordinates: func(data interface{}) (c []db.Coordinates) {
			if v, ok := data.(*[]Saskatchewan); ok {
				for _, a := range *v {
					c = append(c, db.Coordinates{a.Latitude, a.Longitude})
				}
			}
			return c
		},
	}
)

type Saskatchewan struct {
	Latitude  float32 `gorm:"column:latitude"`
	Longitude float32 `gorm:"column:longitude"`
}
