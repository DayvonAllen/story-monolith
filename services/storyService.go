package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/domain"
	"story-app-monolith/repo"
)

type StoryService interface {
	Create(dto *domain.CreateStoryDto) error
	UpdateById(primitive.ObjectID, string, string, string, *domain.Tag, bool) error
	FindAll(string, bool) (*domain.StoryList, error)
	FeaturedStories() (*[]domain.FeaturedStoryDto, error)
	LikeStoryById(primitive.ObjectID, string) error
	DisLikeStoryById(primitive.ObjectID, string) error
	FindById(primitive.ObjectID, string) (*domain.StoryDto, error)
	UpdateFlagCount(flag *domain.Flag) error
	DeleteById(primitive.ObjectID, string) error
}

type DefaultStoryService struct {
	repo repo.StoryRepo
}

func (s DefaultStoryService) Create(story *domain.CreateStoryDto) error {
	err := s.repo.Create(story)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) UpdateById(id primitive.ObjectID, newContent string, newTitle string, username string, tag *domain.Tag, updated bool) error {
	err := s.repo.UpdateById(id, newContent, newTitle, username, tag, updated)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) FindAll(page string, newStoriesQuery bool) (*domain.StoryList, error) {
	story, err := s.repo.FindAll(page, newStoriesQuery)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) FeaturedStories() (*[]domain.FeaturedStoryDto, error) {
	story, err := s.repo.FeaturedStories()
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) FindById(id primitive.ObjectID, username string) (*domain.StoryDto, error) {
	story, err := s.repo.FindById(id, username)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) LikeStoryById(id primitive.ObjectID, username string) error {
	err := s.repo.LikeStoryById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) DisLikeStoryById(id primitive.ObjectID, username string) error {
	err := s.repo.DisLikeStoryById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) UpdateFlagCount(flag *domain.Flag) error {
	err := s.repo.UpdateFlagCount(flag)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) DeleteById(id primitive.ObjectID, username string) error {
	err := s.repo.DeleteById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewStoryService(repository repo.StoryRepo) DefaultStoryService {
	return DefaultStoryService{repository}
}
