package server

import (
	"path/filepath"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
)

type Server struct {
	RootDir     string
	Engine      *gin.Engine
	Env         *Env
	DB          *DB
	TierConfigs map[models.Tier]TierConfig
}

func (server *Server) LoadServer() {

	// gin.DisableConsoleColor()
	server.Engine = gin.Default()

	server.Engine.LoadHTMLGlob(filepath.Join(server.RootDir, "public/*"))

	server.loadEnv()

	server.loadDB()

	server.loadPayment()
}

func (server *Server) KillServer() {
	server.killDB()
}
