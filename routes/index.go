package routes

import (
	"net/http"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
)

type indexLoad struct {
	pageLoad
	Scanners *[]models.Scanner
}

func (routes *Routes) index() {

	routes.server.Engine.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"data": indexLoad{
				pageLoad: pageLoad{
					Title: "Home",
				},
				Scanners: &[]models.Scanner{
					{
						ID:      123,
						Address: "my address",
					},
					{
						ID:      125,
						Address: "my addqwd12s",
					},
					{
						ID:      122,
						Address: "my adasdwqdess",
					}},
			},
		})

	})
}
