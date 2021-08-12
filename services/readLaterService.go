package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
	"story-app-monolith/repo"
)

type ReadLaterService interface {
	Create(username string, storyId primitive.ObjectID) error
	GetByUsername(username string, page string) (*domain.ReadLaterDto, error)
	Delete(id primitive.ObjectID, username string) error
}

type DefaultReadLaterService struct {
	repo repo.ReadLaterRepo
}

func (s DefaultReadLaterService) Create(username string, storyId primitive.ObjectID) error {
	err := s.repo.Create(username, storyId)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultReadLaterService) GetByUsername(username string, page string) (*domain.ReadLaterDto, error) {
	readLaterItems, err := s.repo.GetByUsername(username, page)
	if err != nil {
		return nil, err
	}
	return readLaterItems, nil
}

func (s DefaultReadLaterService) Delete(id primitive.ObjectID, username string) error {
	err := s.repo.Delete(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewReadLaterService(repository repo.ReadLaterRepo) DefaultReadLaterService {
	return DefaultReadLaterService{repository}
}