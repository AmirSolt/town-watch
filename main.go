package main

import (
	"fmt"

	"github.com/AmirSolt/town-watch/routes"
	"github.com/AmirSolt/town-watch/server"
)

func main() {

	server := server.LoadServer("./")
	routes.LoadRoutes(server)

	fmt.Println("=======")
	fmt.Println("http://localhost:8080")
	fmt.Println("=======")

	server.Engine.Run()
}
