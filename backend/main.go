package main

import (
	"log"
	"net/http"

	"github.com/Misq25/Fluxx/backend/config"
	"github.com/Misq25/Fluxx/backend/database"
)

func main() {
	config.LoadEnvVariables()

	database.ConnectDB()

	log.Println("Le serveur démarre sur le port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur", err)
	}
}
