// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOTP = `-- name: CreateOTP :one
INSERT INTO otps (
    expires_at,
    is_active,
    user_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, created_at, expires_at, is_active, user_id
`

type CreateOTPParams struct {
	ExpiresAt pgtype.Timestamptz
	IsActive  bool
	UserID    int32
}

func (q *Queries) CreateOTP(ctx context.Context, arg CreateOTPParams) (Otp, error) {
	row := q.db.QueryRow(ctx, createOTP, arg.ExpiresAt, arg.IsActive, arg.UserID)
	var i Otp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsActive,
		&i.UserID,
	)
	return i, err
}

type CreateReportsParams struct {
	OccurAt       pgtype.Timestamptz
	ExternalSrcID string
	Neighborhood  pgtype.Text
	LocationType  pgtype.Text
	CrimeType     CrimeType
	Region        Region
	Lat           float64
	Long          float64
}

const createScannerNotifs = `-- name: CreateScannerNotifs :many
SELECT scanner_notifs($1, $2, $3)
`

type CreateScannerNotifsParams struct {
	FromDate         pgtype.Timestamptz
	ToDate           pgtype.Timestamptz
	ScanReportsLimit int32
}

func (q *Queries) CreateScannerNotifs(ctx context.Context, arg CreateScannerNotifsParams) ([]pgtype.UUID, error) {
	rows, err := q.db.Query(ctx, createScannerNotifs, arg.FromDate, arg.ToDate, arg.ScanReportsLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var scanner_notifs pgtype.UUID
		if err := rows.Scan(&scanner_notifs); err != nil {
			return nil, err
		}
		items = append(items, scanner_notifs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    email
) VALUES (
    $1
)
RETURNING id, member, autho_id, created_at, email
`

func (q *Queries) CreateUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, createUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Member,
		&i.AuthoID,
		&i.CreatedAt,
		&i.Email,
	)
	return i, err
}

const getOTP = `-- name: GetOTP :one
SELECT id, created_at, expires_at, is_active, user_id FROM otps
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOTP(ctx context.Context, id pgtype.UUID) (Otp, error) {
	row := q.db.QueryRow(ctx, getOTP, id)
	var i Otp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsActive,
		&i.UserID,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, member, autho_id, created_at, email FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Member,
		&i.AuthoID,
		&i.CreatedAt,
		&i.Email,
	)
	return i, err
}

const getUserByAuthoId = `-- name: GetUserByAuthoId :one
SELECT id, member, autho_id, created_at, email FROM users
WHERE autho_id = $1 LIMIT 1
`

func (q *Queries) GetUserByAuthoId(ctx context.Context, authoID pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByAuthoId, authoID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Member,
		&i.AuthoID,
		&i.CreatedAt,
		&i.Email,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, member, autho_id, created_at, email FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Member,
		&i.AuthoID,
		&i.CreatedAt,
		&i.Email,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, member, autho_id, created_at, email FROM users
WHERE id = ANY($1::int[])
`

func (q *Queries) GetUsers(ctx context.Context, dollar_1 []int32) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Member,
			&i.AuthoID,
			&i.CreatedAt,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const scanReports = `-- name: ScanReports :many
SELECT id, created_at, occur_at, external_src_id, neighborhood, location_type, crime_type, region, point, lat, long
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
LIMIT $7
`

type ScanReportsParams struct {
	StPoint   interface{}
	StPoint_2 interface{}
	StDwithin interface{}
	Region    Region
	OccurAt   pgtype.Timestamptz
	OccurAt_2 pgtype.Timestamptz
	Limit     int32
}

func (q *Queries) ScanReports(ctx context.Context, arg ScanReportsParams) ([]Report, error) {
	rows, err := q.db.Query(ctx, scanReports,
		arg.StPoint,
		arg.StPoint_2,
		arg.StDwithin,
		arg.Region,
		arg.OccurAt,
		arg.OccurAt_2,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Report
	for rows.Next() {
		var i Report
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.OccurAt,
			&i.ExternalSrcID,
			&i.Neighborhood,
			&i.LocationType,
			&i.CrimeType,
			&i.Region,
			&i.Point,
			&i.Lat,
			&i.Long,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
