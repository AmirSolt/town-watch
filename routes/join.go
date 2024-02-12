package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type joinLoad struct {
	pageLoad
}

func (routes *Routes) joinRoutes() {
	routes.join()
	routes.joinVerify()
	routes.testJoin()
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

func (routes *Routes) testJoin() {
	if !routes.server.Env.IS_PROD {
		routes.server.Engine.POST("/join/test/singin", func(c *gin.Context) {

			c.HTML(http.StatusOK, "verify.tmpl", gin.H{
				"data": verifyLoad{
					pageLoad: pageLoad{
						Title: "Verify",
					},
				},
			})

		})

		routes.server.Engine.POST("/join/test/singout", func(c *gin.Context) {

			c.HTML(http.StatusOK, "verify.tmpl", gin.H{
				"data": verifyLoad{
					pageLoad: pageLoad{
						Title: "Verify",
					},
				},
			})

		})
	}
}
