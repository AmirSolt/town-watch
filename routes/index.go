package routes

import (
	"net/http"

	"github.com/AmirSolt/town-watch/server"
	"github.com/gin-gonic/gin"
)

func (routes *Routes) index(server *server.Server) {

	server.Engine.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", nil)
	})

}
