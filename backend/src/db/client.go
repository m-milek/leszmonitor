package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logging"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
	"time"
)

type Client struct {
	uri      string
	client   *mongo.Client
	database *mongo.Database
}

var ErrNotFound = errors.New("document not found")

func (*Client) getDatabase() *mongo.Database {
	return dbClient.client.Database(databaseName)
}

func (*Client) getUsersCollection() *mongo.Collection {
	return dbClient.getDatabase().Collection(usersCollectionName)
}

func (*Client) getMonitorsCollection() *mongo.Collection {
	return dbClient.getDatabase().Collection(monitorsCollectionName)
}

func (*Client) getTeamsCollection() *mongo.Collection {
	return dbClient.getDatabase().Collection(teamsCollectionName)
}

type dbResult[T any] struct {
	Duration time.Duration
	Result   T
}

type collectionAlreadyExistsError string

func (err collectionAlreadyExistsError) Error() string {
	return "collection already exists: " + string(err)
}

var dbClient Client

const timeoutDuration = 5 * time.Second

const (
	databaseName           = "leszmonitor"
	usersCollectionName    = "users"
	monitorsCollectionName = "monitors"
	teamsCollectionName    = "teams"
)

func InitDbClient(ctx context.Context) error {
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logging.Db.Info().Msg("Connecting to MongoDB...")

	uri := os.Getenv(env.MongoDbUri)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	dbClient = Client{
		uri:      uri,
		client:   client,
		database: nil,
	}

	_, err = ping(ctx)
	if err != nil {
		logging.Db.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	database := client.Database(databaseName)
	dbClient.database = database

	err = initSchema(ctx)
	if err != nil {
		logging.Db.Fatal().Err(err).Msg("Failed to initialize database schema")
	}

	logging.Db.Info().Msg("MongoDB connection established.")

	return nil
}

func initSchema(ctx context.Context) error {
	database := dbClient.getDatabase()

	err := initUsersCollection(ctx, database)
	if err != nil {
		logging.Db.Error().Err(err).Msg("Failed to initialize users collection")
		return err
	}

	err = initMonitorsCollection(ctx, database)
	if err != nil {
		logging.Db.Error().Err(err).Msg("Failed to initialize monitors collection")
		return err
	}

	err = initTeamsCollection(ctx, database)
	if err != nil {
		logging.Db.Error().Err(err).Msg("Failed to initialize teams collection")
	}

	return err
}

func collectionExists(ctx context.Context, database *mongo.Database, collectionName string) (bool, error) {
	collections, err := database.ListCollections(ctx, bson.D{{"name", collectionName}})
	if err != nil {
		return false, err
	}
	defer collections.Close(ctx)

	for collections.Next(ctx) {
		var result bson.M
		err := collections.Decode(&result)
		if err != nil {
			return false, err
		}
		if result["name"] == collectionName {
			return true, nil
		}
	}

	return false, nil
}

func createCollection(ctx context.Context, database *mongo.Database, collectionName string) error {
	// Check if the collection already exists
	exists, err := collectionExists(ctx, database, collectionName)
	if err != nil {
		return err
	}
	if exists {
		return collectionAlreadyExistsError(collectionName)
	}

	// Create the collection
	err = database.CreateCollection(ctx, collectionName)
	if err != nil {
		if !errors.Is(err, mongo.CommandError{}) {
			return err
		}
	}

	return nil
}

// withTimeout creates a child context with timeout and handles cancellation.
func withTimeout[T any](timeoutCtx context.Context, operation func(context.Context) (T, error)) (dbResult[T], error) {
	timeoutCtx, cancel := context.WithTimeout(timeoutCtx, timeoutDuration)
	defer cancel()

	start := time.Now()
	result, err := operation(timeoutCtx)
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
	result, err := withTimeout(ctx, func(context.Context) (int64, error) {
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
