package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type joinLoad struct {
	pageLoad
}

func (routes *Routes) join() {

	routes.server.Engine.GET("/join", func(c *gin.Context) {

		c.HTML(http.StatusOK, "join.tmpl", gin.H{
			"data": joinLoad{
				pageLoad: pageLoad{
					Title: "Join",
				},
			},
		})

	})
}

type verifyLoad struct {
	pageLoad
}

func (routes *Routes) joinVerify() {

	routes.server.Engine.GET("/join/verify", func(c *gin.Context) {

		c.HTML(http.StatusOK, "verify.tmpl", gin.H{
			"data": verifyLoad{
				pageLoad: pageLoad{
					Title: "Verify",
				},
			},
		})

	})
}
