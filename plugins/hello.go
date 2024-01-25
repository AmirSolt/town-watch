package plugins

import (
	"fmt"

	"github.com/AmirSolt/town-watch/server"
)

func (plugins *Plugins) loadHelloPlugin(server *server.Server) {
	fmt.Println("THIS IS THE HELLO PLUGIN")
}
