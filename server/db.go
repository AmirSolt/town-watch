package server

import (
	"context"
	"log"

	"github.com/AmirSolt/town-watch/models"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	Queries *models.Queries
	conn    *pgx.Conn
}

func (server *Server) loadDB() {
	conn, dbErr := pgx.Connect(context.Background(), server.Env.DATABASE_URL)
	if dbErr != nil {
		log.Fatalln("Error db:", dbErr)
	}

	queries := models.New(conn)
	server.DB = &DB{
		conn:    conn,
		Queries: queries,
	}
}

func (server *Server) killDB() {
	server.DB.conn.Close(context.Background())
}
