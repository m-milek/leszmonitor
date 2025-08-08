package db

import (
	"context"
	"errors"
	"github.com/m-milek/leszmonitor/common"
	"github.com/m-milek/leszmonitor/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initTeamsCollection(database *mongo.Database) error {
	err := createCollection(dbClient.baseCtx, database, teamsCollectionName)
	if err != nil {
		if errors.Is(err, collectionAlreadyExistsError(teamsCollectionName)) {
			logger.Db.Debug().Msg("Teams collection already exists.")
			return nil
		}
		return err
	}

	// unique index on the "id" field
	teamsCollection := database.Collection(teamsCollectionName)
	indexName, err := teamsCollection.Indexes().CreateOne(
		dbClient.baseCtx,
		mongo.IndexModel{
			Keys: bson.D{
				{"id", 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logger.Db.Fatal().Err(err).Msg("Failed to create index on teams collection")
	} else {
		logger.Db.Info().Msgf("Index created: %s", indexName)
	}

	return nil
}

func AddTeam(team *common.Team) (*mongo.InsertOneResult, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*mongo.InsertOneResult, error) {
		res, err := dbClient.getTeamsCollection().InsertOne(ctx, team)
		if err != nil {
			return nil, err
		}
		return res, nil
	})

	logDbOperation("InsertTeam", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetTeamById(id string) (*common.Team, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*common.Team, error) {
		var team common.Team
		err := dbClient.getTeamsCollection().FindOne(ctx, bson.M{"id": id}).Decode(&team)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, nil
			}
			return nil, err
		}
		return &team, nil
	})

	logDbOperation("GetTeamById", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetAllTeams() ([]*common.Team, error) {
	dbRes, err := withTimeout(func(ctx context.Context) ([]*common.Team, error) {
		cursor, err := dbClient.getTeamsCollection().Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		teamsList := make([]*common.Team, 0)

		for cursor.Next(ctx) {
			var team common.Team
			if err := cursor.Decode(&team); err != nil {
				return nil, err
			}
			teamsList = append(teamsList, &team)
		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}

		return teamsList, nil
	})

	logDbOperation("GetAllTeams", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func UpdateTeam(team *common.Team) (bool, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (bool, error) {
		result, err := dbClient.getTeamsCollection().UpdateOne(ctx, bson.M{"id": team.Id}, bson.M{"$set": team})
		if err != nil {
			return false, err
		}
		if result.MatchedCount == 0 {
			return false, ErrNotFound
		}
		return result.ModifiedCount > 0, nil
	})

	logDbOperation("UpdateTeam", dbRes, err)

	if err != nil {
		return false, err
	}
	return dbRes.Result, nil
}

func DeleteTeam(id string) (bool, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (bool, error) {
		result, err := dbClient.getTeamsCollection().DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			return false, err
		}
		if result.DeletedCount == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})

	logDbOperation("DeleteTeam", dbRes, err)

	if err != nil {
		return dbRes.Result, err
	}
	return dbRes.Result, nil
}
