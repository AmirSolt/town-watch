package server

import (
	"database/sql"
	"log"

	"github.com/AmirSolt/town-watch/models"
	_ "github.com/jackc/pgx/v5"
)

type DB struct {
	queries    *models.Queries
	connection *sql.DB
}

func (server *Server) loadDB() {
	conn, dbErr := sql.Open("postgres", server.Env.DATABASE_URL)
	if dbErr != nil {
		log.Fatalln("Error db:", dbErr)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalln("Error db:", err)
	}

	queries := models.New(conn)

	server.DB.connection = conn
	server.DB.queries = queries
}

func (server *Server) killDB() {
	server.DB.connection.Close()
}
