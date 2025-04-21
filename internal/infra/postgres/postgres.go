package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/undersleep7x/cryo-project/internal/config"
	_ "github.com/lib/pq" // postgres driver
)

// NewPostgresClient sets up a raw *sql.DB instance
func NewPostgresClient(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	log.Println("Postgres connected successfully")
	return db, nil
}
