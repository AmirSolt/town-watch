package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AmirSolt/town-watch/server"
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

	server := server.LoadServer("../")

	for i, tc := range tests {
		got, err := server.FetchReports(tc.inputFromDate, tc.inputToDate)

		fmt.Println(">>> Features Len:", len(got.Features))
		if (err == nil) != tc.ok || (len(got.Features) > 0) != tc.want {
			t.Fatalf("test %d: expected: %v, got: %v, ok:%v error:%v", i+1, tc.want, len(got.Features), tc.ok, err)
		}

		time.Sleep(1 * time.Second)
	}
}

func TestConvertArcgisResponseToReports(t *testing.T) {
	tests := []struct {
		inputFromDate time.Time
		inputToDate   time.Time
		want          bool
		ok            bool
	}{
		{inputFromDate: time.Now().UTC().AddDate(0, 0, -2), inputToDate: time.Now().UTC().AddDate(0, 0, -1), want: true, ok: true},
	}

	server := server.LoadServer("../")

	for i, tc := range tests {
		arcReports, err := server.FetchReports(tc.inputFromDate, tc.inputToDate)

		// new code
		got := server.ConvertArcgisResponseToReports(arcReports)

		fmt.Println(">>> Reports Len:", len(*got))
		if (err == nil) != tc.ok || (len(*got) > 0) != tc.want {
			t.Fatalf("test %d: expected: %v, got: %v, ok:%v error:%v", i+1, tc.want, len(*got), tc.ok, err)
		}

		time.Sleep(1 * time.Second)
	}
}
