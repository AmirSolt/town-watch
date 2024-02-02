package server

import (
	"context"
	"fmt"

	"github.com/AmirSolt/town-watch/models"
)

const scanReportsLimit = 50

func (server *Server) Scan(scanParams models.ScanReportsParams) (*[]models.Report, error) {
	// query db
	reports, err := server.DB.queries.ScanReports(context.Background(), scanParams)
	if err != nil {
		return nil, fmt.Errorf("error ScanReports: %w", err)
	}

	return &reports.(*[]models.Report), nil
}
