package server

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5"
)

func (server *Server) loadDB() {
	dbEngine, dbErr := sql.Open("postgres", server.Env.DATABASE_URL)
	if dbErr != nil {
		log.Fatalln("Error db:", dbErr)
	}

	if err := dbEngine.Ping(); err != nil {
		log.Fatalln("Error db:", err)
	}

	server.DB = dbEngine
}

func (server *Server) killDB() {
	server.DB.Close()
}
