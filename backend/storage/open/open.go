package open

import (
	"database/sql"

	"github.com/rahulkumarpahwa/traer/config"
	_ "modernc.org/sqlite"
)


type OpenDB struct {
	Config *config.Config
}

func (o *OpenDB) Open() (*sql.DB, error) {
	sql.Open("sqlite", o.Config.Database.Path)

	return nil, nil
}