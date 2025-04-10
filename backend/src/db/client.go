package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
	"time"
)

type Client struct {
	uri     string
	client  *mongo.Client
	baseCtx context.Context
}

type dbResult[T any] struct {
	Duration time.Duration
	Result   T
}

var dbClient Client

const timeoutDuration = 10 * time.Second

func InitDbClient(baseCtx context.Context) error {
	ctx, cancel := context.WithTimeout(baseCtx, 10*time.Second)
	defer cancel()

	uri := os.Getenv(env.MONGODB_URI)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}
	logger.Db.Info().Msg("MongoDB connection established")

	dbClient = Client{
		uri:     uri,
		client:  client,
		baseCtx: baseCtx,
	}

	return nil
}

// withTimeout creates a child context with timeout and handles cancellation
func withTimeout[T any](operation func(ctx context.Context) (T, error)) (dbResult[T], error) {
	ctx, cancel := context.WithTimeout(dbClient.baseCtx, timeoutDuration)
	defer cancel()

	start := time.Now()
	result, err := operation(ctx)
	elapsed := time.Since(start)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = fmt.Errorf("operation timed out after %v", elapsed)
		} else if ctxErr := ctx.Err(); errors.Is(ctxErr, context.Canceled) {
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
		logger.Db.Error().Err(err).Msgf("DB operation %s failed", operationName)
		return
	}
	logger.Db.Info().Msgf("DB operation %s took %v", operationName, result.Duration)
}

func Ping() (int64, error) {
	result, err := withTimeout(func(ctx context.Context) (int64, error) {
		start := time.Now()
		err := dbClient.client.Ping(ctx, nil)
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
