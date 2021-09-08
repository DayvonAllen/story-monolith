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
	UserCollection         *mongo.Collection
	FlagCollection         *mongo.Collection
	StoryCollection        *mongo.Collection
	ReadLaterCollection    *mongo.Collection
	CommentsCollection     *mongo.Collection
	RepliesCollection      *mongo.Collection
	ConversationCollection *mongo.Collection
	MessageCollection      *mongo.Collection
	NotificationCollection *mongo.Collection
	IdentityCollection     *mongo.Collection
	*mongo.Database
}

var MongoConn *Connection

func ConnectToDB() {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(n+h+p))

	if err != nil {
		panic(err)
	}
	// create database
	db := client.Database("story-service")

	// create collection
	userCollection := db.Collection("users")
	flagCollection := db.Collection("flags")
	storiesCollection := db.Collection("stories")
	commentsCollection := db.Collection("comments")
	repliesCollection := db.Collection("replies")
	readLaterCollection := db.Collection("readLater")
	conversationCollection := db.Collection("conversation")
	messageCollection := db.Collection("message")
	notificationCollection := db.Collection("notification")
	identityCollection := db.Collection("identity")

	dbConnection := &Connection{client, userCollection, flagCollection, storiesCollection,
		commentsCollection, repliesCollection, readLaterCollection,
		conversationCollection, messageCollection, notificationCollection, identityCollection,
		db}

	MongoConn = dbConnection
}
