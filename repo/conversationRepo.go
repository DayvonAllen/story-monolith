package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
)

type ConversationRepo interface {
	Create(message domain.Message) error
	FindByOwner(message domain.Message) (*domain.Conversation, error)
	FindConversation(owner, to string) (*domain.Conversation, error)
	GetConversationPreviews(owner string) (*[]domain.ConversationPreview, error)
	UpdateConversation(conversation domain.Conversation, message domain.Message) error
	DeleteByID(conversationId primitive.ObjectID, username string) error
	DeleteAllByUsername(username string) error
}
