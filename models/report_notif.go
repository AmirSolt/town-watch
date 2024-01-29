package models

import "time"

type ReportNotif struct {
	Id        int8      `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	Notif  Notif
	Report Report
}

const ReportNotifSchema string = `
CREATE TABLE IF NOT EXISTS report_notif (
	PRIMARY KEY (report_id, notif_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notif_id uuid NOT NULL,
    report_id INTEGER NOT NULL,

	CONSTRAINT fk_notif FOREIGN KEY(notif_id) REFERENCES notif(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_report FOREIGN KEY(report_id) REFERENCES report(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS report_notif_notif_id_key ON report_notif("notif_id");
CREATE UNIQUE INDEX IF NOT EXISTS report_notif_report_id_key ON report_notif("report_id");
`
