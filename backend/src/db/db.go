package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-milek/leszmonitor/config"
	"github.com/m-milek/leszmonitor/log"
)

var ErrNotFound = errors.New("document not found")
var ErrAlreadyExists = errors.New("resource already exists")

func pgErrIs(err error, pgErrCode string) bool {
	var e *pgconn.PgError
	return errors.As(err, &e) && e.Code == pgErrCode
}

const dbSchemaFilePath = "db/schema.sql"

const timeoutDuration = 1000 * time.Second

// DB defines the database access surface. It returns repository interfaces for easy mocking.
type DB interface {
	Users() IUserRepository
	Monitors() IMonitorRepository
	Projects() IProjectRepository
	Close()
}

// DBClient implements DB using a pgx connection pool.
type DBClient struct {
	dbPool
	// cached repositories to avoid re-allocation on every getter call
	users    IUserRepository
	monitors IMonitorRepository
	projects IProjectRepository
}

type dbPool struct {
	pool *pgxpool.Pool
}

type dbResult[T any] struct {
	Duration time.Duration
	Result   T
}

type baseRepository struct {
	dbPool
}

func newBaseRepository(pool *pgxpool.Pool) baseRepository {
	return baseRepository{
		dbPool: dbPool{pool: pool},
	}
}

// New creates a new DB client using the provided DSN. It pings the DB and ensures the schema exists.
func New(ctx context.Context, dsn string) (*DBClient, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	c := &DBClient{
		dbPool: dbPool{pool: pool},
	}

	if err := c.initSchema(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	// initialize and cache repositories once
	c.users = newUserRepository(newBaseRepository(pool))
	c.monitors = newMonitorRepository(newBaseRepository(pool))
	c.projects = newProjectRepository(newBaseRepository(pool))

	return c, nil
}

// initSchema reads the database schema from a file and executes it to set up the database structure.
func (c *DBClient) initSchema(ctx context.Context, pool *pgxpool.Pool) error {
	schemaBytes, err := os.ReadFile(dbSchemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read DB schema file: %w", err)
	}
	schema := string(schemaBytes)

	status, err := pool.Exec(ctx, schema)
	if err != nil {
		return err
	}

	log.Db.Info().Msgf("Database schema initialized: %s", status.String())
	return nil
}

// Close closes the underlying connection pool.
func (c *DBClient) Close() {
	c.pool.Close()
}

// dbWrap creates a child context with timeout and handles cancellation.
func dbWrap[T any](timeoutCtx context.Context, operationName string, operation func() (T, error)) (T, error) {
	fun := func() (dbResult[T], error) {
		timeoutCtx, cancel := context.WithTimeout(timeoutCtx, timeoutDuration)
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

	if err != nil {
		log.Db.Error().Err(err).Msgf("DB operation %s failed", operationName)
	} else {
		log.Db.Trace().Dur("duration", result.Duration).Any("result", result.Result).Msgf("DB operation %s completed", operationName)
	}

	return result.Result, err
}

// Repository getters (return interfaces for mocking)
func (c *DBClient) Users() IUserRepository       { return c.users }
func (c *DBClient) Monitors() IMonitorRepository { return c.monitors }
func (c *DBClient) Projects() IProjectRepository { return c.projects }

// --------------------------
// Singleton management (unexported global within the db package for convenience)
// --------------------------
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Db.Info().Msg("Connecting to PostgreSQL...")

	uri := os.Getenv(config.PostgresURI)
	c, err := New(ctx, uri)
	if err != nil {
		return err
	}

	Set(c)
	log.Db.Info().Msg("PostgreSQL connection established.")
	return nil
}
