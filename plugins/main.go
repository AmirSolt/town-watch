package plugins

import "github.com/AmirSolt/town-watch/server"

type Plugins struct {
	Env      *Env
	Reporter *Reporter
}

func LoadPlugins(server *server.Server) *Plugins {
	plugins := Plugins{}

	plugins.Env.loadEnv(&plugins)

	return &plugins
}
