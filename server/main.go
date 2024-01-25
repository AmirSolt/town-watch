package server

import "github.com/gin-gonic/gin"

type Server struct {
	Engine *gin.Engine
}

func LoadServer() *Server {

	// gin.DisableConsoleColor()
	server := Server{
		Engine: gin.Default(),
	}
	server.Engine.LoadHTMLGlob("public/*")

	return &server
}
