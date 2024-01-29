package models

type Region string
type CrimeType string

const (
	TORONTO Region = "TORONTO"
)

const (
	ASSAULT            CrimeType = "Assault"
	AUTO_THEFT         CrimeType = "Auto Theft"
	THEFT_FROM_VEHICLE CrimeType = "Theft From Vehicle"
	BREAK_AND_ENTER    CrimeType = "Break And Enter"
	SEXUAL_VIOLATION   CrimeType = "Sexual Violation"
	ROBBERY            CrimeType = "Robbery"
	THEFT_OVER         CrimeType = "Theft Over"
	BIKE_THEFT         CrimeType = "Bike Theft"
	SHOOTING           CrimeType = "Shooting"
	HOMICIDE           CrimeType = "Homicide"
)
