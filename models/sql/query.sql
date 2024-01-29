-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateReports :exec
INSERT INTO reports (
    occur_at,
    external_src_id,
    neighborhood,
    location_type,
    crime_type,
    region,
    point,
    lat,
    long
) VALUES ($1,$2,$3,$4,$5,$6,ST_POINT($7, $8 ,3857),$7,$8);

