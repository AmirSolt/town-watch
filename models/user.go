package models

import "time"

type User struct {
	Id        int8      `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Email     string    `db:"email"`

	Scanners []Scanner
}

const UserSchema string = `
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email TEXT NOT NULL
);
`
