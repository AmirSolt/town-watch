package server

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB sqlx.DB

const schema string = `CREATE TABLE place (
	id SERIAL
    country text,
    city text NULL,
    telcode integer);`

func (server *Server) loadDB() {
	dbEngine := sqlx.MustConnect("postgres", server.Env.DATABASE_URL)

	if err := dbEngine.Ping(); err != nil {
		log.Fatalln(err)
	}

	dbEngine.MustExec(schema)

	server.DB = (*DB)(dbEngine)
}

func (server *Server) killDB() {
	server.DB.Close()
}
