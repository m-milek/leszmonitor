package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/common"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitor"
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
	baseCtx  context.Context
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
)

func InitDbClient(baseCtx context.Context) error {
	_, cancel := context.WithTimeout(baseCtx, 10*time.Second)
	defer cancel()

	logger.Db.Info().Msg("Connecting to MongoDB...")

	uri := os.Getenv(env.MongoDbUri)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	dbClient = Client{
		uri:      uri,
		client:   client,
		database: nil,
		baseCtx:  baseCtx,
	}

	_, err = ping()
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	database := client.Database(databaseName)
	dbClient.database = database

	err = initSchema()
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to initialize database schema")
	}

	logger.Db.Info().Msg("MongoDB connection established.")

	return nil
}

func initSchema() error {
	database := dbClient.getDatabase()

	err := initUsersCollection(database)
	if err != nil {
		logger.Db.Error().Err(err).Msg("Failed to initialize users collection")
		return err
	}

	err = initMonitorsCollection(database)
	if err != nil {
		logger.Db.Error().Err(err).Msg("Failed to initialize monitors collection")
		return err
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

func initUsersCollection(database *mongo.Database) error {
	err := createCollection(dbClient.baseCtx, database, usersCollectionName)
	if err != nil {
		if errors.Is(err, collectionAlreadyExistsError(usersCollectionName)) {
			logger.Db.Debug().Msg("Users collection already exists.")
			return nil
		}
		return err
	}

	usersCollection := database.Collection(usersCollectionName)
	indexName, err := usersCollection.Indexes().CreateOne(
		dbClient.baseCtx,
		mongo.IndexModel{
			Keys: bson.D{
				{"username", 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to create index on users collection")
	} else {
		logger.Db.Info().Msgf("Index created: %s", indexName)
	}
	return nil
}

func initMonitorsCollection(database *mongo.Database) error {
	err := createCollection(dbClient.baseCtx, database, monitorsCollectionName)
	if err != nil {
		if errors.Is(err, collectionAlreadyExistsError(monitorsCollectionName)) {
			logger.Db.Debug().Msg("Monitors collection already exists.")
			return nil
		}
		return err
	}

	// unique index on the "id" field
	monitorsCollection := database.Collection(monitorsCollectionName)
	indexName, err := monitorsCollection.Indexes().CreateOne(
		dbClient.baseCtx,
		mongo.IndexModel{
			Keys: bson.D{
				{"id", 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to create index on monitors collection")
	} else {
		logger.Db.Info().Msgf("Index created: %s", indexName)
	}

	return nil
}

// withTimeout creates a child context with timeout and handles cancellation.
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
	logger.Db.Debug().Dur("duration", result.Duration).Any("result", result.Result).Msgf("DB operation %s completed", operationName)
}

func ping() (int64, error) {
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

func AddUser(user *common.User) (*mongo.InsertOneResult, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*mongo.InsertOneResult, error) {
		res, err := dbClient.getUsersCollection().InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}
		return res, nil
	})

	logDbOperation("InsertUser", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetUser(username string) (*common.User, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*common.User, error) {
		var user common.User
		err := dbClient.getUsersCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	})

	logDbOperation("GetUser", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

// AddMonitor adds a new monitor to the database and returns its ID (short ID).
func AddMonitor(monitor monitors.IMonitor) (string, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (string, error) {
		_, err := dbClient.getMonitorsCollection().InsertOne(ctx, monitor)
		if err != nil {
			return "", err
		}
		return monitor.GetId(), nil
	})

	logDbOperation("InsertMonitor", dbRes, err)

	if err != nil {
		return "", err
	}

	// Broadcast that a monitor has been added
	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		Id:      monitor.GetId(),
		Status:  monitors.Created,
		Monitor: &monitor,
	})

	return dbRes.Result, nil
}

func GetAllMonitors() ([]monitors.IMonitor, error) {
	dbRes, err := withTimeout(func(ctx context.Context) ([]monitors.IMonitor, error) {
		cursor, err := dbClient.getMonitorsCollection().Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		monitorsList := make([]monitors.IMonitor, 0)

		// First decode into a map to determine the monitor type
		for cursor.Next(ctx) {
			// Decode into a raw document first
			var rawDoc bson.M
			if err := cursor.Decode(&rawDoc); err != nil {
				return nil, err
			}

			monitor, err := monitors.FromRawBsonDoc(rawDoc)
			if err != nil {
				return nil, fmt.Errorf("failed to map monitor from BSON: %w", err)
			}

			monitorsList = append(monitorsList, monitor)
		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}

		return monitorsList, nil
	})

	logDbOperation("GetAllMonitors", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func DeleteMonitor(id string) (bool, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (bool, error) {
		result, err := dbClient.getMonitorsCollection().DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			return false, err
		}
		return result.DeletedCount > 0, nil
	})

	logDbOperation("DeleteMonitor", dbRes, err)

	if err != nil {
		return false, err
	}

	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		Id:      id,
		Status:  monitors.Deleted,
		Monitor: nil,
	})

	return dbRes.Result, nil
}

func GetMonitorById(id string) (monitors.IMonitor, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (monitors.IMonitor, error) {
		var rawDoc bson.M
		err := dbClient.getMonitorsCollection().FindOne(ctx, bson.M{"id": id}).Decode(&rawDoc)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		monitor, err := monitors.FromRawBsonDoc(rawDoc)

		if err != nil {
			return nil, fmt.Errorf("failed to map monitor from BSON: %w", err)
		}

		return monitor, nil
	})

	logDbOperation("GetMonitorById", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func UpdateMonitor(newMonitor monitors.IMonitor) (bool, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (bool, error) {
		result, err := dbClient.getMonitorsCollection().UpdateOne(ctx, bson.M{"id": newMonitor.GetId()}, bson.M{"$set": newMonitor})
		if err != nil {
			return false, err
		}
		if result.MatchedCount == 0 {
			return false, ErrNotFound
		}
		wasUpdated := result.ModifiedCount > 0
		return wasUpdated, nil
	})

	logDbOperation("UpdateMonitor", dbRes, err)

	if err != nil {
		return false, err
	}

	// Broadcast that a monitor has been updated
	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		Id:      newMonitor.GetId(),
		Status:  monitors.Edited,
		Monitor: &newMonitor,
	})

	return dbRes.Result, nil
}
