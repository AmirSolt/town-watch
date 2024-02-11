package server

import (
	"context"
	"fmt"

	"github.com/AmirSolt/town-watch/models"
)

const scanReportsLimit = 50

func (server *Server) Scan(scanParams models.ScanReportsParams) (*[]models.Report, error) {
	// query db
	reports, err := server.DB.Queries.ScanReports(context.Background(), scanParams)
	if err != nil {
		return nil, fmt.Errorf("error ScanReports: %w", err)
	}

	// Casting reports results
	reportCasted := make([]models.Report, len(reports))
	for i, n := range reports {
		report, ok := n.(models.Report)
		if !ok {
			return nil, fmt.Errorf("error reports type assertion failed")
		}
		reportCasted[i] = report
	}

	return &reportCasted, nil
}
