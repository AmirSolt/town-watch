package server

import (
	"time"

	"github.com/AmirSolt/town-watch/models"
)

type ScanParams struct {
	Lat      float64
	Long     float64
	Radius   float64
	Region   models.Region
	FromDate time.Time
	ToDate   time.Time
}

func (server *Server) Scan(scanParams ScanParams) *[]models.Report {
	// query db

	`
                SELECT *
                FROM "Report"
                WHERE 
                ST_DWithin(
                    point,
                    ST_Point(${lat}, ${long}, 3857),
                    ${radius}
                )
                AND region = ${region} 
                AND reported_at > CURRENT_DATE - INTERVAL '${afterDays} day'
                ${toDaysAgoQuery}
                ORDER BY
                    reported_at
                LIMIT ${REPORTS_LIMIT};
            `

	return &[]models.Report{}
}
