package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"golang.org/x/crypto/bcrypt"
)

const expirationDurationSeconds = 60 * 60 * 24 * 15

type JWT struct {
	JWT_ID string `json:"jwt_id"`
	IP     string `json:"ip"`
	EXP    int64  `json:"exp"`
}

func (server *Server) ParseJWT(jwtEncrypted string) (*JWT, error) {
	jwt, err := dencryptJWT([]byte(jwtEncrypted), server.Env.JWE_SECRET_KEY)
	if err != nil {
		return nil, fmt.Errorf("jwt authorization failed: %w", err)
	}
	return jwt, nil
}

func (server *Server) ValidateUserByJWT(ginContext *gin.Context, jwt *JWT) (*models.User, error) {

	if jwt.EXP < time.Now().Unix() {
		return nil, fmt.Errorf("error jwt is expired")
	}

	if jwt.IP == ginContext.ClientIP() {
		return nil, fmt.Errorf("error jwt is from an invalid IP")
	}

	userJWTId := pgtype.UUID{
		Bytes: stringToByte16(jwt.JWT_ID),
		Valid: true,
	}
	user, err := server.DB.queries.GetUserByJWTId(context.Background(), userJWTId)
	if err != nil {
		return nil, fmt.Errorf("error jwt user not found")
	}

	return &user, nil
}

// ==============================================================

func setJWT(ginContext *gin.Context, user *models.User, jwtSecret string) error {

	jwt := JWT{
		JWT_ID: string(user.JwtID.Bytes[:]),
		IP:     ginContext.ClientIP(),
		EXP:    time.Now().Add(time.Second * expirationDurationSeconds).Unix(),
	}

	jwtEncrypted, err := encryptJWT(jwt, jwtSecret)
	if err != nil {
		return fmt.Errorf("jwt authorization failed: %w", err)
	}

	// attach to cookie
	ginContext.SetSameSite(http.SameSiteLaxMode)
	ginContext.SetCookie("Authorization", string(jwtEncrypted), expirationDurationSeconds, "/", "", true, true)

	return nil
}

func hashPassword(password string, salt string) ([]byte, error) {
	saltedPassword := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), 10)
	if err != nil {
		return nil, fmt.Errorf("password hashing error: %w", err)
	}
	return hashedPassword, nil
}

func compareToHashedPassword(user *models.User, password string, salt string) error {
	saltedPassword := password + salt

	errHashed := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(saltedPassword))
	if errHashed != nil {
		return fmt.Errorf("login error: password failed compare")
	}
	return nil
}

func encryptJWT(jwt JWT, jwe_secret_key string) ([]byte, error) {

	payload, err := json.Marshal(jwt)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal payload: %w", err)
	}
	encrypted, err := jwe.Encrypt(payload, jwe.WithKey(jwa.A128GCM, jwe_secret_key))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt jwt payload: %w", err)
	}
	return encrypted, nil
}
func dencryptJWT(encryptedJWT []byte, jwe_secret_key string) (*JWT, error) {
	decrypted, err := jwe.Decrypt(encryptedJWT, jwe.WithKey(jwa.A128GCM, jwe_secret_key))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt payload: %w", err)
	}

	var jwt JWT
	errJson := json.Unmarshal(decrypted, &jwt)
	if errJson != nil {
		return nil, fmt.Errorf("failed to json unmarshal payload: %w", err)
	}

	return &jwt, nil
}

func stringToByte16(str string) [16]byte {
	var arr [16]byte
	byteSlice := []byte(str)
	copy(arr[:], byteSlice)
	return arr
}
