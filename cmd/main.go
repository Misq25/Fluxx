package main

import (
	"log"
	"net/http"
	"os"

	"fluxx/internal/api"
	"fluxx/internal/store"
	"fluxx/internal/websocket"

	"github.com/joho/godotenv"
)

func main() {
	// --- Étape 0 : Charger le fichier .env (Local uniquement) ---
	_ = godotenv.Load() // On ignore l'erreur car Render gère ça en interne

	// --- Étape 1 : Connexion à Supabase ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Variable DATABASE_URL manquante.")
	}

	s, err := store.NewStore(dbURL)
	if err != nil {
		log.Fatalf("Erreur connexion BDD: %v", err)
	}
	log.Println("✅ Connexion Supabase réussie !")
	defer s.Close()

	// --- Étape 2 : Initialisation du Hub ---
	hub := websocket.NewHub(s)
	go hub.Run()

	// --- Étape 3 : Router & Middlewares ---
	r := api.NewRouter(hub)

	// --- Étape 4 : Lancement ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Fluxx API lancée sur le port %s", port)

	// 🚨 IMPORTANT : On enveloppe le router 'r' avec enableCORS
	if err := http.ListenAndServe(":"+port, enableCORS(r)); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// enableCORS est le garde du corps qui autorise ton HTML à parler au serveur
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// On autorise toutes les origines pour le moment (plus simple pour le dev)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Si c'est une requête de pré-vérification (OPTIONS), on répond OK tout de suite
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
