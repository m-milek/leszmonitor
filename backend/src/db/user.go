package db

import (
	"context"
	"errors"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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

func AddUser(user *models.RawUser) (*mongo.InsertOneResult, error) {
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

func GetUserByUsername(username string) (*models.UserResponse, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*models.UserResponse, error) {
		var user models.RawUser
		err := dbClient.getUsersCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			return nil, err
		}
		return user.IntoUser(), nil
	})

	logDbOperation("GetUserByUsername", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetRawUserByUsername(username string) (*models.RawUser, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*models.RawUser, error) {
		var user models.RawUser
		err := dbClient.getUsersCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	})

	logDbOperation("GetUserByUsername", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetAllUsers() ([]models.UserResponse, error) {
	dbRes, err := withTimeout(func(ctx context.Context) ([]models.UserResponse, error) {
		cursor, err := dbClient.getUsersCollection().Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		usersList := make([]models.UserResponse, 0)

		for cursor.Next(ctx) {
			var user models.RawUser
			if err := cursor.Decode(&user); err != nil {
				return nil, err
			}
			usersList = append(usersList, *user.IntoUser())
		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}

		return usersList, nil
	})

	logDbOperation("GetAllUsers", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}
