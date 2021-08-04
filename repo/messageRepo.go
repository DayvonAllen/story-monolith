package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
)

type MessageRepo interface {
	Create(message *domain.Message) (*domain.Conversation, error)
	DeleteByID(owner string, id primitive.ObjectID) error
	DeleteAllByIDs(owner string, messages []domain.DeleteMessage) error
}

