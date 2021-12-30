package water_wells

import db "github.com/cleanflo/open_data/water_wells/db_layer"

var (
	BritishColumbiaR = db.Retreiver{
		DBName: "british-columbia",
		Table:  "well",
		Data:   &[]BritishColumbia{},
		Query: db.ClientRequestData{
			Completed: db.TimeOption{
				Column: "construction_end_date",
				Layout: "2006-01-02 3:04:00.000 PM", //2002-05-16 12:00:00.000 AM
			},
			Abandoned: db.TimeOption{
				Column: "construction_end_date",
				Layout: "2006-01-02 3:04:00.000 PM", //2002-01-15 12:00:00.000 AM
				Required: map[string]interface{}{
					"well_status_code": "ABANDONED",
				},
			},
			/*
				well_status_code	COUNT
				''					9
				ABANDONED			1311
				ALTERATION			1363
				CLOSURE				1683
				NEW					116625
				OTHER				24
			*/
			Status: db.NumberListOption{
				Column: "well_status_code",
				Items: map[string][]interface{}{
					"supply":     {"NEW"},
					"research":   {},
					"geothermal": {},
					"abandoned":  {"ABANDONED", "CLOSURE"},
					"other":      {"ALTERATION", "OTHER"},
					"unknown":    {""},
				},
				Multiple: true,
			},
			/*
				intended_water_use_code		COUNT
				COM							3147
				DOM							63882
				DWS							3977
				IRR							2512
				NA							9658
				OBS							136
				OP_LP_GEO					30
				OTHER						1388
				TST							171
				UNK							36114
			*/
			Use: db.NumberListOption{
				Column: "intended_water_use_code",
				Items: map[string][]interface{}{
					"domestic":    {"DOM"},
					"commercial":  {"COM", "DWS"},
					"industial":   {},
					"municipal":   {},
					"irrigation":  {"IRR"},
					"agriculture": {},
					"research":    {"TST", "OBS", "OP_LP_GEO"},
					"other":       {"OTHER", "NA"},
					"unknown":     {"UNK"},
				},
				Multiple: true,
			},
			/*
				lithology.lithology_colour_code			COUNT
				(null)									496081
				0 nothing entered						5003
				black									5649
				blue									2869
				brown									44077
				dark									642
				green									4684
				grey									40275
				light									315
				purple									156
				red										930
				rust-coloured							136
				salt & pepper							287
				speckled								8
				tan										1760
				vari-coloured							4742
				white									924
				yellow									273
			*/
			Colour: db.NumberListOption{
				Column: "lithology.lithology_colour_code",
				Items: map[string][]interface{}{
					"clear":   {"NULL", "grey", "white"},
					"cloudy":  {"salt & pepper", "speckled", "tan"},
					"light":   {"light", "vari-coloured", "yellow"},
					"dark":    {"black", "blue", "brown", "dark", "green", "purple", "red", "rust-coloured"},
					"other":   {},
					"unknown": {"0 nothing entered"},
				},
				Multiple: true,
				Joins: []db.JoinStatement{
					{
						Diff:    db.INNER_JOIN,
						Left:    db.JoinTable{"well", "well_tag_number"},
						Right:   db.JoinTable{"lithology", "well_tag_number"},
						GroupBy: "well.well_tag_number",
						OrderBy: "MAX(latitude_Decdeg)",
						Select:  "MAX(latitude_Decdeg) as latitude_Decdeg, MAX(latitude_Decdeg) as latitude_Decdeg",
					},
				},
			},
			Taste: db.NumberListOption{},
			/*
				water_quality_odour				COUNT
				''								120746
				A LITTLE ODOUR					1
				brown cloudiness				1
				brown cloudy					1
				clear							3
				CLEAR / NONE					1
				COLOUR OR ODOUR: SLIGHT			1
				CRISP, CLEAN					1
				faint sulfur					1
				Fecal coliform count of 5.2		1
				grey							2
				grey cloudy						1
				HIGHER CONDUCTIVITY				1
				iron							10
				IRON 7; HARDNESS 114 GRAINS		1
				Iron, sulfur					2
				Irony							1
				LITTLE ODOUR					1
				Nil								2
				No								21
				NO IRON, NO MAGNESIUM			1
				no iron, no manganese			1
				NO ODOR							1
				NO ODOU							1
				No odour						5
				NO ORDER						1
				none							148
				orange / Iron					1
				orange/ light iron smell		1
				pump slow to start				1
				SLIGHT							7
				SLIGHT GAS						1
				Slight odor.					1
				SLIGHT ODOUR					5
				SLIGHT ODOUR FIRST DAY			1
				SLIGHT SMELL					1
				SLIGHT SULFUR					1
				SLIGHT SULPHUR					2
				SLIGHT SULPHUR SMELL			1
				slite							1
				SMELL							1
				SOME ODOUR						2
				SOME ODOUR, THEN FRESHENED UP	1
				SOME SULPHUR					1
				Sulfer							2
				sulfur							9
				SULFURIC						1
				SULPHUR							9
				SULPHUR SMELL					1
				Sulphur smell noted				3
				SULPHURY						1
				suphur							1
				SWEET							1
				SWEET ODOUR						1
			*/
			Odour: db.NumberListOption{},

			/*
				Rate
				!NULL: Count(109021)
				Min: 0,
				Max: 97920
			*/
			Rate: db.NumberOption{
				Column: "well_yield_usgpm",
			},

			/*
				Depth
				!NULL: Count(118756)
				Min: -30,
				Max: 55030
			*/
			Depth: db.NumberOption{
				Column: "finished_well_depth_ft-bgl",
			},

			/*
				Bedrock
				!NULL: Count(37796)
				Min: -3,
				Max: 7422
			*/
			Bedrock: db.NumberOption{
				Column: "bedrock_depth_ft-bgl",
			},
		},
		ResponseChunks: map[int]int{
			10000:  1,
			50000:  2,
			150000: 4,
			500000: 5,
		},
		Coordinates: func(data interface{}) (c []db.Coordinates) {
			if v, ok := data.(*[]BritishColumbia); ok {
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

type BritishColumbia struct {
	Latitude  float32 `gorm:"column:latitude_Decdeg"`
	Longitude float32 `gorm:"column:longitude_Decdeg"`
}
