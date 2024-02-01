-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByAuthoId :one
SELECT * FROM users
WHERE autho_id = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE id = ANY($1::int[]);

-- name: CreateUser :one
INSERT INTO users (
    email
) VALUES (
    $1
)
RETURNING *;



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
SELECT *
FROM reports
WHERE 
ST_DWithin(
    point,
    ST_Point($1, $2, 3857),
    $3
)
AND region = $4
AND occur_at >= $5
AND occur_at <= $6
ORDER BY occur_at
LIMIT $7;

-- name: CreateScannerNotifs :many
SELECT scanner_notifs($1, $2, $3);

