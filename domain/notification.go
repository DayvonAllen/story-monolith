package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	For       string             `bson:"for" json:"-"`
	Content   string             `bson:"content" json:"content"`
	Path      string             `bson:"path" json:"path"`
	ReadStatus    bool               `bson:"readStatus" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}
