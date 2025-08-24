package register

import (
	"database/sql"

	"github.com/arsura/goose"
)

func init() {
	goose.AddMigrationNoTx(
		func(_ *sql.DB) error { return nil },
		func(_ *sql.DB) error { return nil },
	)
}
