package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"story-app-monolith/config"
	"time"
)

type Connection struct {
	*mongo.Client
	UserCollection      *mongo.Collection
	FlagCollection      *mongo.Collection
	StoryCollection     *mongo.Collection
	ReadLaterCollection *mongo.Collection
	CommentsCollection  *mongo.Collection
	RepliesCollection   *mongo.Collection
	*mongo.Database
}

func ConnectToDB() (*Connection, error) {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(n+h+p))
	if err != nil {
		return nil, err
	}

	// create database
	db := client.Database("user-service")

	// create collection
	userCollection := db.Collection("users")
	flagCollection := db.Collection("flags")
	storiesCollection := db.Collection("stories")
	commentsCollection := db.Collection("comments")
	repliesCollection := db.Collection("replies")
	readLaterCollection := db.Collection("readLaterCollection")

	dbConnection := &Connection{client, userCollection, flagCollection, storiesCollection,
		commentsCollection, repliesCollection, readLaterCollection, db}

	return dbConnection, nil
}
