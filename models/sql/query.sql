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



-- name: GetCustomerByStripeID :one
SELECT * 
FROM customers
WHERE stripe_customer_id = $1
LIMIT 1;

-- name: GetCustomerByUserID :one
SELECT * 
FROM customers
WHERE user_id = $1
LIMIT 1;





-- name: CreateSubscription :one
INSERT INTO subscriptions( 
    stripe_subscription_id,
    tier_id,
    is_active,
    customer_id
) VALUES ($1,$2,$3,$4)
RETURNING *;

-- name: GetActiveSubscriptionByCustomer :one
SELECT * 
FROM subscriptions
WHERE customer_id = $1 
AND is_active = TRUE
LIMIT 1;

-- name: DeactivateSubscriptionByStripeID :exec
UPDATE subscriptions
SET is_active = FALSE
WHERE stripe_subscription_id=$1;

-- name: DeactivateSubscriptionByCustomerID :one
UPDATE subscriptions
SET is_active = FALSE
WHERE customer_id=$1
RETURNING *;




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

