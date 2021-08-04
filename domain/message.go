package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id        primitive.ObjectID `bson:"_id" json:"messageId"`
	Content   string             `bson:"content" json:"content"`
	From      string             `bson:"from" json:"-" `
	To        string             `bson:"to" json:"to"`
	Read	  bool 				 `bson:"read" json:"read"`
	CreatedAt time.Time          `bson:"createdAt" json:"sentAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}

type DeleteMessage struct {
	Id        primitive.ObjectID `bson:"_id" json:"messageId"`
}

type DeleteMessages struct {
	MessageIds []DeleteMessage `json:"ids"`
}

