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

	return &[]models.Report{}
}
