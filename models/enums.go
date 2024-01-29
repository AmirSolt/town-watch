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

const EnumSchema string = `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'region') THEN
		CREATE TYPE region AS ENUM ('TORONTO');
    END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'crime_type') THEN
		CREATE TYPE crime_type AS ENUM ('Assault', 'Auto Theft', 'Theft From Vehicle', 'Break And Enter', 'Sexual Violation', 'Robbery', 'Theft Over', 'Bike Theft', 'Shooting', 'Homicide');
    END IF;
END$$;

`
