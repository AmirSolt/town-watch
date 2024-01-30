package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	IS_PROD                 string `validate:"boolean"`
	DATABASE_URL            string `validate:"url"`
	ARCGIS_TORONTO_URL      string `validate:"url"`
	EMAIL_CF_WORKER_URL     string `validate:"url"`
	EMAIL_CF_WORKER_API_KEY string `validate:"required"`
}

func (server *Server) loadEnv() {

	if err := godotenv.Load(filepath.Join(server.RootDir, ".env")); err != nil {
		log.Fatal("Error .env:", err)
	}

	env := Env{
		IS_PROD:                 os.Getenv("IS_PROD"),
		DATABASE_URL:            os.Getenv("DATABASE_URL"),
		ARCGIS_TORONTO_URL:      os.Getenv("ARCGIS_TORONTO_URL"),
		EMAIL_CF_WORKER_URL:     os.Getenv("EMAIL_CF_WORKER_URL"),
		EMAIL_CF_WORKER_API_KEY: os.Getenv("EMAIL_CF_WORKER_API_KEY"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(env)
	if err != nil {
		log.Fatal("Error .env:", err)
	}

	server.Env = &env
}
