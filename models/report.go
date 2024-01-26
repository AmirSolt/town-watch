package models

import "time"

type Region string
type CrimeType string

const (
	TORONTO Region = "TORONTO"
)
const (
	ASSAULT            = "ASSAULT"
	AUTO_THEFT         = "AUTO_THEFT"
	THEFT_FROM_VEHICLE = "THEFT_FROM_VEHICLE"
	BREAK_AND_ENTER    = "BREAK_AND_ENTER"
	SEXUAL_VIOLATION   = "SEXUAL_VIOLATION"
	ROBBERY            = "ROBBERY"
	THEFT_OVER         = "THEFT_OVER"
	BIKE_THEFT         = "BIKE_THEFT"
	SHOOTING           = "SHOOTING"
	HOMICIDE           = "HOMICIDE"
)

type Report struct {
	OccurAt      *time.Time "db:created_at"
	Neighborhood string     "db:neighborhood"
	LocationType string     "db:location_type"
	CrimeType    CrimeType  "db:crime_type"
	Region       Region     "db:region"
	Lat          string     "db:lat"
	Long         string     "db:long"
}
