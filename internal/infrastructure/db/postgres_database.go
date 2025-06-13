package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close() 
		return nil, fmt.Errorf("error en la conexión a la base de datos: %w", err)
	}

	log.Println("✅ Conectado a PostgreSQL")
	return db, nil
}

