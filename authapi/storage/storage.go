package storage

import (
	"fmt"

	"github.com/matrix-org/dendrite/authapi/storage/sqlite3"
	"github.com/matrix-org/dendrite/internal/config"
)

// NewDatabase opens a new Postgres or Sqlite database (based on dataSourceName scheme)
// and sets postgres connection parameters
func NewDatabase(dbProperties *config.DatabaseOptions) (Database, error) {
	switch {
	case dbProperties.ConnectionString.IsSQLite():
		return sqlite3.NewDatabase(dbProperties)
	case dbProperties.ConnectionString.IsPostgres():
		return nil, fmt.Errorf("unimplemented")
	default:
		return nil, fmt.Errorf("unexpected database type")
	}
}
