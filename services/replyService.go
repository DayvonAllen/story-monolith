package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
	"story-app-monolith/repo"
	"time"
)

type ReplyService interface {
	Create(comment *domain.Reply) error
	FindAllRepliesByResourceId(id primitive.ObjectID, username string) (*[]domain.Reply, error)
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error
	LikeReplyById(primitive.ObjectID, string) error
	DisLikeReplyById(primitive.ObjectID, string) error
	UpdateFlagCount(flag *domain.Flag) error
	DeleteById(id primitive.ObjectID, username string) error
}

type DefaultReplyService struct {
	repo repo.ReplyRepo
}

func (r DefaultReplyService) Create(comment *domain.Reply) error {
	err := r.repo.Create(comment)
	if err != nil {
		return err
	}
	return nil
}

func (r DefaultReplyService) FindAllRepliesByResourceId(id primitive.ObjectID,  username string) (*[]domain.Reply, error) {
	comment, err := r.repo.FindAllRepliesByResourceId(id, username)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r DefaultReplyService) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error {
	err := r.repo.UpdateById(id, newContent, edited, updatedTime, username)
	if err != nil {
		return err
	}
	return nil
}

func (r DefaultReplyService) LikeReplyById(id primitive.ObjectID, username string) error {
	err := r.repo.LikeReplyById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (r DefaultReplyService) DisLikeReplyById(id primitive.ObjectID, username string) error {
	err := r.repo.DisLikeReplyById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (r DefaultReplyService) UpdateFlagCount(flag *domain.Flag) error {
	err := r.repo.UpdateFlagCount(flag)
	if err != nil {
		return err
	}
	return nil
}

func (r DefaultReplyService) DeleteById(id primitive.ObjectID, username string) error {
	err := r.repo.DeleteById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewReplyService(repository repo.ReplyRepo) DefaultReplyService {
	return DefaultReplyService{repository}
}
