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

	inMemoryLock     sync.RWMutex
	inMemoryClients  map[string]*api.Client
	inMemorySessions map[string]*api.Session
}

func NewDatabase(dbProperties *config.DatabaseOptions) (*Database, error) {
	db, err := sqlutil.Open(dbProperties)
	if err != nil {
		return nil, err
	}

	d := &Database{
		db:               db,
		writer:           sqlutil.NewExclusiveWriter(),
		inMemoryClients:  make(map[string]*api.Client),
		inMemorySessions: make(map[string]*api.Session),
	}

	return d, nil
}

func (d *Database) GetClient(
	ctx context.Context, id string,
) (*api.Client, error) {
	d.inMemoryLock.RLock()
	defer d.inMemoryLock.RUnlock()

	client, ok := d.inMemoryClients[id]
	if !ok {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (d *Database) CreateClient(
	ctx context.Context, client *api.Client,
) error {
	d.inMemoryLock.Lock()
	defer d.inMemoryLock.Unlock()

	if _, ok := d.inMemoryClients[client.GetID()]; ok {
		return errors.New("client exists")
	}

	d.inMemoryClients[client.GetID()] = client
	return nil
}

func (d *Database) UpdateClient(
	ctx context.Context, client *api.Client,
) error {
	d.inMemoryLock.Lock()
	defer d.inMemoryLock.Unlock()

	if _, ok := d.inMemoryClients[client.GetID()]; !ok {
		return errors.New("client not found")
	}

	d.inMemoryClients[client.GetID()] = client
	return nil
}

func (d *Database) GetSession(
	ctx context.Context, id string,
) (*api.Session, error) {
	d.inMemoryLock.RLock()
	defer d.inMemoryLock.RUnlock()

	session, ok := d.inMemorySessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (d *Database) CreateSession(
	ctx context.Context, id string, session *api.Session,
) error {
	d.inMemoryLock.Lock()
	defer d.inMemoryLock.Unlock()

	if _, ok := d.inMemorySessions[id]; ok {
		return errors.New("session exists")
	}

	d.inMemorySessions[id] = session
	return nil
}
