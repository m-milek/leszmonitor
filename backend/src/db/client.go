package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logging"
	"os"
	"time"
)

type Client struct {
	conn *pgxpool.Pool
}

var ErrNotFound = errors.New("document not found")
var ErrAlreadyExists = errors.New("resource already exists")

const DB_SCHEMA_FILE = "db/schema.sql"

type dbResult[T any] struct {
	Duration time.Duration
	Result   T
}

var dbClient Client

const timeoutDuration = 1000 * time.Second

func InitDbClient(ctx context.Context) error {
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logging.Db.Info().Msg("Connecting to MongoDB...")

	uri := os.Getenv(env.PostgresURI)

	client, err := pgxpool.New(ctx, uri)
	if err != nil {
		return err
	}

	dbClient = Client{
		conn: client,
	}

	_, err = ping(ctx)
	if err != nil {
		logging.Db.Fatal().Err(err).Msg("Failed to ping PostgreSQL")
	}

	err = dbClient.initSchema(ctx)
	if err != nil {
		logging.Db.Fatal().Err(err).Msg("Failed to initialize database schema")
	}

	logging.Db.Info().Msg("MongoDB connection established.")

	return nil
}

// initSchema reads the database schema from a file and executes it to set up the database structure.
func (c *Client) initSchema(ctx context.Context) error {
	schemaBytes, err := os.ReadFile(DB_SCHEMA_FILE)
	if err != nil {
		return fmt.Errorf("failed to read DB schema file: %w", err)
	}
	schema := string(schemaBytes)

	status, err := c.conn.Exec(ctx, schema)
	if err != nil {
		return err
	}

	logging.Db.Info().Msgf("Database schema initialized: %s", status.String())
	return nil
}

// withTimeout creates a child context with timeout and handles cancellation.
func withTimeout[T any](timeoutCtx context.Context, operation func() (T, error)) (dbResult[T], error) {
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

func logDbOperation[T any](operationName string, result dbResult[T], err error) {
	if err != nil {
		logging.Db.Error().Err(err).Msgf("DB operation %s failed", operationName)
		return
	}
	logging.Db.Trace().Dur("duration", result.Duration).Any("result", result.Result).Msgf("DB operation %s completed", operationName)
}

func ping(ctx context.Context) (int64, error) {
	result, err := withTimeout(ctx, func() (int64, error) {
		start := time.Now()
		err := dbClient.conn.Ping(ctx)
		if err != nil {
			return 0, err
		}
		duration := time.Since(start).Milliseconds()
		return duration, nil
	})

	logDbOperation("Ping", result, err)

	if err != nil {
		return 0, err
	}
	return result.Result, nil
}
