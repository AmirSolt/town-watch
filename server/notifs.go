package server

import (
	"context"
	"fmt"
	"time"

	"github.com/AmirSolt/town-watch/models"
	"github.com/jackc/pgx/v5/pgtype"
)

const scannerNotifScanReportsLimit = 20
const scannerNotifFromSeconds = 60 * 60 // 60 minutes

func (server *Server) GetNotifs(currentTime time.Time) (*[]models.Notif, error) {
	notifs, err := server.DB.queries.CreateScannerNotifs(context.Background(), models.CreateScannerNotifsParams{
		FromDate:         pgtype.Timestamptz{Time: currentTime.Add(-time.Second * scannerNotifFromSeconds)},
		ToDate:           pgtype.Timestamptz{Time: currentTime},
		ScanReportsLimit: scannerNotifScanReportsLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("error ScanReports: %w", err)
	}

	return &notifs.(*[]models.Notif), nil
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

	var failEmailCount int = 0
	for _, user := range users {
		err := server.SendEmail(user.Email, "User", "Town Watch", "Reports Found Notification", "content")
		if err != nil {
			failEmailCount++
			// log error
			fmt.Errorf("error failed to send email: %w", err)
		}
	}

	if failEmailCount >= len(users) {
		return fmt.Errorf(">>> error all notification emails failed: %w", err)
	}

	return nil
}
