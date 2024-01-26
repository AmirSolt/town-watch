package routes

import (
	"github.com/AmirSolt/town-watch/plugins"
	"github.com/AmirSolt/town-watch/server"
)

type Routes struct{}

func LoadRoutes(server *server.Server, plugins *plugins.Plugins) *Routes {
	routes := Routes{}
	routes.index(server, plugins)

	return &routes
}
