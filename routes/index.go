package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
)

type indexLoad struct {
	pageLoad
	Scanners *[]models.Scanner
}

func (routes *Routes) index() {

	routes.server.Engine.GET("/",
		routes.server.ValidateUser,
		func(c *gin.Context) {

			tempUser, _ := c.Get("user")
			user := (tempUser).(*models.User)

			scanners, err := routes.server.DB.Queries.GetScanners(context.Background(), user.ID)
			if err != nil && err != sql.ErrNoRows {
				fmt.Errorf("scanner lookup failed: %w", err)
				return
			}

			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"data": indexLoad{
					pageLoad: pageLoad{
						Title: "Home",
						User:  user,
					},
					Scanners: &scanners,
				},
			})

		})
}
