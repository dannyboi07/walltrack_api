package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var userCollection string = "user"

func getDbContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func InitDb() error {
	var (
		client *mongo.Client
		err    error
	)

	var opts *options.ClientOptions = options.Client()
	opts.SetTimeout(5 * time.Second)
	opts.ApplyURI(os.Getenv("MONGO_URI"))

	ctx, cancel := getDbContext()
	defer cancel()

	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	db = client.Database(os.Getenv("MONGO_DB"))
	return nil
}

func RunConfig() error {
	ctx, cancel := getDbContext()
	defer cancel()
	var unique bool = true

	_, err := db.Collection(userCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		}, Options: &options.IndexOptions{
			Unique: &unique,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
