package main

import (
	"log"
	"net/http"
	"os"

	"fluxx/internal/api"
	"fluxx/internal/database" // 👈 Nouvel import
	"fluxx/internal/store"
	"fluxx/internal/websocket"

	"github.com/joho/godotenv"
)

func main() {
	// --- Étape 0 : Charger le fichier .env (Local uniquement) ---
	_ = godotenv.Load()

	// --- Étape 1 : Connexion technique (database) ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Variable DATABASE_URL manquante.")
	}

	// On utilise le package database pour créer la connexion brute
	dbConn, err := database.NewPostgres(dbURL)
	if err != nil {
		log.Fatalf("Impossible de se connecter à Postgres: %v", err)
	}

	// --- Étape 2 : Store métier ---
	// On injecte la connexion brute dans ton Store
	s := store.NewStore(dbConn)
	log.Println("✅ Connexion Supabase via Store réussie !")
	defer s.Close()

	// --- Étape 3 : Hub WebSocket ---
	// ICI : On donne le store au Hub pour qu'il puisse sauver les messages
	hub := websocket.NewHub(s)
	go hub.Run()

	// --- Étape 4 : Router & Serveur ---
	r := api.NewRouter(hub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Fluxx API lancée sur le port %s", port)

	// Lancement avec support CORS pour ton HTML
	if err := http.ListenAndServe(":"+port, enableCORS(r)); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Middleware CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
