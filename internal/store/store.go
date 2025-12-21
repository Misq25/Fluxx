package store

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Pilote PostgreSQL
)

// Store est la structure qui maintient la connexion à la BDD ouverte
type Store struct {
	db *sqlx.DB
}

// NewStore est la fonction qui établit la connexion
func NewStore(dbURL string) (*Store, error) {
	// sqlx.Connect ouvre la connexion en utilisant le driver "postgres" (lib/pq)
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Printf("Erreur de connexion à la BDD: %v", err)
		return nil, err
	}

	// S'assurer que la connexion est bien vivante
	if err = db.Ping(); err != nil {
		log.Printf("Ping de la BDD échoué: %v", err)
		return nil, err
	}
	log.Println("Connexion à la base de données Supabase réussie.")

	return &Store{db: db}, nil
}

// 🚨 NOUVELLE MÉTHODE : Close ferme la connexion à la base de données.
func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
		log.Println("Connexion à la base de données fermée.")
	}
}

// 🚨 NOUVELLE MÉTHODE : SaveMessage enregistre un nouveau message dans la BDD.
func (s *Store) SaveMessage(clientID string, content string) error {
	// Requête SQL pour insérer les données. Supabase (PostgreSQL) gérera les IDs et timestamps.
	query := `INSERT INTO messages (client_id, content) VALUES ($1, $2)`

	_, err := s.db.Exec(query, clientID, content)

	if err != nil {
		log.Printf("Erreur SQL lors de l'insertion: %v", err)
	}

	return err
}
