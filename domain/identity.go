package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Identity struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
	Identifier []byte      `bson:"identifier" json:"identifier"`
	StoryId primitive.ObjectID `bson:"storyId" json:"storyId"`
}