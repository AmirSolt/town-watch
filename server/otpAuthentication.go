package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

const otpExpirationDurationSeconds = 60 * 5  // 5 minutes
const otpRetryExpirationDurationSeconds = 10 // 10 sec

func (server *Server) InitOTP(email string) error {

	var user models.User
	var err error

	user, err = server.DB.queries.GetUserByEmail(context.Background(), email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error user email lookup: %w", err)
	}

	if err == sql.ErrNoRows {
		user, err = server.DB.queries.CreateUser(context.Background(), email)
		if err != nil {
			return fmt.Errorf("error user creation: %w", err)
		}
	}

	otp, err := server.DB.queries.CreateOTP(context.Background(), models.CreateOTPParams{
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(time.Second * otpExpirationDurationSeconds).UTC(), Valid: true},
		IsActive:  true,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error OTP creation: %w", err)
	}

	errEmail := server.SendOTP(&user, &otp)
	if errEmail != nil {
		return errEmail
	}
	return nil
}

func (server *Server) SendOTP(user *models.User, otp *models.Otp) error {
	content := "content" + string(otp.ID.Bytes[:])
	errEmail := server.SendEmail(user.Email, "User", "Town Watch", "Email Verification Link", content)
	if errEmail != nil {
		return fmt.Errorf("error OTP email could not be sent: %w", errEmail)
	}
	return nil
}

func (server *Server) ResendOTP(email string) error {
	user, err := server.DB.queries.GetUserByEmail(context.Background(), email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error user email lookup: %w", err)
	}

	lastOTP, err := server.DB.queries.GetLatestOTPByUser(context.Background(), user.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error latest otp lookup: %w", err)
	}

	if time.Now().Add(-time.Second * otpRetryExpirationDurationSeconds).UTC().Before(lastOTP.CreatedAt.Time) {
		return fmt.Errorf("you have to wait %v after sending OTP: %w", otpRetryExpirationDurationSeconds, err)
	}

	errOTP := server.InitOTP(email)
	if errOTP != nil {
		return fmt.Errorf("error resend InitOTP: %w", errOTP)
	}

	return nil
}

func (server *Server) ValidateOTP(ginContext *gin.Context, otpId string) error {

	otp, err := server.DB.queries.GetOTP(context.Background(), pgtype.UUID{Bytes: stringToByte16(otpId), Valid: true})
	if err != nil {
		return fmt.Errorf("error OTP lookup: %w", err)
	}

	if !otp.IsActive {
		return fmt.Errorf("error OTP is not active: %w", err)
	}

	if time.Now().UTC().After(otp.ExpiresAt.Time) {
		return fmt.Errorf("error OTP is expired: %w", err)
	}

	user, err := server.DB.queries.GetUser(context.Background(), otp.UserID)
	if err != nil {
		return fmt.Errorf("error user not found by OTP: %w", err)
	}

	server.SetJWT(ginContext, &user)

	return nil
}

// func (server *Server) Signup(ginContext *gin.Context, email string, password string) (*models.User, error) {

// 	if len(password) < minPasswordLen {
// 		return nil, fmt.Errorf("signup error: password must bigger than %v charachters", minPasswordLen)
// 	}

// 	hashedPassword, err := hashPassword(password, server.Env.PASSWORD_HASHING_SALT)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := server.DB.queries.CreateUser(context.Background(), models.CreateUserParams{Email: email, HashedPassword: string(hashedPassword)})
// 	if err != nil {
// 		return nil, fmt.Errorf("signup error: %w", err)
// 	}

// 	errJWT := setJWT(ginContext, &user, server.Env.JWE_SECRET_KEY)
// 	if errJWT != nil {
// 		return nil, errJWT
// 	}

// 	return &user, nil
// }

// func (server *Server) Login(ginContext *gin.Context, email string, password string) (*models.User, error) {

// 	user, err := server.DB.queries.GetUserByEmail(context.Background(), email)
// 	if err != nil {
// 		return nil, fmt.Errorf("login error: Invalid email/password")
// 	}

// 	errCompare := compareToHashedPassword(&user, password, server.Env.PASSWORD_HASHING_SALT)
// 	if errCompare != nil {
// 		return nil, fmt.Errorf("login error: Invalid email/password")
// 	}

// 	errJWT := setJWT(ginContext, &user, server.Env.JWE_SECRET_KEY)
// 	if errJWT != nil {
// 		return nil, errJWT
// 	}

// 	return &user, nil

// }
