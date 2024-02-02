-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE id = ANY($1::text[]);

-- name: CreateUser :one
INSERT INTO users (
    email
) VALUES (
    $1
)
RETURNING *;

-- name: UpdateUserCustomerID :exec
UPDATE users
SET customer_id = $1
WHERE id = $2;




-- name: CreateOTP :one
INSERT INTO otps (
    expires_at,
    is_active,
    user_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetOTP :one
SELECT * FROM otps
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
SELECT scan($1, $2, $3, $4, $5, $6, $7);

-- name: CreateScannerNotifs :many
SELECT scanner_notifs($1, $2, $3);

