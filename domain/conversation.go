package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Conversation struct {
	Id        primitive.ObjectID `bson:"_id" json:"-"`
	Owner     string             `bson:"owner" json:"-"`
	From      string             `bson:"from" json:"-"`
	To        string             `bson:"to" json:"to"`
	Messages  []Message          `bson:"messages" json:"messages"`
	UnreadCount	   int 				  `bson:"unreadCount" json:"unreadCount"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}

type ConversationPreview struct {
	Id             primitive.ObjectID `bson:"_id" json:"-"`
	Owner          string             `bson:"-" json:"-"`
	From           string             `bson:"-" json:"-" `
	To             string             `bson:"-" json:"-"`
	Messages       []Message          `bson:"-" json:"-"`
	UnreadCount	   int 				  `bson:"-" json:"unreadCount"`
	PreviewMessage Message            `bson:"-" json:"previewMessage"`
	CreatedAt      time.Time          `bson:"-" json:"-"`
	UpdatedAt      time.Time          `bson:"-" json:"updatedAt"`
}
