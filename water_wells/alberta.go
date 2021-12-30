package water_wells

import db "github.com/cleanflo/open_data/water_wells/db_layer"

var (
	AlbertaR = db.Retreiver{
		DBName: "alberta",
		Table:  "Wells",
		Data:   &[]Alberta{},
		Query: db.ClientRequestData{
			Completed: db.TimeOption{
				Column: "Well_Reports.Drilling_End_Date",
				Layout: "2006-01-02 3:04:00.000 PM", //2002-05-16 12:00:00.000 AM
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "Well_ID"},
						Right:   db.JoinTable{"Well_Reports", "Well_ID"},
						GroupBy: "Wells.Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
					},
				},
			},
			Abandoned: db.TimeOption{
				Column: "Well_Reports.Plug_Date",
				Layout: "2006-01-02 3:04:00.000 PM", //2002-01-15 12:00:00.000 AM
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "Well_ID"},
						Right:   db.JoinTable{"Well_Reports", "Well_ID"},
						GroupBy: "Wells.Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
					},
				},
			},
			/*
				Well_Reports.Type_of_Work		COUNT
				(null)							397
												2
				Cathodic Protection				142
				Chemistry						78488
				Coal Test Hole					1916
				Core Hole						3982
				Deepened						3289
				Drill Stem Test Hole			841
				Dry Hole						4555
				Dry Hole-Decommissioned			4614
				Existing Well-Decommissioned	9200
				Federal Well Survey				20341
				Flowing Shot Hole				29905
				New Well						222475
				New Well-Decommissioned			7180
				Oil Exploratory					1820
				Old Well-Yield					1759
				Other							359
				Piezometer						1212
				Reconditioned					1800
				Seismic Shot Hole				719
				Spring							3122
				Structure Test Hole				7388
				Test Hole						10678
				Test Hole-Decommissioned		12017
				Unknown							1688
				Well Inventory					14023
			*/
			Status: db.NumberListOption{
				Column: "Well_Reports.Type_of_Work",
				Items: map[string][]interface{}{
					"supply":     {"New Well", "Deepened", "Reconditioned", "Spring"},
					"research":   {"Test Hole", "Coal Test Hole", "Core Hole", "Federal Well Survey", "Chemistry", "Structure Test Hole", "Well Inventory", "Drill Stem Test Hole", "Piezometer", "Seismic Shot Hole"},
					"geothermal": {},
					"abandoned":  {"Dry Hole", "Old Well-Yeild", "Dry Hole-Decommissioned", "New Well-Decommissioned", "Test Hole-Decommissioned", "Existing Well-Decommissioned"},
					"other":      {"Other", "Flowing Shot Hole", "Oil Exploratory", "Cathodic Protection"},
					"unknown":    {"Unknown"},
				},
				Multiple: true,
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "Well_ID"},
						Right:   db.JoinTable{"Well_Reports", "Well_ID"},
						GroupBy: "Wells.Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
					},
				},
			},
			/*
				Well_Reports.Well_Use			COUNT
				(null)							1234
												1
				Aggregate Washing				2
				Chemistry						3
				Commercial						448
				Contamination Invest.			104
				Co-ops (Colonies)				119
				Dewatering						834
				Domestic						220989
				Domestic & Industrial			478
				Domestic & Irrigation			327
				Domestic & Stock				59312
				Dry Hole - Abandoned			4
				Geothermal						28
				Golf Courses					25
				Heat Transfer					9
				Hydrostatic Testing				3
				Industrial						65167
				Industrial & Stock				91
				Industrial Camp					212
				Injection						745
				Intensive Livestock Operation	36
				Investigation					7142
				Irrigation						490
				Monitoring						1955
				Municipal						5669
				Municipal & Industrial			37
				Municipal & Observation			36
				New Well						8
				Observation						6472
				Old Well - Abandoned			7
				Old Well - Test					4
				Other							3234
				Rural Subdivision				36
				Standby							103
				Stock							35368
				Test Hole						1
				Test Hole - Abandoned			2
				Test Hole-Abandoned				1
				Unknown							33173
				Water Hauling					2
				Well Inventory					1
			*/
			Use: db.NumberListOption{
				Column: "Well_Reports.Well_Use",
				Items: map[string][]interface{}{
					"domestic":    {"NULL", "Domestic", "New Well", "Standby", "Water Hauling"},
					"commercial":  {"Industial", "Dewatering", "Geothermal", "Heat Transfer"},
					"industial":   {"Industial", "Domestic & Industrial", "Industrial Camp", "Injection"},
					"municipal":   {"Municipal", "Co-ops (Colonies)", "Municipal & Industrial", "Rural Subdivision"},
					"irrigation":  {"Irrigation", "Domestic & Irrigation", "Golf Courses"},
					"agriculture": {"Stock", "Domestic & Stock", "Industrial & Stock", "Intensive Livestock Operation"},
					"research":    {"Observation", "Contamination Invest.", "Hydrostatic Testing", "Investigation", "Monitoring"},
					"other":       {"Other", "Dry Hole - Abandoned", "Old Well - Abandoned", "Test Hole - Abandoned", "Test Hole-Abandoned"},
					"unknown":     {"Unknown"},
				},
				Multiple: true,
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "Well_ID"},
						Right:   db.JoinTable{"Well_Reports", "Well_ID"},
						GroupBy: "Wells.Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
					},
				},
			},
			/*
				Lithologies.Colour	COUNT
				(null)				1026451
									74157
				Black				25903
				Blue				129277
				Blue Gray			11910
				Bluish Green		1166
				Brown				279248
				Brownish Gray		26531
				Brownish Green		1858
				Brownish Yellow		3273
				Dark				13070
				Dark Blue			1731
				Dark Brown			7190
				Dark Gray			35271
				Dark Green			1221
				Dark Red			72
				Dark Yellow			114
				Gray				675733
				Gray Salt & Pepper	2373
				Green				59555
				Greenish Gray		15262
				Greenish Yellow		169
				Light				4090
				Light Blue			1959
				Light Brown			10874
				Light Gray			29403
				Light Green			2088
				Light Red			117
				Light Yellow		211
				Red					2369
				Salt & Pepper		4284
				See Comments		1705
				Tan					3996
				Unknown				594
				Unreadable			77
				White				5515
				Yellow				32027
			*/
			Colour: db.NumberListOption{
				Column: "Lithologies.Colour",
				Items: map[string][]interface{}{
					"clear":   {"NULL", "Gray", "White"},
					"cloudy":  {"Yellow", "Brown", "Green", "Blue", "Tan", "Red", "Gray Salt & Pepper", "Salt & Pepper"},
					"light":   {"Light", "Light Gray", "Light Red", "Greenish Gray", "Light Yellow", "Light Green", "Light Blue", "Light Brown", "Blue Gray", "Greenish Yellow", "Greenish Gray"},
					"dark":    {"Dark", "Black", "Dark Gray", "Dark Brown", "Dark Red", "Dark Blue", "Dark Green", "Dark Yellow", "Bluish Green", "Brownish Gray", "Brownish Green", "Brownish Yellow"},
					"other":   {"See Comments"},
					"unknown": {"Unreadable", "Unknown"},
				},
				Multiple: true,
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "GIC_Well_ID"},
						Right:   db.JoinTable{"Lithologies", "GIC_Well_ID"},
						GroupBy: "Wells.GIC_Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
					},
				},
			},
			Taste: db.NumberListOption{},
			Odour: db.NumberListOption{},

			/*
				Rate
				!NULL: Count(270424)
				Min: 0,
				Max: 1000
			*/
			Rate: db.NumberOption{
				Column: "Well_Reports.Recommended_Rate",
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "Well_ID"},
						Right:   db.JoinTable{"Well_Reports", "Well_ID"},
						GroupBy: "Wells.Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
					},
				},
			},

			/*
				Depth
				!NULL: Count(437450)
				Min: 0,
				Max: 9999
			*/
			Depth: db.NumberOption{
				Column: "Well_Reports.Total_Depth_Drilled",
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"Wells", "Well_ID"},
						Right:   db.JoinTable{"Well_Reports", "Well_ID"},
						GroupBy: "Wells.Well_ID",
						OrderBy: "MAX(Latitude)",
						Select:  "MAX(Latitude) as Latitude, MAX(Longitude) as Longitude",
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
			if v, ok := data.(*[]Alberta); ok {
				for _, a := range *v {
					v := &db.Coordinates{}
					v.Add(float32(a.Latitude), float32(a.Longitude))
					c = append(c, *v)
				}
			}
			return c
		},
	}
)

type Alberta struct {
	Latitude  float32 `gorm:"column:Latitude"`
	Longitude float32 `gorm:"column:Longitude"`
}
