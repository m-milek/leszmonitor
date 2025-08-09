package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/logger"
	monitors "github.com/m-milek/leszmonitor/uptime/monitor"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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

// CreateMonitor adds a new monitor to the database and returns its ID (short ID).
func CreateMonitor(monitor monitors.IMonitor) (string, error) {
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
