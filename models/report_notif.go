package models

import "time"

type ReportNotif struct {
	Id        int8      `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	Notif  Notif
	Report Report
}
