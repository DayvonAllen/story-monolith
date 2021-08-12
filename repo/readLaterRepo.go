package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
)

type ReadLaterRepo interface {
	Create(username string, storyId primitive.ObjectID) error
	GetByUsername(username string, page string) (*domain.ReadLaterDto, error)
	Delete(id primitive.ObjectID, username string) error
}
