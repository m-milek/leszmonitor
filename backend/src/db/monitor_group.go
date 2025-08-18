package db

import (
	"context"
	"errors"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initGroupsCollection(ctx context.Context, database *mongo.Database) error {
	err := createCollection(ctx, database, groupsCollectionName)
	if err != nil {
		if errors.Is(err, collectionAlreadyExistsError(groupsCollectionName)) {
			logging.Db.Debug().Msg("Groups collection already exists.")
			return nil
		}
		return err
	} else {
		logging.Db.Info().Msg("Groups collection created successfully.")
	}

	// unique index on the "id" field
	groupsCollection := database.Collection(groupsCollectionName)
	indexName, err := groupsCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{ID_FIELD, 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logging.Db.Error().Err(err).Msg("Failed to create index on groups collection")
		return err
	}

	logging.Db.Info().Msgf("Index created: %s", indexName)
	return nil
}

func CreateMonitorGroup(ctx context.Context, group *models.MonitorGroup) (*mongo.InsertOneResult, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*mongo.InsertOneResult, error) {
		res, err := dbClient.getGroupsCollection().InsertOne(timeoutCtx, group)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}
		return res, nil
	})

	logDbOperation("CreateMonitorGroup", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetMonitorGroupById(ctx context.Context, teamId string, groupId string) (*models.MonitorGroup, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*models.MonitorGroup, error) {
		var group models.MonitorGroup
		err := dbClient.getGroupsCollection().FindOne(timeoutCtx, bson.M{ID_FIELD: groupId}).Decode(&group)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &group, nil
	})

	logDbOperation("GetMonitorGroupById", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetMonitorGroupsForTeam(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) ([]models.MonitorGroup, error) {
		cursor, err := dbClient.getGroupsCollection().Find(timeoutCtx, bson.M{"teamId": team.ObjectId})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(timeoutCtx)

		groups := make([]models.MonitorGroup, 0)
		for cursor.Next(timeoutCtx) {
			var group models.MonitorGroup
			if err := cursor.Decode(&group); err != nil {
				return nil, err
			}
			groups = append(groups, group)
		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}

		return groups, nil
	})

	logDbOperation("GetTeamMonitorGroups", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func UpdateMonitorGroup(ctx context.Context, teamId string, group *models.MonitorGroup) (*mongo.UpdateResult, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*mongo.UpdateResult, error) {
		res, err := dbClient.getGroupsCollection().UpdateOne(
			timeoutCtx,
			bson.M{OBJECT_ID_FIELD: group.ObjectId},
			bson.M{"$set": group},
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	})

	logDbOperation("UpdateMonitorGroup", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func DeleteMonitorGroup(ctx context.Context, teamId string, groupId string) (bool, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (bool, error) {
		res, err := dbClient.getGroupsCollection().DeleteOne(timeoutCtx, bson.M{ID_FIELD: groupId})
		if err != nil {
			return false, err
		}
		return res.DeletedCount != 0, nil
	})

	logDbOperation("DeleteMonitorGroup", dbRes, err)

	if err != nil {
		return false, err
	}
	return dbRes.Result, nil
}
