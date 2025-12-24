package store

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Pilote PostgreSQL
)

type Store struct {
	db *sqlx.DB
}

// NewStore reçoit désormais un *sqlx.DB déjà connecté
// C'est ce qu'on appelle l'injection de dépendance
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
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
	query := `
		INSERT INTO profiles (id, username, display_name, avatar_url)
		VALUES (:id, :username, :display_name, :avatar_url)
		ON CONFLICT (id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = now();
	`

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
