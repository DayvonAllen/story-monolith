package services

import (
	"context"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"story-app-monolith/domain"
	"story-app-monolith/repo"
	"strings"
	"time"
)

type UserService interface {
	GetAllUsers(primitive.ObjectID, string, context.Context, string) (*domain.UserResponse, error)
	GetAllBlockedUsers(primitive.ObjectID,  context.Context, string) (*[]domain.UserDto, error)
	CreateUser(*domain.User) error
	GetUserByID(primitive.ObjectID, context.Context) (*domain.UserDto, error)
	GetUserByUsername(string, context.Context) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility, context.Context) error
	UpdateMessageAcceptance(primitive.ObjectID, *domain.UpdateMessageAcceptance, context.Context) error
	UpdateCurrentBadge(primitive.ObjectID, *domain.UpdateCurrentBadge, context.Context) error
	UpdateProfilePicture(primitive.ObjectID, *domain.UpdateProfilePicture, context.Context) error
	UpdateProfileBackgroundPicture(primitive.ObjectID, *domain.UpdateProfileBackgroundPicture, context.Context) error
	UpdateCurrentTagline(primitive.ObjectID, *domain.UpdateCurrentTagline, context.Context)  error
	UpdateDisplayFollowerCount(primitive.ObjectID, *domain.UpdateDisplayFollowerCount) error
	UpdateVerification(primitive.ObjectID, *domain.UpdateVerification) error
	UpdatePassword(primitive.ObjectID, string) error
	UpdateFlagCount(*domain.Flag) error
	FollowUser(username string, currentUser string) error
	UnfollowUser(username string, currentUser string) error
	BlockUser(primitive.ObjectID, string, context.Context, string) error
	UnblockUser(primitive.ObjectID, string, context.Context, string) error
	DeleteByID(primitive.ObjectID, context.Context, string) error
}

// DefaultUserService the service has a dependency of the repo
type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetAllUsers(id primitive.ObjectID, page string, ctx context.Context, username string) (*domain.UserResponse, error) {
	u, err := s.repo.FindAll(id, page, ctx, username)
	if err != nil {
		return nil, err
	}
	return  u, nil
}

func (s DefaultUserService) GetAllBlockedUsers(id primitive.ObjectID, ctx context.Context, username string) (*[]domain.UserDto, error) {
	u, err := s.repo.FindAllBlockedUsers(id, ctx, username)
	if err != nil {
		return nil, err
	}
	return  u, nil
}

func (s DefaultUserService) CreateUser(user *domain.User) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	a := new(domain.Authentication)
	h := utils.UUIDv4()
	signedHash, err := a.SignToken([]byte(h))

	if err != nil {
		return err
	}

	hash := h + "-" + string(signedHash)
	user.VerificationCode = hash

	err = s.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) GetUserByID(id primitive.ObjectID, ctx context.Context) (*domain.UserDto, error) {
	u, err := s.repo.FindByID(id, ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) GetUserByUsername(username string, ctx context.Context) (*domain.UserDto, error) {
	u, err := s.repo.FindByUsername(username, ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfileVisibility(id, user, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateDisplayFollowerCount(id primitive.ObjectID, user *domain.UpdateDisplayFollowerCount) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateDisplayFollowerCount(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateMessageAcceptance(id, user, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateCurrentBadge(id, user, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfilePicture(id, user, ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s DefaultUserService) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfileBackgroundPicture(id, user, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateCurrentTagline(id, user, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdatePassword(id primitive.ObjectID, password string) error {
	err := s.repo.UpdatePassword(id, password)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateVerification(id primitive.ObjectID, user *domain.UpdateVerification) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateVerification(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateFlagCount(flag *domain.Flag) error {
	err := s.repo.UpdateFlagCount(flag)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) DeleteByID(id primitive.ObjectID, ctx context.Context, username string) error {
	err := s.repo.DeleteByID(id, ctx, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) FollowUser(username string, currentUser string) error {
	err := s.repo.FollowUser(username, currentUser)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UnfollowUser(username string, currentUser string) error {
	err := s.repo.UnfollowUser(username, currentUser)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) BlockUser(id primitive.ObjectID, username string, ctx context.Context, currentUsername string) error {
	err := s.repo.BlockUser(id, username, ctx, currentUsername)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UnblockUser(id primitive.ObjectID, username string, ctx context.Context, currentUsername string) error {
	err := s.repo.UnblockUser(id, username, ctx, currentUsername)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}