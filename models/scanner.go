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

const ScannerSchema string = `
CREATE TABLE IF NOT EXISTS scanner (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    address TEXT,
    region region NOT NULL,
    radius DOUBLE PRECISION NOT NULL,
    point geometry(Point, 3857) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,

	user_id INT NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE
);
`
