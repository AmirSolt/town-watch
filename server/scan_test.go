package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AmirSolt/town-watch/models"
	"github.com/AmirSolt/town-watch/server"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockDB struct {
	mockScanReports func(ctx context.Context, params models.ScanReportsParams) ([]interface{}, error)
}

func (mdb *mockDB) ScanReports(ctx context.Context, params models.ScanReportsParams) ([]interface{}, error) {
	return mdb.mockScanReports(ctx, params)
}

func TestServer_Scan(t *testing.T) {
	// Example report to be used in tests
	exampleReport := models.Report{
		// Initialize necessary fields
	}

	tests := []struct {
		name         string
		scanParams   models.ScanReportsParams
		mockScanFunc func(ctx context.Context, params models.ScanReportsParams) ([]interface{}, error)
		want         *[]models.Report
		wantErr      bool
	}{
		{
			name: "Success scan",
			scanParams: models.ScanReportsParams{
				Lat:      40.712776,
				Long:     -74.005974,
				Radius:   10,
				FromDate: pgtype.Timestamptz{Time: time.Now().AddDate(0, -1, 0)},
				ToDate:   pgtype.Timestamptz{Time: time.Now()},
			},
			mockScanFunc: func(ctx context.Context, params models.ScanReportsParams) ([]interface{}, error) {
				// Simulate successful database retrieval
				return []interface{}{exampleReport}, nil
			},
			want: &[]models.Report{exampleReport},
		},
		{
			name: "Scan error",
			scanParams: models.ScanReportsParams{
				Lat:      40.712776,
				Long:     -74.005974,
				Radius:   10,
				FromDate: pgtype.Timestamptz{},
				ToDate:   pgtype.Timestamptz{},
			},
			mockScanFunc: func(ctx context.Context, params models.ScanReportsParams) ([]interface{}, error) {
				// Simulate a database error
				return nil, errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &server.Server{
				DB: &mockDB{
					mockScanReports: tc.mockScanFunc,
				},
			}

			got, err := s.Scan(tc.scanParams)
			if (err != nil) != tc.wantErr {
				t.Errorf("Server.Scan() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil && !isEqual(got, tc.want) {
				t.Errorf("Server.Scan() = %v, want %v", got, tc.want)
			}
		})
	}
}

// isEqual helps to compare if two slices of reports are equal.
// Adjust comparison according to your model's fields.
func isEqual(a, b *[]models.Report) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}
	for i, v := range *a {
		if v != (*b)[i] {
			return false
		}
	}
	return true
}
