package register

import (
	"database/sql"

	"github.com/arsura/goose"
)

func init() {
	goose.AddMigration(
		func(_ *sql.Tx) error { return nil },
		func(_ *sql.Tx) error { return nil },
	)
}
