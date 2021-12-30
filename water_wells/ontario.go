package water_wells

import (
	db "github.com/cleanflo/open_data/water_wells/db_layer"
	"github.com/icholy/utm"
)

var (
	OntarioR = db.Retreiver{
		DBName: "ontario",
		Table:  "qryWaterWellRecord",
		Data:   &[]Ontario{},
		Joins: []db.JoinStatement{
			{
				Diff:    db.INNER_JOIN,
				Left:    db.JoinTable{"qryWaterWellRecord", "WELL_ID"},
				Right:   db.JoinTable{"gryVBUTM", "Well_ID"},
				OrderBy: "northing",
				Select:  "northing, easting, ZONE, code",
			},
		},
		Query: db.ClientRequestData{
			Completed: db.TimeOption{
				Column: "Received",
				Layout: "2006-01-02 3:04:00.000 PM", // 1963-10-22 12:00:00.000 AM
			},
			Abandoned: db.TimeOption{
				Column: "Received",
				Layout: "2006-01-02 3:04:00.000 PM", //2002-01-15 12:00:00.000 AM
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"qryWaterWellRecord", "WELL_ID"},
						Right:   db.JoinTable{"qryAbandoned", "Well_ID"},
						GroupBy: "qryWaterWellRecord.Well_ID",
						OrderBy: "northing",
						Select:  "northing, easting, ZONE, code",
					},
				},
			},
			/*
				qryWaterWellRecord.Final_Status		COUNT
				(null)								46680
				Abandoned Monitoring and Test Hole	751
				Abandoned-Other						47453
				Abandoned-Quality					6889
				Abandoned-Supply					27759
				Alteration							1502
				Dewatering							2281
				Monitoring and Test Hole			23810
				Not A Well							518
				Observation Wells					59956
				Other Status						1139
				Recharge Well						1205
				Replacement Well					1058
				Test Hole							30197
				Unfinished							1710
				Water Supply						608951
			*/
			Status: db.NumberListOption{
				Column: "Final_Status",
				Items: map[string][]interface{}{
					"supply":     {"NULL", "Water Supply"},
					"research":   {"Test Hole", "Monitoring and Test Hole", "Observation Wells"},
					"geothermal": {},
					"abandoned":  {"Abandoned Monitoring and Test Hole", "Abandoned-Other", "Abandoned-Quality", "Abandoned-Supply"},
					"other":      {"Alteration", "Dewatering", "Not A Well", "Other Status", "Recharge Well", "Replacement Well"},
					"unknown":    {"Unfinished"},
				},
				Multiple: true,
			},
			/*
				qryWaterWellRecord.Use1		COUNT
				(null)						110472
				Commerical					12577
				Cooling And A/C				810
				Dewatering					1500
				Domestic					534050
				Industrial					3771
				Irrigation					5643
				Livestock					46283
				Monitoring					41285
				Monitoring and Test Hole	32778
				Municipal					4891
				Not Used					33345
				Other						2207
				Public						11151
				Test Hole					21096
			*/
			/*
				qryWaterWellRecord.Use2		COUNT
				(null)						802327
				Commerical					1230
				Cooling And A/C				489
				Dewatering					166
				Domestic					33518
				Industrial					303
				Irrigation					801
				Livestock					8902
				Monitoring					11487
				Monitoring and Test Hole	4
				Municipal					389
				Not Used					918
				Other						405
				Public						490
				Test Hole					430
			*/
			Use: db.NumberListOption{
				Column: "Use1",
				Items: map[string][]interface{}{
					"domestic":    {"NULL", "Domestic"},
					"commercial":  {"Commerical", "Cooling And A/C", "Dewatering"},
					"industial":   {"Industrial"},
					"municipal":   {"Municipal", "Public"},
					"irrigation":  {"Irrigation"},
					"agriculture": {"Livestock"},
					"research":    {"Monitoring", "Monitoring and Test Hole", "Test Hole"},
					"other":       {"Other", "Not Used"},
					"unknown":     {},
				},
				Multiple: true,
			},
			Colour: db.NumberListOption{},
			Taste:  db.NumberListOption{},
			Odour:  db.NumberListOption{},

			/*
				Rate
				!NULL: Count(270424)
				Min: 0,
				Max: 1000
			*/
			Rate: db.NumberOption{
				Column: "tblPump_Test.Recom_rate",
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"qryWaterWellRecord", "WELL_ID"},
						Right:   db.JoinTable{"tblPump_Test", "Well_ID"},
						GroupBy: "qryWaterWellRecord.Well_ID",
						OrderBy: "northing",
						Select:  "northing, easting, ZONE, code",
					},
				},
			},

			/*
				qryWellDepth.Well_Depth_m
				!NULL: Count(437450)
				Min: 0,
				Max: 9999
			*/
			Depth: db.NumberOption{
				Column: "qryWellDepth.Well_Depth_m",
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"qryWaterWellRecord", "WELL_ID"},
						Right:   db.JoinTable{"qryWellDepth", "Well_ID"},
						GroupBy: "qryWaterWellRecord.Well_ID",
						OrderBy: "northing",
						Select:  "northing, easting, ZONE, code",
					},
				},
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
			if v, ok := data.(*[]Ontario); ok {
				for _, a := range *v {
					zone := utm.Zone{
						Number: a.Zone,
						Letter: []rune(a.Code)[0],
						North:  true,
					}
					latitude, longitude := zone.ToLatLon(float64(a.Easting), float64(a.Northing))
					v := &db.Coordinates{}
					v.Add(float32(latitude), float32(longitude))
					c = append(c, *v)
				}
			}
			return c
		},
	}
)

type Ontario struct {
	Northing float32 `gorm:"column:northing"`
	Easting  float32 `gorm:"column:easting"`
	Zone     int     `gorm:"column:ZONE"`
	Code     string  `gorm:"column:code"`
}
