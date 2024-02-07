package server_test

import "github.com/AmirSolt/town-watch/server"

func loadTestServer() *server.Server {
	server := server.Server{
		RootDir: "../../",
	}
	server.LoadServer()

	return &server
}
