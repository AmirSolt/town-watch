package main

import (
	"github.com/AmirSolt/town-watch/plugins"
	"github.com/AmirSolt/town-watch/routes"
	"github.com/AmirSolt/town-watch/server"
)

func main() {

	server := server.LoadServer()
	plugins := plugins.LoadPlugins(server)
	routes.LoadRoutes(server, plugins)

	server.Engine.Run()
}
