package server

import (
	"context"
	"fmt"
	"time"

	"github.com/AmirSolt/town-watch/models"
)

func (server *Server) GetNotifs(currentTime time.Time) *[]models.Notif {
	return &[]models.Notif{}
}

func (server *Server) SendNotifs(notifs *[]models.Notif) error {

	var userIds []int32
	for _, notif := range *notifs {
		userIds = append(userIds, notif.UserID)
	}

	users, err := server.DB.queries.GetUsers(context.Background(), userIds)

	if err != nil {
		return fmt.Errorf("failed to fetch users from GetUsers(): %w", err)
	}

	for _, user := range users {
		err := server.SendEmail(user.Email, "User", "Town Watch", "Reports Found Notification", "content")
		if err != nil {

		}
	}

	return nil
}
