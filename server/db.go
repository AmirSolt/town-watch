package server

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB sqlx.DB

func (server *Server) loadDB() {
	dbEngine, dbErr := sqlx.Connect("postgres", server.Env.DATABASE_URL)

	if dbErr != nil {
		log.Fatalln("Error db:", dbErr)
	}

	// initQueriues := []string{
	// 	models.ExtSchema,
	// 	models.EnumSchema,
	// 	models.UserSchema,
	// 	models.ScannerSchema,
	// 	models.NotifSchema,
	// 	models.ReportSchema,
	// 	models.ReportNotifSchema,
	// }

	// init := strings.Join(initQueriues, "\n")
	// _, err := dbEngine.Exec(init)
	// if err != nil {
	// 	log.Fatalln("Error db:", err)
	// }

	server.DB = (*DB)(dbEngine)
}

func (server *Server) killDB() {
	server.DB.Close()
}
