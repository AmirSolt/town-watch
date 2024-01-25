package plugins

import "github.com/AmirSolt/town-watch/server"

type Plugins struct{}

func LoadPlugins(server *server.Server) *Plugins {
	plugins := Plugins{}
	plugins.loadHelloPlugin(server)

	return &plugins
}
