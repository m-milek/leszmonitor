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

func AddUser(user *common.RawUser) (*mongo.InsertOneResult, error) {
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

func GetUser(username string) (*common.UserResponse, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*common.UserResponse, error) {
		var user common.RawUser
		err := dbClient.getUsersCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			return nil, err
		}
		return user.IntoUser(), nil
	})

	logDbOperation("GetUser", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetRawUser(username string) (*common.RawUser, error) {
	dbRes, err := withTimeout(func(ctx context.Context) (*common.RawUser, error) {
		var user common.RawUser
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

func GetAllUsers() ([]*common.UserResponse, error) {
	dbRes, err := withTimeout(func(ctx context.Context) ([]*common.UserResponse, error) {
		cursor, err := dbClient.getUsersCollection().Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		usersList := make([]*common.UserResponse, 0)

		for cursor.Next(ctx) {
			var user common.RawUser
			if err := cursor.Decode(&user); err != nil {
				return nil, err
			}
			usersList = append(usersList, user.IntoUser())
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
