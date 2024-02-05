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
	UserID    pgtype.UUID
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

func (q *Queries) CreateScannerNotifs(ctx context.Context, arg CreateScannerNotifsParams) ([]interface{}, error) {
	rows, err := q.db.Query(ctx, createScannerNotifs, arg.FromDate, arg.ToDate, arg.ScanReportsLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var scanner_notifs interface{}
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
RETURNING id, tier, created_at, email, stripe_customer_id, stripe_subscription_id
`

func (q *Queries) CreateUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, createUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Tier,
		&i.CreatedAt,
		&i.Email,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
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
SELECT id, tier, created_at, email, stripe_customer_id, stripe_subscription_id FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Tier,
		&i.CreatedAt,
		&i.Email,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, tier, created_at, email, stripe_customer_id, stripe_subscription_id FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Tier,
		&i.CreatedAt,
		&i.Email,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
	)
	return i, err
}

const getUserByStripeCustomerID = `-- name: GetUserByStripeCustomerID :one
SELECT id, tier, created_at, email, stripe_customer_id, stripe_subscription_id FROM users
WHERE stripe_customer_id = $1 LIMIT 1
`

func (q *Queries) GetUserByStripeCustomerID(ctx context.Context, stripeCustomerID pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByStripeCustomerID, stripeCustomerID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Tier,
		&i.CreatedAt,
		&i.Email,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, tier, created_at, email, stripe_customer_id, stripe_subscription_id FROM users
WHERE id = ANY($1::text[])
`

func (q *Queries) GetUsers(ctx context.Context, dollar_1 []string) ([]User, error) {
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
			&i.Tier,
			&i.CreatedAt,
			&i.Email,
			&i.StripeCustomerID,
			&i.StripeSubscriptionID,
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
SELECT scan($1, $2, $3, $4, $5, $6, $7)
`

type ScanReportsParams struct {
	Lat        float64
	Long       float64
	Radius     float64
	Region     Region
	FromDate   pgtype.Timestamptz
	ToDate     pgtype.Timestamptz
	CountLimit int32
}

func (q *Queries) ScanReports(ctx context.Context, arg ScanReportsParams) ([]interface{}, error) {
	rows, err := q.db.Query(ctx, scanReports,
		arg.Lat,
		arg.Long,
		arg.Radius,
		arg.Region,
		arg.FromDate,
		arg.ToDate,
		arg.CountLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var scan interface{}
		if err := rows.Scan(&scan); err != nil {
			return nil, err
		}
		items = append(items, scan)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserStripeCustomerID = `-- name: UpdateUserStripeCustomerID :exec
UPDATE users
SET stripe_customer_id = $1
WHERE id = $2
`

type UpdateUserStripeCustomerIDParams struct {
	StripeCustomerID pgtype.Text
	ID               pgtype.UUID
}

func (q *Queries) UpdateUserStripeCustomerID(ctx context.Context, arg UpdateUserStripeCustomerIDParams) error {
	_, err := q.db.Exec(ctx, updateUserStripeCustomerID, arg.StripeCustomerID, arg.ID)
	return err
}

const updateUserSubAndTier = `-- name: UpdateUserSubAndTier :exec
UPDATE users
SET 
stripe_subscription_id = $1,
tier = $2
WHERE id = $3
`

type UpdateUserSubAndTierParams struct {
	StripeSubscriptionID pgtype.Text
	Tier                 Tier
	ID                   pgtype.UUID
}

func (q *Queries) UpdateUserSubAndTier(ctx context.Context, arg UpdateUserSubAndTierParams) error {
	_, err := q.db.Exec(ctx, updateUserSubAndTier, arg.StripeSubscriptionID, arg.Tier, arg.ID)
	return err
}
