package database

import (
	"expense-tracker/internal/config"

	"github.com/jmoiron/sqlx"
)

func Connect(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	switch cfg.Type {
	case "postgres":
		pgConfig := PostgresConfig{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			DBName:   cfg.DBName,
			SSLMode:  cfg.SSLMode,
		}
		return NewPostgresDB(pgConfig)
	case "sqlite":
		sqliteConfig := SQLiteConfig{
			DBPath: cfg.DBName,
		}
		return NewSQLiteDB(sqliteConfig)
	default:
		// Default to SQLite
		sqliteConfig := SQLiteConfig{
			DBPath: "expense_tracker.db",
		}
		return NewSQLiteDB(sqliteConfig)
	}
}