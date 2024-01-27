package server

type DB struct {
}

const schema string = ``

func (server *Server) loadDB() {
	// db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// db.MustExec(schema)
}
