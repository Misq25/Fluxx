package main

import (
	"log"
	"net/http"
	"os"

	"fluxx/internal/api"
	"fluxx/internal/store"
	"fluxx/internal/websocket"

	"github.com/joho/godotenv" // 👈 AJOUTE CET IMPORT
)

func main() {
	// --- Étape 0 : Charger le fichier .env ---
	// Cette ligne lit le fichier .env à la racine et injecte les variables
	err := godotenv.Load()
	if err != nil {
		// On met un Println et pas un Fatal car sur Render (en production),
		// le fichier .env n'existe pas, les variables sont injectées directement.
		log.Println("Note: Aucun fichier .env trouvé, utilisation des variables d'environnement système.")
	}

	// --- Étape 1 : Connexion à Supabase (PostgreSQL) ---

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Variable d'environnement DATABASE_URL manquante. Le serveur ne peut pas se connecter à la BDD.")
	}

	// 1.2. Établir la connexion à la base de données via le Store.
	s, err := store.NewStore(dbURL)
	if err != nil {
		log.Fatalf("Impossible d'initialiser la connexion à Supabase: %v", err)
	}
	log.Println("Connexion à la base de données Supabase réussie !") // Petit message de confort

	defer s.Close()

	// --- Étape 2 : Initialisation du Hub WebSocket ---
	hub := websocket.NewHub(s)
	go hub.Run()

	// --- Étape 3 : Démarrage du Serveur HTTP ---
	r := api.NewRouter(hub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Fluxx API starting on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
