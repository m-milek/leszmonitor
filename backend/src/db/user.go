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

func initUsersCollection(ctx context.Context, database *mongo.Database) error {
	err := createCollection(ctx, database, usersCollectionName)
	if err != nil {
		if errors.Is(err, collectionAlreadyExistsError(usersCollectionName)) {
			logging.Db.Debug().Msg("Users collection already exists.")
			return nil
		}
		return err
	} else {
		logging.Db.Info().Msg("Users collection created successfully.")
	}

	usersCollection := database.Collection(usersCollectionName)
	indexName, err := usersCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{ID_FIELD, 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logging.Db.Error().Err(err).Msg("Failed to create index on users collection")
		return err
	} else {
		logging.Db.Info().Msgf("Index created: %s", indexName)
	}
	return nil
}

func CreateUser(ctx context.Context, user *models.RawUser) (*mongo.InsertOneResult, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*mongo.InsertOneResult, error) {
		res, err := dbClient.getUsersCollection().InsertOne(timeoutCtx, user)
		if err != nil {
			return nil, err
		}
		return res, nil
	})

	logDbOperation("CreateUser", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetUserByUsername(ctx context.Context, username string) (*models.UserResponse, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*models.UserResponse, error) {
		var user models.RawUser
		err := dbClient.getUsersCollection().FindOne(timeoutCtx, bson.M{ID_FIELD: username}).Decode(&user)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
			}
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

func GetRawUserByUsername(ctx context.Context, username string) (*models.RawUser, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) (*models.RawUser, error) {
		var user models.RawUser
		err := dbClient.getUsersCollection().FindOne(timeoutCtx, bson.M{ID_FIELD: username}).Decode(&user)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})

	logDbOperation("GetRawUserByUsername", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
	dbRes, err := withTimeout(ctx, func(timeoutCtx context.Context) ([]models.UserResponse, error) {
		cursor, err := dbClient.getUsersCollection().Find(timeoutCtx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(timeoutCtx)

		usersList := make([]models.UserResponse, 0)

		for cursor.Next(timeoutCtx) {
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
