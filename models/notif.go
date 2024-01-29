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

const NotifSchema string = `
CREATE TABLE notif (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_sent BOOLEAN NOT NULL DEFAULT false,
    is_opened BOOLEAN NOT NULL DEFAULT false,

	scanner_id INT NOT NULL,
	CONSTRAINT fk_scanner FOREIGN KEY(scanner_id) REFERENCES scanner(id) ON DELETE CASCADE ON UPDATE CASCADE;
);
`
