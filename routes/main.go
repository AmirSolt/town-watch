package routes

import (
	"github.com/AmirSolt/town-watch/server"
)

type Routes struct{}

func LoadRoutes(server *server.Server) *Routes {
	routes := Routes{}
	routes.index(server)

	return &routes
}
