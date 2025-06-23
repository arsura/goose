package gomigrations

import (
	"context"
	"database/sql"

	"github.com/arsura/goose"
)

func init() {
	goose.AddMigrationNoTxContext(up014, nil)
}

func up014(ctx context.Context, db *sql.DB) error {
	return createTable(db, "hotel")
}
