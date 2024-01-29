package models

import "time"

type User struct {
	Id        int8      `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Email     string    `db:"email"`

	Scanners []Scanner
}
