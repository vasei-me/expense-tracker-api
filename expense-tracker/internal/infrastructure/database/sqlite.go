package database

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConfig struct {
	DBPath string
}

func NewSQLiteDB(cfg SQLiteConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite doesn't support multiple writers
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to SQLite database:", cfg.DBPath)
	return db, nil
}