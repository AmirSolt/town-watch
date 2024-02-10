package routes

import (
	"github.com/AmirSolt/town-watch/models"
	"github.com/AmirSolt/town-watch/server"
)

type pageLoad struct {
	Title  string
	Domain string
	User   *models.User
}

type Routes struct {
	server *server.Server
}

func LoadRoutes(server *server.Server) *Routes {
	routes := Routes{
		server: server,
	}
	routes.index()

	return &routes
}
