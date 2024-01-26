package plugins

import "github.com/AmirSolt/town-watch/server"

type Plugins struct {
	Env *Env
}

func LoadPlugins(server *server.Server) *Plugins {
	plugins := Plugins{}

	return &plugins
}
