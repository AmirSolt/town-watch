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

func (routes *Routes) indexRoutes() {
	routes.index()
}

func (routes *Routes) index() {

	routes.server.Engine.GET("/",
		routes.server.OptionalUserMiddleware,
		func(c *gin.Context) {

			tempUser, exists := c.Get("user")
			if !exists {
				fmt.Errorf("'user' key/value doesn't exist")
			}
			user := (tempUser).(*models.User)

			fmt.Println("===========================")
			fmt.Println("'user':%w", user == nil)
			fmt.Println("===========================")

			var scanners []models.Scanner
			var err error

			if user != nil {
				scanners, err = routes.server.DB.Queries.GetScanners(context.Background(), user.ID)
				if err != nil && err != sql.ErrNoRows {
					// fmt.Errorf("scanner lookup failed: %w", err)
					return
				}
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
