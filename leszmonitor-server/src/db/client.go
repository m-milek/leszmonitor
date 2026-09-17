package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/m-milek/leszmonitor/platform/audit"
	platformdb "github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

var ErrNotFound = platformdb.ErrNotFound
var ErrAlreadyExists = platformdb.ErrAlreadyExists

// DB defines the database access surface. It returns DAO interfaces for easy mocking.
type DB interface {
	Users() IUserDAO
	Monitors() IMonitorDAO
	MonitorResults() IMonitorResultDAO
	MonitorStatusChanges() IMonitorStatusChangeDAO
	MonitorStats() IMonitorStatsDAO
	AuditLog() audit.IAuditLogDAO
	Querier() platformdb.Querier
	WithTx(ctx context.Context, fn func(q platformdb.Querier) error) error
	Close()
}

// Client implements DB on top of a platform DB client.
type Client struct {
	root *platformdb.Client
	pool platformdb.Querier
	// cached DAOs to avoid re-allocation on every getter call
	users                IUserDAO
	monitors             IMonitorDAO
	monitorResults       IMonitorResultDAO
	monitorStatusChanges IMonitorStatusChangeDAO
	monitorStats         IMonitorStatsDAO
	auditLog             audit.IAuditLogDAO
}

type baseDAO struct {
	pool platformdb.Querier
}

func newBaseDAO(pool platformdb.Querier) baseDAO {
	return baseDAO{pool: pool}
}

// newClient creates a Client wired to the given Querier.
// root is nil for transaction-scoped clients.
func newClient(root *platformdb.Client, pool platformdb.Querier) *Client {
	base := newBaseDAO(pool)
	return &Client{
		root:                 root,
		pool:                 pool,
		users:                newUserDAO(base),
		monitors:             newMonitorDAO(base),
		monitorResults:       newMonitorResultDAO(base),
		monitorStatusChanges: newMonitorStatusChangeDAO(base),
		monitorStats:         newMonitorStatsDAO(base),
		auditLog:             audit.NewAuditLogDAO(pool),
	}
}

// New creates a new DB client using the provided DSN.
func New(ctx context.Context, dsn string) (*Client, error) {
	root, err := platformdb.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return newClient(root, root.Querier()), nil
}

// Querier returns the Querier the client is bound to.
func (c *Client) Querier() platformdb.Querier {
	return c.pool
}

// WithTx executes fn inside a database transaction. Nested transactions are not supported.
func (c *Client) WithTx(ctx context.Context, fn func(q platformdb.Querier) error) error {
	if c.root == nil {
		return fmt.Errorf("nested transactions are not supported")
	}
	return c.root.WithTx(ctx, fn)
}

// Close closes the underlying connection pool. No-op for transaction-scoped clients.
func (c *Client) Close() {
	if c.root != nil {
		c.root.Close()
	}
}

// DAO getters (return interfaces for mocking)

func (c *Client) Users() IUserDAO                               { return c.users }
func (c *Client) Monitors() IMonitorDAO                         { return c.monitors }
func (c *Client) MonitorResults() IMonitorResultDAO             { return c.monitorResults }
func (c *Client) MonitorStatusChanges() IMonitorStatusChangeDAO { return c.monitorStatusChanges }
func (c *Client) MonitorStats() IMonitorStatsDAO                { return c.monitorStats }
func (c *Client) AuditLog() audit.IAuditLogDAO                  { return c.auditLog }

// WithAuditedTx executes fn inside a database transaction and, if successful, records an audit log.
// It returns the result from fn and any error that occurred.
func WithAuditedTx[T any](
	ctx context.Context,
	client DB,
	fn func(tx DB) (T, *audit.AuditLogParams, error),
) (T, error) {
	return audit.WithAuditedTx(ctx, client, func(q platformdb.Querier) (T, *audit.AuditLogParams, error) {
		return fn(newClient(nil, q))
	})
}

// WithAuditedVoidTx executes fn inside a database transaction and, if successful, records an audit log.
// This is a convenience wrapper for operations that do not return a result.
func WithAuditedVoidTx(
	ctx context.Context,
	client DB,
	fn func(tx DB) (*audit.AuditLogParams, error),
) error {
	return audit.WithAuditedVoidTx(ctx, client, func(q platformdb.Querier) (*audit.AuditLogParams, error) {
		return fn(newClient(nil, q))
	})
}

// --------------------------
// Singleton management (unexported global within the db package for convenience)
// --------------------------.
var (
	instance DB
	instMu   sync.RWMutex
)

// Get returns the current DB singleton (maybe nil if not initialized).
func Get() DB {
	instMu.RLock()
	defer instMu.RUnlock()
	return instance
}

// Set sets the DB singleton. Useful for tests to inject a mock.
func Set(db DB) {
	instMu.Lock()
	defer instMu.Unlock()
	if instance != nil {
		// Close previous instance if it was a real client
		instance.Close()
	}
	instance = db
}

// InitFromEnv initializes the DB singleton using the DSN from environment.
func InitFromEnv(ctx context.Context) error {
	logger := log.FromContext(ctx)

	root, err := platformdb.NewFromEnv(ctx)
	if err != nil {
		return err
	}

	Set(newClient(root, root.Querier()))
	logger.Info().Msg("SQLite connection established.")
	return nil
}
