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

func LoadServer(rootDir string) *Server {

	// gin.DisableConsoleColor()
	server := Server{
		RootDir: rootDir,
		Engine:  gin.Default(),
	}

	server.Engine.LoadHTMLGlob(filepath.Join(server.RootDir, "public/*"))

	server.loadEnv()

	return &server
}
