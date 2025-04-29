package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Misq25/Fluxx/backend/config"
	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	var err error
	DB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données.", err)
	}

	log.Println("Connxion à la base de données établie avec succès.")
}
