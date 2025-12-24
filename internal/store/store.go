package store

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Pilote PostgreSQL
)

type Store struct {
	db *sqlx.DB
}

func NewStore(dbURL string) (*Store, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Printf("Erreur de connexion à la BDD: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("Ping de la BDD échoué: %v", err)
		return nil, err
	}
	log.Println("Connexion à la base de données Supabase réussie.")

	return &Store{db: db}, nil
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
		log.Println("Connexion à la base de données fermée.")
	}
}

// --- PARTIE AUTH / PROFIL ---

// SyncUserProfile utilise ta struct Profile pour créer/modifier un utilisateur
func (s *Store) SyncUserProfile(p Profile) error {
	// Avec sqlx, on utilise ":" devant les noms pour qu'il mappe tout seul avec la struct
	query := `
		INSERT INTO profiles (id, username, display_name, avatar_url)
		VALUES (:id, :username, :display_name, :avatar_url)
		ON CONFLICT (id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = now();
	`

	// NamedExec est magique : il lit les tags `db:"..."` de ton models.go
	_, err := s.db.NamedExec(query, p)
	if err != nil {
		log.Printf("Erreur lors de la synchro du profil %s: %v", p.Username, err)
	}
	return err
}

// --- PARTIE MESSAGES ---

func (s *Store) SaveMessage(clientID string, content string) error {
	query := `INSERT INTO messages (client_id, content) VALUES ($1, $2)`
	_, err := s.db.Exec(query, clientID, content)
	if err != nil {
		log.Printf("Erreur SQL lors de l'insertion du message: %v", err)
	}
	return err
}
