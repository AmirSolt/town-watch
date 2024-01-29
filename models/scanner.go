package models

import "time"

type Scanner struct {
	Id        int8      `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	IsActive  bool      `db:"is_active"`
	Address   string    `db:"address"`
	Region    Region    `db:"region"`
	Radius    float32   `db:"radius"`
	Lat       float32   `db:"lat"`
	Long      float32   `db:"long"`

	User   User
	Notifs []Notif
}
