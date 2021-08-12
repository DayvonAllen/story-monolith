package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReadLater struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Username       string 			  `bson:"username" json:"username"`
	Story		   StoryDto `bson:"story" json:"story"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"-"`
}

type ReadLaterDto struct {
	ReadLaterItems  []ReadLater `json:"readLaterItems"`
	NumberOfStories int64		`json:"numberOfStories"`
	CurrentPage int			`json:"currentPage"`
	NumberOfPages int		`json:"numberOfPages"`
}
