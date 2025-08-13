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

func initTeamsCollection(ctx context.Context, database *mongo.Database) error {
	err := createCollection(ctx, database, teamsCollectionName)
	if err != nil {
		if errors.Is(err, collectionAlreadyExistsError(teamsCollectionName)) {
			logging.Db.Debug().Msg("Teams collection already exists.")
			return nil
		}
		return err
	}

	// unique index on the "id" field
	teamsCollection := database.Collection(teamsCollectionName)
	indexName, err := teamsCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{"id", 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logging.Db.Error().Err(err).Msg("Failed to create index on teams collection")
		return err
	} else {
		logging.Db.Info().Msgf("Index created: %s", indexName)
	}

	return nil
}

func CreateTeam(ctx context.Context, team *models.Team) (*mongo.InsertOneResult, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*mongo.InsertOneResult, error) {
		res, err := dbClient.getTeamsCollection().InsertOne(timeoutCtx, team)
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

func GetTeamById(ctx context.Context, id string) (*models.Team, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*models.Team, error) {
		var team models.Team
		err := dbClient.getTeamsCollection().FindOne(timeoutCtx, bson.M{"id": id}).Decode(&team)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
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

func GetAllTeams(ctx context.Context) ([]models.Team, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) ([]models.Team, error) {
		cursor, err := dbClient.getTeamsCollection().Find(timeoutCtx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(timeoutCtx)

		teamsList := make([]models.Team, 0)

		for cursor.Next(timeoutCtx) {
			var team models.Team
			if err := cursor.Decode(&team); err != nil {
				return nil, err
			}
			teamsList = append(teamsList, team)
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

func UpdateTeam(ctx context.Context, team *models.Team) (bool, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (bool, error) {
		result, err := dbClient.getTeamsCollection().UpdateOne(timeoutCtx, bson.M{"id": team.Id}, bson.M{"$set": team})
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

func DeleteTeam(ctx context.Context, id string) (bool, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (bool, error) {
		result, err := dbClient.getTeamsCollection().DeleteOne(timeoutCtx, bson.M{"id": id})
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
