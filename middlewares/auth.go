package middlewares

import (
	"net/http"

	"github.com/AmirSolt/town-watch/server"
	"github.com/gin-gonic/gin"
)

func RequireAuth(ginContext *gin.Context, server *server.Server) {
	// get it from cookie
	tokenString, err := ginContext.Cookie("Authorization")
	if err != nil {
		ginContext.AbortWithStatus(http.StatusUnauthorized)
	}

	// parse and validate token
	mapClaims, err := server.ParseAuthorizationToken(tokenString)
	if err != nil {
		ginContext.AbortWithStatus(http.StatusUnauthorized)
	}

	// find user and check exp
	user, err := server.ValidateUserAuthorization(mapClaims)
	if err != nil {
		ginContext.AbortWithStatus(http.StatusUnauthorized)
	}

	// attach user to locals
	ginContext.Set("user", user)

	ginContext.Next()
}
