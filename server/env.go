package server

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	IS_PROD string `validate:"boolean"`
}

func (server *Server) loadEnv() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	env := Env{
		IS_PROD: os.Getenv("IS_PROD"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(env)
	if err != nil {
		log.Fatal("Error a variable is missing from .env")
	}

	server.Env = &env
}
