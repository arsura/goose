package gomigrations

import (
	"github.com/arsura/goose"
)

func init() {
	goose.AddMigrationNoTx(nil, nil)
}
