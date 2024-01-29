package server

import (
	"log"
	"strings"

	"github.com/AmirSolt/town-watch/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB sqlx.DB

func (server *Server) loadDB() {
	dbEngine := sqlx.MustConnect("postgres", server.Env.DATABASE_URL)

	if err := dbEngine.Ping(); err != nil {
		log.Fatalln("Error db:", err)
	}

	initQueriues := []string{
		models.EnumSchema,
		models.NotifSchema,
		models.ReportNotifSchema,
		models.ReportSchema,
		models.ScannerSchema,
		models.UserSchema,
	}

	init := strings.Join(initQueriues, "\n")
	dbEngine.MustExec(init)

	server.DB = (*DB)(dbEngine)
}

func (server *Server) killDB() {
	server.DB.Close()
}
