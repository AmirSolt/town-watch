package server

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Server struct {
	RootDir string
	Engine  *gin.Engine
	Env     *Env
	DB      *DB
}

func (server *Server) LoadServer() {

	// gin.DisableConsoleColor()
	server.Engine = gin.Default()

	server.Engine.LoadHTMLGlob(filepath.Join(server.RootDir, "public/*"))

	server.loadEnv()

	server.loadDB()
}

func (server *Server) KillServer() {
	server.killDB()
}
