package routes

import (
	"net/http"

	"github.com/AmirSolt/town-watch/models"
	"github.com/AmirSolt/town-watch/server"
	"github.com/gin-gonic/gin"
)

type walletLoad struct {
	pageLoad
	TierDisplays *[]TierDisplay
}

type TierDisplay struct {
	TierConfig       server.TierConfig
	UserSubscription *UserSubscription
}

type UserSubscription struct {
	Tier                     models.Tier
	PaymentMethodLast4Digits string
	IsAuto                   bool
}

func (routes *Routes) wallet() {

	userSubscription := UserSubscription{
		Tier:                     models.TierT2,
		PaymentMethodLast4Digits: "1234",
		IsAuto:                   true,
	}

	var tierDisplays []TierDisplay
	for tier, tierConfig := range routes.server.TierConfigs {
		if tier == models.TierT0 {
			continue
		}

		tierDisplays = append(tierDisplays, TierDisplay{
			TierConfig:       tierConfig,
			UserSubscription: &userSubscription,
		})
	}

	routes.server.Engine.GET("/wallet", func(c *gin.Context) {

		c.HTML(http.StatusOK, "wallet.tmpl", gin.H{
			"data": walletLoad{
				pageLoad: pageLoad{
					Title: "Wallet",
				},
				TierDisplays: &tierDisplays,
			},
		})

	})
}
