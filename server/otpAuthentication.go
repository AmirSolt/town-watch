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

	user, err := server.findOrCreateUser(email)
	if err != nil {
		return err
	}
	otp, err := server.createOTP(user)
	if err != nil {
		return err
	}

	errEmail := server.sendOTPEmail(user, otp)
	if errEmail != nil {
		return errEmail
	}
	return nil
}

func (server *Server) ResendOTP(email string) error {

	// find user
	user, err := server.DB.Queries.GetUserByEmail(context.Background(), email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error user email lookup: %w", err)
	}

	// =======================
	// make sure last otp happened before otpRetryExpirationDurationSeconds ago
	lastOTP, err := server.DB.Queries.GetLatestOTPByUser(context.Background(), user.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error latest otp lookup: %w", err)
	}
	if time.Now().Add(-time.Second * otpRetryExpirationDurationSeconds).UTC().Before(lastOTP.CreatedAt.Time) {
		return fmt.Errorf("you have to wait %v after sending OTP: %w", otpRetryExpirationDurationSeconds, err)
	}
	// =======================

	otp, err := server.createOTP(&user)
	if err != nil {
		return err
	}

	errEmail := server.sendOTPEmail(&user, otp)
	if errEmail != nil {
		return errEmail
	}

	return nil
}

func (server *Server) ValidateOTP(ginContext *gin.Context, otpId string) error {

	// Find OTP
	otp, err := server.DB.Queries.GetOTP(context.Background(), pgtype.UUID{Bytes: stringToByte16(otpId), Valid: true})
	if err != nil {
		return fmt.Errorf("error OTP lookup: %w", err)
	}

	if !otp.IsActive {
		return fmt.Errorf("error OTP is not active: %w", err)
	}
	defer server.deactivateOTP(&otp)

	if time.Now().UTC().After(otp.ExpiresAt.Time) {
		return fmt.Errorf("error OTP is expired: %w", err)
	}

	user, err := server.DB.Queries.GetUser(context.Background(), otp.UserID)
	if err != nil {
		return fmt.Errorf("error user not found by OTP: %w", err)
	}

	lastOTP, err := server.DB.Queries.GetLatestOTPByUser(context.Background(), user.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error latest otp lookup: %w", err)
	}

	if lastOTP.ID != otp.ID {
		return fmt.Errorf("otp does not match latest user otp: %w", err)
	}

	server.SetJWTCookie(ginContext, &user)

	return nil
}

func (server *Server) Signout(ginContext *gin.Context) {
	server.removeJWTCookie(ginContext)
}

// =====================================================================

func (server *Server) findOrCreateUser(email string) (*models.User, error) {
	var user models.User
	var err error

	user, err = server.DB.Queries.GetUserByEmail(context.Background(), email)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error user email lookup: %w", err)
	}

	if err == sql.ErrNoRows {
		user, err = server.DB.Queries.CreateUser(context.Background(), email)
		if err != nil {
			return nil, fmt.Errorf("error user creation: %w", err)
		}
	}

	return &user, nil
}

func (server *Server) createOTP(user *models.User) (*models.Otp, error) {
	otp, err := server.DB.Queries.CreateOTP(context.Background(), models.CreateOTPParams{
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(time.Second * otpExpirationDurationSeconds).UTC(), Valid: true},
		IsActive:  true,
		UserID:    user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("error OTP creation: %w", err)
	}
	return &otp, nil
}

func (server *Server) sendOTPEmail(user *models.User, otp *models.Otp) error {
	content := "content" + string(otp.ID.Bytes[:])
	errEmail := server.SendEmail(user.Email, "User", "Town Watch", "Email Verification Link", content)
	if errEmail != nil {
		return fmt.Errorf("error OTP email could not be sent: %w", errEmail)
	}
	return nil
}

func (server *Server) deactivateOTP(otp *models.Otp) error {
	err := server.DB.Queries.DeactivateOTP(context.Background(), otp.ID)
	if err != nil {
		return fmt.Errorf("deactivating otp failed: %w", err)
	}
	return nil
}
