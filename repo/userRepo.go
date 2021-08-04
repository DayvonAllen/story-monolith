package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
)

type UserRepo interface {
	FindAll(primitive.ObjectID, string, context.Context, string) (*domain.UserResponse, error)
	FindAllBlockedUsers(primitive.ObjectID, context.Context, string) (*[]domain.UserDto, error)
	Create(*domain.User) error
	FindByID(primitive.ObjectID, context.Context) (*domain.UserDto, error)
	FindByUsername(string, context.Context) (*domain.UserDto, error)
	UpdateByID(primitive.ObjectID, *domain.User) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility, context.Context) error
	UpdateMessageAcceptance(primitive.ObjectID, *domain.UpdateMessageAcceptance, context.Context) error
	UpdateCurrentBadge(primitive.ObjectID, *domain.UpdateCurrentBadge, context.Context) error
	UpdateProfilePicture(primitive.ObjectID, *domain.UpdateProfilePicture, context.Context) error
	UpdateProfileBackgroundPicture(primitive.ObjectID, *domain.UpdateProfileBackgroundPicture, context.Context) error
	UpdateCurrentTagline(primitive.ObjectID, *domain.UpdateCurrentTagline, context.Context)  error
	UpdateVerification(primitive.ObjectID, *domain.UpdateVerification) error
	UpdateDisplayFollowerCount(primitive.ObjectID, *domain.UpdateDisplayFollowerCount) error
	FollowUser(username string, currentUser string) error
	UnfollowUser(username string, currentUser string) error
	UpdatePassword(primitive.ObjectID, string) error
	UpdateFlagCount(*domain.Flag) error
	BlockUser(primitive.ObjectID, string, context.Context, string) error
	UnblockUser(primitive.ObjectID, string, context.Context, string) error
	DeleteByID(primitive.ObjectID, context.Context, string) error
}
