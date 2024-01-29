package server

import (
	"time"

	"github.com/AmirSolt/town-watch/models"
)

func (server *Server) getNotifs(currentTime time.Time) *[]models.Notif {
	return &[]models.Notif{}
}

func (server *Server) sendNotifs(notifs *[]models.Notif) {

}
