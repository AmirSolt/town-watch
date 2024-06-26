package server_test

import (
	"fmt"
	"testing"
	"time"
)

func TestFetchReports(t *testing.T) {

	tests := []struct {
		inputFromDate time.Time
		inputToDate   time.Time
		want          bool
		ok            bool
	}{
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -2), inputToDate: time.Now().UTC().AddDate(0, 0, -1), want: true, ok: true},
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -1), inputToDate: time.Now().UTC().AddDate(0, 0, -1), want: false, ok: true},
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -1), inputToDate: time.Now().UTC().AddDate(0, 0, 0), want: false, ok: true},
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -30), inputToDate: time.Now().UTC().AddDate(0, 0, 0), want: true, ok: true},
		{inputFromDate: time.Now().UTC().AddDate(0, 0, 2), inputToDate: time.Now().UTC().AddDate(0, 0, 0), want: false, ok: true},
	}

	server := loadTestServer()

	for i, tc := range tests {

		got, err := server.FetchArcgisReports(tc.inputFromDate, tc.inputToDate)

		fmt.Println(">>> Features Len:", len(got.Features))
		if (err == nil) != tc.ok || (len(got.Features) > 0) != tc.want {
			t.Fatalf("test %d: expected: %v, got: %v, ok:%v error:%v", i+1, tc.want, len(got.Features), tc.ok, err)
		}

		time.Sleep(1 * time.Second)
	}
}

func TestConvertArcgisResponseToReportsParams(t *testing.T) {
	tests := []struct {
		inputFromDate time.Time
		inputToDate   time.Time
		want          bool
		ok            bool
	}{
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -2), inputToDate: time.Now().UTC().AddDate(0, 0, -1), want: true, ok: true},
	}

	server := loadTestServer()

	for i, tc := range tests {
		arcReports, err := server.FetchArcgisReports(tc.inputFromDate, tc.inputToDate)

		got := server.ConvertArcgisResponseToReportsParams(arcReports)

		fmt.Println(">>> Reports Len:", len(*got))
		if (err == nil) != tc.ok || (len(*got) > 0) != tc.want {
			t.Fatalf("test %d: expected: %v, got: %v, ok:%v error:%v", i+1, tc.want, len(*got), tc.ok, err)
		}

		time.Sleep(1 * time.Second)
	}
}

func TestCreateReports(t *testing.T) {
	tests := []struct {
		inputFromDate time.Time
		inputToDate   time.Time
		want          bool
		ok            bool
	}{
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -2), inputToDate: time.Now().UTC().AddDate(0, 0, -1), ok: true},
	}

	server := loadTestServer()

	for i, tc := range tests {
		arcReports, err := server.FetchArcgisReports(tc.inputFromDate, tc.inputToDate)

		if (err == nil) != tc.ok {
			t.Fatalf("test %d: ok:%v error:%v", i+1, tc.ok, err)
		}

		reportsParams := server.ConvertArcgisResponseToReportsParams(arcReports)

		server.CreateReports(reportsParams)

		// fmt.Println(">>> Reports Len:", len(*got))
		// if (err == nil) != tc.ok || (len(*got) > 0) != tc.want {
		// 	t.Fatalf("test %d: expected: %v, got: %v, ok:%v error:%v", i+1, tc.want, len(*got), tc.ok, err)
		// }

		// time.Sleep(1 * time.Second)
	}
}
