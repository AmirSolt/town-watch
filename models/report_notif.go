package models

import "time"

type ReportNotif struct {
	Id        int8      `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	Notif  Notif
	Report Report
}

const ReportNotifSchema string = `
CREATE TABLE report_notif (
	PRIMARY KEY (article_id, tag_id)
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notif_id INTEGER NOT NULL,
    report_id INTEGER NOT NULL,

	CONSTRAINT fk_notif FOREIGN KEY(notif_id) REFERENCES notif(id)
	CONSTRAINT fk_report FOREIGN KEY(report_id) REFERENCES report(id)
);
`
