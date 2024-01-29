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

const ReportSchema string = `
CREATE TABLE report (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    occur_at TIMESTAMP(3) NOT NULL,
    external_src_id TEXT NOT NULL,
    neighborhood TEXT,
    location_type TEXT,
    crime_type "CrimeType" NOT NULL,
    region "Region" NOT NULL,
    point geometry(Point, 3857) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,
);

CREATE INDEX report_occ_at_idx ON report ("occur_at");
CREATE INDEX report_point_idx ON report USING GIST ("point");
`
