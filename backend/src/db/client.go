package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/common"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
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

func (*Client) getDatabase() *mongo.Database {
	return dbClient.client.Database(DatabaseName)
}

func (*Client) getUsersCollection() *mongo.Collection {
	return dbClient.getDatabase().Collection(UsersCollectionName)
}

func (*Client) getMonitorsCollection() *mongo.Collection {
	return dbClient.getDatabase().Collection(MonitorsCollectionName)
}

type dbResult[T any] struct {
	Duration time.Duration
	Result   T
}

type CollectionAlreadyExistsError string

func (err CollectionAlreadyExistsError) Error() string {
	return "collection already exists: " + string(err)
}

var dbClient Client

const timeoutDuration = 5 * time.Second

const (
	DatabaseName           = "leszmonitor"
	UsersCollectionName    = "users"
	MonitorsCollectionName = "monitors"
)

func InitDbClient(baseCtx context.Context) error {
	_, cancel := context.WithTimeout(baseCtx, 10*time.Second)
	defer cancel()

	logger.Db.Info().Msg("Connecting to MongoDB...")

	uri := os.Getenv(env.MONGODB_URI)

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

	_, err = Ping()
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	database := client.Database(DatabaseName)
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
		return CollectionAlreadyExistsError(collectionName)
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
	err := createCollection(dbClient.baseCtx, database, UsersCollectionName)
	if err != nil {
		if errors.Is(err, CollectionAlreadyExistsError(UsersCollectionName)) {
			logger.Db.Info().Msg("Users collection already exists.")
			return nil
		}
		return err
	}

	usersCollection := database.Collection(UsersCollectionName)
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
	err := createCollection(dbClient.baseCtx, database, MonitorsCollectionName)
	if err != nil {
		if errors.Is(err, CollectionAlreadyExistsError(MonitorsCollectionName)) {
			logger.Db.Info().Msg("Monitors collection already exists.")
			return nil
		}
		return err
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
	logger.Db.Debug().Dur("duration", result.Duration).Any("result", result.Result).Msgf("DB operation %s completed", operationName)
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

func AddUser(user *model.User) (*mongo.InsertOneResult, error) {
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

func GetUser(username string) (*model.User, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*model.User, error) {
		var user model.User
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

func AddMonitor(monitor monitors.IMonitor) (*mongo.InsertOneResult, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*mongo.InsertOneResult, error) {
		res, err := dbClient.getMonitorsCollection().InsertOne(ctx, monitor)
		if err != nil {
			return nil, err
		}
		return res, nil
	})

	logDbOperation("InsertMonitor", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetAllMonitors() ([]monitors.IMonitor, error) {
	dbRes, err := withTimeout(func(ctx context.Context) ([]monitors.IMonitor, error) {
		cursor, err := dbClient.getMonitorsCollection().Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var monitorsList []monitors.IMonitor

		// First decode into a map to determine the monitor type
		for cursor.Next(ctx) {
			// Decode into a raw document first
			var rawDoc bson.M
			if err := cursor.Decode(&rawDoc); err != nil {
				return nil, err
			}

			// Get the monitor type
			typeValue, ok := rawDoc["type"]
			if !ok {
				return nil, fmt.Errorf("monitor type not found in document")
			}

			monitorType, ok := typeValue.(string)
			if !ok {
				return nil, fmt.Errorf("monitor type is not a string")
			}

			// Create the appropriate concrete type based on the monitor type
			var monitor monitors.IMonitor

			switch monitors.MonitorType(monitorType) {
			case monitors.Http:
				httpMonitor := &monitors.HttpMonitor{}
				// Re-encode and decode to the concrete type
				data, err := bson.Marshal(rawDoc)
				if err != nil {
					return nil, err
				}
				if err := bson.Unmarshal(data, httpMonitor); err != nil {
					return nil, err
				}
				monitor = httpMonitor

			default:
				return nil, fmt.Errorf("unknown monitor type: %s", monitorType)
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
