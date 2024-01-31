package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLen = 6
const expirationDurationSeconds = 60 * 60 * 24 * 15

func (server *Server) Signup(ginContext *gin.Context, email string, password string) (*models.User, error) {

	if len(password) < minPasswordLen {
		return nil, fmt.Errorf("signup error: password must bigger than %v charachters", minPasswordLen)
	}

	hashedPassword, err := encryptPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := server.DB.queries.CreateUser(context.Background(), models.CreateUserParams{Email: email, HashedPassword: string(hashedPassword)})
	if err != nil {
		return nil, fmt.Errorf("signup error: %w", err)
	}

	errAuth := authorizeUser(ginContext, &user, server.Env.JWT_SECRET)
	if errAuth != nil {
		return nil, errAuth
	}

	return &user, nil
}

func (server *Server) Login(ginContext *gin.Context, email string, password string) (*models.User, error) {

	user, err := server.DB.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("login error: Invalid email/password")
	}

	errHashed := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if errHashed != nil {
		return nil, fmt.Errorf("login error: Invalid email/password")
	}

	errAuth := authorizeUser(ginContext, &user, server.Env.JWT_SECRET)
	if errAuth != nil {
		return nil, errAuth
	}

	return &user, nil

}

func (server *Server) ParseAuthorizationToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(server.Env.JWT_SECRET), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parsing error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt getting claims error")
	}

	return claims, nil
}

func (server *Server) ValidateUserAuthorization(mapClaims jwt.MapClaims) (*models.User, error) {
	// check the exp
	if mapClaims["exp"].(int64) < time.Now().Unix() {
		return nil, fmt.Errorf("jwt is expired error")
	}

	// find user
	userJWTId := pgtype.UUID{
		Bytes: mapClaims["sub"].([16]byte),
		Valid: true,
	}
	user, err := server.DB.queries.GetUserByJWTId(context.Background(), userJWTId)
	if err != nil {
		return nil, fmt.Errorf("jwt parsed error: user not found")
	}

	return &user, nil

}

func encryptPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf("password hashing error: %w", err)
	}
	return hashedPassword, nil
}

func authorizeUser(ginContext *gin.Context, user *models.User, jwtSecret string) error {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.JwtID.Bytes,
		"exp": time.Now().Add(time.Second * expirationDurationSeconds).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return fmt.Errorf("jwt creation error: %w", err)
	}

	// attach to cookie
	ginContext.SetSameSite(http.SameSiteLaxMode)
	ginContext.SetCookie("Authorization", tokenStr, expirationDurationSeconds, "/", "", true, true)

	return nil
}
