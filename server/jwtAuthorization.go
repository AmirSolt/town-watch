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
)

const jwtExpirationDurationSeconds = 60 * 60 * 24 * 15 // 15 days

type JWT struct {
	ID  string `json:"id"`
	IP  string `json:"ip"`
	EXP int64  `json:"exp"`
}

func (server *Server) ValidateUser(ginContext *gin.Context) {
	// get it from cookie
	tokenString, err := ginContext.Cookie("Authorization")
	if err != nil {
		// fmt.Errorf("jwt not found on cookie: %w", err)
		ginContext.Redirect(302, "/")
		return
	}

	// parse and validate token
	jwt, err := server.ParseJWT(tokenString)
	if err != nil {
		// fmt.Errorf("jwt parse failed: %w", err)
		ginContext.Redirect(302, "/")
		return
	}

	// find user and check exp
	user, err := server.ValidateJWTByUser(ginContext, jwt)
	if err != nil {
		// fmt.Errorf("jwt validation by user failed: %w", err)
		ginContext.Redirect(302, "/")
		return
	}

	ginContext.Set("user", user)

	ginContext.Next()
}

// func (server *Server) ValidateUser(ginContext *gin.Context) (*models.User, error) {
// 	// get it from cookie
// 	tokenString, err := ginContext.Cookie("Authorization")
// 	if err != nil {
// 		return nil, fmt.Errorf("jwt not found on cookie: %w", err)
// 	}

// 	// parse and validate token
// 	jwt, err := server.ParseJWT(tokenString)
// 	if err != nil {
// 		return nil, fmt.Errorf("jwt parse failed: %w", err)
// 	}

// 	// find user and check exp
// 	user, err := server.ValidateJWTByUser(ginContext, jwt)
// 	if err != nil {
// 		return nil, fmt.Errorf("jwt validation by user failed: %w", err)
// 	}
// 	return user, nil
// }

func (server *Server) SetJWT(ginContext *gin.Context, user *models.User) error {

	jwt := JWT{
		ID:  string(user.ID.Bytes[:]),
		IP:  ginContext.ClientIP(),
		EXP: time.Now().Add(time.Second * jwtExpirationDurationSeconds).Unix(),
	}

	jwtEncrypted, err := encryptJWT(jwt, server.Env.JWE_SECRET_KEY)
	if err != nil {
		return fmt.Errorf("jwt authorization failed: %w", err)
	}

	// attach to cookie
	ginContext.SetSameSite(http.SameSiteLaxMode)
	ginContext.SetCookie("Authorization", string(jwtEncrypted), jwtExpirationDurationSeconds, "/", "", true, true)

	return nil
}

func (server *Server) ParseJWT(jwtEncrypted string) (*JWT, error) {
	jwt, err := dencryptJWT([]byte(jwtEncrypted), server.Env.JWE_SECRET_KEY)
	if err != nil {
		return nil, fmt.Errorf("jwt authorization failed: %w", err)
	}
	return jwt, nil
}

func (server *Server) ValidateJWTByUser(ginContext *gin.Context, jwt *JWT) (*models.User, error) {

	if jwt.EXP < time.Now().Unix() {
		return nil, fmt.Errorf("error jwt is expired")
	}

	if jwt.IP == ginContext.ClientIP() {
		return nil, fmt.Errorf("error jwt is from an invalid IP")
	}

	UserID := pgtype.UUID{
		Bytes: stringToByte16(jwt.ID),
		Valid: true,
	}
	user, err := server.DB.Queries.GetUser(context.Background(), UserID)
	if err != nil {
		return nil, fmt.Errorf("error jwt user not found")
	}

	return &user, nil
}

// ==============================================================

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
