package server

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB sqlx.DB

const schema string = ``

func (server *Server) loadDB() {
	dbEngine := sqlx.MustConnect("postgres", server.Env.DATABASE_URL)

	if err := dbEngine.Ping(); err != nil {
		log.Fatalln("Error db:", err)
	}

	dbEngine.MustExec(schema)

	server.DB = (*DB)(dbEngine)
}

func (server *Server) killDB() {
	server.DB.Close()
}
