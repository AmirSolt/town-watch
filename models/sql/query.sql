-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateReports :copyfrom
INSERT INTO reports (
    occur_at,
    external_src_id,
    neighborhood,
    location_type,
    crime_type,
    region,
    lat,
    long
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);


-- name: ScanReports :many

-- name: CreateScannerNotifs :many
