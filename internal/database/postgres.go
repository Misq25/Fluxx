package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgres crée une connexion brute à la base de données
func NewPostgres(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Connexion brute à PostgreSQL établie")
	return db, nil
}
