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

	server.DB = (*DB)(dbEngine)
}

func (server *Server) killDB() {
	server.DB.Close()
}
