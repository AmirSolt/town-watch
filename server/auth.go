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

const minPasswordLen = 6
const expirationDurationSeconds = 60 * 60 * 24 * 15

type JWT struct {
	JWT_ID string `json:"jwt_id"`
	IP     string `json:"ip"`
	EXP    int64  `json:"exp"`
}

func (server *Server) Signup(ginContext *gin.Context, email string, password string) (*models.User, error) {

	if len(password) < minPasswordLen {
		return nil, fmt.Errorf("signup error: password must bigger than %v charachters", minPasswordLen)
	}

	hashedPassword, err := hashPassword(password, server.Env.PASSWORD_HASHING_SALT)
	if err != nil {
		return nil, err
	}

	user, err := server.DB.queries.CreateUser(context.Background(), models.CreateUserParams{Email: email, HashedPassword: string(hashedPassword)})
	if err != nil {
		return nil, fmt.Errorf("signup error: %w", err)
	}

	errJWT := setJWT(ginContext, &user, server.Env.JWE_SECRET_KEY)
	if errJWT != nil {
		return nil, errJWT
	}

	return &user, nil
}

func (server *Server) Login(ginContext *gin.Context, email string, password string) (*models.User, error) {

	user, err := server.DB.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("login error: Invalid email/password")
	}

	errCompare := compareToHashedPassword(&user, password, server.Env.PASSWORD_HASHING_SALT)
	if errCompare != nil {
		return nil, fmt.Errorf("login error: Invalid email/password")
	}

	errJWT := setJWT(ginContext, &user, server.Env.JWE_SECRET_KEY)
	if errJWT != nil {
		return nil, errJWT
	}

	return &user, nil

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
