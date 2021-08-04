package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
	"time"
)

type ReplyRepo interface {
	Create(comment *domain.Reply) error
	FindAllRepliesByResourceId(id primitive.ObjectID, username string) (*[]domain.Reply, error)
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error
	LikeReplyById(primitive.ObjectID, string) error
	DisLikeReplyById(primitive.ObjectID, string) error
	UpdateFlagCount(flag *domain.Flag) error
	DeleteById(id primitive.ObjectID, username string) error
}

