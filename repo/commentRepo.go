package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
	"time"
)

type CommentRepo interface {
	Create(comment *domain.Comment) error
	FindAllCommentsByResourceId(id primitive.ObjectID, username string) (*[]domain.CommentDto, error)
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error
	LikeCommentById(primitive.ObjectID, string) error
	DisLikeCommentById(primitive.ObjectID, string) error
	UpdateFlagCount(flag *domain.Flag) error
	DeleteById(id primitive.ObjectID, username string) error
	DeleteManyById(id primitive.ObjectID, username string) error
}
