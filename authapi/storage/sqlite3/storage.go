package sqlite3

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/internal/config"
	"github.com/matrix-org/dendrite/internal/sqlutil"
)

// Database represents an account database
type Database struct {
	db     *sql.DB
	writer sqlutil.Writer

	sqlutil.PartitionOffsetStatements

	clientMu sync.Mutex

	inMemoryClients map[string]*api.Client
}

func NewDatabase(dbProperties *config.DatabaseOptions) (*Database, error) {
	db, err := sqlutil.Open(dbProperties)
	if err != nil {
		return nil, err
	}

	d := &Database{
		db:              db,
		writer:          sqlutil.NewExclusiveWriter(),
		inMemoryClients: make(map[string]*api.Client),
	}

	return d, nil
}

func (d *Database) GetClient(
	ctx context.Context, id string,
) (*api.Client, error) {
	client, ok := d.inMemoryClients[id]
	if !ok {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (d *Database) CreateClient(
	ctx context.Context, client *api.Client,
) error {
	if _, ok := d.inMemoryClients[client.GetID()]; ok {
		return errors.New("client exists")
	}

	d.inMemoryClients[client.GetID()] = client
	return nil
}

func (d *Database) UpdateClient(
	ctx context.Context, client *api.Client,
) error {
	if _, ok := d.inMemoryClients[client.GetID()]; !ok {
		return errors.New("client not found")
	}

	d.inMemoryClients[client.GetID()] = client
	return nil
}
