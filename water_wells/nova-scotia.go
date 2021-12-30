package water_wells

import (
	db "github.com/cleanflo/open_data/water_wells/db_layer"
	"github.com/icholy/utm"
)

var (
	NovaScotiaR = db.Retreiver{
		DBName: "nova-scotia",
		Table:  "tblWellLogs",
		Data:   &[]NovaScotia{},
		Query: db.ClientRequestData{
			Completed: db.TimeOption{
				Column: "DateWellCompleted",
				Layout: "2006-01-02 15:04:00.000", //2000-08-17 00:00:00.000
			},
			Abandoned: db.TimeOption{
				Column: "DateWellCompleted",
				Layout: "2006-01-02 15:04:00.000", //2000-08-17 00:00:00.000
				Required: map[string]interface{}{
					"FinalStatusOfWellL": []int{5, 6, 7, 9, 16, 27},
				},
			},
			/*
				   FinalStatusOfWellL	FinalStatusOfWell						Count
				   NULL 				__										33037
					1					Water Supply Well						86629
					2					Observation Well						303
					3					Test Hole								382
					4					GEOTHERMAL, OPEN LOOP, RECHARGE WELL	110
					5					Abandoned, Dry							692
					6					Abandoned, Poor Quality					128
					7					Abandoned, Salt Water					180
					8					Unfinished								74
					9					ABANDONED, OTHER						48
					10					DRILL NEW & STIMULATE WELL				56
					13					Unknown									67
					16					Dry Hole								100
					19					Well Repaired/Reconstructed/Rehabed		301
					23					Other									48
					27					Abandoned/Decommissioned				300
					28					WELL STIMULATION						384
					31					DRILL/DIG NEW & ABANDON OLD				45
					35					Deepened								1513
					38					GEOTHERMAL, general						238
					39					GEOTHERMAL, OPEN LOOP, SUPPLY WELL		15
					40					GEOTHERMAL, CLOSED LOOP					133
			*/
			Status: db.NumberListOption{
				Column: "FinalStatusOfWellL",
				Items: map[string][]interface{}{
					"supply":     {1},
					"research":   {2, 3, 18},
					"geothermal": {4, 38, 39, 40},
					"abandoned":  {5, 6, 7, 8, 9, 16, 27, 36},
					"other":      {10, 19, 23, 28, 31, 35},
					"unknown":    {13},
				},
				Multiple: true,
			},
			/*
				   WaterUseL	WaterUse					Count
				   NULL 		__							11507
					1			Domestic					108695
					2			Industrial					775
					3			Commercial					819
					4			Municipal					372
					5			Other						349
					6			Irrigation					129
					7			Public (not municipal)		793
					8			Agricultural (not irriga)	474
					9			Heat Pump (source or dis)	473
					11			Dewatering					2
					12			Standby						5
					13			Domestic & Irrigation		6
					14			Domestic & Industrial		8
					17			Monitoring					77
					18			Unknown						118
					19			Heat Transfer				63
					22			Observation					13
					26			Aquaculture					5
					29			Domestic & Heat Pump		98
			*/
			Use: db.NumberListOption{
				Column: "WaterUseL",
				Items: map[string][]interface{}{
					"domestic":    {1, 29},
					"commercial":  {3},
					"industial":   {2},
					"municipal":   {4},
					"irrigation":  {6},
					"agriculture": {8, 26},
					"research":    {17, 22},
					"other":       {5, 7, 9, 10, 12, 13, 14, 19},
					"unknown":     {0, 18, 28},
				},
				Multiple: true,
			},
			/*
				   wqColourL	wqColour	Count
				   NULL 		__ 			77540
					1			Clear		34764
					2			Slight		229
					3			Coloured	70
					4			Cloudy		1625
					5			Turbid		3
					6			Yellow		19
					7			Brown		1215
					8			Tea			7
					10			Other		204
					11			Unknown		5660
					12			Red			518
					13			Gray		275
					14			None		2640
					15			Clearing	11
					16			REDDISH		3
			*/
			Colour: db.NumberListOption{
				Column: "wqColourL",
				Items: map[string][]interface{}{
					"clear":   {1, 2, 14},
					"cloudy":  {4, 5, 15},
					"light":   {3, 6, 8},
					"dark":    {7, 12, 13, 16},
					"other":   {10},
					"unknown": {9, 11},
				},
				Multiple: true,
			},
			/*
				   wqTasteL	wqTaste		Count
				   NULL		__			117228
					1		SULFUR		3
					2		METALLIC	7
					3		SALTY		47
					4		HARD		55
					5		GOOD		2094
					6		MINERAL		2
					9		OTHER		103
					10		UNKNOWN		6
					11		FRESH		430
					12		NONE		4808
			*/
			Taste: db.NumberListOption{
				Column: "wqTasteL",
				Items: map[string][]interface{}{
					"fresh":   {5, 7, 11, 12},
					"mineral": {2, 4, 6},
					"sulfur":  {1},
					"salt":    {3},
					"other":   {9},
					"unknown": {8, 10},
				},
				Multiple: true,
			},
			/*
				wqOdourL	wqOdour		Count
				NULL		__			92530
				1			ROTTEN EGG	18
				2			SULFUR		32
				3			NONE		25025
				4			METALLIC	2
				5			OILY		2
				9			OTHER		13
				10			UNKNOWN		7027
				11			GAS			73
				12			SLIGHT		38
				13			YES			23
			*/
			Odour: db.NumberListOption{
				Column: "wqOdourL",
				Items: map[string][]interface{}{
					"fresh":   {3},
					"mineral": {4, 12, 13},
					"sulfur":  {1, 2},
					"organic": {5, 6, 7, 11},
					"other":   {9},
					"unknown": {8, 10},
				},
				Multiple: true,
			},
			/*
				Rate
				NULL: Count(6991)
				Min: 0,
				Max: 3000
			*/
			Rate: db.NumberOption{
				Column: "wyRate",
			},
			/*
				Depth
				NULL: Count(854)
				Min: 2,
				Max: 1220
			*/
			Depth: db.NumberOption{
				Column: "TotalOrFinishedDepth",
			},

			/*
				Bedrock
				NULL: Count(26846)
				Min: 0,
				Max: 3000
			*/
			Bedrock: db.NumberOption{
				Column: "DepthToBedrock",
			},
		},
		ResponseChunks: map[int]int{
			10000:  1,
			50000:  2,
			150000: 4,
			500000: 5,
		},
		Coordinates: func(data interface{}) (c []db.Coordinates) {
			if v, ok := data.(*[]NovaScotia); ok {
				for _, a := range *v {
					zone := utm.Zone{
						Number: 20,
						Letter: 'T',
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

type NovaScotia struct {
	// WellNumber string `gorm:"column:WellNumber;primaryKey"`
	// DillerCompany string  `gorm:"column:DrillerCompany"`
	// Completed string `gorm:"column:DateWellCompleted"`
	// Owner         string  `gorm:"column:WellDrilledForLast"`
	// Community string `gorm:"column:NearestCommunity"`
	// Address   string `gorm:"column:CivicAddress"`
	// WaterUse  string `gorm:"column:WaterUseL"`
	// Colour        string  `gorm:"column:wqColourL"`
	// Taste         string  `gorm:"column:wqTasteL"`
	// Odour         string  `gorm:"column:wqOdourL"`
	// Other         string  `gorm:"column:wqOtherL"`
	// Rate          float64 `gorm:"column:wyRate"`
	// Drawdown      float64 `gorm:"column:wyTotalDrawdown"`
	// Recovered     float64 `gorm:"column:wyRecoveredTo"`
	// RecoveryHours float64 `gorm:"column:wyRecoveryByHours"`
	// RecoveryMins  float64 `gorm:"column:wyRecoveryByMin"`
	// BedrockDepth float64 `gorm:"column:DepthToBedrock"`
	// FinalDepth   float64 `gorm:"column:TotalOrFinishedDepth"`
	Easting  int `gorm:"column:Easting"`
	Northing int `gorm:"column:Northing"`
}
