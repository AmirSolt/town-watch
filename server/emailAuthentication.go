package server

import (
	"context"
	"fmt"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
)

const minPasswordLen = 6

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
