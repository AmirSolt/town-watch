package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (routes *Routes) index() {

	routes.server.Engine.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", nil)
	})

}
