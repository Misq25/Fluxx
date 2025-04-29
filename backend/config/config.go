package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur de chargement du fichier .env")
	}

	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
}
