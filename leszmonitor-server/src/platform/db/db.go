package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	config "github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/log"

	// Blank import to initialize the SQLite driver.
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var dbSchema embed.FS

var ErrNotFound = errors.New("resource not found")
var ErrAlreadyExists = errors.New("resource already exists")

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// SQLite reports unique constraint errors as text like: "UNIQUE constraint failed: table.column".
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint failed")
}

const timeoutDuration = 1000 * time.Second

// Querier is the query surface shared by *sqlx.DB and *sqlx.Tx.
type Querier = sqlx.ExtContext

// DB is the database handle used by services.
type DB interface {
	Querier() Querier
	WithTx(ctx context.Context, fn func(q Querier) error) error
}

// Client owns the SQLite connection pool.
type Client struct {
	sqlxDB *sqlx.DB
}

type dbResult[T any] struct {
	Duration time.Duration
	Result   T
}

// New creates a new DB client using the provided DSN.
func New(ctx context.Context, dsn string) (*Client, error) {
	pool, err := sqlx.ConnectContext(ctx, "sqlite", dsn)
	if err != nil {
		return nil, err
	}

	c := &Client{sqlxDB: pool}

	if err := c.initSchema(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	pool.SetMaxOpenConns(1)

	return c, nil
}

// NewFromEnv creates a new DB client using the DSN from environment.
func NewFromEnv(ctx context.Context) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logger := log.FromContext(ctx)

	logger.Info().Msg("Connecting to SQLite...")

	uri := os.Getenv(config.SqliteDBPath)
	if uri == "" {
		logger.Fatal().Msg("SQLite DB path is not defined")
	}
	uri = setupSQLiteConfig(uri)
	return New(ctx, uri)
}

// Querier returns the connection pool as a Querier.
func (c *Client) Querier() Querier {
	return c.sqlxDB
}

// initSchema reads the database schema from a file and executes it to set up the database structure.
func (c *Client) initSchema(ctx context.Context, pool *sqlx.DB) error {
	logger := log.FromContext(ctx)
	dbSchemaContent, err := dbSchema.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read database schema: %w", err)
	}

	_, err = pool.ExecContext(ctx, string(dbSchemaContent))
	if err != nil {
		return err
	}

	logger.Info().Msg("Database schema initialized successfully")
	return nil
}

// Close closes the underlying connection pool.
func (c *Client) Close() {
	if c.sqlxDB != nil {
		c.sqlxDB.Close()
	}
}

// Wrap creates a child context with timeout and handles cancellation.
func Wrap[T any](ctx context.Context, operationName string, operation func() (T, error)) (T, error) {
	logger := log.FromContext(ctx)
	fun := func() (dbResult[T], error) {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		start := time.Now()
		result, err := operation()
		elapsed := time.Since(start)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				err = fmt.Errorf("operation timed out after %v", elapsed)
			} else if ctxErr := timeoutCtx.Err(); errors.Is(ctxErr, context.Canceled) {
				err = fmt.Errorf("operation canceled: %w", err)
			} else {
				err = fmt.Errorf("operation failed: %w", err)
			}
			return dbResult[T]{
				Duration: elapsed,
				Result:   result,
			}, err
		}

		return dbResult[T]{
			Duration: elapsed,
			Result:   result,
		}, nil
	}
	result, err := fun()

	if err != nil && !errors.Is(err, ErrNotFound) {
		logger.Error().Err(err).Dur("duration_ms", result.Duration).Msgf("DB operation %s failed", operationName)
	} else if err == nil {
		logger.Trace().Dur("duration_ms", result.Duration).Msgf("DB operation %s completed", operationName)
	}

	return result.Result, err
}

// --------------------------
// Singleton management (unexported global within the db package for convenience)
// --------------------------.
var (
	instance *Client
	instMu   sync.RWMutex
)

// Get returns the current DB singleton (maybe nil if not initialized).
func Get() *Client {
	instMu.RLock()
	defer instMu.RUnlock()
	return instance
}

// Set sets the DB singleton. Useful for tests to inject a different client.
func Set(client *Client) {
	instMu.Lock()
	defer instMu.Unlock()
	if instance != nil {
		// Close previous instance if it was a real client
		instance.Close()
	}
	instance = client
}

// InitFromEnv initializes the DB singleton using the DSN from environment.
func InitFromEnv(ctx context.Context) error {
	logger := log.FromContext(ctx)

	client, err := NewFromEnv(ctx)
	if err != nil {
		return err
	}

	Set(client)
	logger.Info().Msg("SQLite connection established.")
	return nil
}

func setupSQLiteConfig(dsn string) string {
	pragmas := []string{
		"_pragma=journal_mode(WAL)",
		"_pragma=synchronous(NORMAL)",
		"_pragma=busy_timeout(5000)",
		"_pragma=foreign_keys(1)",
	}

	for _, p := range pragmas {
		if !strings.Contains(dsn, p) {
			if strings.Contains(dsn, "?") {
				dsn += "&" + p
			} else {
				dsn += "?" + p
			}
		}
	}

	// _txlock=immediate prevents SQLITE_BUSY on write transactions in WAL mode
	if !strings.Contains(dsn, "_txlock=") {
		if strings.Contains(dsn, "?") {
			dsn += "&_txlock=immediate"
		} else {
			dsn += "?_txlock=immediate"
		}
	}

	return dsn
}
