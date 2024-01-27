package server_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/AmirSolt/town-watch/server"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		inputFromDate time.Time
		inputToDate   time.Time
		want          int32
		ok            bool
	}{
		{inputFromDate: time.Now(), inputToDate: time.Now(), want: 20, ok: true},
	}

	server := server.LoadServer("../")

	for i, tc := range tests {
		got, err := server.FetchReports(tc.inputFromDate, tc.inputToDate)
		if !reflect.DeepEqual(tc.want, got) || ((err == nil) == tc.ok) {
			t.Fatalf("test %d: expected: %v, got: %v, ok:%v error:%v", i+1, tc.want, got, tc.ok, err)
		}
	}
}
