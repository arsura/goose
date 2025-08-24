package gomigrations

import (
	"database/sql"

	"github.com/arsura/goose"
)

func init() {
	goose.AddMigrationNoTx(up006, nil)
}

func up006(db *sql.DB) error {
	return createTable(db, "delta")
}
