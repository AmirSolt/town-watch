package routes

import (
	"net/http"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
)

type notifLoad struct {
	pageLoad
	Scanner *models.Scanner
	Reports *[]models.Report
}

func (routes *Routes) notif() {

	routes.server.Engine.GET("/notif/:id", func(c *gin.Context) {

		c.HTML(http.StatusOK, "notif.tmpl", gin.H{
			"data": notifLoad{
				pageLoad: pageLoad{
					Title: "Notif",
				},
				Scanner: &models.Scanner{
					ID:       123,
					Address:  "my address",
					IsActive: true,
				},
				Reports: &[]models.Report{
					{
						ID: 13,
					},
				},
			},
		})

	})
}
