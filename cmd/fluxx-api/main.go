package main

import (
	"log"
	"net/http"

	"fluxx/internal/api"
	"fluxx/internal/websocket"
)

func main() {
	// 1. Initialiser le Hub WebSocket
	// C'est la structure centrale qui gère toutes les connexions temps réel.
	hub := websocket.NewHub()

	// Démarrer la fonction Run() dans une goroutine.
	// Le Hub tourne en parallèle pour traiter les messages sans bloquer le serveur HTTP.
	go hub.Run()

	// 2. Initialiser le Routeur (gestion des URLs)
	// On passe l'instance du Hub au routeur pour que les handlers puissent y accéder.
	r := api.NewRouter(hub)

	// 3. Démarrer le serveur HTTP
	log.Println("Fluxx API starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
