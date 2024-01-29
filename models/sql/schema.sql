-- Amirali Soltani
-- 2024-01-29
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE region AS ENUM ('TORONTO');


CREATE TYPE crime_type AS ENUM (
    'Assault',
    'Auto Theft',
    'Theft from Motor Vehicle',
    'Break and Enter',
    'Sexual Violation',
    'Robbery',
    'Theft Over',
    'Bike Theft',
    'Shooting',
    'Homicide'
);

-- ======

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email TEXT NOT NULL
);

-- ======

CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    occur_at TIMESTAMPTZ NOT NULL,
    external_src_id TEXT NOT NULL,
    neighborhood TEXT,
    location_type TEXT,
    crime_type crime_type NOT NULL,
    region region NOT NULL,
    point geometry(Point, 3857) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX report_occ_at_idx ON reports ("occur_at");
CREATE INDEX report_point_idx ON reports USING GIST ("point");
CREATE FUNCTION report_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.lat, NEW.long, 3857)
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_report_insert BEFORE INSERT OR UPDATE ON reports
    FOR EACH ROW EXECUTE FUNCTION report_insert();

-- ======

CREATE TABLE scanners (
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
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE FUNCTION scanner_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.lat, NEW.long, 3857)
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_scanner_insert BEFORE INSERT OR UPDATE ON scanners
    FOR EACH ROW EXECUTE FUNCTION scanner_insert();

-- ======

CREATE TABLE notifs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_sent BOOLEAN NOT NULL DEFAULT false,
    is_opened BOOLEAN NOT NULL DEFAULT false,
    scanner_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_scanner FOREIGN KEY(scanner_id) REFERENCES scanners(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- ======

CREATE TABLE report_notifs (
    PRIMARY KEY (report_id, notif_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notif_id uuid NOT NULL,
    report_id INTEGER NOT NULL,
    CONSTRAINT fk_notif FOREIGN KEY(notif_id) REFERENCES notifs(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_report FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE UNIQUE INDEX report_notif_notif_id_key ON report_notifs("notif_id");
CREATE UNIQUE INDEX report_notif_report_id_key ON report_notifs("report_id");