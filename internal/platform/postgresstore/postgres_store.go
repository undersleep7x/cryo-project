package postgresstore

import (
	"context"
	"database/sql"
)

type pgClientImpl struct {
	db *sql.DB
}

// New returns a new Store that implements PostgresStore interface.
func NewPgClientWrapper(db *sql.DB) *pgClientImpl {
	return &pgClientImpl{db: db}
}

func (p *pgClientImpl) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *pgClientImpl) GetDB() *sql.DB {
	return p.db
}

func (p *pgClientImpl) Close() error {
	return p.db.Close()
}
