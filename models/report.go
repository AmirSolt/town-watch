package models

import "time"

type Region string
type CrimeType string

const (
	TORONTO Region = "TORONTO"
)
const (
	ASSAULT            = "Assault"
	AUTO_THEFT         = "Auto Theft"
	THEFT_FROM_VEHICLE = "Theft From Vehicle"
	BREAK_AND_ENTER    = "Break And Enter"
	SEXUAL_VIOLATION   = "Sexual Violation"
	ROBBERY            = "Robbery"
	THEFT_OVER         = "Theft Over"
	BIKE_THEFT         = "Bike Theft"
	SHOOTING           = "Shooting"
	HOMICIDE           = "Homicide"
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
