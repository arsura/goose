package gomigrations

import (
	"context"
	"database/sql"

	"github.com/arsura/goose"
)

func init() {
	goose.AddMigrationNoTxContext(nil, down015)
}

func down015(ctx context.Context, db *sql.DB) error {
	return dropTable(db, "hotel")
}
