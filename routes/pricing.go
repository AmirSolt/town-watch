package routes

import (
	"net/http"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
)

type pricingLoad struct {
	pageLoad
	Tier1 models.Tier
	Tier2 models.Tier
}

func (routes *Routes) pricing() {

	routes.server.Engine.GET("/pricing", func(c *gin.Context) {

		c.HTML(http.StatusOK, "pricing.tmpl", gin.H{
			"data": pricingLoad{
				pageLoad: pageLoad{
					Title: "Pricing",
				},
				Tier1: models.TierT1,
				Tier2: models.TierT2,
			},
		})

	})
}

func (routes *Routes) checkout() {

	routes.server.Engine.GET("/checkout/:tier", func(c *gin.Context) {
		c.String(200, c.Param("tier"))
	})
}
