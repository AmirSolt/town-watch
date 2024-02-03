package routes

import (
	"github.com/AmirSolt/town-watch/server"
)

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
