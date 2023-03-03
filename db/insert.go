package db

import (
	"walltrack/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(user model.User) (model.User, error) {
	var coll *mongo.Collection = db.Collection(userCollection)

	ctx, cancel := getDbContext()
	defer cancel()

	insertResult, err := coll.InsertOne(ctx, user)
	if err == nil {
		user.Id = insertResult.InsertedID.(primitive.ObjectID)
	}

	return user, err
}
