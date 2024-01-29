package models

import "time"

type Notif struct {
	Id        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Is_sent   bool      `db:"is_sent"`
	Is_opened bool      `db:"is_opened"`

	Scanner Scanner
	Reports []Report
}
