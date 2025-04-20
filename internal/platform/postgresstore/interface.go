package postgresstore

import (
	"context"
	"database/sql"
)

type PostgresClient interface {
    Ping(ctx context.Context) error
	GetDB() *sql.DB
    Close() error
}