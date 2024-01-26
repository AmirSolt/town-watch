package main

import (
	"fmt"

	"github.com/AmirSolt/town-watch/plugins"
	"github.com/AmirSolt/town-watch/routes"
	"github.com/AmirSolt/town-watch/server"
)

func main() {

	server := server.LoadServer()
	plugins := plugins.LoadPlugins(server)
	routes.LoadRoutes(server, plugins)

	fmt.Println("=======")
	fmt.Println("http://localhost:8080")
	fmt.Println("=======")

	server.Engine.Run()
}
