package gomigrations

import (
	"github.com/arsura/goose"
)

func init() {
	goose.AddMigrationNoTxContext(nil, nil)
}
