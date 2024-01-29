package models

import "time"

type Report struct {
	Id           int8      `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	OccurAt      time.Time `db:"occur_at"`
	Neighborhood string    `db:"neighborhood"`
	LocationType string    `db:"location_type"`
	CrimeType    CrimeType `db:"crime_type"`
	Region       Region    `db:"region"`
	Lat          float32   `db:"lat"`
	Long         float32   `db:"long"`
}
